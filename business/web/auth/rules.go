package auth

import (
	_ "embed"
)

// É o conjunto de regras que temos para autenticação/autorização
// que estão nos scripts rego (usados para aplicar OPA usando go)
// OPA = Open Policy Agent = Uma forma de definir políticas de forma padronizada
const (
	RuleAuthenticate   = "auth"
	RuleAny            = "ruleAny"
	RuleAdminOnly      = "ruleAdminOnly"
	RuleUserOnly       = "ruleUserOnly"
	RuleAdminOrSubject = "ruleAdminOrSubject"
)

// Nome do pacote definido nos arquivos rego
// usado para localizar na estrutura de pacotes do rego a validação que
// queremos executar
const (
	opaPackage string = "vitor.rego"
)

// Embute no sistema de arquivos do binário os scripts de validação
var (
	//go:embed rego/authentication.rego
	opaAuthentication string

	//go:embed rego/authorization.rego
	opaAuthorization string
)
