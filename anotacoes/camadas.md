O projeto utiliza camamadas de código (não grupos), de forma que o modelo 
mental do projeto seja composto de 5 camadas (estudos que mostram que conseguimos
manter 5 coisas na cabeça num mesmo momento (??)).

Podem existir subcamadas, mas também ficam restritas a 5.

Camadas (stair case)

```
app/ ----------------- referente à aplicação que será executada
|
+---+- services/
    |         |
    |         +-- metrics
    |         |
    |         +-- sales-api
    |
    +- tooling/
             |
             +-- logfmt
             |
             +-- sales-admin

business/ --------------------- regras de negócio e capacidades que são decisão
|                               do projeto específico
+---+- core/
    |
    +- cview/

foundation/ -------------------- stdlib do projeto
     +- docker                   não estão estritamente ligados
     |- vault                    às regras de negócio
     +- logger                   eventualmente poderiam estar
     |- web                      em seus próprios repos. Comportamento geral
                                 

vendor   --------------- código de dependências de terceiros

zarf    ----- zarf (capas em que colocamos os copos de reci
              pientes quentes para não nos queimarmos
              tudo relacionado a deploys, configuração de 
              deploy (docker, k8s...)
```

As camadas apenas importam camadas abaixo ou do lado, nunca acima

Pacotes e as camadas permitem criar barreiras (firewalls)
entre o código e sermos específicos em termos de domínio (DDD)



