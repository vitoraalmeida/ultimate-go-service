package mid

import (
	"context"
	"net/http"

	"github.com/vitoraalmeida/service/business/sys/validate"
	"github.com/vitoraalmeida/service/business/web/auth"
	v1 "github.com/vitoraalmeida/service/business/web/v1"
	"github.com/vitoraalmeida/service/foundation/web"
	"go.uber.org/zap"
)

// Errors lida com erros que aparecem na cadeia de chamadas de funções.
// Detecta com erros de aplicação que são usados para responder ao cliente de
// forma uniforme.
// Erros inesperados (status >= 500) são loggados
func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// executamos o handler e vemos se há erro
			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				var er v1.ErrorResponse
				var status int

				// Se for um erro conhecido, buscamos esse erro e criamos um
				// erro de resposta e definimos o status code da resposta
				switch {
				// erros de validação
				case validate.IsFieldErrors(err):
					fieldErrors := validate.GetFieldErrors(err)
					er = v1.ErrorResponse{
						Error:  "data validation error",
						Fields: fieldErrors.Fields(),
					}
					status = http.StatusBadRequest
				case v1.IsRequestError(err):
					reqErr := v1.GetRequestError(err)
					er = v1.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					er = v1.ErrorResponse{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				// se for um erro inesperado, criamos o erro de resposta como
				// internalServerError (500)
				default:
					er = v1.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				// enviamos a resposta de erro
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// se temos erro ao tentar enviar a resposta ao cliente, temos
				// problema de integridade do sistema e devemos desligar
				// retornamos o erro para o handler base (definido no framework)
				// que é o mais próximo da main()
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
