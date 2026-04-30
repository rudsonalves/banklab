# Anatomia Interna dos Módulos

Este capítulo detalha como cada módulo da API é organizado internamente. Depois de compreender o mapa geral do projeto, o próximo passo é entender a estrutura recorrente usada dentro dos módulos de negócio e por que essa estrutura é importante para o desenvolvimento diário.

## Estrutura base dos módulos

Os módulos principais da aplicação seguem, de forma geral, a seguinte estrutura:

```shell
internal/<module>/
├── delivery
├── application
├── domain
└── infrastructure
```

Essa estrutura não é apenas uma convenção estética. Ela representa a separação de responsabilidades entre entrada HTTP, casos de uso, regras de negócio e implementação técnica.

Cada módulo deve ser entendido como uma pequena unidade funcional dentro do sistema. Ele possui:

- sua forma de entrada;
- seus fluxos de negócio;
- seus conceitos e regras;
- suas implementações concretas de persistência e integração.

## Objetivo da separação interna

A separação interna dos módulos existe para evitar que o código se organize em torno de detalhes técnicos antes de se organizar em torno de responsabilidade.

Sem essa divisão, seria comum encontrar:

- regra de negócio dentro de handlers HTTP;
- decisões de domínio espalhadas em queries SQL;
- lógica de autorização misturada com serialização JSON;
- fluxo de caso de uso escondido dentro de repositórios.

Essa mistura aumenta o acoplamento e dificulta leitura, manutenção e testes. A estrutura em camadas reduz esse problema ao distribuir o código de acordo com o papel que cada parte desempenha.

## Camada `delivery`

A camada `delivery` representa a borda HTTP do módulo.

Ela é responsável por receber a requisição e adaptá-la para a aplicação. Isso inclui:

- leitura de método, rota e parâmetros;
- decodificação de body JSON;
- leitura de query params;
- extração do usuário autenticado do contexto;
- construção da entrada do caso de uso;
- transformação do resultado em resposta HTTP;
- mapeamento de erros para status code e payload.

Essa camada trabalha com o protocolo HTTP e com o contrato REST da API.

Ela **não é dona da regra de negócio**. Seu papel é adaptar entrada e saída. Quando um handler verifica se um UUID é inválido, se um body JSON está malformado ou se um query param inesperado foi enviado, ele está tratando formato e contrato de transporte, não regra de domínio.

Por outro lado, decisões como:

- um usuário pode criar conta;
- uma conta pode receber depósito;
- há saldo suficiente para saque;
- a conta de destino pode participar da transferência;

não pertencem à `delivery`. Essas decisões são tomadas nas camadas internas.

## Camada `application`

A camada `application` contém os casos de uso do módulo.

Um caso de uso representa uma ação do sistema. É nessa camada que o sistema expressa operações como:

- registrar usuário;
- aprovar usuário;
- listar contas;
- criar conta;
- consultar saldo;
- depositar;
- sacar;
- transferir;
- consultar extrato.

A `application` coordena o fluxo da operação. Suas responsabilidades incluem:

- validar pré-condições da operação;
- consultar contratos de repositório;
- combinar regras do domínio;
- abrir e controlar transações quando necessário;
- orquestrar interações entre mais de uma dependência;
- devolver um resultado para a camada de entrega.

Essa camada não deve ser confundida com o domínio. O domínio expressa regras e conceitos; a aplicação expressa execução coordenada de um fluxo.

Por exemplo:

- a regra de que uma conta inativa não pode receber depósito pertence ao domínio;
- a decisão de carregar a conta, validar a operação, atualizar saldo e registrar ledger dentro de uma transação pertence ao caso de uso na aplicação.

## Camada `domain`

A camada `domain` contém o núcleo conceitual do módulo.

Ela é responsável por representar:

- entidades;
- invariantes;
- regras de negócio;
- contratos;
- erros de domínio;
- tipos centrais do contexto funcional.

No módulo `account`, por exemplo, o domínio inclui conceitos como:

- conta;
- transação;
- status da conta;
- saldo insuficiente;
- conta inativa;
- contrato de repositório de contas.

O domínio não lida com HTTP, JSON, SQL, `pgx`, middleware ou drivers. Ele representa o que o sistema considera verdadeiro do ponto de vista do negócio.

É importante perceber que o domínio não precisa conter toda a lógica do sistema em isolamento absoluto. Algumas decisões de fluxo pertencem naturalmente à aplicação. Ainda assim, o domínio é o lugar onde ficam as regras que definem o significado das operações.

Exemplos de regra de domínio:

- uma transferência não pode ter a mesma conta como origem e destino;
- um saque exige saldo suficiente;
- uma conta inativa não pode operar;
- determinados valores são inválidos para uma operação financeira.

