# Relatório Completo da API Implementada

Data de referência: 12/05/2026  
Escopo analisado: `api/` (código-fonte, migrations e documentação técnica)

## 1. Resumo Executivo

A API do projeto BankLab está implementada como um monólito modular em Go, com separação em camadas (delivery, application, domain e infrastructure), autenticação híbrida (`X-App-Token` no onboarding e JWT nas rotas protegidas), persistência em PostgreSQL e operações financeiras com foco em consistência forte.

Principais capacidades já implementadas:

- cadastro de usuário com criação de cliente no mesmo fluxo transacional
- login com emissão de `access_token` + `refresh_token`
- refresh de sessão com rotação de token
- recuperação do usuário autenticado (`/auth/me`)
- aprovação administrativa de usuário pendente com criação atômica de conta
- consulta do perfil do cliente autenticado (`/customers/me`)
- criação e listagem de contas
- operações monetárias: depósito, saque e transferência interna
- recibo por referência da transferência
- extrato com paginação baseada em cursor e filtros por período
- consulta de saldo da conta

Observação importante: os endpoints atuais de criação de conta, depósito e saque
representam provisionamento operacional ou operações diretas sobre o ledger. Eles
não devem ser tratados como funcionalidades voltadas ao mobile ou a um app web de
cliente final. A direção recomendada é mover a criação de conta para uma
superfície administrativa de onboarding/provisionamento e mover depósito/saque
para uma superfície protegida da estrutura (admin/backoffice/sandbox/integração
operacional) ou substituí-los por fluxos reais de cash-in/cash-out quando estes
forem definidos.

## 2. Arquitetura e Organização Implementadas

Stack atual:

- Go 1.26.1
- PostgreSQL 16
- `pgx/v5`
- `net/http` (roteamento nativo do Go)

Estrutura da API:

- `cmd/api/main.go`: bootstrap da aplicação e roteamento
- `internal/bootstrap`: inicialização, carga de ambiente e registro global de erros
- `internal/database`: conexão/pool e contexto de transação
- `internal/auth`: autenticação, sessões e segurança
- `internal/customer`: domínio e consulta de perfil do cliente
- `internal/account`: conta bancária, transações e extrato
- `internal/admin`: aprovação de usuário
- `internal/shared`: contexto de autenticação, envelope HTTP e catálogo de erros
- `migrations`: evolução de esquema do banco

Direção de dependências observada:

- delivery -> application -> domain
- infrastructure -> domain

## 3. Inicialização e Configuração em Runtime

Comportamento de bootstrap:

1. carrega `.env` em caminhos candidatos (`.env`, `api/.env`, diretório do executável)
2. valida fail-fast variáveis obrigatórias:
   - `APP_TOKEN`
   - `JWT_SECRET`
3. cria pool PostgreSQL
4. instancia repositórios
5. instancia serviços de segurança:
   - hash de senha com `bcrypt`
   - JWT (`15m` de validade para access token)
6. instancia casos de uso
7. instancia handlers HTTP
8. registra middlewares e rotas
9. inicia servidor em `:8080`

## 4. Contrato HTTP Implementado

