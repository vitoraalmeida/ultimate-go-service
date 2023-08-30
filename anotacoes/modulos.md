`go env`
    GOMODCACHE="/home/viandrade/go/pkg/mod"

GOMODCACHE armazena os módulos que já foram baixados anteriormente

`go mod init github.com/vitoraalmeida/service`-> cria um módulo com o nome passado

Um repositório representa um projeto, um projeto representa um módulo, um módulo
cria um namespace único para acessar código dentro do projeto

O arquivo go.mod é uma âncora do ferramental Go para toda informação relativa
a nomes de modulos

Possui o nome do módulo e a versão do Go que gerou aquele módulo.
A versão diz ao ferramental do Go a versão mínima do go que é necessária para
construir o projeto

## Adicionando dependências
Para adicionar uma dependência no projeto, podemos adicionar o endereço do
repositório no bloco de imports do programa em que estamos trabalhando,
chamar o módulo como se estivéssemos usando uma função dele (não precisa ser
algo que esteja de fato no módulo, é apenas para que o intelisense do go não
remova a dependencia que não está sendo utilizada no código) e executar o comando
go mod tidy para que o ferramental do Go busque a dependência e adicione no
nosso go.mod automaticamente, além de salvar a versão exata que foi baixada
no go.sum com um hash relativo ao projeto e ao go.mod da dependência. 
O código fonte do modulo é baixado para nossa máquina e fica no GOMODCACHE

```go
package main

import(
    "github.com/ardanlabs/conf"
)

func main() {
    conf.New()
}

```

```
$ go mod tidy
# o servidor de módulos do google busca no github (ou no vcs utilizado) se existe um repositório com esse nome
go: finding module for package github.com/ardanlabs/conf 
# se achou, faz o donwload como zip do repositório
go: downloading github.com/ardanlabs/conf vX.X.X # se não passamos uma versão explicita, busca a mais atual
# busca um módulo dentro do zip baixado
go: found github.com/ardanlabs/conf in github.com/ardanlabs/conf vX.X.X
# baixa dependências indiretas para que a direta funcione
go: downloading <dependência-indireta> vY.Y.Y
```

go.mod
```
module github.com/vitoraalmeida/service

go 1.21

require github.com/ardanlabs/conf v1.5.0 //dependência direta
```

O servidor do google que serve os módulos é o que fica na variável
GOPROXY (proxy do google). Podemos não utilizar o proxy do google (por motivos
de privacidade talvez) e buscar diretamente (GOPROXY="direct"). Assim, caso 
já tenhamos feito esse processo anteriormente utilizando o proxy do google,
será feito o download novamente, e gerado um novo hash no go.sum. Quando estamos
construindo módulos para terceiros, é importante que sejamos os primeiros a 
ir até ao servidor do proxy para que o hash que ficará salvo seja de fato do 
nosso código verificado, de forma que alterações no código sem alteração de 
versão não passem desapercebidas

É possível que hospedemos nosso próprio servidor de proxy para que para alguns
casos irmos num VCS local, e para outros casos (opcional) ir no proxy do google

Podemos configurar para que modulos publicos sejam buscados do proxy e modulos 
privados sejam buscados no VCS da empresa.
`go PRIVATE="company.gitlab.com"`
Qualquer import que começe com company.gitlab.com não passará pelo proxy do google.

Existem modulos que mantém cada major version separada em "submodulos"
ex.: github.com/ardanlabs/conf/v3. Para sabermos se é o caso, podemos ir até o 
código do módulo e ver se a versão major dele está usando essa delcaração e seu
go.mod

Para não precisar de chamadas de rede para buscar dependências, podemos utilziar
"vendoring" que é a prática de baixar todo o código da dependência e manter
como parte do código do projeto. Agora possuimos a dependência. Criamos uma pasta
"vendoring" no projeto e deixamos o código lá, também fazermos commits do vendoring
