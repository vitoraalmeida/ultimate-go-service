// Package userdb contains user related CRUD functionality.
package userdb

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/data/order"
	database "github.com/vitoraalmeida/service/business/sys/database/pgx"
	"github.com/vitoraalmeida/service/business/sys/database/pgx/dbarray"
	"go.uber.org/zap"
)

// Store gerencia o conjunt de API que usamos para interagir com o banco de dados
type Store struct {
	log *zap.SugaredLogger // para fazer logs de queries e erros
	// db sqlx.ExtContext // para casos em que queremos usar transactions
	db *sqlx.DB
	// intran bool // para cassos em que queremos usar transactions
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) WithinTran(ctx context.Context, fn func(s user.Storer) error) error {
	if s.inTran {
		return fn(s)
	}
	f := func(tx *sqlx.Tx) error {
		s := &Store{
			log:    s.log,
			db:     tx,
			inTran: true,
		}
		return fn(s)
	}
}

// Create insere um novo usuário no banco
func (s *Store) Create(ctx context.Context, usr user.User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, department, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :department, :date_created, :date_updated)`

	// o user.User é passado pela Aplicação e convertemos para o user do modelo que usamos na camada de acesso aos dados
	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update substitui o usuário no banco de dados
func (s *Store) Update(ctx context.Context, usr user.User) error {
	const q = `
	UPDATE
		users
	SET 
		"name" = :name,
		"email" = :email,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"department" = :department,
		"date_updated" = :date_updated
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	// a função de query espera um struct para saber sobre quais dados opera
	// então construimos o struct aqui para possibilitar que o usuário da função
	// possa passar o usuário inteiro
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: usr.ID.String(),
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query busca uma lista de usuários existentes no banco
func (s *Store) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]user.User, error) {
	// mapa que salva os dados que foram passados na query para que possam
	// ser logados quando a consulta for realizada
	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		users`

	buf := bytes.NewBufferString(q)
	// adiciona na query a parte do WHERE e suas clausulas conforme
	// os campos do objeto user que não forem núlos
	s.applyFilter(filter, data, buf)

	// gera a parte do ORDER BY (campo e direção [crescente/descrescente])
	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	// adiciona a parte do orderby
	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	// a lista que será recebida e preenchida pelo banco
	var dbUsrs []dbUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbUsrs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreUserSlice(dbUsrs), nil
}

// Count retorna o total de usuário no banco
func (s *Store) Count(ctx context.Context, filter user.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `
	SELECT
		count(1)
	FROM
		users`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := database.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID busca o usuário especificado
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = :user_id`

	var dbUsr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr), nil
}

// QueryByIDs busca os usuário especificados
func (s *Store) QueryByIDs(ctx context.Context, userIDs []uuid.UUID) ([]user.User, error) {
	// gera o slice de us
	ids := make([]string, len(userIDs))
	for i, userID := range userIDs {
		ids[i] = userID.String()
	}

	// possibilita passar um array para o sqlx
	data := struct {
		UserID interface {
			driver.Valuer
			sql.Scanner
		} `db:"user_id"`
	}{
		UserID: dbarray.Array(ids),
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		user_id = ANY(:user_id)`

	var usrs []dbUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &usrs); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUserSlice(usrs), nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`

	var dbUsr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toCoreUser(dbUsr), nil
}
