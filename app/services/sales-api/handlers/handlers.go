package handlers

import (
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/vitoraalmeida/service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/vitoraalmeida/service/app/services/sales-api/handlers/v1/usergrp"
	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/core/user/stores/userdb"
	"github.com/vitoraalmeida/service/business/web/auth"
	"github.com/vitoraalmeida/service/business/web/v1/mid"
	"github.com/vitoraalmeida/service/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contém o que é necessário para consturir um handler
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth // Objeto que armazena informções referentes à autenticação
	DB       *sqlx.DB
}

// APIMux contrói um mux ( que implementa http.Handler) com todas as rotas
// definidas para a aplicação.
// web.App é um envoltório para o mux httptreemux.ContextMux, adicionando
// possível configurações e contexto que precisarmos
// o httptreemux.ContextMux implementa o http.Handler.
func APIMux(cfg APIMuxConfig) *web.App {
	// registramos o middleware de logs e erros em toda a aplicação,
	// ou seja, todo handler a ser executado ocorrerá depois de passar
	// pelo middleware de logs e depois de erros, de forma que se o handler
	// retornar um erro, será lidado pelo mid de erros
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	// Registra um handleFunc que irá prcessar requisições get em /test
	app.Handle(http.MethodGet, "/test", testgrp.Test)
	// Registra handler para testar autenticação
	app.Handle(http.MethodGet, "/test/auth", testgrp.Test, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))

	// -------------------------------------------------------------------------

	usrCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))

	ugh := usergrp.New(usrCore)

	app.Handle(http.MethodGet, "/users", ugh.Query)

	// o objeto App implementa a internface http.Handler que é necessário para
	// construir um http.Server
	return app
}
