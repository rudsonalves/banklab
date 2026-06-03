# Backlog: cadastro de senha transacional no mobile

## 1. Contexto

O backlog `api/006` fechou o contrato do MVP ZTA na API para criação da senha
transacional:

- o usuário autenticado cria uma senha transacional em
  `POST /security/transaction-password`;
- a chamada não exige step-up prévio, porque a credencial ainda não existe;
- a API retorna apenas metadados da credencial criada, sem PIN, hash ou
  qualquer material sensível.

Este backlog trata somente do cadastro da senha transacional no mobile.

O uso dessa credencial para autorizar operações sensíveis fica fora deste
escopo e será tratado na backlog `012 - step-up-transferencia-interna.md`.

Decisão de produto para o mobile:

- após login bem-sucedido, o app deve verificar se o usuário possui senha
  transacional ativa;
- se não possuir, o app deve direcionar o usuário para o cadastro da senha
  transacional antes de liberar a Home;
- se possuir, o app segue para a Home normalmente.

## 2. Objetivo

Permitir que o usuário autenticado cadastre uma senha transacional de 6 dígitos
pelo app.

O fluxo mobile deve:

- verificar, após o login, se o usuário já possui senha transacional ativa;
- bloquear o avanço para a Home enquanto a credencial obrigatória não estiver
  criada;
- solicitar PIN transacional de 6 dígitos;
- solicitar confirmação local do PIN;
- validar formato e confirmação antes de chamar a API;
- criar a credencial em `POST /security/transaction-password`;
- tratar erros por `error.code`;
- não armazenar a senha transacional no dispositivo;
- não enviar senha de login para esse endpoint;
- respeitar a arquitetura mobile atual.

## 3. Contratos de API

### 3.1 Sessão autenticada já disponível

Dependência já implementada:

O status da senha transacional já é lido do snapshot de sessão autenticada em
`GET /auth/session`, implementado pela backlog de API
`docs/backlogs/api/009 - auth-session-bootstrap.md` e já consumido pelo fluxo
mobile de login.

O mobile não deve descobrir a existência da senha tentando criá-la e tratando
`TRANSACTION_PASSWORD_ALREADY_SET`, porque isso mistura verificação de estado
com mutação.

Formato disponível:

```http
GET /auth/session
Authorization: Bearer <access_token>
```

Resposta parcial esperada:

```json
{
  "data": {
    "readiness": {
      "transaction_password_status": "active",
      "can_access_home": true
    }
  },
  "error": null
}
```

Detalhes do contrato estão na backlog concluída
`docs/backlogs/api/009 - auth-session-bootstrap.md`.

Nesta backlog, esse contrato deve ser usado apenas como fonte do gate
pós-login. Não há novo endpoint de status a implementar.

### 3.2 Criar senha transacional

Finalidade:

Criar a primeira senha transacional do usuário autenticado.

Autenticação:

- exige usuário logado, com JWT válido no header `Authorization`;
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

## 4. Endpoints empregados na implementação

A implementação mobile deve começar pela criação do pacote de API:

```text
mobile/lib/data/services/apis/transaction_password/
```

Esse pacote deve concentrar as chamadas de senha transacional. A UI, os view
models e os use cases não devem montar paths, headers ou payloads diretamente.

### 4.1 Mapa de endpoints

| Fluxo    | Método | Path                             | Uso no mobile                                     | API service sugerido              |
| -------- | ------ | -------------------------------- | ------------------------------------------------- | --------------------------------- |
| Cadastro | `POST` | `/security/transaction-password` | Criar a senha transacional com PIN e confirmação. | `TransactionPasswordApi.create()` |

Endpoint existente hoje:

- `POST /security/transaction-password`.

### 4.2 Endpoint: criar senha transacional

Uso:

- chamado somente a partir do fluxo de cadastro;
- exige usuário logado;
- recebe PIN e confirmação;
- não recebe senha de login;
- depende apenas do JWT injetado pelo interceptor de autenticação;
- não usa `X-Step-Up-Token`.

Resultado esperado no mobile:

- sucesso com `status = active`: considerar credencial criada e navegar para a
  Home;
- `TRANSACTION_PASSWORD_ALREADY_SET`: considerar que a credencial já existe e
  seguir o comportamento definido para usuário apto a ir para Home;
- `INVALID_DATA`: mostrar erro de PIN inválido ou confirmação divergente;
- `INVALID_REQUEST`: tratar erro de payload/contrato;
- `UNAUTHORIZED` / `INVALID_TOKEN`: tratar como sessão inválida;
- `FORBIDDEN`: bloquear o avanço e exibir erro apropriado.

DTOs sugeridos:

