package handlers

import (
	"net/http"
	"os"

	"github.com/vitoraalmeida/service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/vitoraalmeida/service/business/web/v1/mid"
	"github.com/vitoraalmeida/service/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contém o que é necessário para consturir um handler
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux contrói um mux ( que implementa http.Handler) com todas as rotas
// definidas para a aplicação.
// web.App é um envoltório para o mux httptreemux.ContextMux, adicionando
// possível configurações e contexto que precisarmos
// o httptreemux.ContextMux implementa o http.Handler.
func APIMux(cfg APIMuxConfig) *web.App {
	// registramos o middleware de logs em toda a aplicação
	// ou seja, a função
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log))

	// Registra um handleFunc que irá prcessar requisições get em /test
	app.Handle(http.MethodGet, "/test", testgrp.Test)

	// o objeto App implementa a internface http.Handler que é necessário para
	// construir um http.Server
	return app
}
