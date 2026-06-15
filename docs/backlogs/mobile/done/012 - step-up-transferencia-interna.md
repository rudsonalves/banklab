# Backlog: step-up na transferência interna mobile

## 1. Contexto

O backlog `api/006` fechou o contrato do MVP ZTA na API para autorização de
operações sensíveis:

- o app autoriza uma operação sensível em
  `POST /security/step-up/authorize`;
- a API retorna um `step_up_token` curto, de uso único, emitido para o endpoint
  público `POST /accounts/internal-transfers`;
- `POST /accounts/internal-transfers` exige o header `X-Step-Up-Token`;
- a senha transacional é enviada somente para o endpoint de autorização
  step-up, nunca no payload da transferência.

Este backlog trata somente do uso da senha transacional para autorizar a
transferência interna.

O cadastro da senha transacional foi implementado pela backlog concluída
`done/011 - cadastro-senha-transacional.md`.

O fluxo existente de cadastro navega para a Home após sucesso. Para reutilizá-lo
a partir da transferência, este backlog deve torná-lo consciente de uma intenção
de retorno não sensível, sem duplicar as páginas de criação.

O acesso à transferência deve verificar o estado atual da senha transacional
quando o usuário tocar em **Transferir**. A transferência só pode ser aberta
quando
`AppSection.currentSession.readiness.transactionPasswordStatus` for `active`.

## 2. Objetivo

Permitir que o usuário autorize uma transferência interna com senha
transacional antes da execução final.

O fluxo mobile deve:

- manter a senha transacional fora do endpoint de transferência;
- verificar a existência de senha transacional ativa antes de abrir o fluxo de
  transferência;
- reutilizar o cadastro existente quando o estado for `not_set`;
- solicitar step-up apenas no momento de confirmar a operação sensível;
- chamar `POST /security/step-up/authorize` para obter um `step_up_token`;
- enviar `X-Step-Up-Token` em `POST /accounts/internal-transfers`;
- descartar o token após sucesso, erro ou cancelamento;
- tratar erros ZTA por `error.code`;
- preservar a `idempotency_key` durante retries da mesma tentativa lógica;
- impedir que erros funcionais de step-up com HTTP 401 acionem indevidamente a
  renovação da sessão;
- impedir que senha transacional, token step-up e headers de autenticação sejam
  registrados em logs;
- respeitar a arquitetura mobile atual.

## 3. Contratos de API

### 3.1 Autorizar step-up

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

O mobile deve pedir autorização para a superfície pública que irá acessar.

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

- `TRANSACTION_PASSWORD_INVALID`: PIN transacional incorreto;
- `TRANSACTION_PASSWORD_LOCKED`: PIN bloqueado temporariamente por tentativas
  inválidas;
- `STEP_UP_ENDPOINT_NOT_ALLOWED`: operação pública não permitida pela API para
  step-up;
- `INVALID_DATA`: PIN com formato inválido;
- `INVALID_REQUEST`: JSON inválido ou campo inesperado;
- `UNAUTHORIZED` / `INVALID_TOKEN`: sessão ausente ou inválida;
- `FORBIDDEN`: usuário autenticado não está apto para a operação.

`TRANSACTION_PASSWORD_NOT_SET` continua existindo no contrato defensivo da API,
mas não é um resultado esperado deste fluxo mobile. Se ocorrer após o usuário
ter entrado na transferência, representa inconsistência entre a verificação de
entrada e o estado atual da credencial.

### 3.2 Executar transferência interna

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

1. gerar a `idempotency_key` ao criar a tentativa confirmada de transferência;
2. solicitar senha transacional;
3. obter `step_up_token`;
4. executar transferência com `X-Step-Up-Token`;
5. descartar `step_up_token` em sucesso, erro ou cancelamento.

A mesma `idempotency_key` deve permanecer associada à tentativa enquanto o
usuário estiver na confirmação e precisar repetir apenas o step-up. Uma nova
chave deve ser gerada somente quando o usuário abandonar a tentativa ou criar
uma nova transferência.

## 4. Fluxo de produto

### 4.1 Entrada da transferência

Ao tocar em **Transferir**, antes de abrir a seleção de destinatário:

1. O app consulta o snapshot de sessão já carregado no login em
   `AppSection.currentSession`.
2. Se `transaction_password_status=active`, abre
   `TransferRoutes.recipient`.
3. Se `transaction_password_status=not_set`, abre o fluxo existente de criação
   de senha transacional com intenção de retorno para a transferência.
4. Se o status for `locked` ou `unknown`, não abre a transferência e apresenta
   mensagem apropriada.
5. Se o snapshot estiver ausente, mantém o usuário na Home e apresenta erro de
   sessão, sem iniciar cadastro ou transferência.

