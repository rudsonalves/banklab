# Persistência, Transações e Consistência

Este capítulo descreve como a API trata persistência e consistência no estágio atual. O foco não é apenas listar tabelas ou operações de banco, mas explicar como o sistema usa PostgreSQL, transações, locks e ledger para preservar a integridade das operações de negócio, especialmente as operações financeiras.

## Papel da persistência na aplicação

Na arquitetura atual, a persistência não é apenas um meio de armazenamento. Ela participa diretamente da preservação de regras importantes do sistema.

Isso é especialmente verdadeiro em fluxos que envolvem:

- criação de usuário e customer;
- aprovação de usuário;
- criação de conta;
- depósito;
- saque;
- transferência;
- leitura de extrato.

Em todos esses casos, o banco não é apenas um repositório passivo de dados. Ele também é o ambiente onde atomicidade, bloqueio concorrente, histórico e integridade relacional são garantidos.

## Banco de dados principal

O backend utiliza PostgreSQL como banco principal.

Essa escolha é coerente com o estágio atual do projeto porque PostgreSQL oferece:

- transações robustas;
- locks em nível de linha;
- constraints relacionais;
- boa capacidade para modelagem transacional;
- suporte adequado a queries, índices e consistência forte.

Para um sistema que já lida com saldo e operações monetárias, essas características são mais importantes do que apenas capacidade de armazenar registros.

## Repositórios e implementação concreta

Na aplicação, os casos de uso dependem de contratos de repositório. A implementação concreta desses contratos fica na camada `infrastructure`, principalmente em implementações PostgreSQL com `pgx`.

Isso significa que a persistência está organizada em torno de dois níveis:

- contratos definidos pelas camadas internas;
- implementação concreta executada pela infraestrutura.

Essa separação permite que os casos de uso operem sobre capacidades do domínio, como buscar conta, listar contas, atualizar saldo ou persistir transações, sem depender diretamente da query SQL usada para isso.

## Principais áreas persistidas

No estágio atual, o modelo persistido da API inclui, entre outros, os seguintes conjuntos de dados:

- usuários;
- customers;
- contas;
- sessões de usuário;
- transações financeiras.

Essas áreas sustentam os principais fluxos do sistema:

- `users` e `user_sessions` participam da autenticação e continuidade da sessão;
- `customers` conectam identidade autenticada ao domínio de negócio;
- `accounts` representam o estado atual das contas;
- `transactions` registram o histórico imutável de movimentações.

## Persistência transacional em fluxos compostos

Uma das decisões mais importantes do sistema é usar transações explícitas em fluxos que não podem deixar estado parcial.

Isso aparece claramente em operações como:

- registro de usuário com criação de customer;
- aprovação de usuário com criação de conta;
- depósito;
- saque;
- transferência.

### O que a transação protege

Uma transação protege a operação como unidade lógica.

Isso significa que, quando um fluxo exige vários passos dependentes entre si, todos eles devem ser confirmados juntos ou revertidos juntos.

Exemplos:

- no registro, não deve existir customer sem usuário nem usuário sem customer dentro do fluxo proposto;
- na aprovação, não deve existir usuário ativo sem a conta criada pelo mesmo fluxo;
- no depósito, saldo e ledger precisam permanecer coerentes;
- na transferência, débito, crédito e duas pernas de ledger precisam ser persistidos de forma consistente.

Sem transação, o sistema ficaria exposto a estados parciais em caso de falha intermediária.

## Locks e concorrência

Em operações financeiras, uma das maiores fontes de risco não é apenas falha técnica, mas concorrência.

Se duas operações alterarem o saldo da mesma conta ao mesmo tempo sem coordenação adequada, podem surgir inconsistências como:

- leitura do mesmo saldo anterior por duas operações simultâneas;
- débito concorrente além do disponível;
- perda de atualização;
- divergência entre saldo atual e ledger.

Para reduzir esse risco, a API usa bloqueio de linha com:

```sql
SELECT ... FOR UPDATE
```

Esse padrão é usado quando uma conta precisa ser lida e depois modificada dentro de uma transação.

O lock impede que outra transação concorrente altere a mesma linha até que a operação atual seja concluída por commit ou rollback.

## Ordenação determinística de locks em transferência

Transferência é um caso especial porque envolve duas contas.

Quando duas transações concorrentes tentam bloquear dois recursos em ordens diferentes, existe risco de deadlock. Para reduzir esse risco, a API adota ordenação determinística ao adquirir os locks das contas envolvidas.

Isso significa que as contas são bloqueadas segundo um critério estável, e não segundo a ordem acidental em que apareceram na requisição.

Essa decisão é importante porque:

- reduz risco de bloqueio circular;
- aumenta previsibilidade do fluxo transacional;
- preserva segurança em uma operação naturalmente mais complexa.

## Saldo atual e histórico persistido

O sistema trabalha com duas representações complementares do estado financeiro:

- saldo atual em `accounts.balance`;
- histórico imutável em `transactions`.

### Saldo atual

O saldo atual serve para representar o estado operacional presente da conta.

Ele é usado em fluxos como:

- consulta de saldo;
- depósito;
- saque;
- transferência.

