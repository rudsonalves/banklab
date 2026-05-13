# Superfície REST da Aplicação

Este capítulo apresenta a superfície REST exposta pela API no estágio atual. O objetivo é descrever o contrato externo do sistema, isto é, os endpoints que clientes podem utilizar, como esses endpoints se organizam por responsabilidade e o que essa superfície já revela sobre o modelo de uso da aplicação.

## O que significa a superfície REST

A superfície REST é o conjunto de rotas HTTP disponibilizadas pela API para interação externa.

Ela representa:

- as capacidades visíveis do sistema;
- a forma como o cliente acessa essas capacidades;
- os grupos funcionais da aplicação;
- as exigências de autenticação e autorização de cada fluxo.

Do ponto de vista de um novo desenvolvedor, a superfície REST é um bom ponto de entrada para entender o que a aplicação faz. Antes mesmo de aprofundar em código interno, a lista de endpoints já mostra quais operações o sistema oferece, como elas se agrupam e quais áreas do domínio já possuem implementação concreta.

## Organização por grupos funcionais

A API está organizada, na prática, em quatro grupos principais de endpoints:

- autenticação;
- administração;
- customer;
- conta.

Essa divisão acompanha os módulos de negócio do sistema e ajuda a manter o contrato HTTP alinhado à estrutura interna da aplicação.

## Endpoints de autenticação

O grupo de autenticação é responsável por estabelecer e renovar a identidade do usuário dentro do sistema.

Os endpoints atuais são:

```text
POST /auth/register
POST /auth/login
POST /auth/refresh
GET  /auth/me
```

### `POST /auth/register`

Esse endpoint registra um novo usuário e cria o customer associado.

Ele marca a entrada inicial do usuário no sistema. O registro já cria a base de identidade e associação com a entidade de cliente, mas isso não significa automaticamente que o usuário já esteja apto a executar todos os fluxos disponíveis. O processo de aprovação ainda pode ser necessário para liberar determinadas operações.

### `POST /auth/login`

Esse endpoint autentica o usuário e devolve os tokens necessários para uso das rotas protegidas.

O login retorna dados importantes para o restante da aplicação, como:

- `access_token`;
- `refresh_token`;
- `user_id`;
- `role`;
- `customer_id`.

Esses dados estabelecem o contexto autenticado que será usado pelos demais módulos.

### `POST /auth/refresh`

Esse endpoint renova o `access_token` a partir de um `refresh_token`.

No projeto atual, o fluxo implementa rotação de refresh token, o que significa que a renovação bem-sucedida revoga o token anterior e emite um novo. Isso permite maior controle sobre sessão e reduz riscos de reutilização indevida de tokens antigos.

### `GET /auth/me`

Esse endpoint retorna os dados do usuário autenticado atual.

Ele é útil para que o cliente confirme sua identidade ativa no sistema, incluindo informações como `role` e `customer_id`. Também funciona como forma explícita de consultar o contexto autenticado estabelecido a partir do token JWT.

## Endpoints administrativos

O grupo administrativo concentra operações restritas a usuários com role elevada.

No estágio atual, o principal endpoint dessa área é:

```text
POST /admin/users/{id}/approve
```

Esse endpoint aprova um usuário pendente e cria sua conta de forma atômica.

Do ponto de vista do contrato HTTP, ele é uma operação administrativa porque:

- exige autenticação;
- exige autorização por role;
- atua sobre o estado de outro usuário.

Ao mesmo tempo, esse endpoint mostra uma característica importante da API: uma rota pode pertencer ao grupo administrativo sem que todas as regras envolvidas sejam de propriedade do módulo `admin`. O fluxo interage com regras de `auth`, `customer` e `account`, mas o acesso continua sendo administrativamente restrito.

## Endpoints de customer

O grupo `customer` concentra informações sobre o cliente de negócio associado ao usuário autenticado.

No estágio atual, o principal endpoint é:

```text
GET /customers/me
```

Esse endpoint retorna o customer vinculado ao usuário autenticado.

Um aspecto importante desse grupo é que a API não depende de um `customer_id` enviado pelo cliente para descobrir qual customer consultar. Essa identificação é derivada a partir do contexto autenticado da requisição.

Isso reforça um princípio importante da aplicação: o cliente não escolhe livremente o escopo de ownership em operações protegidas; o backend deriva esse escopo a partir da identidade autenticada.

## Endpoints de conta

O grupo `account` concentra o maior número de endpoints e cobre tanto operações cadastrais quanto operações financeiras.

Os endpoints atuais são:

```text
GET  /accounts
POST /accounts
GET  /accounts/{id}/balance
POST /terminal/accounts/{id}/deposit
POST /terminal/accounts/{id}/withdraw
GET  /accounts/internal-transfers/recipients
POST /accounts/internal-transfers
GET  /accounts/{id}/statement
```

### `GET /accounts`

Lista as contas pertencentes ao usuário autenticado.

A operação usa o `customer_id` do contexto autenticado e não aceita filtros neste momento. O objetivo do endpoint é listar contas, não consultar saldo.

### `POST /accounts`

Cria uma nova conta para o usuário autenticado.

Essa operação depende de condições de negócio, como o estado do usuário e a existência do customer associado. O cliente não informa `customer_id`; esse valor é derivado do token autenticado.

### `GET /accounts/{id}/balance`

Consulta o saldo atual de uma conta específica.

Esse endpoint existe separadamente da listagem de contas para manter o propósito da operação explícito e permitir tratamento próprio do saldo.

### `POST /terminal/accounts/{id}/deposit`

Executa um depósito em uma conta.

O endpoint expressa diretamente uma ação de negócio, em vez de expor uma atualização genérica de saldo.

