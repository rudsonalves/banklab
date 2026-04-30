# Testes e Estratégia de Validação

Este capítulo descreve como a API é validada por meio de testes e por que a estratégia de testes acompanha a arquitetura da aplicação. O objetivo é ajudar novos desenvolvedores a entender que tipos de teste existem, o que cada um valida e como essa distribuição contribui para segurança de evolução do sistema.

## Por que a estratégia de testes acompanha a arquitetura

A aplicação está organizada em camadas e módulos com responsabilidades distintas. Como consequência, a validação também tende a ser mais eficaz quando acompanha essa separação.

Em vez de depender apenas de testes de ponta a ponta, o projeto distribui a verificação entre diferentes níveis, cada um adequado a uma pergunta específica.

Isso permite testar:

- regras puras de negócio;
- coordenação dos casos de uso;
- adaptação HTTP;
- persistência e integração real com banco;
- wiring e autorização em cenários concretos.

Essa abordagem traz mais precisão, reduz custo de execução e melhora a capacidade de localizar problemas quando um comportamento se desvia do esperado.

## Principais níveis de teste

No estágio atual, a aplicação já trabalha com uma combinação de níveis de teste, incluindo:

- testes de domínio;
- testes de aplicação;
- testes de delivery;
- testes de integração.

Cada um desses níveis responde a um tipo diferente de pergunta.

## Testes de domínio

Os testes de domínio verificam regras puras e invariantes centrais do sistema.

Eles costumam ser os mais simples de executar, porque não dependem de servidor HTTP nem de banco de dados real. Seu foco é validar se o modelo de negócio se comporta corretamente quando submetido a diferentes entradas.

Exemplos de comportamento que se encaixam bem nesse nível:

- validação de regras de conta;
- rejeição de valores inválidos;
- saldo insuficiente para saque;
- conta inativa;
- transferência para a mesma conta;
- invariantes de criação de entidades.

Esses testes são importantes porque isolam o significado do domínio de detalhes de protocolo e infraestrutura.

## Testes de aplicação

Os testes de aplicação verificam os casos de uso.

Aqui o foco não é o HTTP nem a query SQL exata, mas a coordenação da operação. Esse tipo de teste ajuda a confirmar que o fluxo do caso de uso está correto.

Perguntas comuns respondidas por esse nível:

- o caso de uso exige autenticação ou contexto adequado?
- ele valida ownership corretamente?
- ele usa os repositórios necessários?
- ele devolve o erro correto no cenário esperado?
- ele abre transação quando precisa?
- ele combina corretamente regras e dependências?

Exemplos relevantes no projeto incluem testes de:

- registro de usuário;
- aprovação de usuário;
- criação de conta;
- listagem de contas;
- consulta de saldo;
- depósito;
- saque;
- transferência;
- extrato;
- políticas de acesso.

Normalmente, esse tipo de teste usa mocks, stubs ou implementações falsas para representar repositórios e dependências.

## Testes de delivery

Os testes de delivery validam a camada HTTP.

O objetivo aqui é verificar se a API está interpretando e devolvendo corretamente o contrato REST.

Esses testes costumam cobrir:

- leitura de path params;
- leitura de query params;
- decodificação de body JSON;
- validação de formato;
- extração do usuário autenticado do contexto;
- mapeamento de erro para status code;
- montagem do envelope de resposta.

Esse nível é importante porque o handler é a borda visível da aplicação. Um fluxo interno pode estar correto do ponto de vista de negócio, mas ainda assim falhar em contrato externo se o handler ler mal a entrada ou devolver a resposta errada.

Exemplos úteis nesse nível:

- rejeição de UUID inválido;
- rejeição de body malformado;
- rejeição de query param inesperado;
- retorno de `403` em caso de acesso proibido;
- retorno de `404` em caso de recurso não encontrado;
- retorno de envelope padronizado.

## Testes de integração

Os testes de integração verificam a interação entre partes concretas do sistema.

Seu foco é validar se componentes reais trabalham corretamente em conjunto, especialmente quando isso depende de:

- PostgreSQL;
- queries SQL;
- transações;
- locks;
- mapeamento de erro do banco;
- middleware de autenticação;
- wiring da aplicação.

Esse tipo de teste é importante porque mocks não conseguem provar tudo. Um caso de uso pode parecer correto contra um repositório falso, mas ainda falhar quando encontra a implementação real de SQL, constraints, tipos persistidos ou concorrência de banco.

