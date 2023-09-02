package handlers

import (
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/vitoraalmeida/service/app/services/sales-api/handlers/v1/testgrp"
	"go.uber.org/zap"
)

// APIMuxConfig contém o que é necessário para consturir um handler
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux contrói um mux (http.Handler) com todas as rotas definidas para a aplicação
// o httptreemux.ContextMux implementa o http.Handler. Retornamos o tipo concreto
// ao invés de retornar uma interface (http.Handler) para que o usuário decida
// o que fazer com o que a api retorna.
func APIMux(cfg APIMuxConfig) *httptreemux.ContextMux {
	// NewContextMux retorna um mux que implementa http.Handler
	mux := httptreemux.NewContextMux()

	// Registra um handleFunc que irá prcessar requisições get em /test
	mux.Handle(http.MethodGet, "/test", testgrp.Test)

	return mux
}
