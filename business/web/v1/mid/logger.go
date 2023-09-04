package mid

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vitoraalmeida/service/foundation/web"
	"go.uber.org/zap"
)

// Logger adiciona a capacidade de realizar logging antes e depois do processamento
// de uma requisição
func Logger(log *zap.SugaredLogger) web.Middleware {
	// cria o middleware no formato esperado pelo mux (implementa handleFunc)
	m := func(handler web.Handler) web.Handler {
		// cria o handler que é executado pelo nosso framework
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Executa a lógica para coletar os dados necessários para logging
			v := web.GetValues(ctx)

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Infow("request started", "trace_id", v.TraceID, "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr)

			// executa de fato o handler específico da requisição
			err := handler(ctx, w, r)

			// loga o fim do request
			log.Infow("request completed", "trace_id", v.TraceID, "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))
			// retorna o erro para ser tratado por quem deve tratar
			return err
		}

		return h
	}

	return m
}
