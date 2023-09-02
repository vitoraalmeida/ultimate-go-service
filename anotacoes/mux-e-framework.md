Um servidor http no go http.Server possui um campo Handler, que é o responsável
por responser a uma requisição. Se nenhum handler específico por passado,
é utilizado o http.DefaultServeMux.

Um handler deve implementar o método ServeHTTP(ReponseWriter, \*Request), que 
escreve headers e dados no ResponseWriter e então retorna. O retorno indica que
a requisoção finalizou.

Podemos utilizar handlers personalizados que implementam a interface Handler,
atribuindo esse handler num http.Server.Handler.

Existem diversos Handlers de terceiros (mux, routers) e podemos usar qualquer um
que seja bem feito, pois a diferença é de nanosegundos

Aqui será usado o httptreemux

Idealmente, queremos que os handlers validem os dados que estão chegando na 
requisição, invoquem a camada de negócio para processar os dados, retorne
erros caso ocorram e lidem com o caso de tudo estar correto. Devemos retornar
erros, não lidar com erros dentro do handler, pois assim não deixamos o usuário
da nossa API lidar da forma que ele quer com os erros.

Para isso, precisamos que existam passos antes dos handlers que possam lidar
com os erros retornados por ele.

A implementação básica de um web server em go é:

```
                              Goroutine                    Test Handler
                                  |  +-----------+   +----------------+
              +-------+ServeHTTP  |  |  mux      |   |                |
http req ---> |listen |-----------|--|---E->/test------> Processa     |
              +-------+              +-----------+   |       |        |
                                                     +-------|--------+
           <-------------------------------------------------+ retorna
```
Porém qualquer erro que aconteça no handler terá que ser lidado no handler,
mas para uma API ser mais óbvia e clara, é melhor que retornemos erros que 
ocorram para o usuário da API lidar da forma que achar melhor, e para isso
temos que adicionar camadas no entorno do handler para que passos intermediários
possam ocorrer. Middlewares

```
  1=validar                                                  Log
  2=chamar pacote de negócios                        +-----------------------+
  3=retornar erros                                   |       Auth            |
  4=retornar OK                                      | +--------------------+|
                              Goroutine              | |     Erros          ||
                                  |  +-----------+   | | +----------------+ ||
              +-------+ServeHTTP  |  |  mux /... |   | | |   Test Handler | ||
http req ---> |listen |-----------|--|---E->/test--------->  1,2,3,4      | ||
              +-------+              |      /... |   | | +--------|-------+ ||
                                     +-----------+   | +----------|---------+|
                                                     +------------|----------+
           <---------------------------------------------------+ retorna
```

Assim podemos separar as atribuições melhor.

Para isso vamos adicionar um contexto para as requisições onde podemos adicionar
informações, mas não no contexto embutido do pacote http, pois assim fica escondido
de quem usa a API, vamos passar contexto como parte da assinatura das funções dos
handlers, para tornar obvio quando estamos usando contexto.

Nossos handlers possuem a assinatura seguinte

```go
func HandlerX(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
}
``` 

Para isso vamos adicionar funcionalidades no httptreemux, construindo um framework
ao redor dele para termos suporte a essa necessidade
