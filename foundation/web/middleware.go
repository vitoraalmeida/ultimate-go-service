package web

// Middleware são função designadas para executar código antes e/ou depois de
// outro handler. Permite que não precisemos adicionar o mesmo código em todo
// lugar que desejamos que um mesmo comportamento ocorra.
// Ex.: queremos que cada requisição gere um log. sem middleware precisariamos
// adicionar o mesmo código em cada handler
// então o middleware recebe como argumento um Handler, e adiciona uma camada
// ao redor do handler.
type Middleware func(Handler) Handler

// wrapMiddleware cria um novo handler que engloba o handler final que executa
// a lógica principal (lidar com requisições em endpoints)
// Os middlewares serão executados na ordem em que forem passados, de forma que
// a função principal é executada no final
func wrapMiddleware(mw []Middleware, handler Handler) Handler {

	// Para que os middlewares sejam executados na ordem que forem registrados,
	// devemos fazer com que o primeiro middleware registrado seja a camada
	// mais externa, então iteramos de trás para frente na lista
	// assim, a função que lida com a request é a última a chamada e a
	// primeira a retornar
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
