# Backlog: senha transacional e step-up no mobile

## 1. Contexto

O backlog `api/006` fechou o contrato do MVP ZTA na API:

- o usuário cria uma senha transacional em `POST /security/transaction-password`;
- o app autoriza uma operação sensível em
  `POST /security/step-up/authorize`;
- a API retorna um `step_up_token` curto, de uso único, emitido para o
  endpoint público `POST /accounts/internal-transfers`;
- `POST /accounts/internal-transfers` exige o header `X-Step-Up-Token`.

Este backlog define os elementos básicos para implementar esse fluxo no mobile,
sem alterar o contrato da API.

## 2. Objetivo

Permitir que o usuário cadastre sua senha transacional e use essa credencial
para autorizar uma transferência interna antes da execução final.

O fluxo mobile deve:

- manter a senha transacional fora do endpoint de transferência;
- solicitar step-up apenas no momento de confirmar uma operação sensível;
- enviar `X-Step-Up-Token` na transferência interna;
- tratar os erros ZTA por `error.code`;
- respeitar a arquitetura mobile atual.

## 3. Contratos de API

### 3.1 Criar senha transacional

Finalidade:

Criar a primeira senha transacional do usuário autenticado. Esta chamada é
usada no setup inicial da credencial e não exige step-up prévio, porque a
credencial ainda não existe.

Autenticação:

- exige JWT válido no header `Authorization`;
- não usa `X-Step-Up-Token`;
- não deve receber senha de login.

```http
POST /security/transaction-password
Authorization: Bearer <access_token>
```

Payload:

- `transaction_password`: PIN numérico de 6 dígitos;
- `transaction_password_confirmation`: confirmação local do mesmo PIN.

Request:

```json
{
  "transaction_password": "123456",
  "transaction_password_confirmation": "123456"
}
```

Resposta:

- retorna somente metadados da credencial criada;
- não retorna PIN, hash ou qualquer material sensível;
- o mobile deve tratar sucesso como confirmação de que o usuário já possui
  senha transacional ativa.

Resposta de sucesso:

```json
{
  "data": {
    "user_id": "<user_id>",
    "status": "active",
    "created_at": "2026-05-29T10:00:00Z"
  },
  "error": null
}
```

Erros relevantes para o mobile:

- `TRANSACTION_PASSWORD_ALREADY_SET`: usuário já possui senha transacional;
- `INVALID_DATA`: PIN inválido ou confirmação divergente;
- `INVALID_REQUEST`: JSON inválido ou campo inesperado;
- `UNAUTHORIZED` / `INVALID_TOKEN`: sessão ausente ou inválida;
- `FORBIDDEN`: usuário autenticado não está apto para a operação.

### 3.2 Autorizar step-up

Finalidade:

Validar a senha transacional para uma operação HTTP pública sensível e receber
um `step_up_token` curto. No MVP, a operação sensível usada pelo mobile é
`POST /accounts/internal-transfers`.

Autenticação:

- exige JWT válido no header `Authorization`;
- não usa `X-Step-Up-Token`;
- recebe a senha transacional apenas nesta chamada;
- o mobile não deve tentar interpretar o JWT retornado.

```http
POST /security/step-up/authorize
Authorization: Bearer <access_token>
```

Payload:

- `method`: verbo HTTP público da operação sensível;
- `path`: path público da operação sensível;
- `transaction_password`: PIN transacional informado no prompt de step-up.

Observação para o mobile:

O mobile deve pedir autorização para a superfície pública que irá acessar. Ele
não deve conhecer chaves internas de policy da API, como
`internal_transfer.create`.
`endpoint_key` pode continuar existindo internamente no backend/JWT, mas não é
campo de input público do contrato mobile.

Request:

```json
{
  "method": "POST",
  "path": "/accounts/internal-transfers",
  "transaction_password": "123456"
}
```

Resposta:

- `step_up_token`: token opaco para o mobile;
- `expires_in`: tempo de validade em segundos, atualmente `120`;
- o token é de uso único e deve ficar somente em memória até a tentativa da
  operação sensível.