### 4.1 Envelope padrão de resposta

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
    "message": "mensagem",
    "details": {}
  }
}
```

Observação: `details` existe no tipo de resposta, mas normalmente não é preenchido pelos handlers atuais.

### 4.2 Endpoints expostos atualmente

#### Autenticação

- `POST /auth/register` (requer `X-App-Token`)
- `POST /auth/login` (requer `X-App-Token`)
- `POST /auth/refresh` (requer `refresh_token` no body)
- `GET /auth/me` (requer JWT Bearer)

#### Administração

- `POST /admin/users/{id}/approve` (JWT + role `admin`)

#### Cliente

- `GET /customers/me` (JWT)

#### Contas

- `GET /accounts` (JWT)
- `POST /accounts` (JWT; atual, mas candidato a admin/onboarding)
- `GET /accounts/internal-transfers/recipients` (JWT)
- `GET /accounts/{id}/balance` (JWT)
- `GET /accounts/{id}/statement` (JWT)

#### Transações

- `POST /terminal/accounts/{id}/deposit` (JWT)
- `POST /terminal/accounts/{id}/withdraw` (JWT)
- `POST /accounts/internal-transfers` (JWT)
- `GET /accounts/transfer/{transaction_reference}/receipt` (JWT)

Nota: `POST /accounts`, `deposit` e `withdraw` permanecem expostos no contrato
atual, mas devem ser considerados endpoints operacionais/provisórios. Eles não
devem ser consumidos por mobile ou web client de usuário final até serem
reposicionados atrás de proteção estrutural adequada.

## 5. Segurança e Autorização Implementadas

### 5.1 App Token (onboarding)

Middleware dedicado valida `X-App-Token` com comparação em tempo constante (`crypto/subtle`), protegendo:

- `POST /auth/register`
- `POST /auth/login`

### 5.2 JWT

Middleware JWT:

- extrai `Authorization: Bearer <token>`
- valida assinatura/claims via `TokenService`
- injeta principal autenticado no contexto (`user_id`, `role`, `customer_id`)

Rotas protegidas retornam erro padronizado para ausência/token inválido.

### 5.3 Modelo de autorização aplicado

- operações administrativas exigem `role=admin`
- fluxos de conta/transação exigem usuário autenticado e validação de posse de conta
- fluxos de cliente exigem `customer_id` válido no principal

## 6. Módulos Funcionais Implementados

### 6.1 Auth (`internal/auth`)

Implementado:

- `RegisterUserUseCase`
  - valida email/senha
  - cria `customer` e `user` de forma atômica
  - aplica invariantes de estado
- `LoginUserUseCase`
  - valida credenciais
  - emite JWT + refresh token
  - persiste sessão em `user_sessions`
- `RefreshAccessTokenUseCase`
  - rotação de refresh token (revoga token antigo e persiste novo)
- `GetCurrentUserUseCase`
  - resolve usuário atual via contexto autenticado

Entidades/chaves de domínio:

- `User` com `Role` (`admin`, `customer`)
- `UserStatus` (`pending`, `active`, `blocked`)
- `RoleCustomer` exige `customer_id` não nulo (invariante de domínio)

### 6.2 Admin (`internal/admin`)

Implementado:

- `ApproveUserUseCase`
  - carrega usuário com lock (`FindByIDForUpdate`)
  - exige status `pending`
  - atualiza para `active`
  - valida existência do customer associado
  - gera número da conta
  - resolve agência via `BranchPolicy`
  - cria conta
  - tudo em transação única (atomicidade)

### 6.3 Customer (`internal/customer`)

Implementado e exposto:

- `GetCustomerMe` via `GET /customers/me`

Implementado internamente (não roteado no `main.go`):

- fluxo de criação de cliente por handler dedicado (`Create`) existe no pacote, mas não está exposto como endpoint público no roteador principal.

### 6.4 Account - Conta Bancária (`internal/account/bankaccount`)

Implementado:

- `ListAccounts`
- `CreateAccount`
- `GetAccountBalance`
- `LookupInternalTransferRecipients`

Regras relevantes:

- criação de conta deriva `customer_id` do usuário autenticado
- usuário precisa estar elegível para operar (ex.: ativo)
- políticas de acesso impedem atuação fora da própria conta/cliente
- listagem de contas recusa query params inesperados

### 6.5 Account - Transações (`internal/account/transaction`)

Implementado:

- `Deposit`
- `Withdraw`
- `Transfer`
- `GetTransferReceipt`

Características técnicas:

- operações de saldo em transações de banco
- lock de linha em conta(s) durante operação
- transferência trava duas contas em ordem determinística de UUID para reduzir risco de deadlock
- idempotência por `idempotency_key` no `transfer_out`
- em conflito de idempotência, resposta é reconstruída pelo ledger (replay determinístico)
- suporte a `description` em transações

### 6.6 Account - Extrato (`internal/account/statement`)

Implementado:

- `GetStatement`
- endpoint aceita:
  - `limit`
  - `cursor` (timestamp RFC3339)
  - `cursor_id` (UUID)
  - `from`
  - `to`
- retorno com `next_cursor` para paginação

## 7. Regras de Negócio e Invariantes Implementadas

- conta inicia com saldo zero e status `active`
- depósito/saque/transferência exigem valor positivo
- saque/transferência validam saldo suficiente
- transferência entre mesma conta é rejeitada
- operações monetárias só em contas ativas
- usuário `customer` deve manter vínculo com `customer_id`
- aprovação admin é pré-condição para ativação e abertura de conta em onboarding administrativo

## 8. Persistência e Modelo de Dados (PostgreSQL)

Tabelas principais implementadas:

- `customers`
- `users`
- `user_sessions`
- `accounts`
- `transactions`

Decisões já implementadas:

- ledger financeiro imutável (`transactions`) como fonte de verdade
- snapshot de saldo em `accounts.balance` para leitura eficiente
- trigger `prevent_transactions_mutation` bloqueando `UPDATE/DELETE` em `transactions`
- modelagem de transferência com par (`transfer_out`/`transfer_in`) compartilhando `reference_id`
- índice único parcial de idempotência (`account_id`, `idempotency_key`)
- índice único de integridade do par de transferência (`reference_id`, `type`)

## 9. Migrações Existentes

Na pasta `api/migrations`:

- `000001_init_schema.up.sql`
  - baseline com enums, tabelas, índices e trigger de imutabilidade
- `000002_account_number_key.up.sql`
  - troca unicidade para `(branch, number)` em contas
- `000003_transaction_description.up.sql`
  - adiciona coluna `description` em `transactions`

## 10. Catálogo de Erros e Mapeamento HTTP

Implementado em registro central (`bootstrap.RegisterErrors`) com agregação por módulo.

Códigos de erro já presentes:

- `INVALID_REQUEST`
- `INVALID_DATA`
- `ACCOUNT_NOT_FOUND`
- `INSUFFICIENT_FUNDS`
- `INVALID_AMOUNT`
- `FORBIDDEN`
- `CUSTOMER_NOT_FOUND`
- `ACCOUNT_INACTIVE`
- `SAME_ACCOUNT_TRANSFER`
- `USER_ALREADY_EXISTS`
- `INVALID_CREDENTIALS`
- `UNAUTHORIZED`
- `INVALID_TOKEN`
- `INTERNAL_ERROR`
- `INVALID_USER_STATE`
- `USER_NOT_FOUND`
- `USER_ALREADY_ACTIVE`
- `TRANSACTION_NOT_FOUND`

Mapeamentos observados incluem `400`, `401`, `403`, `404`, `409`, `422` e `500`.

## 11. Qualidade e Cobertura de Testes (estado atual)

Métricas rápidas do código da API:

- arquivos Go (`internal/` + `cmd/`): `117`
- arquivos de teste (`*_test.go`): `44`
- pacotes Go (`./...`): `30`

Há testes distribuídos por camadas:

- application (casos de uso)
- delivery (handlers/middlewares)
- infrastructure (repositórios)
- integração pontual (ex.: transações e auth)

## 12. Pontos de Atenção e Limites Atuais

- há documentação de erro com campo `details`, mas uso prático ainda é limitado
- handler de criação de cliente existe no código, porém não está exposto no roteamento principal
- sem versão explícita de API (`/v1`) nas rotas atuais
- observabilidade (correlation-id, tracing) não aparece como padrão obrigatório no fluxo atual

## 13. Conclusão

A API encontra-se funcionalmente robusta para o escopo de core bancário simplificado, com implementação consistente de autenticação, autorização, onboarding com aprovação, operações financeiras críticas e trilha contábil imutável.

O que está implementado hoje já cobre o ciclo essencial do produto:

- entrada do usuário (registro/login)
- aprovação administrativa
- criação e consulta de conta
- movimentação de saldo com garantias transacionais
- consulta de extrato/saldo
- recuperação de comprovante de transferência

Em termos de engenharia, os principais fundamentos esperados para domínio financeiro (atomicidade, controle de concorrência, idempotência, rastreabilidade por ledger e mapeamento padronizado de erros) estão presentes na implementação atual.
