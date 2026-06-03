# Tasks: Cadastro de senha transacional no mobile

Estas tasks dividem o backlog mobile de cadastro de senha transacional em
passos executáveis.

O objetivo é permitir que um usuário logado crie sua primeira senha
transacional de 6 dígitos, usando o snapshot de sessão já carregado por
`GET /auth/session` como gate pós-login e chamando
`POST /security/transaction-password` apenas no fluxo de cadastro.

A implementação inicial deve usar view model + repository, sem use case
dedicado. Nenhum singleton/lazy singleton deve manter PIN ou confirmação em
memória duradoura.

## Task 1/9: Criar DTOs de senha transacional

### Objetivo

Representar no mobile o contrato de criação da senha transacional.

### Escopo

- Criar pacote `mobile/lib/data/services/apis/transaction_password/`.
- Criar DTO de request `CreateTransactionPasswordRequestDto` ou nome
  equivalente.
- Serializar:
  - `transaction_password`;
  - `transaction_password_confirmation`.
- Criar DTO de resposta `TransactionPasswordStatusResponseDto` ou nome
  equivalente.
- Parsear metadados retornados pela API:
  - `user_id`;
  - `status`;
  - `created_at`.
- Não modelar PIN, hash ou qualquer material sensível na resposta.
- Não persistir PIN nem confirmação.

### Critérios de aceite

- Request DTO gera o JSON esperado pela API.
- Response DTO parseia resposta de sucesso válida.
- Datas são parseadas de forma compatível com os padrões atuais do app.
- DTOs não expõem hash ou material sensível.
- Testes cobrem serialização e parse.

### Depende de

- Nenhuma.

## Task 2/9: Criar TransactionPasswordApi

### Objetivo

Permitir que a camada de dados chame `POST /security/transaction-password`.

### Escopo

- Criar `TransactionPasswordApi`.
- Adicionar método `create()` ou nome equivalente.
- Chamar `POST /security/transaction-password`.
- Enviar somente o body com PIN e confirmação.
- Depender da sessão do usuário logado via JWT injetado pelo interceptor.
- Não enviar senha de login.
- Não enviar `X-Step-Up-Token`.
- Parsear resposta com `ApiEnvelope`.
- Retornar `Result` seguindo o padrão atual dos API services.
- Registrar o service no `Services.add`.

### Critérios de aceite

- API service chama o path correto.
- Payload contém `transaction_password` e
  `transaction_password_confirmation`.
- Request não monta headers de step-up.
- Sucesso retorna DTO com metadados.
- Falhas da API retornam `Failure(AppError)`.
- Testes cobrem sucesso, erro de envelope e erro HTTP/client.

### Depende de

- Task 1.

## Task 3/9: Criar helper para código de erro da API

### Objetivo

Centralizar a extração do `error.code` real retornado pela API quando ele vier
encapsulado em `AppError.details`.

### Escopo

- Criar helper compartilhado em local apropriado do core/data.
- Extrair código backend quando `AppError.details` contiver:
  - `code`;
  - ou mapa de erro com `code`.
- Preservar comportamento quando o erro não possuir código backend.
- Evitar strings soltas nas telas para:
  - `TRANSACTION_PASSWORD_ALREADY_SET`;
  - `INVALID_DATA`;
  - `INVALID_REQUEST`;
  - `UNAUTHORIZED`;
  - `INVALID_TOKEN`;
  - `FORBIDDEN`.
- Não alterar o contrato público de `AppError` sem necessidade.

### Critérios de aceite

- Helper retorna o código backend quando disponível.
- Helper retorna `null` ou equivalente quando não houver código backend.
- Helper funciona com os formatos de `details` já usados pelo app.
- Testes cobrem erro com `details.code`, erro com mapa de erro e erro sem
  código backend.

### Depende de

- Nenhuma.

## Task 4/9: Criar TransactionPasswordRepository

### Objetivo

Encapsular a criação da senha transacional e o mapeamento de erros relevantes.

### Escopo

- Criar contrato `TransactionPasswordRepository`.
- Criar implementação `TransactionPasswordRepositoryImpl`.
- Injetar `TransactionPasswordApi`.
- Expor método para criar senha transacional.
- Delegar chamada ao API service.
- Retornar `Result`/`AsyncResult`.
- Tratar `TRANSACTION_PASSWORD_ALREADY_SET` como condição conhecida para a UI.
- Propagar erros de sessão, dados inválidos e autorização sem perder o código
  backend.
- Registrar repository em `Repositories.add`.

### Critérios de aceite

- Repository delega criação para o API service.
- Sucesso retorna metadados da credencial criada.
- `TRANSACTION_PASSWORD_ALREADY_SET` permanece identificável por código.
- Erros relevantes preservam informação suficiente para a UI decidir fluxo.
- Testes cobrem sucesso, senha já existente, dados inválidos, sessão inválida e
  erro inesperado.

### Depende de

- Task 2.
- Task 3.

## Task 5/9: Criar helper de gate pós-login

### Objetivo

Compartilhar a decisão de navegação pós-login entre login completo e login
curto.

### Escopo

- Criar helper para avaliar `AuthRepository.userProfile.readiness`.
- Decidir:
  - `active` + `canAccessHome = true`: navegar para Home;
  - `notSet`: navegar para cadastro da senha transacional;
  - `locked` ou `unknown`: bloquear avanço e retornar estado de bloqueio.
