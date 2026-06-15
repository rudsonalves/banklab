# Tasks: Step-up na transferência interna mobile

Estas tasks dividem o backlog mobile de step-up na transferência interna em
passos executáveis.

O fluxo deve verificar a senha transacional ao entrar na transferência,
reutilizar o cadastro existente quando necessário e exigir autorização step-up
antes de executar `POST /accounts/internal-transfers`.

Senha transacional, PIN e `step_up_token` devem permanecer transitórios. A
`idempotency_key` deve permanecer estável durante retries da mesma tentativa
lógica de transferência.

## Task 1/13: Centralizar e sincronizar a sessão da aplicação

### Objetivo

Manter o snapshot autenticado em um serviço transversal de sessão e sincronizar
localmente o estado da senha transacional após sua criação.

### Escopo

- Usar `AppSection` como fonte única do snapshot `AuthSession` durante a sessão
  autenticada.
- Registrar `AppSection` como singleton da aplicação.
- Fazer o `AuthRepository` preencher o `AppSection` após carregar
  `GET /auth/session` no login e limpá-lo no logout.
- Manter `getAuthSession()` com retorno do snapshot em memória quando já
  carregado.
- Após criação bem-sucedida da senha transacional com resposta `active`, fazer
  `TransactionPasswordRepository` atualizar localmente o
  `transactionPasswordStatus` do `AppSection`.
- Não chamar novamente `GET /auth/session` após a criação bem-sucedida.
- Tratar `TRANSACTION_PASSWORD_ALREADY_SET` como inconsistência defensiva:
  permanecer no cadastro, sem navegar e sem atualizar o snapshot
  otimisticamente.
- Ao receber `TRANSACTION_PASSWORD_NOT_SET` durante o step-up, atualizar
  localmente o `transactionPasswordStatus` do `AppSection` para `not_set`.
- Usar o erro retornado pelo backend como informação autoritativa, sem consultar
  novamente `GET /auth/session`.
- Não consultar novamente `GET /auth/session` na entrada normal da
  transferência.

### Critérios de aceite

- Login preenche o `AppSection` com o snapshot retornado por
  `GET /auth/session`.
- Logout limpa o `AppSection`.
- `getAuthSession()` pode retornar o snapshot já carregado.
- Criação com resposta `active` atualiza o `AppSection` sem uma segunda chamada
  de sessão.
- `TRANSACTION_PASSWORD_ALREADY_SET` não libera Home ou transferência e não
  altera o snapshot.
- `TRANSACTION_PASSWORD_NOT_SET` atualiza localmente o snapshot para `not_set`
  sem nova consulta de sessão.
- A entrada normal da transferência não provoca nova consulta de sessão.
- Testes cobrem ciclo do `AppSection`, atualização local após criação,
  atualização local após `TRANSACTION_PASSWORD_NOT_SET` e ausência de refresh
  na entrada da transferência.

### Depende de

- Nenhuma.

## Task 2/13: Tornar o cadastro de senha transacional reutilizável

### Objetivo

Permitir que o fluxo existente de cadastro conclua na Home ou na transferência,
conforme sua origem.

### Escopo

- Criar enum ou objeto fechado para o destino pós-cadastro, com pelo menos:
  - `postLogin`;
  - `internalTransfer`.
- Não aceitar path ou nome de rota arbitrário como destino.
- Propagar o destino entre as páginas de criação e confirmação.
- Manter o PIN transitório conforme o fluxo atual.
- Após sucesso, usar o estado atualizado pelo
  `TransactionPasswordRepository` antes de navegar.
- Tratar `TRANSACTION_PASSWORD_ALREADY_SET` como falha defensiva, mantendo o
  usuário no cadastro.
- Para `postLogin`, navegar para a Home após confirmar que o `AppSection`
  permite esse acesso.
- Para `internalTransfer`, navegar para `TransferRoutes.recipient` somente se o
  status no `AppSection` for `active`.
- Em falha de criação, permanecer no cadastro.
- Em cancelamento iniciado pela Home, retornar à Home.

### Critérios de aceite

- O fluxo pós-login continua concluindo na Home.
- O fluxo iniciado por **Transferir** conclui na seleção de destinatário.
- Nenhum destino arbitrário pode ser injetado pelas rotas.
- Sucesso com resposta `active` atualiza localmente o `AppSection`.
- `TRANSACTION_PASSWORD_ALREADY_SET` não navega nem altera o `AppSection`.
- Status diferente de `active` não abre a transferência.
- Falha mantém o usuário no cadastro e permite nova tentativa sem reutilizar o
  PIN anterior.
- Testes cobrem os dois destinos, cancelamento, falha e senha já existente.

### Depende de

- Task 1.

## Task 3/13: Verificar senha transacional na entrada da transferência

