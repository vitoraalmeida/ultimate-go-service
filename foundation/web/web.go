package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// Tipo que lida com http requests no nosso framework.
// Sobrescreve o que é um Handler padrão, adicionando o contexto
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App é o entrypoint para nossa aplicação web. Configura o objeto de contexto
// para cada handler e configurações que os handler precisam ter
type App struct {
	// Embutindo o tipo ContextMux como ponteiro, ou seja, é um ContextMux
	// concreto e o App vai possuir todos os seus campos internos e métodos
	// consequentemente implementando qualquer interface que ContextMux
	// implementa
	// Dessa forma, um App é um App e não um ContextMux como ocorre com herança
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// Cria e retorna uma instância de App
// Usar semântica de ponteiro quando estamos lidando com uma API, com algo
// que deve compartilhar estados e recursos
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	// usando
	return &App{
		// httptreemux.NewContextMux() retorna um ponteiro para ContextMux
		// NewContextMux retorna um mux que implementa http.Handler
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// Handle atribui um handler function para uma chamada de determindado método
// em determinado endpoint, utilizando a lógica de roteamento do
// httptreemux.ContextMux internamente.
// Engloba a requisição que será processada com as infos  e capacidades
// a mais que queremos
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// se houver middlewares específico para a rota, engloba
	handler = wrapMiddleware(mw, handler)
	// Mas os middlewares para a aplicação inteira devem ser chamados
	// primeiro, então irão ficar na camada mais externa
	handler = wrapMiddleware(a.mw, handler)

	// é uma função anônima que obedece ao contrate de uma handlerFunc
	// que o httptreemux usa para registrar para uma rota, porém por dentro
	// o que é chamado é nosso Handle customizado que utiliza um contexto
	// que poderemos usar para adicionar camadas ao redor da lógica
	h := func(w http.ResponseWriter, r *http.Request) {
		// pode executar qualquer código antes de chamar o handler final
		// ex.: verificar autenticação, criar um log da requisição etc

		// cria os valores que serão passados no contexto da requisição
		v := Values{
			// As requisições terão um ID único para identificarmos
			// todas as açẽos e passos que fizeram parte do processo
			TraceID: uuid.NewString(),
			// o tempo em que aquela requisição começou para compararmos
			// quando finalizar e termos o tempo total que levou
			Now: time.Now().UTC(),
		}
		// cria o contexto reaproveitando o contexto do request e adicionando nosso
		// dado para os logs
		ctx := context.WithValue(r.Context(), key, &v)

		// chama a cadeia de funções
		if err := handler(ctx, w, r); err != nil {
			fmt.Println(err)
			return
		}

		// pode executar qualquer código depois do handler
		// logs etc
	}

	// como h tem a assinatura que o ContextMux.Handle espera, podemos usar
	// a lógica já implementada de roteamento do httptreemux, porém h possui
	// aa mecânica que usa contexto que desejamos internamente
	a.ContextMux.Handle(method, path, h)
}