- `CreateTransactionPasswordRequestDto`;
- `TransactionPasswordStatusResponseDto`, reutilizando o mesmo formato de
  metadados retornado pelo status.

### 4.3 Fora deste pacote inicial

Não entram em `data/services/apis/transaction_password/` nesta backlog:

- `POST /security/step-up/authorize`;
- envio de `X-Step-Up-Token`;
- qualquer chamada de transferência.

Esses pontos pertencem à backlog
`012 - step-up-transferencia-interna.md`.

Também não entra como endpoint de segurança nesta backlog:

- `GET /auth/session`, porque ele pertence ao bootstrap de sessão pós-login e
  já foi implementado pela backlog de API `009 - auth-session-bootstrap.md`.

## 5. Fluxo de produto

### 5.1 Pós-login

Fluxo obrigatório após login:

1. O usuário faz login com sucesso.
2. O app salva a sessão conforme o fluxo atual.
3. O app usa o snapshot já carregado de `GET /auth/session` para verificar se
   o usuário possui senha transacional ativa.
4. Se possuir senha transacional ativa, o app navega para a Home.
5. Se não possuir senha transacional ativa, o app navega para o cadastro da
   senha transacional.
6. Após cadastro bem-sucedido, o app navega para a Home.

O app não deve liberar a Home para usuário autenticado sem senha transacional
ativa, salvo em rotas técnicas de logout, erro de sessão ou recuperação futura
explicitamente definida.

### 5.2 Cadastro da senha transacional

1. O usuário autenticado acessa a criação de senha transacional após o gate de
   pós-login.
2. O app solicita PIN de 6 dígitos na primeira página do fluxo.
3. O app solicita confirmação do PIN na segunda página do fluxo.
4. O app valida formato e confirmação localmente.
5. O app chama `POST /security/transaction-password`, que é a operação
   responsável por criar a credencial na API.
6. Em sucesso, o app considera a senha transacional ativa.
7. O app navega para a Home.

Decisões iniciais:

- usar helper compartilhado para decidir a navegação pós-login a partir de
  `AuthRepository.userProfile.readiness`;
- bloquear o acesso à Home apenas no fluxo pós-login desta backlog, sem criar
  guard global de rotas neste momento;
- começar a implementação com view model + repository, sem use case dedicado;
- criar helper para extrair o `error.code` real retornado pela API quando ele
  vier encapsulado em `AppError.details`;
- não criar recuperação, troca ou reset de senha transacional neste backlog;
- não integrar este fluxo à execução de transferência neste backlog;
- não usar tentativa de criação como mecanismo principal de verificação de
  existência da credencial;
- não armazenar a senha transacional no dispositivo;
- não exibir senha transacional em logs, analytics ou mensagens de erro.

## 6. Impactos na arquitetura mobile

### 6.1 Data layer

Adicionar API service para o módulo de segurança, seguindo o padrão de
`mobile/lib/data/services/apis`:

- pacote `data/services/apis/transaction_password/`;
- `TransactionPasswordApi` ou equivalente;
- DTO de criação de senha transacional;
- DTO de resposta com metadados da credencial criada.

### 6.2 Repository layer

O status da senha transacional para o gate pós-login deve continuar vindo da
sessão autenticada (`GET /auth/session`), já carregada pelo fluxo de login.

Adicionar um repositório dedicado de senha transacional:

- criar senha transacional;
- mapear erros relevantes para erros de aplicação;
- expor resultado de sucesso sem retornar dados sensíveis.

### 6.3 Domain/use case layer

A implementação inicial não deve criar use case dedicado para este fluxo.

O cadastro deve ser coordenado por view model + repository:

- view model valida formato, confirmação local e estados de UI;
- repository delega a chamada ao API service e mapeia erros relevantes;
- repository retorna `Result`/`AsyncResult`;
- nenhum singleton/lazy singleton deve manter PIN ou confirmação em memória
  duradoura.

Um use case stateless pode ser introduzido depois se o fluxo ganhar regras de
aplicação além de validação local e chamada ao repositório.

### 6.4 UI layer

Elementos mínimos:

- fluxo em duas páginas para cadastrar senha transacional;
- navegação pós-login para o cadastro quando a senha transacional não existir;
- helper compartilhado para aplicar o gate pós-login em login completo e login
  curto;
- primeira página com `TokenInput.visible = true` para entrada do PIN de 6
  dígitos;
- segunda página com `TokenInput.visible = false` para confirmação do PIN;
- validação local antes de chamar a API;
- estados de loading, sucesso e erro;
- feedback claro quando a senha já existe.

