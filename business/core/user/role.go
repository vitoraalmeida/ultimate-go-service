// Define as roles possíveis no sistema para um usuário
package user

import "errors"

// Papeis/responsabilidades possíveis para um usuário
// Esse padrão é o mais próximo que podemos chegar de enums em go
var (
	RoleAdmin = Role{"ADMIN"}
	RoleUser  = Role{"USER"}
)

// Conjunto dos papeis de um usuário
var roles = map[string]Role{
	RoleAdmin.name: RoleAdmin,
	RoleUser.name:  RoleUser,
}

// Role representa um papel/responsabilidade
type Role struct {
	name string // não pode ser modificado/escrito/lido = não exportado
}

// ParseRole recebe um texto e converte para uma role existente
// deve ser chamada na camada de aplicação, não na camada busines. Validação
// não deve ser feita em business
// A camada de aplicação pode aceitar strings e aqui é feita a conversão
func ParseRole(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, errors.New("invalid role")
	}

	return role, nil
}

// MustParseRole chama panic() caso ParseRole retorne erro -> usado em testes
func MustParseRole(value string) Role {
	role, err := ParseRole(value)
	if err != nil {
		panic(err)
	}

	return role
}

// Name retorna o nome da role
func (r Role) Name() string {
	return r.name
}

// UnmarshalText converte json para role
func (r *Role) UnmarshalText(data []byte) error {
	r.name = string(data)
	return nil
}

// MarshalText converte role para json
func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

// Equal provê suporte para o pacote go-cmp e testing
func (r Role) Equal(r2 Role) bool {
	return r.name == r2.name
}
