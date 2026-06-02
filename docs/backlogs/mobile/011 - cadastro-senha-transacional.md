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

### 3.1 Verificar status da senha transacional

Finalidade:

Permitir que o mobile saiba, após login, se o usuário autenticado já possui
senha transacional ativa.

Este contrato ainda não existe na API atual. Hoje, o módulo de segurança expõe
apenas:

- `POST /security/transaction-password`;
- `POST /security/step-up/authorize`.

Portanto, o gate pós-login depende de uma evolução de contrato na API antes de
ser concluído no mobile.

O mobile não deve descobrir a existência da senha tentando criá-la e tratando
`TRANSACTION_PASSWORD_ALREADY_SET`, porque isso mistura verificação de estado
com mutação.

Formato decidido para a evolução da API:

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
      "can_access_home": true,
      "next_required_step": null
    }
  },
  "error": null
}
```

Detalhes do contrato estão na backlog
`docs/backlogs/api/009 - auth-session-bootstrap.md`.

Erros relevantes para o mobile:

- `TRANSACTION_PASSWORD_NOT_SET`: usuário ainda não cadastrou senha
  transacional;
- `UNAUTHORIZED` / `INVALID_TOKEN`: sessão ausente ou inválida;
- `FORBIDDEN`: usuário autenticado não está apto para a operação.

### 3.2 Criar senha transacional

Finalidade:

Criar a primeira senha transacional do usuário autenticado.

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

## 4. Endpoints empregados na implementação

A implementação mobile deve começar pela criação do pacote de API:

```text
mobile/lib/data/services/apis/senha_transacional/
```

Esse pacote deve concentrar as chamadas de senha transacional. A UI, os view
models e os use cases não devem montar paths, headers ou payloads diretamente.

### 4.1 Mapa de endpoints

| Fluxo | Método | Path | Uso no mobile | API service sugerido |
| --- | --- | --- | --- | --- |
| Cadastro | `POST` | `/security/transaction-password` | Criar a senha transacional com PIN e confirmação. | `SenhaTransacionalApi.create()` |
| Gate pós-login | `GET` | `/auth/session` | Verificar readiness pós-login, incluindo senha transacional ativa antes de liberar a Home. | API de sessão/auth pós-login |

Endpoint existente hoje:

- `POST /security/transaction-password`.

Endpoint necessário, ainda pendente de API:

- `GET /auth/session`, conforme backlog
  `docs/backlogs/api/009 - auth-session-bootstrap.md`.

### 4.2 Endpoint: status da senha transacional

Uso:

- ainda depende de contrato na API;
- chamado após login bem-sucedido, depois que a sessão estiver disponível;
- chamado antes de navegar para a Home;
- não envia body;
- depende apenas do JWT injetado pelo interceptor de autenticação;
- não usa `X-App-Token`;
- não usa `X-Step-Up-Token`.

Resultado esperado no mobile:

- sucesso com `status = active`: liberar Home;
- `TRANSACTION_PASSWORD_NOT_SET`: navegar para cadastro de senha transacional;
- `UNAUTHORIZED` / `INVALID_TOKEN`: tratar como sessão inválida;
- `FORBIDDEN`: bloquear o avanço e exibir erro apropriado.

DTOs sugeridos:

- DTO de sessão pós-login no módulo de auth/session;
- campos mínimos para este fluxo:
  `readiness.transactionPasswordStatus`, `readiness.canAccessHome` e
  `readiness.nextRequiredStep`.

Enquanto esse contrato não existir, a implementação mobile pode preparar o
pacote, DTOs e fronteira de repositório, mas o gate pós-login real fica
bloqueado pela API.

### 4.3 Endpoint: criar senha transacional

Uso:

- chamado somente a partir do fluxo de cadastro;
- recebe PIN e confirmação;
- não recebe senha de login;
- depende apenas do JWT injetado pelo interceptor de autenticação;
- não usa `X-App-Token`;
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

- `CreateSenhaTransacionalRequestDto`;
- `SenhaTransacionalStatusResponseDto`, reutilizando o mesmo formato de
  metadados retornado pelo status.

### 4.4 Fora deste pacote inicial

Não entram em `data/services/apis/senha_transacional/` nesta backlog:

- `POST /security/step-up/authorize`;
- envio de `X-Step-Up-Token`;
- qualquer chamada de transferência.

Esses pontos pertencem à backlog
`012 - step-up-transferencia-interna.md`.

Também não entra como endpoint de segurança nesta backlog:

- `GET /auth/session`, porque ele pertence ao bootstrap de sessão pós-login e
  será definido na backlog de API `009 - auth-session-bootstrap.md`.

## 5. Fluxo de produto

### 5.1 Pós-login

Fluxo obrigatório após login:

1. O usuário faz login com sucesso.
2. O app salva a sessão conforme o fluxo atual.
3. O app verifica se o usuário possui senha transacional ativa.
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
2. O app solicita PIN de 6 dígitos.
3. O app solicita confirmação do PIN.
4. O app valida formato e confirmação localmente.
5. O app chama `POST /security/transaction-password`.
6. Em sucesso, o app considera a senha transacional ativa.
7. O app navega para a Home.

Decisões iniciais:

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

- pacote `data/services/apis/senha_transacional/`;
- `SenhaTransacionalApi` ou equivalente;
- DTO de status da senha transacional, caso o status venha de endpoint de
  segurança dedicado;
- DTO de criação de senha transacional;
- DTO de resposta com metadados da credencial criada.

### 6.2 Repository layer

Adicionar repositório de segurança:

- verificar status da senha transacional;
- criar senha transacional;
- mapear erros relevantes para erros de aplicação;
- expor resultado de sucesso sem retornar dados sensíveis.

### 6.3 Domain/use case layer

Criar use case ou método específico para o cadastro quando a coordenação ficar
maior que uma simples delegação de view model para repositório.

O use case, se criado, deve:

- coordenar o gate pós-login quando o usuário não possuir senha transacional;
- validar a intenção de cadastro em termos de aplicação;
- delegar a chamada ao repositório de segurança;
- retornar `Result`/`AsyncResult`;
- não conhecer widgets, navegação ou controllers.

### 6.4 UI layer

Elementos mínimos:

- tela ou fluxo para cadastrar senha transacional;
- navegação pós-login para o cadastro quando a senha transacional não existir;
- componente de entrada de PIN de 6 dígitos;
- confirmação do PIN;
- validação local antes de chamar a API;
- estados de loading, sucesso e erro;
- feedback claro quando a senha já existe.

A senha transacional deve ficar apenas em estado local/transitório da tela ou do
view model responsável pelo formulário. Ela não deve ser persistida, enviada em
route extras, cache local ou estado global duradouro.

## 7. Erros que o mobile deve tratar

Clientes devem depender de `error.code`.

Criação da senha transacional:

- `TRANSACTION_PASSWORD_ALREADY_SET`;
- `INVALID_DATA`;
- `INVALID_REQUEST`;
- `UNAUTHORIZED`;
- `INVALID_TOKEN`;
- `FORBIDDEN`.

Verificação pós-login:

- `TRANSACTION_PASSWORD_NOT_SET`;
- `UNAUTHORIZED`;
- `INVALID_TOKEN`;
- `FORBIDDEN`.

## 8. UX mínima

- PIN de 6 dígitos.
- Confirmação do PIN.
- Entrada não deve mostrar o valor digitado.
- Validação local antes de chamar a API.
- Após login, usuário sem senha transacional deve ser levado diretamente para o
  cadastro.
- Após cadastro bem-sucedido, usuário deve ir para a Home.
- Feedback claro quando a senha já existe.
- Feedback claro para sessão inválida ou usuário não autorizado.

## 9. Segurança e privacidade

- Não persistir senha transacional.
- Não logar senha transacional.
- Não enviar senha transacional em analytics.
- Não colocar senha transacional em route extras.
- Não armazenar senha transacional em secure storage.
- Usar somente `error.code` para decisão de fluxo.

## 10. Fora do escopo inicial

- autorização step-up para transferência interna;
- envio de `X-Step-Up-Token`;
- acesso à Home sem senha transacional ativa;
- recuperação de senha transacional;
- troca de senha transacional;
- reset administrativo;
- biometria local como substituta da senha transacional;
- dispositivo confiável.

## 11. Critérios de aceite

- Usuário autenticado consegue criar senha transacional pelo app.
- Após login, app verifica se o usuário possui senha transacional ativa.
- Verificação pós-login usa `GET /auth/session` quando esse contrato existir na
  API.
- Usuário sem senha transacional ativa é direcionado para cadastro antes da
  Home.
- Usuário com senha transacional ativa segue para a Home.
- Após cadastro bem-sucedido, usuário segue para a Home.
- App valida PIN de 6 dígitos e confirmação localmente.
- App chama `POST /security/transaction-password` com o payload correto.
- Chamadas de status e criação ficam concentradas em
  `data/services/apis/senha_transacional/`.
- App não envia senha de login nem `X-Step-Up-Token` nessa chamada.
- App não persiste nem loga a senha transacional.
- App não usa tentativa de criação como checagem principal de existência da
  senha transacional.
- App trata `TRANSACTION_PASSWORD_ALREADY_SET` por `error.code`.
- App trata `TRANSACTION_PASSWORD_NOT_SET` no gate pós-login por `error.code`.
- App trata erros de sessão e dados inválidos por `error.code`.
- Testes cobrem DTOs, API service, repositório e view model/use case afetados.
- Documentação mobile é atualizada ao final da implementação.

## 12. Referências

- `docs/backlogs/api/006a - transaction-password.md`
- `docs/backlogs/api/006d - zta-contracts-and-docs.md`
- `api/docs/07-api-rest.md`
- `api/docs/implementations/03-zta-step-up-transaction-password.md`
- `docs/backlogs/api/009 - auth-session-bootstrap.md`