## Camada `infrastructure`

A camada `infrastructure` contém as implementações técnicas concretas do módulo.

Ela é responsável por:

- executar queries SQL;
- implementar contratos de repositório;
- integrar com PostgreSQL;
- mapear linhas e registros;
- controlar detalhes de transação;
- aplicar locks de banco;
- integrar com mecanismos técnicos como JWT e hashing, quando pertinente ao módulo.

Enquanto o domínio define contratos e a aplicação depende deles, a infraestrutura fornece a implementação real desses contratos.

No projeto atual, isso aparece principalmente em implementações com PostgreSQL e `pgx`. No módulo `account`, por exemplo, a infraestrutura é responsável por consultar contas, atualizar saldos, bloquear linhas, persistir transações e reconstruir resultados a partir do banco.

Essa camada pode conhecer detalhes do banco. O restante da aplicação não precisa conhecer esses detalhes para executar o fluxo de negócio.

## Relação entre as camadas

As camadas se relacionam por meio de dependências com direção controlada:

```text
delivery -> application -> domain
infrastructure -> domain
```

Em um fluxo típico:

1. `delivery` recebe a chamada HTTP;
2. `delivery` chama um caso de uso da `application`;
3. `application` usa regras e contratos do `domain`;
4. `application` chama implementações que satisfazem esses contratos;
5. `infrastructure` executa a parte técnica;
6. o resultado retorna para `application` e depois para `delivery`.

Essa direção é importante porque reduz contaminação entre camadas. O domínio não depende de HTTP. A aplicação não depende diretamente de SQL. A infraestrutura não define o significado da operação.

## Exemplo prático: módulo `account`

O módulo `account` é um bom exemplo porque possui vários tipos de fluxo.

Sua organização inclui:

```shell
internal/account/
├── delivery
├── application
│   ├── account
│   ├── statement
│   └── transaction
├── domain
└── infrastructure
```

Essa subdivisão na `application` existe porque o módulo concentra responsabilidades diferentes:

- `account/`: criação de conta, listagem de contas e consulta de saldo;
- `transaction/`: depósito, saque e transferência;
- `statement/`: leitura de extrato.

Mesmo com essa subdivisão adicional, a lógica geral continua a mesma: entrada em `delivery`, coordenação em `application`, regras e contratos em `domain`, implementação concreta em `infrastructure`.

## Como essa anatomia ajuda no desenvolvimento

Entender a anatomia interna dos módulos ajuda a decidir onde cada mudança deve ser feita.

Por exemplo:

- se o problema é um status code errado na resposta, a investigação tende a começar em `delivery`;
- se a operação precisa abrir transação ou combinar dependências, a mudança tende a ficar em `application`;
- se a regra do negócio mudou, a mudança tende a ficar em `domain` ou em como a `application` usa o domínio;
- se a query SQL ou o mapeamento do banco está errado, a mudança tende a ficar em `infrastructure`.

Essa previsibilidade reduz o risco de implementar uma solução rápida no lugar errado e acumular acoplamento indevido.

## Como essa anatomia ajuda nos testes

A separação por camadas também organiza melhor a estratégia de testes.

Ela permite:

- testar regras de domínio isoladamente;
- testar casos de uso com mocks ou fakes de repositório;
- testar handlers como adaptação HTTP;
- testar infraestrutura com banco real em integração.

Isso só funciona bem porque as responsabilidades estão razoavelmente separadas no código.

## Limites e cuidado prático

A presença de camadas não garante sozinha que o módulo estará bem organizado. O benefício real depende de preservar a intenção de cada parte.

Alguns sinais de erosão arquitetural seriam:

- handlers contendo a maior parte da regra da operação;
- repositórios definindo decisões de negócio sem justificativa;
- domínio dependendo de biblioteca técnica;
- lógica de autenticação ou autorização espalhada fora do fluxo esperado;
- uso excessivo de `shared/` para código sem dono claro.

Por isso, ao desenvolver novas funcionalidades, o mais importante não é apenas “seguir a pasta certa”, mas manter coerência entre a responsabilidade e o local da implementação.

## Síntese

A anatomia interna dos módulos é uma das bases da organização da API.

Cada módulo reúne sua entrada HTTP, seus casos de uso, seu núcleo de regras e suas implementações técnicas concretas. Essa distribuição torna o código mais legível, melhora a localização de responsabilidades e permite que a aplicação cresça sem misturar demais detalhes de transporte, negócio e persistência.

Para um novo desenvolvedor, compreender essa estrutura é essencial. Ela funciona como o mapa interno de cada módulo e orienta onde procurar comportamento existente e onde inserir mudanças futuras com menor risco de acoplamento indevido.
