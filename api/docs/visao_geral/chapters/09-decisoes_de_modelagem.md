# Decisões Importantes de Modelagem e Domínio

Este capítulo reúne algumas decisões centrais de modelagem já presentes na API. O objetivo é registrar escolhas que influenciam diretamente a implementação de novas funcionalidades e a manutenção das existentes.

Essas decisões não são apenas detalhes internos do código. Elas representam a forma como o sistema interpreta conceitos do domínio, distribui responsabilidade entre módulos e protege operações sensíveis.

## Por que este capítulo é importante

Em muitos projetos, a maior parte da arquitetura visível está nas pastas, nos endpoints e nas dependências entre camadas. No entanto, uma parte igualmente importante da coerência do sistema está em decisões menores, porém estruturantes, como:

- de onde vem a identidade do customer;
- qual módulo é dono de determinada política;
- como o saldo se relaciona com o ledger;
- como o status do usuário se diferencia do status da conta;
- em que ponto uma regra pertence ao domínio ou ao fluxo da aplicação.

Se essas decisões não estiverem claras, é comum que novas implementações “funcionem” do ponto de vista técnico, mas enfraqueçam a coerência do sistema ao longo do tempo.

## `users.status` e `accounts.status` representam coisas diferentes

Uma das decisões mais importantes do modelo atual é manter separados os estados de usuário e de conta.

### `users.status`

O campo `users.status` representa o estado do usuário no fluxo de acesso e onboarding.

Ele responde a perguntas como:

- o usuário já foi aprovado?
- o usuário ainda está pendente?
- o usuário está bloqueado?
- o usuário pode seguir para fluxos que exigem habilitação prévia?

Portanto, esse status está ligado à identidade autenticável e ao ciclo de vida do usuário dentro do sistema.

### `accounts.status`

O campo `accounts.status` representa o estado operacional de uma conta específica.

Ele responde a perguntas como:

- esta conta pode receber depósito?
- esta conta pode realizar saque?
- esta conta pode participar de uma transferência?
- esta conta está ativa, inativa ou bloqueada para operação?

Esse status está ligado ao produto financeiro “conta” e não ao onboarding do usuário.

### Consequência prática

Essa separação é importante porque:

- um usuário pode estar ativo e ainda assim possuir conta bloqueada;
- um usuário pode ter múltiplas contas com estados diferentes;
- a conta pode sofrer mudanças operacionais sem alterar a identidade do usuário;
- o onboarding do usuário não deve ser confundido com a operacionalidade da conta.

Isso significa que esses dois campos não devem ser tratados como duplicação desnecessária. Eles modelam dimensões diferentes do sistema.

## O `customer_id` vem do contexto autenticado

Outra decisão central é que o backend não depende do cliente para informar explicitamente, em várias operações protegidas, a qual customer aquela ação se refere.

Em fluxos como:

- `GET /customers/me`;
- `GET /accounts`;
- `POST /accounts`;

o `customer_id` é derivado a partir do contexto autenticado.

Essa decisão reforça três princípios:

- identidade sensível deve ser derivada pelo backend sempre que possível;
- o cliente não deve definir livremente o escopo de ownership da operação;
- a autenticação deve servir como base de autorização e não apenas como liberação de acesso.

Na prática, isso significa que novas features devem considerar com cuidado se realmente precisam receber identificadores de escopo no payload ou se essas informações já podem ser inferidas de forma mais segura a partir do contexto autenticado.

## Ownership é uma regra estrutural do sistema

Ownership é a ideia de que um usuário autenticado não pode operar livremente sobre qualquer recurso do sistema apenas porque está autenticado.

No domínio atual, essa preocupação aparece especialmente em contas.

Um usuário customer só pode:

- listar suas próprias contas;
- consultar saldo das contas que lhe pertencem;
- consultar extrato das contas que lhe pertencem;
- operar financeiramente sobre contas às quais possui acesso permitido.

Essa decisão é importante porque autenticação sozinha não resolve segurança de domínio. O sistema precisa validar não apenas “quem é o usuário”, mas “qual relação esse usuário possui com o recurso acessado”.

Essa lógica influencia diretamente a implementação de novos fluxos que envolvam recurso identificado por path param, query param ou body.

## A criação de conta pertence ao módulo `account`

Uma decisão arquitetural importante do sistema é que regras de conta permanecem no módulo `account`, mesmo quando a operação que as aciona está localizada em outro contexto funcional.

O melhor exemplo é o fluxo de aprovação de usuário:

```text
POST /admin/users/{id}/approve
```

Esse endpoint é administrativo do ponto de vista da autorização de acesso, mas a criação da conta continua sendo uma responsabilidade de `account`.

Isso significa que:

- `admin` coordena a operação administrativa;
- `account` continua dono das políticas e regras relacionadas à conta;
- a localização da rota não redefine a propriedade da regra.

Essa decisão é importante porque evita que o módulo administrativo se torne dono indevido de lógica de outros contextos.

## A política de branch pertence ao contexto de conta

No estado atual da aplicação, a definição de branch está localizada no contexto de `account/application/account`.