- Não chamar `GET /auth/session` novamente.
- Não criar guard global de rotas.
- Permitir mensagens/feedback para estados bloqueados.
- Cobrir ausência inesperada de `userProfile`.

### Critérios de aceite

- Helper libera Home para sessão pronta.
- Helper direciona para cadastro quando status é `notSet`.
- Helper bloqueia `locked` e `unknown`.
- Helper trata sessão ausente como erro de sessão.
- Testes cobrem todos os status de `TransactionPasswordStatus`.

### Depende de

- Nenhuma.

## Task 6/9: Integrar gate pós-login no LoginPage e ShortLoginPage

### Objetivo

Usar o helper compartilhado para impedir avanço para Home quando a senha
transacional ainda não existir.

### Escopo

- Atualizar `LoginPage`.
- Atualizar `ShortLoginPage`.
- Após sucesso do comando de login, usar o helper de gate pós-login.
- Navegar para Home quando a sessão estiver pronta.
- Navegar para a primeira página de cadastro da senha transacional quando
  `transactionPasswordStatus = notSet`.
- Exibir feedback para `locked`, `unknown` ou sessão ausente.
- Preservar mensagens existentes de erro de login, aprovação pendente e contato
  não verificado.
- Não criar guard global para Home nesta task.

### Critérios de aceite

- Login completo com senha transacional ativa vai para Home.
- Login curto com senha transacional ativa vai para Home.
- Login completo sem senha transacional vai para cadastro.
- Login curto sem senha transacional vai para cadastro.
- Estados `locked` e `unknown` não navegam para Home.
- Testes/widget tests ou testes de helper cobrem os cenários principais.

### Depende de

- Task 5.

## Task 7/9: Criar rotas e páginas do cadastro da senha transacional

### Objetivo

Construir o fluxo em duas páginas para entrada e confirmação do PIN.

### Escopo

- Criar rotas:
  - `TransactionPasswordRoutes.create`;
  - `TransactionPasswordRoutes.confirm`.
- Criar primeira página para PIN de 6 dígitos.
- Usar `TokenInput.visible = true` na primeira página.
- Criar segunda página para confirmação do PIN.
- Usar `TokenInput.visible = false` na segunda página.
- Manter o PIN apenas em estado local/transitório do fluxo ou view model.
- Não persistir PIN.
- Não logar PIN.
- Não colocar PIN em storage ou cache.
- Validar PIN numérico de 6 dígitos antes de avançar.
- Validar confirmação antes de chamar a API.
- Enviar `transaction_password` e `transaction_password_confirmation` para o
  repository mesmo com validação local.
- Navegar para Home após sucesso.
- Tratar `TRANSACTION_PASSWORD_ALREADY_SET` como credencial já existente,
  seguindo o comportamento de usuário apto a ir para Home.

### Critérios de aceite

- Primeira página coleta PIN visível.
- Segunda página coleta confirmação mascarada.
- PIN incompleto não avança.
- Confirmação divergente não chama a API.
- Sucesso chama a API com payload correto e navega para Home.
- `TRANSACTION_PASSWORD_ALREADY_SET` gera feedback claro e não deixa o usuário
  preso no cadastro.
- PIN e confirmação não são persistidos.
- View model limpa estado sensível ao concluir, cancelar ou sair do fluxo.

### Depende de

- Task 4.

## Task 8/9: Registrar dependências e exports

### Objetivo

Integrar os novos serviços, repositories, view models e rotas ao app.

### Escopo

- Registrar `TransactionPasswordApi` em `Services.add`.
- Registrar `TransactionPasswordRepository` em `Repositories.add`.
- Registrar view model do fluxo, se houver.
- Adicionar exports necessários em barrels existentes, quando aplicável.
- Adicionar rotas novas ao router.
- Garantir que login completo e curto conseguem navegar para as novas rotas.
- Não registrar objeto stateful sensível como singleton/lazy singleton.

### Critérios de aceite

- App compila com as novas dependências.
- Rotas novas são resolvidas por nome.
- Fluxo pós-login consegue navegar para cadastro.
- Nenhum singleton/lazy singleton mantém PIN ou confirmação.

### Depende de

- Task 4.
- Task 6.
- Task 7.

## Task 9/9: Validar, testar e atualizar documentação

### Objetivo

Fechar a implementação com cobertura e documentação atualizada.

### Escopo

- Adicionar/atualizar testes de DTOs.
- Adicionar/atualizar testes de API service.
- Adicionar/atualizar testes de repository.
- Adicionar/atualizar testes do helper de erro backend.
- Adicionar/atualizar testes do helper de gate pós-login.
- Adicionar/atualizar testes dos view models afetados.
- Rodar `dart format`.
- Rodar `flutter test`.
- Rodar `flutter analyze`.
- Atualizar documentação mobile ao final, quando necessário.
- Atualizar changelog, se esse for o padrão da entrega.

### Critérios de aceite

- Testes relevantes passam.
- `flutter analyze` não aponta novos problemas.
- Código formatado.
- Backlog e documentação refletem o comportamento implementado.
- Nenhum dado sensível é persistido, logado ou mantido em singleton duradouro.

### Depende de

- Task 8.
