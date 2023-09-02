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


