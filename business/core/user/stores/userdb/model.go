package userdb

import (
	"database/sql"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/sys/database/pgx/dbarray"
)

// dbUser representa a estrutura que precisamos para mover dados entre a aplicação
// e o banco de dados
type dbUser struct {
	ID           uuid.UUID      `db:"user_id"` // tags que informam como o campo está definido em sql
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        dbarray.String `db:"roles"` // usamos o pacote dbarray para fazer a representação de arrays do postgres
	PasswordHash []byte         `db:"password_hash"`
	Enabled      bool           `db:"enabled"`
	Department   sql.NullString `db:"department"` // Quando o dado pode ser nulo, usamos o null específico para sql
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

// converte um User de domínio em UserDb (Modelo de user para interação direta com o banco)
// para inserir dados no banco
func toDBUser(usr user.User) dbUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return dbUser{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department: sql.NullString{
			String: usr.Department,
			Valid:  usr.Department != "",
		},
		Enabled:     usr.Enabled,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateUpdated.UTC(),
	}
}

// converte de UserDB para User de domínio
// para dados que saem do banco
func toCoreUser(dbUsr dbUser) user.User {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roles := make([]user.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		roles[i] = user.MustParseRole(value)
	}

	usr := user.User{
		ID:           dbUsr.ID,
		Name:         dbUsr.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		Enabled:      dbUsr.Enabled,
		Department:   dbUsr.Department.String,
		DateCreated:  dbUsr.DateCreated.In(time.Local),
		DateUpdated:  dbUsr.DateUpdated.In(time.Local),
	}

	return usr
}

// converte o slice de dbusers que vem do banco em slice de usuários de domínio
func toCoreUserSlice(dbUsers []dbUser) []user.User {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		usrs[i] = toCoreUser(dbUsr)
	}
	return usrs
}