Resposta de sucesso:

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 120
  },
  "error": null
}
```

Erros relevantes para o mobile:

- `TRANSACTION_PASSWORD_NOT_SET`: usuário ainda não cadastrou senha
  transacional;
- `TRANSACTION_PASSWORD_INVALID`: PIN transacional incorreto;
- `TRANSACTION_PASSWORD_LOCKED`: PIN bloqueado temporariamente por tentativas
  inválidas;
- `STEP_UP_ENDPOINT_NOT_ALLOWED`: operação pública não permitida pela API para
  step-up;
- `INVALID_DATA`: PIN com formato inválido;
- `INVALID_REQUEST`: JSON inválido ou campo inesperado;
- `UNAUTHORIZED` / `INVALID_TOKEN`: sessão ausente ou inválida;
- `FORBIDDEN`: usuário autenticado não está apto para a operação.

### 3.3 Executar transferência interna

Finalidade:

Executar a transferência interna após o usuário confirmar a intenção com senha
transacional. A transferência continua sendo ID-based; branch e número de conta
podem aparecer no fluxo de busca/confirmação, mas não são identificadores de
execução.

Autenticação e step-up:

- exige JWT válido no header `Authorization`;
- exige `X-Step-Up-Token` emitido para `POST /accounts/internal-transfers`;
- o token deve ser enviado somente nesta chamada protegida;
- a senha transacional nunca deve ser enviada no payload da transferência.

```http
POST /accounts/internal-transfers
Authorization: Bearer <access_token>
X-Step-Up-Token: <step_up_token>
```

O payload da transferência continua usando `from_account_id`, `to_account_id`,
`amount`, `idempotency_key` e `description` opcional.

Payload:

```json
{
  "from_account_id": "<source_account_id>",
  "to_account_id": "<destination_account_id>",
  "amount": 2500,
  "idempotency_key": "<stable_transfer_attempt_key>",
  "description": "Aluguel de maio"
}
```

Resposta:

- retorna os dados públicos da transferência executada;
- inclui `transaction_reference`, usado para abrir o comprovante;
- não retorna o step-up token nem dados da senha transacional.

Erros relevantes para o mobile:

- `STEP_UP_TOKEN_REQUIRED`: header ausente;
- `STEP_UP_TOKEN_INVALID`: token inválido, malformado ou divergente da
  persistência;
- `STEP_UP_TOKEN_EXPIRED`: token expirado;
- `STEP_UP_TOKEN_CONSUMED`: token já usado; o app deve solicitar novo step-up;
- `STEP_UP_ENDPOINT_MISMATCH`: token emitido para outra operação pública;
- erros de negócio da transferência, como `INSUFFICIENT_FUNDS`,
  `ACCOUNT_INACTIVE`, `ACCOUNT_NOT_FOUND`, `FORBIDDEN`, `INVALID_AMOUNT`,
  `INVALID_DATA` e `SAME_ACCOUNT_TRANSFER`.

Ordem esperada no mobile:

1. preservar ou gerar `idempotency_key` da tentativa lógica;
2. solicitar senha transacional;
3. obter `step_up_token`;
4. executar transferência com `X-Step-Up-Token`;
5. descartar `step_up_token` em sucesso, erro ou cancelamento.

## 4. Fluxo de produto

### 4.1 Cadastro da senha transacional

Fluxo inicial:

1. O usuário autenticado acessa a criação de senha transacional.
2. O app solicita PIN de 6 dígitos.
3. O app solicita confirmação do PIN.
4. O app valida formato e confirmação localmente.
5. O app chama `POST /security/transaction-password`.
6. Em sucesso, o app volta ao fluxo anterior ou mostra status de senha criada.

Decisão inicial:

- não criar recuperação, troca ou reset de senha transacional neste backlog;
- não armazenar a senha transacional no dispositivo;
- não exibir senha transacional em logs, analytics ou mensagens de erro.

### 4.2 Step-up na transferência interna

Fluxo esperado:

1. Usuário preenche a transferência como hoje.
2. Na confirmação, antes de executar a transferência, o app solicita a senha
   transacional.
3. O app chama `POST /security/step-up/authorize` para autorizar
   `method=POST` e `path=/accounts/internal-transfers`.
4. Se a autorização retornar `step_up_token`, o app chama
   `POST /accounts/internal-transfers` com `X-Step-Up-Token`.
5. O token deve ser descartado após a tentativa de transferência.
6. Se a transferência falhar depois do enforcement, o mesmo token não deve ser
   reutilizado.

Retry:

- se a API retornar `STEP_UP_TOKEN_CONSUMED`, o app deve solicitar novo step-up;
- mesmo mantendo o mesmo `idempotency_key`, o app pode precisar de novo
  `step_up_token`;
- o `idempotency_key` deve continuar representando a tentativa lógica de
  transferência, não o token de step-up.

## 5. Impactos na arquitetura mobile

### 5.1 Data layer

Adicionar API service para o módulo de segurança, seguindo o padrão de
`mobile/lib/data/services/apis`:

- `SecurityApi` ou equivalente;
- DTO de criação de senha transacional;
- DTO de autorização step-up;
- DTO de resposta com `step_up_token` e `expires_in`.

Atualizar o serviço de transferência:

- permitir enviar headers por chamada;
- incluir `X-Step-Up-Token` somente no `POST /accounts/internal-transfers`;
- não enviar senha transacional no payload da transferência.

### 5.2 Repository layer

Adicionar repositório de segurança:

- criar senha transacional;
- autorizar step-up;
- expor uma operação para autorização de transferência interna usando método e
  path públicos;
- mapear erros ZTA para erros de aplicação.

O repositório de transação deve continuar responsável por executar a
transferência, mas receber o step-up token como parâmetro de execução ou por um
objeto de input explícito.

### 5.3 Domain/use case layer

Criar ou ajustar use case para coordenar transferência protegida:

```text
UI confirma transferencia
-> solicita senha transacional
-> authorizeStepUp(method=POST, path=/accounts/internal-transfers, senha)
-> transfer(..., step_up_token)
-> descarta token
```

O use case deve manter a orquestração fora da UI e evitar que view models
conheçam detalhes de múltiplos repositórios.

O mobile deve conhecer apenas método e path públicos da API. Chaves internas de
policy, como `internal_transfer.create`, são responsabilidade exclusiva da API.

### 5.4 UI layer

Elementos mínimos:

- tela ou fluxo para cadastrar senha transacional;
- componente de entrada de PIN de 6 dígitos;
- confirmação do PIN no cadastro;
- prompt/modal/página de step-up na confirmação da transferência;
- estados de loading, sucesso e erro;
- mensagens específicas para bloqueio, senha ausente e senha inválida.

O prompt de step-up deve ser transitório e não deve preservar a senha em route
extras, cache local ou estado global.

## 6. Erros que o mobile deve tratar

Clientes devem depender de `error.code`.

Criação da senha transacional:

- `TRANSACTION_PASSWORD_ALREADY_SET`;
- `INVALID_DATA`;
- `INVALID_REQUEST`;
- `UNAUTHORIZED`;
- `INVALID_TOKEN`;
- `FORBIDDEN`.

Autorização de step-up:

- `TRANSACTION_PASSWORD_NOT_SET`;
- `TRANSACTION_PASSWORD_INVALID`;
- `TRANSACTION_PASSWORD_LOCKED`;
- `STEP_UP_ENDPOINT_NOT_ALLOWED`;
- `INVALID_DATA`;
- `INVALID_REQUEST`;
- `UNAUTHORIZED`;
- `INVALID_TOKEN`;
- `FORBIDDEN`.

Transferência protegida:

- `STEP_UP_TOKEN_REQUIRED`;
- `STEP_UP_TOKEN_INVALID`;
- `STEP_UP_TOKEN_EXPIRED`;
- `STEP_UP_TOKEN_CONSUMED`;
- `STEP_UP_ENDPOINT_MISMATCH`;
- erros de transferência já existentes, como `INSUFFICIENT_FUNDS`,
  `ACCOUNT_INACTIVE`, `ACCOUNT_NOT_FOUND`, `FORBIDDEN`, `INVALID_AMOUNT`,
  `INVALID_DATA` e `SAME_ACCOUNT_TRANSFER`.

## 7. UX mínima

### Cadastro

- PIN de 6 dígitos.
- Confirmação do PIN.
- Validação local antes de chamar a API.
- Feedback claro quando a senha já existe.

### Step-up

- Entrada curta e focada da senha transacional.
- Não mostrar o valor digitado.
- Permitir cancelar e voltar para a confirmação da transferência.
- Em senha inválida, permitir nova tentativa enquanto a API permitir.
- Em senha bloqueada, informar bloqueio temporário.
- Em senha inexistente, direcionar para cadastro de senha transacional.

## 8. Segurança e privacidade

- Não persistir senha transacional.
- Não logar senha transacional.
- Não logar `step_up_token` completo.
- Não armazenar `step_up_token` em secure storage.
- Manter o token apenas em memória pelo menor tempo possível.
- Descartar o token após sucesso, erro ou cancelamento.
- Usar somente `error.code` para decisão de fluxo.

## 9. Fora do escopo inicial

- recuperação de senha transacional;
- troca de senha transacional;
- reset administrativo;
- biometria local como substituta da senha transacional;
- dispositivo confiável;
- vinculação do step-up token ao payload detalhado da operação;
- suporte a outros endpoints sensíveis além de `POST /accounts/internal-transfers`.

## 10. Critérios de aceite

- Usuário consegue criar senha transacional pelo app.
- App não envia senha transacional no endpoint de transferência.
- App autoriza step-up antes de executar transferência interna.
- Transferência interna envia `X-Step-Up-Token`.
- Token é descartado após a tentativa.
- Retry com `STEP_UP_TOKEN_CONSUMED` solicita novo step-up.
- Mesmo `idempotency_key` pode ser reutilizado com novo step-up token quando
  fizer sentido para a mesma tentativa lógica.
- Erros ZTA são tratados por `error.code`.
- Testes cobrem DTOs, API services, repositórios, use cases e view models
  afetados.
- Documentação mobile é atualizada ao final da implementação.

## 11. Referências

- `docs/backlogs/api/006a - transaction-password.md`
- `docs/backlogs/api/006b - step-up-token.md`
- `docs/backlogs/api/006c - internal-transfer-step-up-enforcement.md`
- `docs/backlogs/api/006d - zta-contracts-and-docs.md`
- `api/docs/07-api-rest.md`
- `api/docs/implementations/03-zta-step-up-transaction-password.md`
