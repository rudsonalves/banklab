# Principais Fluxos Implementados

Este capítulo descreve os fluxos mais importantes já implementados na API. O objetivo é apresentar não apenas quais endpoints existem, mas como certas operações funcionam do ponto de vista do negócio e por que elas são relevantes para a compreensão do sistema.

Os fluxos escolhidos aqui não esgotam tudo o que a API faz, mas representam bem as decisões arquiteturais e de domínio já presentes na aplicação.

## Por que olhar para fluxos

Até este ponto, a documentação explicou:

- o objetivo da aplicação;
- o estilo arquitetural adotado;
- a organização do projeto;
- a anatomia dos módulos;
- o ciclo de uma requisição;
- a superfície REST;
- o papel da autenticação e do contexto do usuário.

O próximo passo natural é observar como tudo isso se combina em operações concretas.

Um fluxo mostra:

- qual problema de negócio está sendo resolvido;
- quais módulos participam da operação;
- quais validações são feitas;
- onde estão os pontos de decisão;
- como o sistema preserva consistência e autorização.

## Fluxo 1: registro de usuário

O fluxo de registro começa em:

```text
POST /auth/register
```

Seu objetivo é criar a identidade inicial do usuário e o customer associado.

Esse fluxo é importante porque mostra que, no sistema atual, o usuário não existe isoladamente como credencial sem vínculo de domínio. O registro já estabelece a conexão entre:

- identidade autenticável;
- dados do customer;
- início do ciclo de vida do usuário.

Do ponto de vista operacional, o fluxo envolve:

1. validação dos dados de entrada;
2. validação de e-mail e senha;
3. verificação de unicidade;
4. criação do customer;
5. criação do usuário associado;
6. persistência transacional das duas partes.

Esse comportamento é importante porque garante que não haja registro parcial do processo. O usuário e o customer são criados como partes do mesmo fluxo lógico.

## Fluxo 2: login e estabelecimento da sessão

O login acontece em:

```text
POST /auth/login
```

Seu objetivo é autenticar o usuário e devolver o conjunto de dados que passa a representar a sessão ativa.

Esse fluxo inclui:

1. validação das credenciais;
2. autenticação do usuário;
3. emissão de `access_token`;
4. emissão de `refresh_token`;
5. devolução dos dados de identidade úteis ao cliente, como `role` e `customer_id`.

O valor desse fluxo está em estabelecer o contexto autenticado reutilizado pelos demais módulos.

Sem ele, a aplicação não teria base para:

- reconhecer o usuário da requisição;
- saber se ele é administrador;
- saber qual customer está associado à identidade atual;
- determinar ownership em operações de conta.

## Fluxo 3: renovação de sessão

A renovação de sessão acontece em:

```text
POST /auth/refresh
```

Seu objetivo é permitir continuidade da sessão sem exigir novo login por credenciais.

No projeto atual, esse fluxo já possui rotação de refresh token. Isso significa que o token apresentado é verificado e, quando aceito, é substituído por um novo token, invalidando o anterior.

Esse comportamento é importante porque:

- reforça o controle do backend sobre sessões;
- reduz reutilização silenciosa de tokens antigos;
- mantém a sessão ativa sem exigir reautenticação completa a cada expiração do access token.

## Fluxo 4: leitura do usuário autenticado

Esse fluxo está preservado como lógica de aplicação, mas não está exposto na
superfície HTTP atual. O path previsto para um canal de terminal é:

```text
GET /auth/me
```

Seu objetivo é expor ao cliente a identidade que a API reconhece no contexto autenticado atual.

Esse endpoint é simples, mas muito importante conceitualmente. Ele mostra que o backend não depende apenas do login no momento inicial; ele consegue reconstruir e reutilizar a identidade autenticada ao longo das requisições protegidas.

Esse fluxo é útil para:

- confirmar a sessão atual;
- reidratar estado no cliente;
- inspecionar `role`;
- inspecionar `customer_id`;
- validar a integridade da autenticação após refresh.

