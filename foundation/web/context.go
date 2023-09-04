// Definição do contexto que será passado em cada requisição
package web

import (
	"context"
	"time"
)

type ctxKey int

// identificador para os valores que ficam no context
// é a chave que passamod para o objeto do pacote context
// para recuperar o valor que inserimos
const key ctxKey = 1

// Values reprensenta o estado de cada requisição
type Values struct {
	TraceID    string // um identificador unico para aquela requisição
	Now        time.Time
	StatusCode int
}

// GetValues retorna o valor atual do contexto
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return time.Now()
	}
	return v.Now
}

func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}
