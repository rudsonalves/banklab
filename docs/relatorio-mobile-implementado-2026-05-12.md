# Relatório Completo do Mobile Implementado

Data de referência: 12/05/2026  
Escopo analisado: `mobile/` (código-fonte Flutter, camada de dados, roteamento, docs e testes)

## 1. Resumo Executivo

O aplicativo mobile (BankFlow) está implementado em Flutter com arquitetura em camadas (core, data, domain e ui), injeção de dependências centralizada com `auto_injector`, navegação por `go_router`, comunicação HTTP via `dio` e persistência de sessão com `flutter_secure_storage`.

As jornadas implementadas hoje cobrem:

- autenticação (login, cadastro e login rápido com identidade em cache)
- bootstrap/splash com decisão de entrada por dados da última sessão
- home com conta selecionada e saldo
- transferência interna completa (busca de destinatário, pagamento, confirmação, status)
- consulta de extrato
- consulta de comprovante por referência

## 2. Arquitetura Implementada

### 2.1 Camadas

- **Core**: infraestrutura transversal (DI, roteamento, Result/Command, HTTP client, interceptadores, storage seguro, configuração de ambiente)
- **Data**: APIs, DTOs, envelopes de resposta e repositórios
- **Domain**: modelos de domínio app-facing e use cases de orquestração
- **UI**: páginas, widgets reutilizáveis, temas e viewmodels

### 2.2 Bootstrap da aplicação

Fluxo de inicialização:

1. `main.dart` chama `setupDependencies()`
2. `setupDependencies()` registra módulos na ordem:
   - `CoreServices`
   - `Services`
   - `Repositories`
   - `Usecases`
   - `Viewmodels`
3. `runApp(const AppWidget())`
4. `AppWidget` monta `MaterialApp.router` com `GoRouter`

### 2.3 Design system e tema

A aplicação utiliza tema customizado com:

- `MaterialTheme`
- tipografia por `createTextTheme(context, "Quicksand", "Google Sans")`
- ajustes de `AppBarTheme` e `InputDecorationTheme`
- variação por brilho de plataforma (light/dark)

## 3. Infraestrutura Técnica Implementada

## 3.1 Configuração de ambiente

`AppEnv` lê parâmetros por `--dart-define`:

- `BASE_URL` (obrigatório)
- `APP_ACCESS_TOKEN`
- `CONNECT_TIMEOUT`
- `RECEIVE_TIMEOUT`
- `APP_MODE` (`dev`, `staging`, `prod`)

Também existem arquivos de ambiente no projeto (`dev.env`, `staging.env`, `prod.env`).

## 3.2 HTTP client

- `RestClient` abstrato
- implementação: `DioRestClient`
- retorno padronizado por `Result<T>` (`Success`/`Failure`)
- mapeamento de erro HTTP em `AppError`

## 3.3 Sessão e autenticação automática

Implementado em `AuthInterceptor`:

- injeta `Authorization: Bearer <token>` quando existe token salvo
- em `401`, tenta refresh automaticamente (exceto quando já está em `/auth/refresh`)
- evita refresh concorrente com controle de operação em voo (`_refreshInFlight`)
- reexecuta requisição original quando refresh funciona
- limpa sessão (access/refresh) quando refresh falha

## 3.4 Armazenamento seguro

- abstração: `LocalSecureStorage`
- implementação: `FlutterSecureStorageLocalStorage`
- chaves usadas em sessão:
  - access token
  - refresh token

## 3.5 Modelo de estado e erro

### Result

A aplicação usa `Result<T>` para evitar exceções atravessando camadas:

- sucesso: `Result.success(value)`
- falha: `Result.failure(AppError)`

### Command

Viewmodels expõem comandos (`Command0`/`Command1`) com estados:

- `idle`
- `running`
- `success`
- `failure`

### Códigos de erro app

`AppErrorCode` já implementa:

- HTTP: `httpError`, `timeout`, `networkError`, `parsingError`, `unauthenticated`
- storage: `storageError`, `storageNotFound`, `storageConflict`, `storageCorrupted`, `storageExpired`
- validação/genérico: `invalidData`, `unexpected`

## 4. Rotas e Navegação Implementadas

Rotas registradas via `GoRouter`:

### Base

- `/splash`
- `/home`
- `/statement`

### Auth

- `/login`
- `/register`
- `/short-login`

### Transferência

- `/recipient`
- `/payment`
- `/confirmation`
- `/status/success`
- `/status/failure`

### Compartilhadas

- `/details`

Recursos de navegação já implementados:

- `extraCodec` para serialização de extras
- `routeObserver` para observação de ciclo de navegação
- `AppCustomTransactionPage` para transições customizadas em rotas de auth

## 5. UI, Páginas e ViewModels Implementados

Páginas implementadas (11):

- Splash
- Login
- Register
- Short Login
- Home
- Statement
- Transfer Recipient
- Transfer Payment
- Transfer Confirmation
- Transfer Status
- Details

Viewmodels implementados (8):