Essa escolha é coerente com o modelo atual porque a branch faz parte da construção e identidade operacional da conta. Portanto, mesmo que uma conta seja criada em um fluxo administrativo, a regra de resolução da branch continua pertencendo ao módulo de conta.

Essa decisão tem uma implicação prática importante: novas políticas relacionadas à abertura e configuração da conta devem, em princípio, ser avaliadas primeiro dentro do contexto de `account`, e não de `admin`.

## Saldo atual e ledger coexistem com papéis diferentes

O sistema utiliza duas perspectivas complementares para representar o estado financeiro:

- saldo atual em `accounts.balance`;
- histórico imutável em `transactions`.

### `accounts.balance`

O campo de saldo na conta representa o estado atual disponível para consulta e operação.

Ele é usado em fluxos como:

- `GET /accounts/{id}/balance`;
- depósito;
- saque;
- transferência.

### `transactions`

A tabela `transactions` funciona como ledger da aplicação.

Ela registra eventos de movimentação como:

- depósito;
- saque;
- `transfer_out`;
- `transfer_in`.

Esse ledger é importante para:

- auditabilidade;
- reconstrução do histórico;
- extrato;
- rastreabilidade da evolução do saldo.

### Consequência prática

Saldo e ledger não são concorrentes; eles cumprem funções diferentes.

O saldo responde:

```text
quanto existe agora?
```

O ledger responde:

```text
como se chegou até aqui?
```

Novas funcionalidades que mexam em dinheiro precisam respeitar essa dualidade. Alterar saldo sem registrar ledger, por exemplo, enfraquece a consistência histórica do sistema.

## Transferência é modelada como duas pernas de ledger

Uma transferência não é tratada como um único evento abstrato e opaco. No modelo atual, ela gera dois registros complementares:

- `transfer_out` para a conta de origem;
- `transfer_in` para a conta de destino.

Esses registros compartilham uma referência comum, o que permite identificar que ambos pertencem à mesma operação.

Essa decisão melhora:

- auditabilidade;
- rastreamento da operação;
- reconstrução de comportamento;
- integridade histórica do extrato.

Também influencia qualquer evolução futura ligada a conciliação, relatórios ou observabilidade operacional de movimentações.

## Idempotência é parte do domínio operacional da transferência

O suporte a `idempotency_key` em transferências não é apenas um detalhe técnico de API. Ele já faz parte da forma como a aplicação modela segurança operacional em um fluxo sensível.

Essa decisão reconhece que uma operação financeira pode ser repetida pelo cliente por:

- timeout;
- falha de rede;
- dúvida sobre o resultado;
- retry automático.

Ao tratar idempotência como parte do fluxo, o sistema protege o domínio contra duplicidade acidental de débito e crédito.

Isso deve ser considerado como referência para futuros fluxos críticos que também possam sofrer retry.

## O cliente não define livremente a identidade da operação

De forma geral, a modelagem atual evita depender do cliente para determinar elementos sensíveis do domínio quando essas informações já podem ser conhecidas pelo backend.

Isso aparece em decisões como:

- derivar `customer_id` do token;
- validar ownership na aplicação;
- restringir operações administrativas por role;
- separar escopo de conta e escopo de usuário;
- não usar o payload como fonte primária de identidade do recurso quando o contexto autenticado já oferece essa informação.

Esse princípio ajuda a reduzir erros de integração e a proteger o sistema contra manipulação indevida do escopo da operação.

## Regras de negócio não devem ser empurradas para adaptação técnica

Outra decisão de modelagem importante é manter regras de negócio fora de lugares cujo papel principal é adaptação técnica.

Isso significa, por exemplo:

- não transformar handler HTTP em centro da regra;
- não deixar query SQL definir sozinha a política do fluxo;
- não colocar a definição de ownership na camada de serialização;
- não deslocar para `shared/` uma regra cujo dono deveria ser um módulo específico.

Esse cuidado é importante porque uma modelagem coerente não depende apenas das entidades. Ela depende também de onde o sistema decide e executa seu comportamento.

## Modelagem atual como base para evolução

As decisões atuais não significam que o modelo esteja fechado ou completo. Elas significam que existe uma linha de coerência que novas evoluções devem considerar.

Ao adicionar uma nova funcionalidade, vale sempre verificar:

- essa regra pertence a qual módulo?
- ela depende de contexto autenticado?
- ela mexe com onboarding do usuário ou com operacionalidade da conta?
- ela altera saldo atual, ledger ou ambos?
- ela exige ownership?
- ela é uma política de conta ou um fluxo administrativo?

Responder a essas perguntas com base no modelo atual ajuda a manter consistência arquitetural ao longo do tempo.

## Síntese

As decisões de modelagem e domínio da API vão além da definição de structs ou tabelas. Elas estabelecem a semântica do sistema: distinguem usuário de conta, contexto autenticado de payload, saldo atual de histórico, operação administrativa de política de conta e autenticação de ownership.

Para novos desenvolvedores, entender essas decisões é essencial. Muitas implementações futuras vão depender menos de descobrir “como fazer tecnicamente” e mais de reconhecer “qual significado de domínio já está sendo preservado pelo sistema”.
