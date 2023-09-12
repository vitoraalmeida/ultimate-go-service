package user

import "errors"

// Papeis/responsabilidades possíveis para um usuário
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
	name string
}

// ParseRole recebe um texto e converte para uma role existente
func ParseRole(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, errors.New("invalid role")
	}

	return role, nil
}

// MustParseRole chama panic() caso ParseRole retorne erro
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
