// Package metrics constrói as métricas que a aplicação vai rastrear
package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// Armazena a única instância de metrics que é necessáira para coletar metricas
// O único lugar em que é usado é no pacote de métricas.
// expvar é baseado em um singleton para diferentes métricas, então temos
// que usar a API que foi fornecida, que é adicionar outras métricas diretamente
// em *expvar
var m *metrics

// =============================================================================

// metrics representa o conjunto de métricas que queremos coletar específicas da
// aplicação. Os campos são seguros de ser acessados concorrentemente graças
// ao pacote expvar.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constrói o valor de metrics que será usado durante a aplicação
// como o único lugar que acessa metrics é o pacote metrics, não há muito
// problema em iniciar o valor aqui
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// =============================================================================

type ctxKey int

const key ctxKey = 1

// Set adiciona dados de metricas no context
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines refreshes the goroutine metric every 100 requests.
func AddGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}

// AddRequests incrementa o número de requests em 1
func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddErrors incrementa o número de erros em 1
func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanics incrementa o número de panics em 1
func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}
