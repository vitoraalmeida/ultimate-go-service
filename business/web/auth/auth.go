// Package auth fornece suporte a autenticação e autorização
// Authentication: Você é quem alega ser
// Authorization:  Você tem permissão para fazer o que deseja fazer
package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/rego"
	"github.com/vitoraalmeida/service/business/core/user"
	"go.uber.org/zap"
)

// ErrForbidden é retornado quando há problemas de autenticação ou autorização
var ErrForbidden = errors.New("attempted action is not allowed")

// Claims representa o conjunto de alegações (claims) que são trasmitidos no JWT
type Claims struct {
	jwt.RegisteredClaims
	Roles []user.Role `json:"roles"`
}

// KeyLookup declara um conjunto de metodos do comportamento de buscar por chaves
// publicas e privadas para uso com JWT. O retorno pode ser uma string codificada
// em PEM ou JWS
// Interface usada para que possamos ter um keystore que possa ser implementado
// de diferentes formas (em memória, banco de dados etc)
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

// Config representa informação necessáira para construir um objeto Auth
type Config struct {
	Log       *zap.SugaredLogger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth usado para autenticar clientes. Pode gerar tokens para um conjunto de
// claims e recriar o claims com base num token
type Auth struct {
	log       *zap.SugaredLogger
	keyLookup KeyLookup // o objeto responsável por consultar o armazenamento de chaves
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string // quem gerou o token
	mu        sync.RWMutex
	cache     map[string]string
}

// New constrói um objeto Auth para autenticação e autorização
func New(cfg Config) (*Auth, error) {
	a := Auth{
		log:       cfg.Log,
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
		cache:     make(map[string]string),
	}

	return &a, nil
}

// GenerateToken gera um token baseado num conjunto de claims
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	// cria o token usando o pacote jwt passando o método de assinatura (ex RSA)
	token := jwt.NewWithClaims(a.method, claims)
	// injeta no token um header "kid" -> key id
	token.Header["kid"] = kid

	// verifica se existe uma chave privada equivalente ao kid no armazenamento
	privateKeyPEM, err := a.keyLookup.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("private key: %w", err)
	}

	// recupera a chave RSA que gerou aquela chave PEM
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("parsing private pem: %w", err)
	}

	// utiliza a chave RSA para assinar o token
	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}

// Authenticate processa o token passado e valida
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	// faz o parse do texto "Bearer <token>"
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	// Faz o parsing do token em um claims
	var claims Claims
	token, _, err := a.parser.ParseUnverified(parts[1], &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	// Checa se o token é valido usando OPA
	kidRaw, exists := token.Header["kid"]
	if !exists {
		return Claims{}, fmt.Errorf("kid missing from header: %w", err)
	}

	kid, ok := kidRaw.(string)
	if !ok {
		return Claims{}, fmt.Errorf("kid malformed: %w", err)
	}

	pem, err := a.publicKeyLookup(kid)
	if err != nil {
		return Claims{}, fmt.Errorf("failed to fetch public key: %w", err)
	}

	input := map[string]any{
		"Key":   pem,
		"Token": parts[1],
		"ISS":   a.issuer,
	}

	if err := a.opaPolicyEvaluation(ctx, opaAuthentication, RuleAuthenticate, input); err != nil {
		return Claims{}, fmt.Errorf("authentication failed : %w", err)
	}

	// Check the database for this user to verify they are still enabled.

	return claims, nil
}

// Authorize tenta autorizar o usuário baseado num determinado papel/responsabilidade
// comparando com o claims que foi passado no token. Se as roles passadas no
// claims não forem compatíveis com a role em questão, retorna erro
func (a *Auth) Authorize(ctx context.Context, claims Claims, rule string) error {
	input := map[string]any{
		"Roles":   claims.Roles,
		"Subject": claims.Subject,
		"UserID":  claims.Subject,
	}

	if err := a.opaPolicyEvaluation(ctx, opaAuthorization, rule, input); err != nil {
		return fmt.Errorf("rego evaluation failed : %w", err)
	}

	return nil
}

// =============================================================================

// publicKeyLookup busca a publickey relativa ao kid passado
func (a *Auth) publicKeyLookup(kid string) (string, error) {
	// verifica primeiro no cache
	pem, err := func() (string, error) {
		a.mu.RLock()
		defer a.mu.RUnlock()

		pem, exists := a.cache[kid]
		if !exists {
			return "", errors.New("not found")
		}
		return pem, nil
	}()
	if err == nil {
		return pem, nil
	}

	pem, err = a.keyLookup.PublicKey(kid)
	if err != nil {
		return "", fmt.Errorf("fetching public key: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[kid] = pem

	return pem, nil
}

// opaPolicyEvaluation asks opa to evaulate the token against the specified token
// policy and public key.
func (a *Auth) opaPolicyEvaluation(ctx context.Context, opaPolicy string, rule string, input any) error {
	// rule definida no script, indicamos qual regra queremos
	// validar
	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

	// registra a validação que queremos fazer
	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", opaPolicy),
	).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	// o input consiste na chave que possui a informação verdadeira
	// e no token que está alegando ser correto
	// a validação vai comparar o token após ser processado com
	// o token que supostamente o gerou
	// ex. autenticação :
	// input := map[string]any{
	// 	"Key":   pem,
	// 	"Token": tokenString,
	// 	"ISS":   issuer,
	// }
	// ex. autorização
	//input := map[string]any{
	//	"Roles":   []string{"ADMIN"},
	//	"Subject": "1234567",
	//	"UserID":  "1234567",
	//}

	// checa se a execução da validação gerou algum resultado valido
	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	// os resultados da validação são booleanos
	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
	}

	return nil
}