A senha transacional deve ficar apenas em estado local/transitório do fluxo. Na
implementação atual, o PIN é passado da página de criação para a página de
confirmação via `GoRouter.extra` como dado transitório entre rotas, sem storage,
cache local, log, analytics ou singleton/lazy singleton segurando o valor em
memória duradoura.

## 7. Erros que o mobile deve tratar

Clientes devem depender de `error.code`.

Criação da senha transacional:

- `TRANSACTION_PASSWORD_ALREADY_SET`;
- `INVALID_DATA`;
- `INVALID_REQUEST`;
- `UNAUTHORIZED`;
- `INVALID_TOKEN`;
- `FORBIDDEN`.

Gate pós-login via snapshot de sessão já existente:

- `readiness.transactionPasswordStatus = active`: liberar Home quando
  `readiness.canAccessHome = true`;
- `readiness.transactionPasswordStatus = notSet`: navegar para o cadastro da
  senha transacional;
- `readiness.transactionPasswordStatus = locked` ou `unknown`: bloquear avanço
  e exibir feedback apropriado definido durante a implementação.

## 8. UX mínima

- Primeira página para PIN de 6 dígitos, usando `TokenInput.visible = true`.
- Segunda página para confirmação do PIN, usando `TokenInput.visible = false`.
- Validação local antes de chamar a API.
- Mesmo com validação local, o app envia `transaction_password` e
  `transaction_password_confirmation` para a API.
- Após login, usuário sem senha transacional deve ser levado diretamente para o
  cadastro.
- Após cadastro bem-sucedido, usuário deve ir para a Home.
- Feedback claro quando a senha já existe.
- Feedback claro para sessão inválida ou usuário não autorizado.

## 9. Segurança e privacidade

- Não persistir senha transacional.
- Não logar senha transacional.
- Não enviar senha transacional em analytics.
- Não persistir senha transacional em route state recuperável, deep link,
  storage, cache ou estado global duradouro.
- Permitir apenas passagem transitória via `GoRouter.extra` entre as páginas de
  criação e confirmação.
- Não armazenar senha transacional em secure storage.
- Usar `error.code`/`AppErrorCode` tipado para decisão de fluxo.
- Usar helper compartilhado para extrair o código de erro da API e evitar
  strings soltas nas telas.

## 10. Fora do escopo inicial

- autorização step-up para transferência interna;
- envio de `X-Step-Up-Token`;
- acesso à Home sem senha transacional ativa;
- guard global de rotas para proteger a Home fora do pós-login;
- recuperação de senha transacional;
- troca de senha transacional;
- reset administrativo;
- biometria local como substituta da senha transacional;
- dispositivo confiável.

## 11. Critérios de aceite

- Usuário autenticado consegue criar senha transacional pelo app.
- Após login, app usa o snapshot de sessão já carregado para decidir se o
  usuário possui senha transacional ativa.
- Usuário sem senha transacional ativa é direcionado para cadastro antes da
  Home.
- Usuário com senha transacional ativa segue para a Home.
- Após cadastro bem-sucedido, usuário segue para a Home.
- App usa duas páginas: PIN visível na primeira e confirmação mascarada na
  segunda.
- App valida PIN de 6 dígitos e confirmação localmente.
- App chama `POST /security/transaction-password` com o payload correto.
- Chamada de criação fica concentrada em
  `data/services/apis/transaction_password/`.
- Status pós-login usa o snapshot de sessão autenticada já carregado por
  `GET /auth/session`.
- App executa a criação usando a sessão do usuário logado, via JWT.
- App não envia senha de login nem `X-Step-Up-Token` nessa chamada.
- App não persiste nem loga a senha transacional.
- App não usa tentativa de criação como checagem principal de existência da
  senha transacional.
- App trata `TRANSACTION_PASSWORD_ALREADY_SET` por `error.code`.
- App mapeia `TRANSACTION_PASSWORD_ALREADY_SET` para
  `AppErrorCode.transactionPasswordAlreadySet` antes da decisão na UI.
- App trata `notSet` no gate pós-login a partir de
  `readiness.transactionPasswordStatus`.
- App usa helper compartilhado para aplicar o gate pós-login no login completo
  e no login curto.
- App usa helper compartilhado para extrair `error.code` real da API.
- App trata erros de sessão e dados inválidos por `error.code`.
- Testes cobrem DTOs, API service, repositório, helper de erro, helper de gate
  pós-login e view models afetados.
- Documentação mobile é atualizada ao final da implementação.

## 12. Referências

- `docs/backlogs/api/006a - transaction-password.md`
- `docs/backlogs/api/006d - zta-contracts-and-docs.md`
- `api/docs/07-api-rest.md`
- `api/docs/implementations/03-zta-step-up-transaction-password.md`
- `docs/backlogs/api/009 - auth-session-bootstrap.md`
