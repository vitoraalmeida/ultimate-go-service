Devemos criar pacotes que fornecem, não que contenham.

Deve ser muito claro o que pacotes fornecem. Pacotes devem fornecer uma
API. Pacotes são firewalls no seu projeto. 

Pacotes como "utils", "helpers" em Go são um mal sinal. 

Toda linha de código que criamos faz 3 coisas: aloca memória, lê memória ou 
escreve na memória.

Toda função que escrevemos realizada transformação de dados (recebe, transforma
e retorna)

Toda API, todo problema que resolvemos com código é um problema de transformação
de dados

Pacotes definem uma API que cria um firewall ao redor de si e a única forma de 
se expor é atravez daquela API publica.

Um sistema de tipos proporciona 2 coisas: que dados entrem e saiam da API

Cada API deve ter seu type system e não devemos criar pacotes que contenham
tipos comuns.

Uma API pode escolher receber um tipo concreto (dados baseados no que são) ou
um tipo abstrato (baseado no que pode fazer/comportamento - polimorfismo).

Polimorfismo - Um pedaço de código muda seu comportamento com base no dado
concreto em que está operando.

Sempre retornamos ao dado concreto, pois apenas dados concretos podem ser
construídos, manipulados e transformados. Interfaces não são reais.

Então o se queremos usar uma interface no pacote, o pacote deve também definir
a interface, pois o pacote define a API e como ele usa os dados que chegam e
saem.

Não queremos pacotes que "contenham" (utils, helpersm type system) pois acarreta na 
chance de que muitos outros pacotes dependam dele e mudanças nele geram muitos efeitos 
colaterais.

Um cheiro de que um pacote é um container é quando não faz sentido ter um 
arquivo nome-do-pacote.go dentro do pacote