Esse valor precisa estar correto do ponto de vista da operação atual da conta.

### Histórico em ledger

A tabela `transactions` funciona como ledger da aplicação.

Ela registra eventos financeiros persistidos, como:

- depósito;
- saque;
- `transfer_out`;
- `transfer_in`.

Esse histórico é usado para:

- extrato;
- auditabilidade;
- reconstrução de contexto;
- rastreamento da evolução do saldo.

### Consequência prática

Persistir apenas o saldo não seria suficiente para o domínio atual. E persistir apenas eventos sem manter saldo operacional também complicaria vários fluxos do escopo atual.

O modelo adotado mantém os dois:

- `accounts.balance` como estado atual;
- `transactions` como trilha histórica.

## Coerência entre saldo e ledger

Um ponto crítico da persistência atual é que saldo e ledger precisam evoluir juntos nos fluxos financeiros.

Isso significa que:

- depósito deve aumentar saldo e registrar `deposit`;
- saque deve reduzir saldo e registrar `withdraw`;
- transferência deve reduzir origem, aumentar destino e registrar `transfer_out` e `transfer_in`.

Se um desses passos for executado sem o outro, o sistema perde coerência entre estado presente e histórico.

Por isso, essas operações são protegidas por transação.

## Transferência como operação composta

Na persistência atual, a transferência é modelada como uma operação com múltiplos efeitos coordenados:

1. validação das contas;
2. lock das linhas envolvidas;
3. débito da conta de origem;
4. crédito da conta de destino;
5. registro de `transfer_out`;
6. registro de `transfer_in`.

As duas pernas do ledger compartilham uma referência comum, o que permite relacionar os dois registros como partes da mesma transferência.

Essa modelagem facilita:

- auditabilidade;
- reconstrução histórica;
- replay de resultado em cenários idempotentes;
- consistência lógica do extrato.

## Idempotência e persistência

No fluxo de transferência, a persistência também participa da proteção contra repetição indevida de operação.

A `idempotency_key` opcional permite que o sistema reconheça que um cliente está repetindo uma requisição logicamente equivalente.

Isso é importante em cenários como:

- timeout;
- falha de rede;
- incerteza do cliente sobre o resultado da chamada anterior;
- retry automático.

Ao tratar idempotência na persistência, a API evita que a repetição da requisição gere um novo débito e um novo crédito indevidos.

## Consultas e consistência de leitura

Nem toda operação da API exige o mesmo grau de proteção transacional.

Fluxos de leitura como:

- `GET /accounts`;
- `GET /accounts/{id}/balance`;
- `GET /accounts/{id}/statement`;
- `GET /customers/me`;
- `GET /auth/me`;

não necessariamente exigem o mesmo tipo de lock usado em operações de mutação financeira.

Ainda assim, cada leitura precisa respeitar:

- autenticação;
- autorização;
- ownership;
- semântica correta do dado retornado.

No caso do saldo, por exemplo, a consulta usa o snapshot atual persistido em `accounts.balance`. No caso do extrato, a leitura recai sobre o ledger persistido.

## Migrações e evolução do schema

A persistência da aplicação evolui por meio de migrations versionadas.

Esse mecanismo é importante porque permite:

- registrar mudanças estruturais do banco;
- evoluir schema de forma reprodutível;
- manter coerência entre código e modelo persistido;
- acompanhar a história técnica da aplicação.

Para um novo desenvolvedor, esse ponto é importante porque alterações de modelo não devem ser tratadas como mudanças isoladas em structs ou queries. A evolução do banco precisa ser explícita, controlada e versionada.

## Persistência como parte da arquitetura, não como detalhe invisível

Em muitos sistemas simples, o banco acaba sendo visto apenas como camada de armazenamento.

Na API atual, isso seria uma visão incompleta. A persistência participa de forma ativa da arquitetura porque sustenta:

- atomicidade;
- locks;
- ownership de dados;
- histórico financeiro;
- integridade relacional;
- continuidade de sessão.

Isso significa que decisões de persistência têm impacto direto no comportamento do sistema, e não apenas em sua eficiência técnica.

## Cuidados práticos para novas implementações

Ao implementar novos fluxos, é importante avaliar cuidadosamente:

- a operação pode deixar estado parcial?
- ela exige transação?
- ela mexe com saldo, ledger ou ambos?
- ela precisa de lock de linha?
- ela pode sofrer concorrência relevante?
- ela precisa de idempotência?
- o dado retornado deve vir do estado atual ou do histórico persistido?

Responder corretamente a essas perguntas ajuda a manter coerência com o modelo atual da aplicação.

## Síntese

A persistência da API está profundamente ligada à consistência do sistema.

PostgreSQL, transações, locks, saldo atual e ledger não aparecem aqui apenas como escolhas técnicas isoladas. Eles formam o mecanismo pelo qual a aplicação protege operações sensíveis, evita estados parciais e mantém coerência entre o presente operacional da conta e o histórico financeiro persistido.

Para novos desenvolvedores, compreender essa relação é essencial. Alterações aparentemente simples em persistência podem afetar autorização, integridade de saldo, comportamento concorrente e rastreabilidade da aplicação como um todo.