### Objetivo

Decidir o destino correto quando o usuário tocar em **Transferir**.

### Escopo

- Expor no `HomeViewmodel` o `TransactionPasswordStatus` armazenado no
  `AppSection`, carregado durante o login.
- Atualizar o botão **Transferir** para tratar o status na UI:
  - abrir transferência quando `active`;
  - abrir cadastro quando `not_set`;
  - permanecer na Home e informar o bloqueio quando `locked`;
  - permanecer na Home e apresentar erro quando `unknown`.
- Não fazer a `HomePage` interpretar diretamente o snapshot da sessão.

### Critérios de aceite

- `active` abre `TransferRoutes.recipient`.
- `not_set` abre o cadastro com destino `internalTransfer`.
- `locked` e `unknown` não abrem cadastro nem transferência.
- Testes cobrem todos os estados.

### Depende de

- Task 2.

## Task 4/13: Restringir refresh automático aos erros de sessão

### Objetivo

Impedir que erros funcionais de step-up com HTTP 401 acionem renovação do access
token.

### Escopo

- Fazer o `AuthInterceptor` ler `error.code` da resposta 401.
- Tentar refresh somente para:
  - `UNAUTHORIZED`;
  - `INVALID_TOKEN`.
- Não tentar refresh para:
  - `TRANSACTION_PASSWORD_INVALID`;
  - `STEP_UP_TOKEN_REQUIRED`;
  - `STEP_UP_TOKEN_INVALID`;
  - `STEP_UP_TOKEN_EXPIRED`;
  - `STEP_UP_TOKEN_CONSUMED`.
- Encaminhar respostas 401 funcionais ao chamador sem alterar o erro.
- Não usar refresh automático para códigos 401 desconhecidos.
- Preservar serialização de refresh concorrente já existente.

### Critérios de aceite

- `UNAUTHORIZED` e `INVALID_TOKEN` continuam tentando refresh.
- Erros funcionais de step-up não chamam `/auth/refresh`.
- Código e body originais chegam ao cliente chamador.
- Refresh concorrente continua realizando uma única chamada.
- Testes cobrem todos os códigos listados e código desconhecido.

### Depende de

- Nenhuma.

## Task 5/13: Preservar códigos backend e sanitizar logs HTTP

### Objetivo

Permitir decisões por `error.code` sem expor credenciais nos logs.

### Escopo

- Garantir que erros retornados pelo `RestClient` preservem o código backend em
  `AppError.details`.
- Manter compatibilidade com `backendErrorCode(error)`.
- Não exigir um `AppErrorCode` para cada erro ZTA.
- Remover ou sanitizar logging genérico de headers em requests.
- Mascarar ou omitir:
  - `Authorization`;
  - `X-Step-Up-Token`;
  - outros headers reconhecidos como credenciais.
- Não registrar body de `POST /security/step-up/authorize`.
- Não registrar senha transacional em API services, repositories ou view
  models.
- Manter logs técnicos úteis sem dados sensíveis.

### Critérios de aceite

- `backendErrorCode(error)` identifica códigos ZTA propagados pelo cliente.
- Logs não contêm access token, refresh token, step-up token ou PIN.
- Logs de transferência não exibem `X-Step-Up-Token`.
- Logs de autorização step-up não exibem o request body.
- Testes cobrem extração de código e sanitização.

### Depende de

- Nenhuma.

## Task 6/13: Criar DTOs e chamada de autorização step-up

### Objetivo

Representar e executar `POST /security/step-up/authorize`.

### Escopo

- Adicionar DTO de request de autorização step-up.
- Serializar:
  - `method`;
  - `path`;
  - `transaction_password`.
- Usar representação canônica:
  - `method=POST`;
  - `path=/accounts/internal-transfers`.
- Adicionar DTO de resposta com:
  - `step_up_token`;
  - `expires_in`.
- Ampliar `TransactionPasswordApi` com `authorizeStepUp()` ou nome equivalente.
- Parsear resposta com `ApiEnvelope`.
- Preservar `error.code` nos erros.
- Não interpretar o JWT retornado.
- Não persistir ou logar PIN e token.

### Critérios de aceite

- Request gera o JSON esperado pela API.
- API service chama `/security/step-up/authorize`.
- Sucesso retorna token e expiração.
- Erros de senha, sessão e policy permanecem identificáveis por código.
- Nenhum dado sensível aparece em logs.
- Testes cobrem serialização, sucesso, envelope inválido e falhas relevantes.

### Depende de

- Task 4.
- Task 5.

## Task 7/13: Ampliar o repositório de senha transacional para step-up

### Objetivo

Expor autorização de transferência interna sem apresentar detalhes de API à
camada de domínio.

### Escopo

