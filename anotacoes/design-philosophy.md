Não fazer coisas fáceis de fazer, mas sim fáceis de entender.
Não adicionar complexidade até que seja necessário.
    Ex.: não começar com microservices
Tudo que fazemos deve ser preciso

Pacotes e as camadas permitem criar barreiras (firewalls)
entre o código e sermos específicos em termos de domínio (DDD)

Logging é a janela para a saúde dos serviços, a janela para os 
problemas que estão acontecendo, a primeira chance de identificar
problemas e consertar. Utilizar debuggers em desenvolvimento pode 
tornar a vida em produção mais difícil, pois em produção utilizar
debuggers é mais complicado. Se não somos capazes de identificar
problemas utilizando logs, então talvez tenhamos que refatorar.

Debuggers não acham bugs, apenas executam eles devagar.

Logs são para serem consumidos por humanos

Toda abstração adiciona complexidade

Rob Pike - Não fazemos design com interfaces, nós as descobrimos

Depois de descobrirmos abstrações, adicionamos
