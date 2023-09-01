O único código que deve ler configurações (seja arquivos, variáveis etc) deve
ficar em main.go e a config é passada para os componentes.

Devemos conseguir digitar help no binário do serviço e ver as opções de 
configuração para o serviço, inclusive os defaults.

As configs padrão devem poder ser sobrescritas por variáveis de ambientes ou 
flags de cli do binário.

Padrões precisam funcionar completamente no ambiente de desenvolvimento. Se 
algum dado tiver de ser fornecido de forma complementar (chaves aws etc) deve
ser claro no readme como conseguir e como adicionar.