O fluxo de cadastro deve aceitar uma origem/destino tipado e não sensível:

- quando iniciado pelo gate pós-login, sucesso mantém o comportamento de seguir
  para a Home;
- quando iniciado pelo botão **Transferir**, sucesso com resposta `active`
  atualiza localmente o `AppSection` e abre `TransferRoutes.recipient`;
- se o `AppSection` não confirmar status `active`, não abre a transferência;
- `TRANSACTION_PASSWORD_ALREADY_SET` é tratado como inconsistência defensiva:
  permanece no cadastro, sem navegar e sem atualizar localmente a sessão;
- cancelamento do cadastro iniciado pela Home retorna à Home;
- falha na criação mantém o usuário no cadastro e permite nova tentativa.

O destino não deve ser recebido como path ou nome de rota arbitrário. Deve ser
representado por enum ou objeto fechado, por exemplo `postLogin` e
`internalTransfer`, e propagado entre as páginas sem incluir PIN, token ou
draft de transferência.

### 4.2 Autorização e execução

Fluxo esperado:

1. Após a verificação de entrada, o usuário preenche a transferência como hoje.
2. Na confirmação, antes de executar a transferência, o app solicita a senha
   transacional.
3. O app chama `POST /security/step-up/authorize` para autorizar
   `method=POST` e `path=/accounts/internal-transfers`.
4. Após a autorização retornar o `step_up_token`, o app chama obrigatoriamente
   `POST /accounts/internal-transfers` com o header `X-Step-Up-Token`. Se a
   autorização falhar, a transferência não deve ser chamada.
5. O token deve ser descartado após a tentativa de transferência.
6. Se a transferência falhar depois do enforcement, o mesmo token não deve ser
   reutilizado.

Se, defensivamente, a autorização retornar
`TRANSACTION_PASSWORD_NOT_SET`, o app deve:

- interromper a transferência sem chamar o endpoint protegido;
- descartar PIN, token, draft e `idempotency_key`;
- atualizar localmente o `AppSection` para `not_set`, usando o erro do backend
  como informação autoritativa;
- retornar à Home;
- permitir que um novo toque em **Transferir** encaminhe ao cadastro caso o
  estado atualizado seja `not_set`;
- não abrir o cadastro diretamente a partir da confirmação da transferência.

### Retry

- se a API retornar `STEP_UP_TOKEN_CONSUMED`, o app deve solicitar novo step-up;
- se a API retornar `STEP_UP_TOKEN_EXPIRED`, o app deve solicitar novo step-up;
- mesmo mantendo o mesmo `idempotency_key`, o app pode precisar de novo
  `step_up_token`;
- o `idempotency_key` deve continuar representando a tentativa lógica de
  transferência, não o token de step-up.
- `STEP_UP_TOKEN_CONSUMED` e `STEP_UP_TOKEN_EXPIRED` não devem levar o usuário
  à tela genérica de falha: o app deve permanecer ou retornar à confirmação,
  abrir novo prompt de PIN e repetir a autorização;
- o retry não deve reenviar automaticamente a senha transacional nem reutilizar
  o token anterior.

## 5. Impactos na arquitetura mobile

### 5.1 Data layer

Ampliar o módulo existente de senha transacional, seguindo o padrão de
`mobile/lib/data/services/apis/transaction_password`:

- adicionar `authorizeStepUp` ao `TransactionPasswordApi` ou nome equivalente
  dentro do mesmo módulo;
- DTO de autorização step-up;
- DTO de resposta com `step_up_token` e `expires_in`.

Atualizar o serviço de transferência:

- permitir enviar headers por chamada;
- incluir `X-Step-Up-Token` somente no `POST /accounts/internal-transfers`;
- não enviar senha transacional no payload da transferência.

O parsing de erros deve preservar o `error.code` original da API em
`AppError.details`, compatível com o helper existente `backendErrorCode(error)`.
Não é requisito criar um valor de `AppErrorCode` para cada erro ZTA.

### 5.2 Cliente HTTP

O cliente HTTP compartilhado deve ser ajustado para o contrato de step-up:

- o `AuthInterceptor` só deve tentar refresh para respostas HTTP 401 cujo
  `error.code` seja `UNAUTHORIZED` ou `INVALID_TOKEN`;
- `TRANSACTION_PASSWORD_INVALID`, `STEP_UP_TOKEN_REQUIRED`,
  `STEP_UP_TOKEN_INVALID`, `STEP_UP_TOKEN_EXPIRED` e
  `STEP_UP_TOKEN_CONSUMED` devem chegar ao fluxo chamador sem refresh;
- respostas 401 desconhecidas não devem ser usadas para mascarar erros
  funcionais com uma tentativa automática de refresh;
