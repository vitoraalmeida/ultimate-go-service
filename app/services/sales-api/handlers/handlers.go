package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
)

// APIMuxConfig contém o que é necessário para consturir um handler
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux contrói um http.Handler com todas as rotas definidas para a aplicação
func APIMux(cfg APIMuxConfig) http.Handler {
	// NewContextMux retorna um mux que implementa http.Handler
	mux := httptreemux.NewContextMux()

	// http.HandlerFunc
	h := func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string
		}{
			Status: "OK",
		}

		json.NewEncoder(w).Encode(status)
	}

	// Registra um handleFunc que irá prcessar requisições get em /test
	mux.Handle(http.MethodGet, "/test", h)

	return mux
}