## Fluxo 5: aprovação de usuário

A aprovação de usuário acontece em:

```text
POST /admin/users/{id}/approve
```

Esse é um dos fluxos mais importantes do sistema porque conecta autenticação, autorização, customer e account.

Seu objetivo é:

- mudar `users.status` de `pending` para `active`;
- validar o customer associado;
- criar a conta correspondente;
- garantir que tudo isso ocorra de forma atômica.

O fluxo envolve:

1. autenticação do chamador;
2. validação de role administrativa;
3. carregamento do usuário alvo com proteção transacional;
4. validação do status atual do usuário;
5. ativação do usuário;
6. validação do customer vinculado;
7. geração de número de conta;
8. resolução de branch;
9. criação da conta;
10. commit da transação.

Esse fluxo é importante porque mostra uma operação transversal. Embora o acesso seja administrativo, a lógica não pertence exclusivamente ao módulo `admin`. Ele orquestra a operação, mas usa regras e dependências de outros módulos.

## Fluxo 6: leitura do customer atual

Esse fluxo acontece em:

```text
GET /customers/me
```

Seu objetivo é consultar o customer associado à identidade autenticada.

Ele reforça um princípio importante da API: o cliente não escolhe explicitamente qual customer consultar. A aplicação deriva isso a partir do contexto autenticado.

Esse fluxo ajuda a manter coerência entre:

- autenticação;
- modelo de cliente;
- ownership das operações futuras.

## Fluxo 7: listagem de contas do usuário autenticado

Esse fluxo acontece em:

```text
GET /accounts
```

Seu objetivo é listar as contas pertencentes ao customer associado ao usuário autenticado.

Ele é um bom exemplo de fluxo relativamente simples, mas arquiteturalmente importante, porque mostra:

- uso do contexto autenticado como fonte de `customer_id`;
- separação entre listagem de conta e consulta de saldo;
- validação de pré-condições no caso de uso;
- dependência de contrato de repositório;
- resposta padronizada ao cliente.

O fluxo envolve:

1. autenticação do usuário;
2. leitura do usuário autenticado do contexto;
3. validação da presença de `customer_id`;
4. busca das contas por `customer_id`;
5. retorno de lista resumida.

Esse endpoint não retorna saldo. Essa omissão é intencional e reforça a separação de propósitos entre listagem e consulta financeira.

## Fluxo 8: criação de conta

Esse fluxo acontece em:

```text
POST /admin/customers/{customer_id}/accounts
```

Seu objetivo é criar uma conta adicional para um customer existente por ação administrativa.

Ele depende de algumas condições importantes:

- o operador deve estar autenticado como admin;
- o `customer_id` alvo deve ser informado na rota;
- o customer alvo precisa existir.

A criação de conta não é uma ação self-service do cliente. A primeira conta é criada automaticamente na aprovação do onboarding; este fluxo existe para provisionamento administrativo de contas adicionais.

## Fluxo 9: consulta de saldo

Esse fluxo acontece em:

```text
GET /accounts/{id}/balance
```

Seu objetivo é devolver o saldo atual de uma conta específica.

Esse endpoint existe separadamente da listagem de contas para manter o contrato mais claro e permitir tratamento próprio de saldo.

O fluxo envolve:

- autenticação;
- verificação de ownership ou permissão;
- leitura da conta;
- retorno do saldo atual.

Mesmo sendo uma leitura, ele é importante no sistema porque o saldo é um dado sensível e precisa ser consultado dentro do escopo correto de acesso.

## Fluxo 10: depósito

Esse fluxo acontece em:

```text
POST /terminal/accounts/{id}/deposit
```

Seu objetivo é incrementar o saldo de uma conta e registrar a movimentação
correspondente. A rota permanece comentada no wiring até existir um canal de
terminal com autenticação e controles próprios.

Esse fluxo já mostra que operações financeiras não são tratadas como simples atualização de registro. Ele envolve:

