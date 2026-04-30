# Objetivo e Escopo Atual da Aplicação

Este capítulo descreve o propósito da API e delimita o que ela implementa no estágio atual do projeto. O foco é ajudar o leitor a entender rapidamente qual problema a aplicação resolve, quais responsabilidades já estão cobertas pelo backend e quais áreas do domínio já possuem implementação concreta.

## Objetivo da API

A API do projeto `banklab` implementa um core bancário simplificado. Seu papel é oferecer as capacidades backend necessárias para autenticação de usuários, gestão de identidade de cliente, criação e consulta de contas e execução de operações financeiras básicas.

Embora o projeto ainda não represente um sistema bancário completo, ele já incorpora preocupações típicas desse domínio:

- autenticação e controle de sessão;
- autorização baseada em identidade e papel do usuário;
- onboarding com aprovação administrativa;
- vínculo entre usuário autenticado e entidade de cliente;
- criação e consulta de contas;
- operações monetárias com consistência transacional;
- registro histórico de movimentações.

Portanto, a API não deve ser entendida como um simples sistema de cadastro. Mesmo em um escopo reduzido, ela já lida com regras de autorização, ownership, integridade de saldo e separação de responsabilidades entre módulos.

## Papel da API dentro do projeto

Dentro do monorepo `banklab`, a API funciona como backend principal consumido pela aplicação mobile. Ela expõe um contrato HTTP REST que permite ao cliente:

- registrar usuários;
- autenticar sessão;
- renovar tokens;
- consultar o usuário autenticado;
- consultar dados do customer associado;
- listar e criar contas;
- consultar saldo;
- executar depósito, saque e transferência;
- consultar extrato;
- executar fluxos administrativos de aprovação.

Esse backend também concentra decisões sensíveis que não devem ficar no cliente, como:

- derivação de `customer_id` a partir do token autenticado;
- verificação de ownership sobre contas;
- checagem de permissões administrativas;
- aplicação de regras de status de usuário e conta;
- proteção transacional de operações financeiras.

## Escopo funcional implementado

No estágio atual, a aplicação já possui implementação para quatro áreas principais do domínio.

### 1. Autenticação e sessão

O módulo `auth` cobre:

- registro de usuário;
- criação do customer associado no registro;
- login com emissão de `access_token` e `refresh_token`;
- refresh de sessão com rotação de token;
- leitura do usuário autenticado atual em `GET /auth/me`.

Essa área é responsável por estabelecer a identidade do usuário dentro do sistema e disponibilizar contexto autenticado para os demais módulos.

### 2. Customer

O módulo `customer` cobre:

- criação do customer durante o fluxo de registro;
- consulta do customer vinculado ao usuário autenticado em `GET /customers/me`.

Esse módulo separa a identidade de acesso do conceito de cliente de negócio. O usuário autentica no sistema; o customer representa a entidade de domínio associada àquele acesso.

### 3. Administração

O módulo `admin` cobre, neste momento:

- aprovação de usuário por meio de `POST /admin/users/{id}/approve`.

Esse fluxo é restrito a administradores e coordena uma operação composta: transição de `users.status` para `active` e criação da conta associada. Trata-se de uma funcionalidade administrativa porque depende de permissão elevada, mas ela interage diretamente com regras de `auth`, `customer` e `account`.

### 4. Contas e operações financeiras

O módulo `account` cobre:

- listagem de contas do usuário autenticado;
- criação de conta;
- consulta de saldo;
- depósito;
- saque;
- transferência;
- extrato.

Essa é a área mais sensível da API do ponto de vista transacional, porque lida com saldo, ledger, concorrência e idempotência em operações críticas.

## Escopo técnico implementado

Além do escopo funcional, a aplicação já possui algumas decisões técnicas relevantes implementadas.

### Arquitetura

A API segue um modelo de monolito modular com separação em camadas:

- `delivery`;
- `application`;
- `domain`;
- `infrastructure`.

Essa estrutura é repetida nos módulos principais do sistema e organiza a aplicação por responsabilidade.

### Persistência

O backend utiliza PostgreSQL como banco principal. A persistência já contempla:

- tabela de usuários;
- tabela de customers;
- tabela de contas;
- tabela de sessões;
- tabela de transações como ledger.

O projeto também utiliza migrations para evolução controlada do schema.

### Consistência financeira

Operações que alteram saldo já contam com mecanismos de proteção, incluindo:

- transações explícitas de banco;
- locks com `SELECT ... FOR UPDATE`;
- ordenação determinística de bloqueios em transferência;
- ledger imutável em `transactions`;
- suporte a idempotência em transferência.

Esses mecanismos fazem parte do escopo atual porque a aplicação já trata consistência financeira como requisito de arquitetura, e não apenas como detalhe de implementação.

### Testes

O projeto também já possui cobertura distribuída entre:

- testes de domínio;
- testes de aplicação;
- testes de delivery;
- testes de integração.

Isso indica que o escopo atual não se limita à implementação funcional. Já existe também uma preocupação concreta com validação de comportamento e segurança de evolução.

## O que a aplicação ainda não busca resolver

Para entender corretamente o estágio atual da API, também é importante delimitar o que ela ainda não pretende resolver como foco principal.

No momento, a aplicação não está orientada a:

- arquitetura distribuída por microsserviços;
- deploy independente por módulo;
- processamento assíncrono como base do fluxo principal;
- event sourcing completo;
- CQRS amplo;
- múltiplos produtos bancários complexos;
- observabilidade operacional avançada;
- escalabilidade horizontal altamente especializada por domínio.

Essas ausências não representam falhas do projeto. Elas indicam apenas que a aplicação foi desenhada para resolver primeiro um conjunto mais fundamental de problemas: organização do domínio, consistência das operações principais e clareza arquitetural.

## Limite atual do domínio

Hoje, a aplicação cobre um núcleo bancário simplificado. Isso significa que ela já possui noções de:

- usuário;
- customer;
- conta;
- saldo;
- movimentação financeira;
- histórico de transações;
- aprovação administrativa.

Por outro lado, ela ainda não modela áreas mais amplas que poderiam existir em um sistema financeiro maior, como:

- cartões;
- empréstimos;
- limites de crédito;
- tarifas;
- investimento;
- notificações multicanal;
- conciliação externa;
- antifraude especializado;
- auditoria operacional expandida.

Essas possíveis expansões ajudam a entender que o sistema atual não é um banco completo, mas sim uma base arquitetural e funcional consistente para um recorte específico do domínio.

## Síntese do estágio atual

No estágio atual, a API já deve ser entendida como um backend funcional com regras reais de domínio e não apenas como um protótipo de endpoints.

Ela já possui:

- identidade autenticada propagada para os fluxos;
- autorização por role e ownership;
- onboarding com aprovação;
- criação e consulta de contas;
- operações financeiras básicas;
- proteção de consistência em saldo;
- documentação e testes em evolução.

Essa combinação torna a aplicação adequada tanto para estudo arquitetural quanto para continuação de desenvolvimento incremental.

O ponto central deste capítulo é que a API já possui um escopo suficientemente definido para exigir disciplina arquitetural. Mesmo sem cobrir um domínio bancário completo, ela já implementa fluxos que precisam preservar regras de negócio, controle de acesso e consistência de estado.
