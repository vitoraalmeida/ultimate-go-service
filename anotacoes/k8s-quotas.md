Um cluster K8s possui poder computacional que é fornecido por diversos nodes.

No caso do ambiente local (dev) rodando Kind, o cluster ter o poder computacional
dá máquina em que ele executa.

É possível configurar o serviço em go para utilizar apenas o recurso disponível
na plataforma em que ele executa. E é possível determinar os limites que os
elementos do k8s podem usar do cluster.

Normalmente um SO linux possui um time-slice de 100ms (milisegundos).
Se SO roda numa máquina (ou cluster) com 4 cores, cada core consegue atribuir 
100ms para cada processo. 

No k8s utilizamos 1000m para dizer que queremos 1 core,  então  se queremos que 
um container seja executado consumindo um core para ele, podemos definir 1000m de cpu,
o que faz com que ele possua 100% dos 100ms disponíveis que o SO disponibilizou
para o runtime

O atributo "request" da definição de recursos para um container ajuda o  k8s a 
determinar qual node consegue hospedar aquele pod em que o container vai rodar,
garantindo que a soma total dos pods não vão usar mais do que o disponível num
nó. Se temos 8 pods requisitando 1 core (1000m) cada, então um nó com 8 cores
pode executar esses 8 pods. Se temos 16 pods requisitando cada um 1/2 core (500m),
então esse nó de 8 cores pode rodar os 16 pods.

Já o "limit" ajuda o runtime de containers a determinar quais container podem
usar um CPU e por quanto tempo. Se o time-slice é de 100ms e queremos que aquele
container utilize todos os 100ms, atribuimos 1000m no "Limit". Se queremos que 
ele ocupe metade do tempo disponível, atrbuimos 500m.


Se definimos que o "limit" é igual ao "request" dizemos que o container pode 
apenas o tempo diponível. A desvantagem é que se ele começar a receber um
aumento brusco de requisições, pode ser que ele não tenha tempo suficiente e
precise esperar pelos próximos 100ms que o SO disponibilizar. Então para casos
em que nosso programa visa lidar com momentos de pico muito grander, não ter 
limites pode ser melhor. Mas caso o trafego seja constante, podemos limitar

    Nick Stogner
    Requests are separate from Limits b/c setting requests lower than limits allows
    for more efficient resource utilization when you start packing multiple PODs
    onto a given Node. Limits are there to prevent noisy neighbor conditions with
    colocated PODs.

## consumo de recursos pelo go

Se um computador tem 4 cores, o go runtime vai requisitar 1 thread para cada core
e vai executar diversas go routines em cada thread (green threads).

Uma troca de contexto para threads de sistema operacional leva 1 microssegundo,
mas para goroutines apenas 200nanossegundos (1/5). 

Se estamos fazendo apenas tarefas CPU bound (computação) não temos troca de contexto
e mantemos os cores funcionando 100% do tempo para as goroutines. Fora o SO
precisar de um core, nada mais deve afetar a execução.

Se temos apenas nosso binário rodando num sistema, dificilmente o poder computacional
será utlizado para outras coisas.

Se temos 2 programas Go executando, os 2 consideram que tem os 4 cores disponíveis,
então são 2 threads de SO em cada core, e aí terão de ocorrer trocas de contexto
entre threads.

Se podemos pagar por garbage collection (tempo para limpar memória -> menos tempo
executando o que queremos), Go é possivelmente a melhor tecnologia.

Para fazermos programas Go o mais eficientes possíveis, devemos definir que ele
deve usar tantas trheads quanto o número de cores disponíveis.
