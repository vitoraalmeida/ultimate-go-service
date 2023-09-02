package web

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
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
}

// Cria e retorna uma instância de App
// Usar semântica de ponteiro quando estamos lidando com uma API, com algo
// que deve compartilhar estados e recursos
func NewApp(shutdown chan os.Signal) *App {
	// usando
	return &App{
		// httptreemux.NewContextMux() retorna um ponteiro para ContextMux
		// NewContextMux retorna um mux que implementa http.Handler
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}

// Handle atribui um handler function para uma chamada de determindado método
// em determinado endpoint, utilizando a lógica de roteamento do
// httptreemux.ContextMux internamente.
func (a *App) Handle(method string, path string, handler Handler) {
	// é uma função anônima que obedece ao contrate de uma handlerFunc
	// que o httptreemux usa para registrar para uma rota, porém por dentro
	// o que é chamado é nosso Handle customizado que utiliza um contexto
	// que poderemos usar para adicionar camadas ao redor da lógica
	h := func(w http.ResponseWriter, r *http.Request) {
		// pode executar qualquer código antes de chamar o handler final
		// ex.: verificar autenticação, criar um log da requisição etc

		// chama
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
			return
		}

		// pode executar qualquer código depois do handler
		// logs etc
	}

	// como h tem a assinatura que o ContextMux.Handler espera, podemos usar
	// a lógica já implementada de roteamento do httptreemux, porém h possui
	// aa mecânica que usa contexto que desejamos
	a.ContextMux.Handle(method, path, h)
}