- Ampliar `TransactionPasswordRepository`.
- Adicionar operação específica para autorizar transferência interna.
- Encapsular `POST` e `/accounts/internal-transfers` no repository ou em input
  de domínio fechado.
- Delegar a autorização ao `TransactionPasswordApi`.
- Retornar token e expiração somente em memória.
- Preservar códigos backend para decisão do fluxo.
- Não armazenar o último token no repository.

### Critérios de aceite

- Chamadores não precisam conhecer chaves internas como
  `internal_transfer.create`.
- Repository usa método e path públicos corretos.
- Sucesso retorna token transitório.
- Falha preserva o código backend.
- Repository não mantém PIN nem token após a chamada.
- Testes cobrem sucesso e erros ZTA relevantes.

### Depende de

- Task 6.

## Task 8/13: Enviar o step-up token na transferência

### Objetivo

Cumprir o contrato obrigatório de `POST /accounts/internal-transfers`.

### Escopo

- Alterar a operação de transferência para receber o step-up token
  explicitamente.
- Enviar `X-Step-Up-Token` somente nessa chamada protegida.
- Não adicionar o token como header global do cliente.
- Não incluir o token no DTO ou body da transferência.
- Não incluir a senha transacional no request de transferência.
- Preservar `from_account_id`, `to_account_id`, `amount`,
  `idempotency_key` e `description`.
- Preservar códigos de erro do endpoint protegido.

### Critérios de aceite

- Transferência envia `X-Step-Up-Token`.
- Body não contém PIN nem token.
- Outras chamadas não recebem o header.
- Token ausente não é substituído silenciosamente.
- Erros `STEP_UP_TOKEN_*` permanecem identificáveis.
- Testes verificam path, body e headers exatos.

### Depende de

- Task 5.

## Task 9/13: Orquestrar a transferência protegida

### Objetivo

Coordenar autorização e execução fora da UI.

### Escopo

- Ajustar `TransferUsecase` ou criar operação dedicada para transferência
  protegida.
- Receber o draft confirmado e a senha transacional.
- Gerar a `idempotency_key` uma vez por tentativa lógica, antes do step-up.
- Chamar autorização step-up.
- Chamar a transferência somente após autorização bem-sucedida.
- Passar o token por parâmetro explícito ao repository de transação.
- Descartar a referência ao token após uma única tentativa de transferência,
  inclusive em erro.
- Preservar a mesma `idempotency_key` ao repetir apenas o step-up.
- Nunca repetir automaticamente usando o PIN anterior.
- Nunca reutilizar o token anterior.

### Critérios de aceite

- Falha de autorização não chama o endpoint de transferência.
- Sucesso de autorização chama a transferência com o token recebido.
- Token é usado em uma única tentativa.
- Retry de step-up mantém a mesma `idempotency_key`.
- Nova transferência gera nova `idempotency_key`.
- PIN não fica armazenado em repository, singleton ou estado duradouro.
- Testes cobrem ordem das chamadas, sucesso e falhas em cada etapa.

### Depende de

- Task 7.
- Task 8.

## Task 10/13: Criar a experiência de entrada do PIN na confirmação

### Objetivo

Solicitar a senha transacional somente após o usuário confirmar os dados da
transferência.

### Escopo

- Adicionar prompt, modal ou página transitória de step-up.
- Reutilizar `TokenInput` com entrada mascarada.
- Solicitar PIN numérico de 6 dígitos.
- Não colocar PIN em route extra, cache, storage ou estado global.
- Permitir cancelamento com retorno à confirmação.
- Bloquear submissões concorrentes.
- Exibir loading durante autorização e transferência.
- Entregar o PIN à operação protegida e limpar o estado local imediatamente
  depois.
- Navegar ao sucesso somente após a transferência.

### Critérios de aceite

- Confirmar transferência abre a entrada de PIN antes da chamada protegida.
- PIN incompleto não submete.
- Cancelamento não chama autorização nem transferência.
- Loading impede toques duplicados.
- PIN é limpo após sucesso, falha ou cancelamento.
- Sucesso mantém o fluxo atual de status e comprovante.
- Testes cobrem estados principais do prompt e da confirmação.

### Depende de

- Task 9.

## Task 11/13: Implementar decisões de erro e retry

### Objetivo

Tratar erros ZTA por código e preservar a tentativa correta.

### Escopo

- Usar `backendErrorCode(error)` para decisões de fluxo.
- Tratar autorização:
  - `TRANSACTION_PASSWORD_INVALID`: informar erro e permitir novo PIN;
  - `TRANSACTION_PASSWORD_LOCKED`: informar bloqueio e impedir nova tentativa
    imediata;
  - `STEP_UP_ENDPOINT_NOT_ALLOWED`: encerrar com erro não recuperável;
  - `UNAUTHORIZED` e `INVALID_TOKEN`: seguir tratamento de sessão;
  - `FORBIDDEN`, `INVALID_DATA` e `INVALID_REQUEST`: apresentar erro adequado.
