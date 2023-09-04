pensar em erros no go como sinais verticais


``` 
app [function]     |    Log, Reponse, Terminate = handling
-----------------  | [ [ [e] ] ]
bus [function]     |
-----------------  | [ [e] ] 
fou [function]     | 
-----------------  | [e]
stdlib [function]  | 
                        
                       /\
                       ||
                      error 
```


O que significa lidar com um erro?

Bill Kennedy: 1- Se você está lidando com o erro, o erro é loggado. 2- Você
deve determinar se podemos nos recuperar desse erro ou se a aplicação deve
parar (ou a goroutine), ou seja, o erro (específico) morre ali. 
Se o erro é propagado, não foi lidado.


Se o erro acontece na camada de foundation, não podemos fazer mais que envolver
o erro com algum contexto e retornar para quem chamou, pois não logamos em foundation

Quando mais cedo podermos lidar com o erro, mais temos chance de poder nos 
recuperarmos daquele erro e não parar a aplicação

No caso dessa aplicação, deixamos que o middleware de erros inspecione o erro
e tome a melhor decisão sobre como lidar com ele.

Definimos tipos diferentes de erros para situações diferentes.

Se o erro não teve um tipo definido, respondemos como 500 (internal server).

Um dos tipos definido é o erro confiável (trusted error) um erro que sabemos
o porquê de existir, sabemos que não vai vazar nenhuma informação importante
e sabemos como descrever.

O outro é um erro de shutdown. Se o sistema está tendo problemas de integridade,
não deveria estar executando. Erros que ocorrem em pacotes mais fundamentais,
mais longe da aplicação final, não devem conseguir desligar a aplicação (panic,
os.Exit()), porém podem sinalizar que a aplicação deveria parar.
