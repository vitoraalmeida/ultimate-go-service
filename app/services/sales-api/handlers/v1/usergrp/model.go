package usergrp

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/cview/user/summary"
	"github.com/vitoraalmeida/service/business/sys/validate"
)

// AppUser representa informação referente a um usuário no contexto de aplicação
// camada de armazenamento de business
// Usa a camada de business e business se encarrega de receber esse usuário
// e converter no formato que ele precsia
// Usamos tipos primitivos para que a conversão de JSON não falhe,
// porém usamos validação após a conversão em que os tipos especificos
// foram definidos
type AppUser struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Department   string   `json:"department"`
	Enabled      bool     `json:"enabled"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

// Converte um usuário de domínio em usuário de aplicação
func toAppUser(usr user.User) AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return AppUser{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

// AppNewUser contem informação necessária para criar um novo usuário
// declaramos a validação que queremos e o pacote validate (business/sys/validate)
// será usado para validar
type AppNewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

// Converte o modelo de usuário para criar novos usuários em usuários de domínio
func toCoreNewUser(app AppNewUser) (user.NewUser, error) {
	roles := make([]user.Role, len(app.Roles))
	for i, roleStr := range app.Roles {
		role, err := user.ParseRole(roleStr)
		if err != nil {
			return user.NewUser{}, fmt.Errorf("parsing role: %w", err)
		}
		roles[i] = role
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parsing email: %w", err)
	}

	usr := user.NewUser{
		Name:            app.Name,
		Email:           *addr,
		Roles:           roles,
		Department:      app.Department,
		Password:        app.Password,
		PasswordConfirm: app.PasswordConfirm,
	}

	return usr, nil
}

// Valida se as informações passadas para criar um novo usuário são validas
func (app AppNewUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

// =============================================================================

// AppUpdateUser contém informação necessáira para atualizar um usuário
type AppUpdateUser struct {
	Name            *string  `json:"name"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	Roles           []string `json:"roles"`
	Department      *string  `json:"department"`
	Password        *string  `json:"password"`
	PasswordConfirm *string  `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
	Enabled         *bool    `json:"enabled"`
}

func toCoreUpdateUser(app AppUpdateUser) (user.UpdateUser, error) {
	var roles []user.Role
	if app.Roles != nil {
		roles = make([]user.Role, len(app.Roles))
		for i, roleStr := range app.Roles {
			role, err := user.ParseRole(roleStr)
			if err != nil {
				return user.UpdateUser{}, fmt.Errorf("parsing role: %w", err)
			}
			roles[i] = role
		}
	}

	var addr *mail.Address
	if app.Email != nil {
		var err error
		addr, err = mail.ParseAddress(*app.Email)
		if err != nil {
			return user.UpdateUser{}, fmt.Errorf("parsing email: %w", err)
		}
	}

	nu := user.UpdateUser{
		Name:            app.Name,
		Email:           addr,
		Roles:           roles,
		Department:      app.Department,
		Password:        app.Password,
		PasswordConfirm: app.PasswordConfirm,
		Enabled:         app.Enabled,
	}

	return nu, nil
}

// Validate checa se as informações passadas para atualizar o usuário são válidas
func (app AppUpdateUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// =============================================================================

// AppSummary representa informação sobre um usuário e seus produtos relacionados
type AppSummary struct {
	UserID     string  `json:"userID"`
	UserName   string  `json:"userName"`
	TotalCount int     `json:"totalCount"`
	TotalCost  float64 `json:"totalCost"`
}

func toAppSummary(smm summary.Summary) AppSummary {
	return AppSummary{
		UserID:     smm.UserID.String(),
		UserName:   smm.UserName,
		TotalCount: smm.TotalCount,
		TotalCost:  smm.TotalCost,
	}
}