Por isso, integração complementa os outros níveis.

## A utilidade prática de cada nível

Cada nível de teste é mais adequado para certos tipos de problema.

### Quando usar teste de domínio

Quando a dúvida principal é:

- a regra de negócio está correta?
- o valor é aceito ou rejeitado como deveria?
- a entidade preserva seus invariantes?

### Quando usar teste de aplicação

Quando a dúvida principal é:

- o caso de uso coordena corretamente a operação?
- a autorização foi aplicada?
- o fluxo chama as dependências na ordem e com os critérios esperados?

### Quando usar teste de delivery

Quando a dúvida principal é:

- a rota está lendo a requisição corretamente?
- o status code está correto?
- o envelope da resposta está conforme o contrato?

### Quando usar teste de integração

Quando a dúvida principal é:

- a persistência real funciona como esperado?
- a query SQL devolve o comportamento correto?
- a transação e os locks estão preservando consistência?
- o sistema real está ligado corretamente de ponta a ponta?

## Por que não testar tudo apenas por endpoint

Uma abordagem comum em projetos menores é concentrar tudo em testes de endpoint. Embora isso pareça mais próximo do comportamento “real”, ela costuma trazer desvantagens:

- execução mais lenta;
- maior dificuldade de localizar a origem do problema;
- custo maior para validar pequenas regras;
- dependência excessiva do ambiente completo.

Na API atual, isso seria especialmente ruim porque muitas regras de negócio podem ser validadas melhor e mais rapidamente em níveis mais próximos do domínio e da aplicação.

Isso não elimina a importância dos testes de endpoint ou integração, mas mostra que eles não precisam carregar toda a responsabilidade da validação.

## Relação entre testes e manutenção

Para novos desenvolvedores, entender a estratégia de testes é tão importante quanto entender a arquitetura.

Ao alterar uma funcionalidade, vale perguntar:

- esta mudança altera uma regra de domínio?
- esta mudança altera o fluxo do caso de uso?
- esta mudança altera o contrato HTTP?
- esta mudança altera SQL ou comportamento persistido?

Essas perguntas ajudam a identificar em que nível de teste a alteração deve ser coberta.

Também ajudam a evitar um erro comum: adicionar apenas um teste distante demais da mudança e deixar desprotegido o nível em que a regra realmente vive.

## Testes como proteção contra erosão arquitetural

Os testes não servem apenas para detectar bugs funcionais. Eles também ajudam a preservar a arquitetura.

Por exemplo:

- testes de domínio reforçam que regras importantes vivem fora de HTTP;
- testes de aplicação reforçam que casos de uso coordenam operações;
- testes de delivery reforçam o contrato REST;
- testes de integração reforçam a coerência entre arquitetura e persistência real.

Quando um projeto possui essa distribuição, fica mais difícil que a erosão arquitetural ocorra silenciosamente.

## Execução e rotina de validação

No fluxo de desenvolvimento, testes também ajudam a organizar a rotina de trabalho.

A aplicação já possui comandos próprios para execução da suíte, incluindo:

```bash
make api-test
```

e, diretamente dentro da pasta da API:

```bash
go test ./...
```

Dependendo da mudança, é comum executar subconjuntos mais focados durante a implementação e depois validar a suíte mais ampla.

## O que novos desenvolvedores devem preservar

Ao contribuir com o projeto, é importante manter alguns princípios:

- novas regras devem vir acompanhadas de testes no nível certo;
- novos endpoints devem ter cobertura de contrato HTTP;
- mudanças de persistência sensíveis devem considerar integração;
- fluxos críticos financeiros merecem validação reforçada;
- testes não devem ser tratados apenas como verificação final, mas como parte da forma de manter a coerência do sistema.

## Síntese

A estratégia de testes da API acompanha a arquitetura em camadas da aplicação.

Domínio, aplicação, delivery e integração são validados em níveis diferentes porque cada um responde a perguntas diferentes sobre o comportamento do sistema. Essa distribuição melhora precisão, acelera feedback, facilita manutenção e reduz o risco de quebrar regras importantes ao evoluir o código.

Para novos desenvolvedores, isso significa que entender onde uma regra vive no código ajuda também a entender onde ela deve ser validada. Arquitetura e testes, neste projeto, caminham juntos.