- logs de requests devem omitir ou mascarar `Authorization`,
  `X-Step-Up-Token` e quaisquer valores de senha transacional;
- nenhum log deve serializar o body de autorização step-up.

### 5.3 Repository layer

Ampliar o repositório existente de senha transacional:

- autorizar step-up;
- expor uma operação para autorização de transferência interna usando método e
  path públicos;
- preservar códigos de erro ZTA para decisão do fluxo.

O repositório de transação deve continuar responsável por executar a
transferência, mas deve receber o step-up token como parâmetro de execução ou
por um objeto de input explícito.

O snapshot autenticado deve ficar centralizado no `AppSection`:

- o `AuthRepository` carrega `GET /auth/session` no login e preenche o
  `AppSection`;
- `getAuthSession()` pode retornar o snapshot já carregado;
- o logout limpa o `AppSection`;
- após criação bem-sucedida com resposta `active`, o
  `TransactionPasswordRepository` atualiza localmente o status no
  `AppSection`, sem nova chamada a `GET /auth/session`;
- `TRANSACTION_PASSWORD_ALREADY_SET` não atualiza o snapshot nem libera
  navegação;
- `TRANSACTION_PASSWORD_NOT_SET` durante o step-up atualiza localmente o
  `AppSection` para `not_set`, sem nova chamada de sessão;
- a entrada normal da transferência não chama novamente `GET /auth/session`.

### 5.4 Domain/use case layer

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

O use case também deve:

- receber ou criar uma tentativa com `idempotency_key` estável antes do
  step-up;
- preservar essa chave ao repetir somente a autorização;
- manter o `step_up_token` em variável local pelo tempo necessário para uma
  única chamada de transferência;
- descartar o token em bloco equivalente a `finally`;
- nunca repetir automaticamente a autorização com o PIN anterior.

A verificação de entrada da transferência também deve ser exposta por uma
operação do view model ou use case, evitando que a `HomePage` consulte
diretamente repositórios ou interprete o snapshot da sessão.

### 5.5 UI layer

Elementos mínimos:

- verificação assíncrona ao tocar em **Transferir**, com bloqueio contra toques
  concorrentes;
- reaproveitamento das páginas existentes de criação e confirmação da senha
  transacional;
- intenção tipada para definir Home ou transferência como destino após o
  cadastro;
- prompt/modal/página de step-up na confirmação da transferência;
- entrada de PIN de 6 dígitos;
- estados de loading, sucesso e erro;
- mensagens específicas para bloqueio, senha ausente e senha inválida;
- cancelamento do step-up com retorno para a confirmação da transferência.

O prompt de step-up deve ser transitório e não deve preservar a senha em route
extras, cache local ou estado global.

A intenção de retorno do cadastro pode trafegar entre rotas por ser não
sensível, mas não deve aceitar paths arbitrários. O PIN atual continua
transitório entre as duas páginas conforme o fluxo existente e deve ser limpo
ao concluir, cancelar ou falhar uma tentativa de criação.

A confirmação deve bloquear submissões concorrentes enquanto autorização ou
transferência estiverem em andamento. Em token consumido ou expirado, ela deve
solicitar novo PIN sem recriar a tentativa e sem navegar para a tela genérica de
falha.

## 6. Erros que o mobile deve tratar

Clientes devem depender de `error.code`.

No mobile, a decisão deve usar o código preservado no `AppError`, por exemplo
por meio de `backendErrorCode(error)`, e não comparar mensagens ou depender
somente do status HTTP.

Autorização de step-up:

- `TRANSACTION_PASSWORD_INVALID`;
- `TRANSACTION_PASSWORD_LOCKED`;
- `STEP_UP_ENDPOINT_NOT_ALLOWED`;
- `INVALID_DATA`;
- `INVALID_REQUEST`;
- `UNAUTHORIZED`;
- `INVALID_TOKEN`;
- `FORBIDDEN`.

Inconsistência defensiva:

- `TRANSACTION_PASSWORD_NOT_SET`: interromper a transferência, executar a
  atualização local do `AppSection` e retornar à Home; o próximo acesso à
  transferência pode iniciar o cadastro a partir da verificação de entrada.

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

- Entrada curta e focada da senha transacional.
- Não mostrar o valor digitado.
- Permitir cancelar e voltar para a confirmação da transferência.
- Em senha inválida, permitir nova tentativa enquanto a API permitir.
- Em senha bloqueada, informar bloqueio temporário.
- Em token expirado ou consumido, manter a transferência confirmada e solicitar
  novamente a senha transacional.
- Não exibir a tela genérica de falha antes de oferecer o novo step-up nesses
  dois casos.
