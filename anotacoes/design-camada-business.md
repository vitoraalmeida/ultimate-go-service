Se não entendemos o Dado não entendemos o problema

## Domain driven design

O banco de dados pode mudar e ser implementado de diferentes formas, mas o
design do dóminio deve resolver o problema de negócio, então faz sentido
primeiro trabalhar no problema de negócio e domínio antes de adicionar  o 
banco de dados

```
               External input
APP      -----O-O-O-----O--O-O--------O-O-O-----O-O-O---  
           +--|-|-|-+ +-|--|-|----+ +-|-|-|-+ +-|-|-|--+ APIs
           |+-----+ | | +--------+| |+-----+| |+------+|
           ||Users| | | |Products|| ||Sales|| ||Orders|| = domínios
           |+-----+ | | +--------+| |+-----+| |+------+|
Business   +--------+ +-----------+ +-------+ +--------+

           ---------------------------------------------
```

Os domínios definem quais dados compõem e representam aquele objeto real que 
estamos abstraindo

Envolvemos o domínio com uma API que permite realizar ações com os objetos
do domínio e a api recebe inputs da aplicação

Existem modelos no nível de business, que representam o domínio e modelos no 
nível da aplicação, que representam dados que são necessários para determinados
fins em cima dos domínios. Digamos que queremos um relatório de vendas e cada
registro de venda registra entre outros dados, o produto, quatidade vendida e 
o id do usuário que comprou. No relatório, também queremos o NOme do usuário, porém
não há esse dado no modelo de business de sales, apenas o ID.

Para resolver isso, a nivel de aplicação, criamos um modelo que contém o dado que
queremos para resolver o problema, e para obter esse dados interagimos com outros
domínios. Assim, em cada domínio fica somente o que faz sentido para representar
aquele domínio. Ou seja, buscamos o dado de vendas e com base no Id do usuário
vinculado, buscamos os dados do usuário no domínio de users e agregamos os dados
para gerar o relatório

Não enviamos dados da camada business diretamente, pois se alterarmos a camada 
business automaticamente quebramos a API. Ao manter modelos especificos de aplicação
evitamos que mudanças no nível de business causem diretamente quebras de contratos 
de api.

No caso de o usuário desejar também que as vendas sejam ordenadas pelo nome do
usuário. Se nosso banco estivesse com um número muito grande de registros de vendas
seria muito custoso buscar todos os registros, colocar em memória, fazer a consulta
de todos os usuários que correspondem àquele ID que está no registro de venda,
unir e depois ordenar todos. 

Podemos criar um novo domínio na camada business que 
atenda a necessidade do relatório, combinando os dados que precisamos, por exemplo,
criando um domínio SalesUsers, que é uma "view" para determinados dados

```
                            External input
APP       --O--O-O----------O-O-O-----O--O-O--------O-O-O-----O-O-O---  
          +-|--|-|------+  +--|-|-|-+ +-|--|-|----+ +-|-|-|-+ +-|-|-|--+ APIs                 
          | +----------+|  |+-----+ | | +--------+| |+-----+| |+------+|                      
          | |SalesUsers||  ||Users| | | |Products|| ||Sales|| ||Orders|| = domínios          
          | +----------+|  |+-----+ | | +--------+| |+-----+| |+------+|                      
Business  +-------------+  +--------+ +-----------+ +-------+ +--------+                      
                                                                                                  
        ----------------------------------------------------------------
```

Quando mantemos as abstrações como domínios puros, cada um reprensentando seu
objeto real de forma precisa, isolados uns dos outros de fato, também devemos
considerar que o armazenamento deles é completamente separado.