- `SplashViewmodel`
- `LoginViewModel`
- `RegisterViewmodel`
- `ShortLoginViewModel`
- `HomeViewmodel`
- `TransferViewmodel`
- `StatementViewmodel`
- `DetailsViewmodel`

Widgets/componentes compartilhados relevantes:

- scaffold seguro (`SafeScaffold`)
- campos e formatadores (`BasicTextFormField`, `CpfInputFormatter`, `MoneyInputFormatter`)
- cartões (`AccountCard`, `BalanceCard`, `RecipientCard`)
- feedback (`AppSnackbar`)
- componentes de extrato (`day_header`, `month_header`, `statement_item_card`)

## 6. Camada de Dados e Integração com API

## 6.1 Serviços de API implementados

### Auth

- `POST /auth/register` (com `X-App-Token`)
- `POST /auth/login` (com `X-App-Token`)
- `GET /customers/me`
- `GET /auth/me`

### Conta

- `GET /accounts`
- `GET /accounts/{accountId}/balance`
- `GET /accounts/{accountId}/statement`

### Transferência

- `POST /accounts/internal-transfers`
- `GET /accounts/internal-transfers/recipients`
- `GET /accounts/transfer/{transactionReference}/receipt`

Obs.: o app mobile atual consome esses endpoints; não há cliente implementado para depósito/saque no recorte atual.

## 6.2 Envelopes e parsing

- uso de `ApiEnvelope<T>` para parse de `data/error`
- DTOs dedicados por contexto de endpoint
- validação de status HTTP (faixa 2xx)
- fallback de erro de parsing (`AppErrorCode.parsingError`)

## 6.3 Repositórios implementados

- `AuthRepositoryImpl`
- `AccountRepositoryImpl`
- `TransactionRepositoryImpl`

Comportamentos importantes:

- cache de conta selecionada e saldo
- stream de saldo (`balance()`) para atualização reativa
- cache do último extrato, último comprovante e última transferência
- persistência de tokens no login/logout
- cache de última identidade de login para fluxo de short login

## 7. Casos de Uso de Domínio Implementados

Use cases registrados:

- `TransferUsecase`
- `DetailsUsecase`

### TransferUsecase

Coordena:

- seleção de conta
- transferência com validações
- idempotência no cliente via geração de `Uuid().v7()` no `TransferViewmodel`
- consulta de destinatários internos
- consulta de comprovante

### DetailsUsecase

Coordena:

- leitura de conta selecionada
- consulta de comprovante por referência

## 8. Jornadas Funcionais Implementadas

## 8.1 Autenticação

- login padrão
- cadastro de usuário
- login rápido (short login) com dados de identidade salvos localmente
- obtenção de perfil combinando `/customers/me` + `/auth/me`

## 8.2 Home e saldo

- listagem de contas
- seleção de conta
- carga de saldo da conta selecionada
- atualização periódica de saldo (timer de 20s no `HomeViewmodel`)

## 8.3 Transferência interna

- busca de destinatário por critérios
- preenchimento de pagamento
- confirmação
- execução da transferência
- tela de status (sucesso/falha)
- recuperação de comprovante por referência

## 8.4 Extrato

- consulta do extrato da conta selecionada
- suporte a query params:
  - `limit`
  - `cursor`
  - `cursor_id`
  - `from`
  - `to`
- parse de `next_cursor` para paginação

## 9. Qualidade e Testes (estado atual)

Métricas levantadas no workspace mobile:

- arquivos Dart (`lib` + `test`): `146`
- arquivos de teste (`*_test.dart`): `16`
- páginas (`*_page.dart`): `11`
- viewmodels (`*_viewmodel.dart`): `8`

Há cobertura de testes em camadas de core, data e UI, com foco em comportamento de parsing, contratos de APIs/repositórios e fluxo de estado de componentes.

## 10. Restrições e Lacunas Atuais

Pontos observados no estado atual da implementação:

- os use cases de domínio estão concentrados em transferência e detalhes; outras jornadas seguem repositório direto via viewmodel
- o `AuthApi.getProfile()` está marcado com observação para extração futura para serviço dedicado de profile
- apesar de a API backend suportar depósito/saque, o cliente mobile atual não expõe essas operações no recorte analisado
- parte da tipagem de erro de APIs converte `error.code` do backend para `AppErrorCode.httpError`, reduzindo granularidade semântica no app

## 11. Conclusão

O mobile está implementado com base arquitetural sólida e consistente para o escopo atual, com boa separação de responsabilidades e fluxo previsível de estado.

O que já está entregue end-to-end no app:

- autenticação e sessão com refresh automático
- carregamento de perfil e contexto de usuário
- seleção de conta e atualização de saldo
- transferência interna completa com comprovante
- consulta de extrato

Em termos de prontidão funcional, o app atende o ciclo principal de uso para conta e transferência, mantendo fundamentos técnicos relevantes: DI estruturada, navegação tipada por rotas conhecidas, tratamento explícito de erros por `Result` e persistência segura de credenciais.
