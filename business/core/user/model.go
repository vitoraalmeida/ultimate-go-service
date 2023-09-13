package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User representa informações sobre um usuário individual
type User struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address // o tipo força que a camada App deva passar o tipo que queremos
	Roles        []Role       // assim podemos mitigar a falta de validações na camada business
	PasswordHash []byte
	Department   string
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contém informação necessária para criar um usuário
type NewUser struct {
	Name            string
	Email           mail.Address
	Roles           []Role
	Department      string
	Password        string
	PasswordConfirm string
}

// UpdateUser  contém informação necessária para atualizar dados de um usuário
type UpdateUser struct {
	Name            *string       // Usamos a semântica de ponteiro aqui
	Email           *mail.Address // para mostrar que alguns desses dados
	Roles           []Role        // podem ser nulos na tentativa de atualizar
	Department      *string       // um usuário, pois podemos querer atualizar
	Password        *string       // apenas algum dos dados do usuário
	PasswordConfirm *string
	Enabled         *bool
}