- Tratar transferência:
  - `STEP_UP_TOKEN_EXPIRED`: solicitar novo PIN;
  - `STEP_UP_TOKEN_CONSUMED`: solicitar novo PIN;
  - `STEP_UP_TOKEN_REQUIRED`, `STEP_UP_TOKEN_INVALID` e
    `STEP_UP_ENDPOINT_MISMATCH`: não reutilizar token e apresentar falha
    apropriada.
- Para token expirado ou consumido, manter o draft e a `idempotency_key`.
- Não navegar para a tela genérica de falha antes de oferecer novo step-up para
  token expirado ou consumido.
- Para `TRANSACTION_PASSWORD_NOT_SET`, descartar a tentativa, manter a
  atualização local do `AppSection`, retornar à Home e exigir nova verificação
  no botão **Transferir**.
- Preservar tratamento dos erros de negócio já existentes.

### Critérios de aceite

- Nenhuma decisão depende de mensagem textual.
- Senha inválida permite nova entrada sem reutilizar o PIN.
- Senha bloqueada não executa transferência.
- Token expirado ou consumido reabre o step-up com a mesma idempotência.
- Token anterior nunca é reutilizado.
- `TRANSACTION_PASSWORD_NOT_SET` não abre cadastro diretamente da confirmação.
- `TRANSACTION_PASSWORD_NOT_SET` mantém o `AppSection` local como `not_set`
  antes do retorno à Home.
- Erros de negócio continuam chegando ao fluxo de status existente.
- Testes cobrem todos os códigos ZTA listados.

### Depende de

- Task 9.
- Task 10.

## Task 12/13: Registrar dependências, rotas e integração do fluxo

### Objetivo

Conectar os novos contratos à composição atual do app.

### Escopo

- Registrar APIs, repositories, use cases e view models adicionados ou
  alterados.
- Preservar os ciclos de vida atuais quando forem compatíveis.
- Não registrar objetos que mantenham PIN ou token como singleton/lazy
  singleton.
- Atualizar rotas do cadastro para transportar destino tipado.
- Atualizar `ExtraCodec` caso o tipo de destino precise ser serializável.
- Integrar a entrada da transferência na Home.
- Integrar o step-up à confirmação da transferência.
- Garantir fallback seguro para route extras ausentes ou inválidos.

### Critérios de aceite

- App compila com todas as dependências resolvidas.
- Home consegue verificar sessão e abrir cadastro ou transferência.
- Cadastro conclui no destino correto.
- Confirmação executa autorização e transferência em sequência.
- Extras inválidos não expõem PIN e não abrem a transferência indevidamente.
- Nenhum singleton mantém credencial transitória.

### Depende de

- Task 2.
- Task 3.
- Task 4.
- Task 5.
- Task 7.
- Task 9.
- Task 10.
- Task 11.

## Task 13/13: Validar, testar e atualizar documentação

### Objetivo

Fechar a implementação com cobertura integrada, análise estática e documentação
coerente.

### Escopo

- Completar testes de DTOs e API services.
- Completar testes de repositories e use cases.
- Completar testes de view models.
- Cobrir atualização local do `AppSection` após criação.
- Cobrir atualização local do `AppSection` após
  `TRANSACTION_PASSWORD_NOT_SET` defensivo, sem nova consulta de sessão.
- Cobrir uso do snapshot já carregado na entrada da transferência, sem nova
  chamada de sessão.
- Cobrir destinos `postLogin` e `internalTransfer`.
- Cobrir `AuthInterceptor` para erros de sessão e step-up.
- Cobrir sanitização de logs sensíveis.
- Cobrir estabilidade da `idempotency_key`.
- Cobrir descarte e não reutilização do token.
- Adicionar widget tests para entrada da transferência, cadastro reutilizado,
  prompt de PIN e retries críticos quando compatível com o padrão do projeto.
- Rodar `dart format`.
- Rodar `flutter test`.
- Rodar `flutter analyze`.
- Atualizar documentação mobile e changelog, se aplicável.

### Critérios de aceite

- Testes relevantes passam.
- `flutter analyze` não aponta novos problemas.
- Código está formatado.
- Fluxo completo funciona para usuário com senha ativa.
- Usuário sem senha é encaminhado ao cadastro e depois à transferência.
- Falha ou cancelamento do cadastro não abre a transferência.
- Transferência nunca é chamada sem step-up token.
- PIN, tokens e headers de autenticação não são persistidos nem logados.
- Backlog e documentação refletem o comportamento implementado.

### Depende de

- Task 12.