- Em `TRANSACTION_PASSWORD_NOT_SET`, encerrar a tentativa, atualizar localmente
  o `AppSection` e retornar à Home. O cadastro só pode ser aberto por uma nova
  verificação no botão **Transferir**.

## 8. Segurança e privacidade

- Não persistir senha transacional.
- Não logar senha transacional.
- Não logar `step_up_token` completo.
- Não logar os headers `Authorization` e `X-Step-Up-Token`.
- Não logar o body de `POST /security/step-up/authorize`.
- Não armazenar `step_up_token` em secure storage.
- Manter o token apenas em memória pelo menor tempo possível.
- Descartar o token após sucesso, erro ou cancelamento.
- Usar somente `error.code` para decisão de fluxo.
- Não executar refresh de sessão para erros funcionais de step-up com HTTP 401.

## 9. Fora do escopo inicial

- reimplementação completa do cadastro da senha transacional, além dos ajustes
  de destino e sincronização do `AppSection` necessários para reutilizá-lo;
- recuperação de senha transacional;
- troca de senha transacional;
- reset administrativo;
- biometria local como substituta da senha transacional;
- dispositivo confiável;
- vinculação do step-up token ao payload detalhado da operação;
- suporte a outros endpoints sensíveis além de `POST /accounts/internal-transfers`.

## 10. Critérios de aceite

- App não envia senha transacional no endpoint de transferência.
- Botão **Transferir** consulta o snapshot já carregado no login em
  `AppSection.currentSession`, sem nova chamada a `/auth/session`.
- Status `active` abre a seleção de destinatário.
- Status `not_set` abre o cadastro existente com destino tipado para
  transferência.
- Status `locked`, `unknown` ou falha de sessão não abre a transferência.
- Cadastro iniciado pelo botão **Transferir** atualiza localmente o
  `AppSection` a partir da resposta `active` e somente então abre a seleção de
  destinatário.
- Cadastro iniciado pelo gate pós-login mantém a Home como destino.
- Falha de criação mantém o usuário no cadastro; cancelamento iniciado pela
  Home retorna à Home.
- `TRANSACTION_PASSWORD_ALREADY_SET` permanece no cadastro, sem atualizar o
  `AppSection` e sem seguir ao destino.
- App autoriza step-up antes de executar transferência interna.
- Transferência interna envia `X-Step-Up-Token`.
- Token é descartado após a tentativa.
- Retry com `STEP_UP_TOKEN_CONSUMED` solicita novo step-up.
- Retry com `STEP_UP_TOKEN_EXPIRED` solicita novo step-up.
- Mesmo `idempotency_key` pode ser reutilizado com novo step-up token quando
  fizer sentido para a mesma tentativa lógica.
- `TRANSACTION_PASSWORD_NOT_SET` durante step-up abandona a tentativa e retorna
  à Home; somente um novo toque em **Transferir** pode iniciar o cadastro após
  a atualização local do `AppSection`.
- Erros ZTA são tratados por `error.code`.
- `AuthInterceptor` não tenta refresh para erros funcionais de step-up com HTTP
  401.
- Requests não registram senha transacional, `Authorization` ou
  `X-Step-Up-Token`.
- `idempotency_key` é criada antes do step-up e permanece estável ao repetir a
  autorização da mesma tentativa.
- `STEP_UP_TOKEN_CONSUMED` e `STEP_UP_TOKEN_EXPIRED` reabrem o step-up na
  confirmação, sem passar pela tela genérica de falha.
- Submissões concorrentes são bloqueadas durante autorização e transferência.
- Testes cobrem DTOs, API services, repositórios, use cases e view models
  afetados.
- Testes cobrem os destinos pós-cadastro `postLogin` e `internalTransfer`,
  inclusive sucesso, cancelamento, falha e o tratamento defensivo de
  `TRANSACTION_PASSWORD_ALREADY_SET`.
- Testes comprovam que a entrada da transferência usa o `AppSection` carregado
  no login e que o repository atualiza seu status após criação bem-sucedida.
- Testes do interceptor cobrem refresh para `UNAUTHORIZED`/`INVALID_TOKEN` e
  ausência de refresh para os códigos funcionais de step-up.
- Testes verificam sanitização de headers sensíveis em logs.
- Documentação mobile é atualizada ao final da implementação.

## 11. Referências

- `docs/backlogs/mobile/done/011 - cadastro-senha-transacional.md`
- `docs/backlogs/api/006b - step-up-token.md`
- `docs/backlogs/api/006c - internal-transfer-step-up-enforcement.md`
- `docs/backlogs/api/006d - zta-contracts-and-docs.md`
- `api/docs/07-api-rest.md`
- `api/docs/implementations/03-zta-step-up-transaction-password.md`