### `POST /terminal/accounts/{id}/withdraw`

Executa um saque em uma conta.

Assim como depósito, ele representa uma intenção de negócio específica e exige validações próprias.

### `GET /accounts/internal-transfers/recipients`

Localiza contas recebedoras elegíveis para transferência interna por agência + conta ou CPF/CNPJ. Retorna apenas dados mínimos de confirmação, como `account_id`, nome do titular, documento mascarado, agência e número da conta.

### `POST /accounts/internal-transfers`

Executa transferência interna entre contas usando `from_account_id` e `to_account_id`.

Esse é um dos fluxos mais críticos da API porque envolve débito, crédito, locks, ledger e suporte a idempotência opcional.

### `GET /accounts/{id}/statement`

Retorna o extrato da conta, com suporte a paginação e filtros de período.

Esse endpoint representa a consulta histórica das movimentações, distinta tanto da listagem de contas quanto da consulta de saldo.

## O que a superfície REST já revela sobre o modelo do sistema

A forma como os endpoints estão organizados já comunica várias decisões arquiteturais e de domínio.

### 1. A API é orientada a ações de negócio

Em vez de expor apenas operações genéricas sobre recursos, a API modela intenções claras, como:

- aprovar usuário;
- depositar;
- sacar;
- transferir;
- consultar extrato.

Isso é especialmente importante em domínio financeiro, porque operações com semânticas diferentes exigem validações, efeitos e regras diferentes.

### 2. Identidade e ownership são derivados do contexto autenticado

Em vários fluxos protegidos, o cliente não fornece explicitamente o identificador da entidade de negócio principal.

Por exemplo:

- `GET /accounts` não recebe `customer_id`;
- `POST /accounts` não recebe `customer_id`;
- `GET /customers/me` não recebe `customer_id`.

Essa decisão reduz risco de manipulação indevida do escopo da operação e reforça a responsabilidade do backend em determinar ownership.

### 3. A API distingue leitura de estado e execução de operação

A superfície REST não mistura todas as capacidades em um único endpoint de conta.

Em vez disso, separa:

- listagem de contas;
- consulta de saldo;
- extrato;
- operações de mutação financeira.

Essa separação melhora clareza de contrato e ajuda a manter a semântica das rotas alinhada ao uso real do sistema.

### 4. Fluxos sensíveis possuem proteção específica

A superfície REST deixa claro que existem diferentes níveis de acesso:

- rotas iniciais de entrada;
- rotas autenticadas;
- rotas restritas por role;
- rotas com validação adicional de ownership.

Essa organização mostra que a API não trata autenticação como um detalhe genérico, mas como parte central do contrato.

## Autenticação na superfície REST

No estágio atual, a API utiliza dois mecanismos principais de proteção no acesso aos endpoints.

### App Token

As rotas de entrada inicial usam App Token:

```text
POST /auth/register
POST /auth/login
```

Esse mecanismo protege o acesso inicial à API sem depender ainda de uma sessão autenticada do usuário.

### JWT Bearer Token

As rotas protegidas usam JWT:

```http
Authorization: Bearer <access_token>
```

Esse token é exigido em endpoints como:

- `POST /auth/refresh`;
- `GET /auth/me`;
- `GET /customers/me`;
- rotas de `account`;
- rotas administrativas protegidas.

Além da autenticação, esses endpoints ainda podem exigir validações adicionais, como role administrativa ou ownership da conta acessada.

## Envelope de resposta

Todos os endpoints seguem um envelope de resposta padronizado:

Sucesso:

```json
{
  "data": {},
  "error": null
}
```

Erro:

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "human readable message"
  }
}
```

Esse padrão torna a superfície REST mais uniforme e facilita integração com clientes como o aplicativo mobile.

## A relação entre a superfície REST e a arquitetura interna

Embora a superfície REST seja externa, ela está fortemente conectada à organização interna da aplicação.

Em geral:

- rotas de `auth` se conectam ao módulo `internal/auth`;
- rotas de `customer` se conectam a `internal/customer`;
- rotas de `account` se conectam a `internal/account`;
- rotas de `admin` se conectam a `internal/admin`.

Essa correspondência não significa que um módulo sempre resolva tudo sozinho. Alguns fluxos atravessam mais de um módulo internamente. Ainda assim, o contrato HTTP preserva uma divisão coerente por responsabilidade principal.

## Consequência prática para desenvolvimento

Ao adicionar ou alterar endpoints, é importante preservar a semântica da superfície REST já adotada.

Isso significa, por exemplo:

- preferir endpoints orientados a ações reais do domínio;
- evitar fazer o cliente enviar identificadores que o backend já conhece por contexto autenticado;
- manter clara a diferença entre listagem, consulta de estado e operação financeira;
- respeitar os níveis de autenticação e autorização já presentes;
- manter a padronização de envelope e códigos de erro.

A superfície REST é uma das partes mais visíveis da aplicação. Mudanças nela afetam não apenas o backend, mas também mobile, testes, documentação e entendimento geral do sistema.

## Síntese

A superfície REST atual da API revela um sistema organizado por áreas funcionais e orientado a operações de negócio claras.

Os grupos `auth`, `admin`, `customer` e `account` expressam os principais contextos do domínio. As rotas protegidas dependem de contexto autenticado, role e ownership. E a separação entre listagem, saldo, extrato e operações financeiras mostra que o contrato externo acompanha as necessidades do domínio, em vez de tratar tudo como simples manipulação genérica de recursos.

Para novos desenvolvedores, estudar a superfície REST é uma forma rápida de compreender o que a aplicação oferece externamente e como essa oferta se conecta às decisões arquiteturais internas do projeto.