- validação do valor;
- verificação da conta;
- verificação do status operacional da conta;
- abertura de transação;
- atualização do saldo;
- registro da movimentação em ledger.

O depósito é um fluxo financeiramente mais simples do que transferência, mas ainda assim exige controle transacional.

## Fluxo 11: saque

Esse fluxo está preservado como lógica de aplicação, mas não está exposto na
superfície HTTP atual. O path previsto para um canal de terminal é:

```text
POST /terminal/accounts/{id}/withdraw
```

Seu objetivo é reduzir o saldo da conta e registrar a movimentação
correspondente. A rota permanece comentada no wiring até existir um canal de
terminal com autenticação e controles próprios.

Além dos cuidados de depósito, ele exige:

- verificação de saldo suficiente;
- proteção contra saque inválido em conta inativa;
- persistência consistente entre saldo e ledger.

Esse fluxo mostra que o sistema já trata saldo insuficiente como regra de negócio explícita.

## Fluxo 12: transferência

Esse fluxo acontece em:

```text
GET  /accounts/internal-transfers/recipients
POST /accounts/internal-transfers
```

Ele é o fluxo financeiro mais crítico da aplicação.

Seu objetivo é transferir valor entre duas contas, mantendo:

- atomicidade;
- proteção contra concorrência;
- histórico consistente;
- possibilidade de retry seguro via idempotência.

Esse fluxo envolve:

1. validação das contas e do valor;
2. validação de que origem e destino são diferentes;
3. abertura de transação;
4. lock das contas em ordem determinística;
5. validação das regras da conta de origem;
6. validação das regras da conta de destino;
7. débito da origem;
8. crédito do destino;
9. persistência de `transfer_out`;
10. persistência de `transfer_in`;
11. tratamento de idempotência quando houver chave.

Esse fluxo sintetiza várias das principais preocupações arquiteturais da API.

## Fluxo 13: extrato

Esse fluxo acontece em:

```text
GET /accounts/{id}/statement
```

Seu objetivo é consultar o histórico de movimentações da conta.

Ele é diferente da consulta de saldo porque não representa o estado atual, mas a trilha histórica das operações.

Esse fluxo inclui:

- autenticação;
- validação de acesso à conta;
- filtros de paginação e período;
- consulta ao ledger;
- retorno dos itens do extrato e cursor, quando aplicável.

Esse endpoint é importante porque mostra que a aplicação já trata o histórico financeiro como informação de primeira classe, e não apenas como efeito colateral da atualização de saldo.

## Como esses fluxos ajudam a entender o sistema

Esses fluxos mostram que a API já combina vários tipos de responsabilidade:

- identidade e autenticação;
- onboarding e aprovação;
- vínculo entre usuário e customer;
- ownership sobre recursos;
- modelagem de conta;
- consulta de estado;
- movimentação financeira;
- histórico persistido.

Observar os fluxos em conjunto ajuda a entender que os módulos não foram separados apenas por conveniência técnica. Eles respondem a partes reais do comportamento da aplicação.

## Consequência prática para desenvolvimento

Ao implementar uma nova funcionalidade, um bom primeiro passo é identificar a qual desses tipos de fluxo ela se assemelha.

Por exemplo:

- se a nova feature depende de identidade e sessão, o ponto de comparação tende a ser `auth`;
- se depende de contexto autenticado e customer, olhar `GET /customers/me` e `GET /accounts` ajuda;
- se depende de mutação financeira, depósito, saque e transferência servem como referência;
- se depende de operação administrativa transversal, aprovação de usuário é o melhor exemplo existente.

Essa comparação ajuda a reaproveitar padrões já consolidados na aplicação.

## Síntese

Os principais fluxos já implementados mostram que a API está além de um conjunto de endpoints isolados.

Ela já possui operações que combinam autenticação, autorização, ownership, consistência financeira, transação, histórico e organização modular. Para novos desenvolvedores, estudar esses fluxos é uma das formas mais eficazes de compreender como a arquitetura se traduz em comportamento concreto do sistema.
