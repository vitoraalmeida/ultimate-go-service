package mid

import (
	"context"
	"net/http"

	"github.com/vitoraalmeida/service/business/web/metrics"
	"github.com/vitoraalmeida/service/foundation/web"
)

// Metrics atualiza os contadores de metricas
func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)

			err := handler(ctx, w, r)

			metrics.AddRequests(ctx)
			// cada requisição é executada em uma goroutine diferente
			metrics.AddGoroutines(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