```
                            External input
APP       --O--O-O------------O-O-O-------O--O-O--------O-O-O--------O-O-O---  
          +-|--|-|------+  +--|-|-|---+ +-|--|-|----+ +-|-|-|----+ +-|-|-|-----+ APIs                 
          | +----------+|  |+-----+   | | +--------+| |+-----+   | |+------+   |                      
          | |SalesUsers||  ||Users|   | | |Products|| ||Sales|   | ||Orders|   | = domínios          
          | +----------+|  |+-----+   | | +--------+| |+-----+   | |+------+   |                      
          | +--------+  |  |+--------+| | +--------+| |+--------+| |+--------+ | 
          | | Storage|  |  || Storage|| | | Storage|| || Storage|| || Storage| | Camada de acesso ao dado
          | +--------+  |  |+--------+| | +--------+| |+--------+| |+--------+ |
          |   +--+      |  | +--+     | |  +--+     | |  +--+    | |  +--+     |
          |   |DB|      |  | |DB|     | |  |DB|     | |  |DB|    | |  |DB|     |
          |   +--+      |  | +--+     | |  +--+     | |  +--+    | |  +--+     |
Business  +-------------+  +----------+ +-----------+ +----------+ +-----------+                      
                                                                                                  
        ----------------------------------------------------------------
```

Dessa forma, ainda que não estejamos de fato usando uma instância de banco de dados
para cada domínio, como os domínios não são dependentes uns dos outros internamente
caso um dia algum deles vier a ser muito exigido ao ponto de fazer sentido criar um
serviço só para ele, poderemos simplesmente separar o serviço com toda a infra
só para ele, pois nunca assumimos que eles estivessem no mesmo banco.

Por isso é interessante que não façamos o relatório simplesmente fazendo joins
de tabelas do banco, pois eventualmente podemos precisar que estejam em bancos
separados. 

    Event driven
    Para que a tabela de SalesUsers seja populada, podemos usar a estratégia de emitir
    eventos quando uma nova compra for inserida, executando a logica de buscar o
    usuário relacionado, fazer a união dos dados necessários e inserir no SalesUsers


Para manter essa consistência e a ideia de cada domínio estar separado e evitar
joins (acoplamento) e ao mesmo tempo não precisar de um servidor de banco de
dados para cada domínio, podemos usar o mesmo banco, com tabelas diferentes
e para nosso domínio SalesUsers usar Views ( no pacote cview = core view), 
que podemos usar para relatórios, agregações etc

```
                            External input
APP       --O--O-O------------O-O-O-------O--O-O--------O-O-O--------O-O-O---  
          +-|--|-|------+  +--|-|-|---+ +-|--|-|----+ +-|-|-|----+ +-|-|-|-----+ APIs                 
          | +----------+|  |+-----+   | | +--------+| |+-----+   | |+------+   |                      
          | |SalesUsers||  ||Users|   | | |Products|| ||Sales|   | ||Orders|   | = domínios          
          | +----------+|  |+-----+   | | +--------+| |+-----+   | |+------+   |                      
          | +--------+  |  |+--------+| | +--------+| |+--------+| |+--------+ | 
          | | Storage|  |  || Storage|| | | Storage|| || Storage|| || Storage| |
          | +---|----+  |  |+----|---+| | +----|---+| |+---|----+| |+---|----+ |
Business  +-----|-------+  +-----|----+ +------|----+ +----|-----+ +----|------+                      
                +----------------+-------------+-----------+------------+
                    VIEW                       |
                                              +--+ 
                                              |DB|
                                              +--+
        ----------------------------------------------------------------
```

Assim mantemos o isolamento com flexibilidade para separar os serviços, bancos e
usar eventos caso necessário, se nos mantivermos disciplinados em manter os
domínios separados


### Validações

Para evitar que toda chamada na camada de business e foundation deva ser validada
e validações muitas vezes precisam bater no banco de dados para buscar uma info,
estabelecemos que a camada de foundation e business confiam nos dados que entram
nelas, mas para isso precisamos reforçar as validações na camada de aplicação.

Além disso, para mitigar a falta de validações fora da camada de aplicação, usamos
o sistema de tipos e métodos "parse" para garantir que os dados que entram são
do tipo que queremos

### 
