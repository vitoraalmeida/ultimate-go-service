Qual o propósito de logs no projeto?

Se não sabemos responder, provavelmente não estamos criando logs tão
efetivamente quanto poderíamos

Logging é a janela para a saúde dos serviços, a janela para os 
problemas que estão acontecendo, a primeira chance de identificar
problemas e consertar. Utilizar debuggers em desenvolvimento pode 
tornar a vida em produção mais difícil, pois em produção utilizar
debuggers é mais complicado. Se não somos capazes de identificar
problemas utilizando logs, então talvez tenhamos que refatorar.

Debuggers não acham bugs, apenas executam eles devagar.

Logs são para serem consumidos por humanos

Neste projeto utilizamos o zap logger (uber). Passamos o logger para
todos os locais que precisamos. Se estamos usando essa lib, vamos nos 
comprometer com ela e não vamos abstrair seu uso. Se for necessário mudar
refatoramos para aderir a outra lib (pois toda mudança acaba gerando
algumas refatorações no caso concreto).

Toda abstração adiciona complexidade
