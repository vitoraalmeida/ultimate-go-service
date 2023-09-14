// Package database provê funções para acessar o banco de dados de diferentes
// formas (queries, configuração, transactions), usando sqlx e pgx como driver
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/vitoraalmeida/service/foundation/web"
	"go.uber.org/zap"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

// Conjunto de erros para operações CRUD
var (
	ErrDBNotFound        = sql.ErrNoRows
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrUndefinedTable    = errors.New("undefined table")
)

// Config representa propriedades necessáiras para usar o banco de dados
type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	Schema       string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

// Open abre uma conexão com base em Config
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	if cfg.Schema != "" {
		q.Set("search_path", cfg.Schema)
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

// StatusCheck retorna nil se conseguiu falar com o bd. Retorna erro caso contrário
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Executa uma query simples para testar a conexão
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// WithinTran executa a função passada em uma transaction e realiza o commit ou
// rollback no fim
func WithinTran(ctx context.Context, log *zap.SugaredLogger, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	traceID := web.GetTraceID(ctx)

	log.Infow("begin tran")
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tran: %w", err)
	}

	// Podemos fazer o defer do rollback, pois o código checa se a transactoin
	// já passou por um commit
	defer func() {
		if err := tx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				return
			}
			log.Errorw("unable to rollback tran", "trace_id", traceID, "ERROR", err)
		}
		log.Infow("rollback tran", "trace_id", traceID)
	}()

	if err := fn(tx); err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == uniqueViolation {
			return ErrDBDuplicatedEntry
		}
		return fmt.Errorf("exec tran: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tran: %w", err)
	}
	log.Infow("commit tran", "trace_id", traceID)

	return nil
}

// ExecContext função helper para executar operações CUD com logging e tracing
func ExecContext(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string) error {
	return NamedExecContext(ctx, log, db, query, struct{}{})
}

// ExecContext função helper para executar operações CUD com logging e tracing
// que necessitam de substituição de campos
func NamedExecContext(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any) error {
	q := queryString(query, data)

	if _, ok := data.(struct{}); ok {
		log.WithOptions(zap.AddCallerSkip(3)).Infow("database.NamedExecContext", "trace_id", web.GetTraceID(ctx), "query", q)
	} else {
		log.WithOptions(zap.AddCallerSkip(2)).Infow("database.NamedExecContext", "trace_id", web.GetTraceID(ctx), "query", q)
	}

	if _, err := sqlx.NamedExecContext(ctx, db, query, data); err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok {
			switch pqerr.Code {
			case undefinedTable:
				return ErrUndefinedTable
			case uniqueViolation:
				return ErrDBDuplicatedEntry
			}
		}
		return err
	}

	return nil
}

// QuerySlice função para executar queries que retornam uma coleção de dados para
// para ser convertido num slice
func QuerySlice[T any](ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, dest *[]T) error {
	// passa scruct vazio, pois a função aceita dados em que os campos podem ser substituidos
	return namedQuerySlice(ctx, log, db, query, struct{}{}, dest, false)
}

// QuerySlice função para executar queries que retornam uma coleção de dados para
// para ser convertido num slice para dados que podem ter campos substituídos
func NamedQuerySlice[T any](ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any, dest *[]T) error {
	return namedQuerySlice(ctx, log, db, query, data, dest, false)
}

// NamedQuerySliceUsingIn is a helper function for executing queries that return
// a collection of data to be unmarshalled into a slice where field replacement
// is necessary. Use this if the query has an IN clause.
func NamedQuerySliceUsingIn[T any](ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any, dest *[]T) error {
	return namedQuerySlice(ctx, log, db, query, data, dest, true)
}

func namedQuerySlice[T any](ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any, dest *[]T, withIn bool) error {
	q := queryString(query, data)

	log.WithOptions(zap.AddCallerSkip(3)).Infow("database.NamedQuerySlice", "trace_id", web.GetTraceID(ctx), "query", q)

	var rows *sqlx.Rows
	var err error

	switch withIn {
	case true:
		rows, err = func() (*sqlx.Rows, error) {
			named, args, err := sqlx.Named(query, data)
			if err != nil {
				return nil, err
			}

			query, args, err := sqlx.In(named, args...)
			if err != nil {
				return nil, err
			}

			query = db.Rebind(query)
			return db.QueryxContext(ctx, query, args...)
		}()

	default:
		rows, err = sqlx.NamedQueryContext(ctx, db, query, data)
	}

	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == undefinedTable {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()

	var slice []T
	for rows.Next() {
		v := new(T)
		if err := rows.StructScan(v); err != nil {
			return err
		}
		slice = append(slice, *v)
	}
	*dest = slice

	return nil
}

// QueryStruct função para executar queries que retornam um único valor para
// ser convertido num struct
func QueryStruct(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, dest any) error {
	return namedQueryStruct(ctx, log, db, query, struct{}{}, dest, false)
}

// QueryStruct função para executar queries que retornam um único valor para
// ser convertido num struct em que substituição de campos é necessário
func NamedQueryStruct(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any, dest any) error {
	return namedQueryStruct(ctx, log, db, query, data, dest, false)
}

// NamedQueryStructUsingIn is a helper function for executing queries that return
// a single value to be unmarshalled into a struct type where field replacement
// is necessary. Use this if the query has an IN clause.
func NamedQueryStructUsingIn(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any, dest any) error {
	return namedQueryStruct(ctx, log, db, query, data, dest, true)
}

func namedQueryStruct(ctx context.Context, log *zap.SugaredLogger, db sqlx.ExtContext, query string, data any, dest any, withIn bool) error {
	q := queryString(query, data)

	log.WithOptions(zap.AddCallerSkip(3)).Infow("database.NamedQueryStruct", "trace_id", web.GetTraceID(ctx), "query", q)

	var rows *sqlx.Rows
	var err error

	switch withIn {
	case true:
		rows, err = func() (*sqlx.Rows, error) {
			named, args, err := sqlx.Named(query, data)
			if err != nil {
				return nil, err
			}

			query, args, err := sqlx.In(named, args...)
			if err != nil {
				return nil, err
			}

			query = db.Rebind(query)
			return db.QueryxContext(ctx, query, args...)
		}()

	default:
		rows, err = sqlx.NamedQueryContext(ctx, db, query, data)
	}

	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == undefinedTable {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// queryString fornece uma versão legível das queries e parametros
func queryString(query string, args any) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case []byte:
			value = fmt.Sprintf("'%s'", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.Trim(query, " ")
}
