package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/vitoraalmeida/service/business/web/metrics"
	"github.com/vitoraalmeida/service/foundation/web"
)

// Panics se recupera de panics e converte o panic em um erro para que possa
// ser reportado no mid de Metrics e lidado no mid Errors
// vai ser sempre posicionado como primeiro mid mais próximo do handler
// principal, do handler que de fato processa a requisição feita pelo cliente
func Panics() web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// invoca defer numa função anônima que se recupera de um panic que
			// ocorrer durante a chamada do handler
			defer func() {
				// recover retorna diferente de nil se tiver ocorrido um panic
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					// err é o nome da variável de retorno na definição de h, então
					// será o valor retornado caso ocorra um panic.
					// funções anônimas são closures (captam o escopo externo ao
					// corpo da função)
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
					// chama o pacote de metricas para incrementar o contador de
					// panics
					metrics.AddPanics(ctx)
				}
			}()

			return handler(ctx, w, r)
			// defer é executado depois que handler() retornar
		}

		return h
	}

	return m
}
