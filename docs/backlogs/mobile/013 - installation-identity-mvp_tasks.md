# Tarefas do MVP de Identidade de Instalação

Estas tarefas dividem a implementação mobile da identidade de instalação em
etapas de tamanho adequado, alinhadas ao contrato da API.

A primeira versão cobre identidade local da instalação, propagação HTTP, login,
registro restrito de instalação e estados de erro aprovados. O gerenciamento de
instalações, com listagem e revogação, fica intencionalmente adiado.

## Tarefa 1/12: Adicionar chaves de armazenamento da identidade de instalação e contrato do marcador local

### Objetivo

Definir as chaves locais e o contrato do marcador local usados para distinguir
atualizações do app de novas instalações.

### Escopo

- Adicionar `StorageKeys.installationId` com o valor
  `banklab.installation.id`.
- Definir o mecanismo de marcador local usado fora do armazenamento seguro
  durável.
- Garantir que o marcador não seja tratado como segredo ou identidade.
- Garantir que uma restauração sem marcador válido resulte em uma nova
  identidade de instalação.
- Documentar que logout e limpeza de credenciais não devem apagar
  `StorageKeys.installationId`.

### Critérios de Aceite

- A chave do id de instalação está disponível por meio de `StorageKeys`.
- O comportamento do marcador local está documentado em código ou testes.
- A implementação consegue distinguir uma atualização do app de um cenário de
  reinstalação/limpeza de dados.
- A limpeza de sessão não mira a chave do id de instalação.

### Depende De

- Nenhuma.

## Tarefa 2/12: Implementar serviço de identidade de instalação

### Objetivo

Criar um serviço que resolva um UUID v4 estável para a instalação atual do app.

### Escopo

- Criar `InstallationIdentityService`.
- Ler o id de instalação existente em `LocalSecureStorage`.
- Validar que o valor armazenado é um UUID v4 canônico.
- Gerar um novo UUID v4 quando o marcador local estiver ausente, o valor estiver
  ausente ou o valor for inválido.
- Persistir o UUID v4 gerado em `flutter_secure_storage`.
- Criar ou atualizar o marcador local de instalação após uma resolução
  bem-sucedida.
- Retornar um `AsyncResult<String>`.
- Não registrar em log o id completo da instalação.

### Critérios de Aceite

- A primeira execução cria e persiste um UUID v4.
- Uma atualização do app com marcador local presente reutiliza o mesmo UUID.
- Reinstalação/limpeza de dados/restauração sem marcador local gera um novo UUID,
  mesmo que o armazenamento seguro ainda retorne um valor antigo.
- Valores armazenados inválidos são substituídos.
- Falhas de leitura/escrita retornam um resultado de falha e não continuam de
  forma silenciosa.
- Logs não incluem o UUID completo.

### Depende De

- Tarefa 1.

## Tarefa 3/12: Bloquear bootstrap quando a identidade de instalação falhar

### Objetivo

Impedir login e chamadas à API quando o app não conseguir resolver um
`installation_id` estável.

### Escopo

- Resolver a identidade de instalação durante o bootstrap do app ou antes do
  primeiro fluxo de API.
- Exibir um erro recuperável quando a resolução da identidade falhar.
- Fornecer um caminho de retentativa.
- Não enviar requisições à API sem `X-Installation-Id`.

### Critérios de Aceite

- O login fica indisponível enquanto a identidade de instalação não estiver
  resolvida.
- Uma falha de armazenamento/leitura/escrita/geração bloqueia o fluxo.
- A retentativa pode resolver a identidade e desbloquear o app.
- Nenhuma requisição é enviada intencionalmente sem `X-Installation-Id`.

### Depende De

- Tarefa 2.

## Tarefa 4/12: Adicionar interceptor de instalação

### Objetivo

Anexar `X-Installation-Id` a todas as requisições tratadas pelo cliente HTTP
principal.

### Escopo

- Criar `InstallationInterceptor`.
- Resolver o id da instalação atual por meio de `InstallationIdentityService`.
- Adicionar `X-Installation-Id` aos headers da requisição.
- Registrar o interceptor na configuração principal de Dio/RestClient.
- Manter o comportamento existente do interceptor de autenticação.
- Remover ou deixar sem uso a referência comentada antiga ao interceptor de
  dispositivo, sem reintroduzir terminologia de dispositivo.

### Critérios de Aceite

- Requisições de login incluem `X-Installation-Id`.
- Requisições autenticadas incluem `X-Installation-Id`.
- O valor do header é um UUID v4 canônico em letras minúsculas.
- O comportamento existente de `Authorization` permanece inalterado.
- Falhas do interceptor bloqueiam a requisição pelo tratamento normal de erro do
  app, em vez de enviar uma requisição sem header.

### Depende De

- Tarefa 2.
- Tarefa 3.

## Tarefa 5/12: Adicionar header de instalação ao refresh de token

### Objetivo

Garantir que `/auth/refresh` envie a mesma identidade de instalação das
requisições normais.

### Escopo

- Atualizar o fluxo de refresh em `AuthInterceptor`.
- Injetar ou acessar de outra forma `InstallationIdentityService` para o
  refresh.
- Adicionar `X-Installation-Id` à requisição Dio dedicada de refresh.
- Preservar o comportamento de deduplicação de refresh e retentativa.
- Preservar a limpeza de tokens em falha de refresh.

### Critérios de Aceite

- Requisições de refresh incluem `X-Installation-Id`.
- O refresh continua gravando access token e refresh token rotacionados.
- A falha de refresh continua limpando os tokens de sessão.
- Nenhuma requisição de refresh é enviada sem id de instalação quando a resolução
  da identidade falhar.

### Depende De

- Tarefa 2.
- Tarefa 4.

## Tarefa 6/12: Modelar respostas operacionais e restritas de login

### Objetivo

Representar respostas de login que podem criar uma sessão operacional ou exigir
registro de instalação.

### Escopo

- Substituir o parsing direto de login apenas como `OperationalAuthState` por
  uma hierarquia de `AuthState` que represente sessão operacional, estado
  anônimo e autorização restrita de instalação.
- Dar suporte à resposta operacional com `access_token` e `refresh_token`.
- Dar suporte à resposta restrita com `restricted_access_token`,
  `restricted_token_type`, `restricted_scope` e `restricted_expires_at` em
  `RestrictedInstallationAuthState`.
- Dar suporte a `INSTALLATION_LIMIT_REACHED` como erro tipado do app.
- Preservar o tratamento atual de aprovação da conta e verificação de contato.
- Não persistir tokens para resultados restritos ou de limite atingido.

### Critérios de Aceite

- Login operacional continua persistindo tokens e carregando o perfil.
- Login restrito não persiste tokens operacionais.
- Limite atingido não persiste tokens.
- Testes existentes de erro de login continuam passando após as atualizações.
- Falhas de parsing retornam `AppErrorCode.parsingError`.

### Depende De

- Tarefa 4.

## Tarefa 7/12: Adicionar API de registro de instalação

### Objetivo

Chamar `POST /security/installations` para trocar uma autorização restrita por
uma sessão operacional.

### Escopo

- Adicionar DTOs de API para requisição/resposta de registro de instalação,
  conforme necessário.
- Enviar `Authorization: Bearer <restricted_access_token>`.
- Enviar `X-Step-Up-Token`.
- Enviar o mesmo `X-Installation-Id` usado pelo app.
- Fazer parsing de `access_token` e `refresh_token` operacionais.
- Fazer parsing dos metadados de instalação retornados pela API.

### Critérios de Aceite

- Registro bem-sucedido retorna tokens operacionais.
- Step-up token ausente ou inválido é mapeado para erro do app.
- Incompatibilidade de instalação e id de instalação inválido são mapeados para
  erros do app.
- O registro não envia a senha transacional.

### Depende De

- Tarefa 5.
- Tarefa 6.

## Tarefa 8/12: Estender operação de step-up para registro de instalação

### Objetivo

Permitir que o fluxo existente de step-up com senha transacional autorize o
registro de instalação.

### Escopo

- Adicionar `StepUpOperation.installationRegistration`.
- Usar método `POST`.
- Usar path `/security/installations`.
- Reutilizar a API existente `/security/step-up/authorize`.
- Garantir que a senha transacional seja enviada apenas ao endpoint de step-up.

### Critérios de Aceite

- O app consegue solicitar um step-up token para `POST /security/installations`.
- O corpo da requisição de step-up corresponde ao contrato público de operação da
  API.
- A senha transacional não é enviada ao login nem ao registro de instalação.
- O comportamento existente de step-up para transferência interna permanece
  inalterado.

### Depende De

- Nenhuma.

## Tarefa 9/12: Implementar fluxo restrito de certificação de instalação

### Objetivo

Concluir o login de um usuário conhecido em uma nova instalação solicitando a
senha transacional e registrando a instalação.

### Escopo

- Quando o login retornar registro restrito de instalação, navegar para um
  estado/tela de certificação.
- Solicitar a senha transacional.
- Autorizar step-up para `POST /security/installations`.
- Chamar o registro de instalação com o token restrito e o step-up token.
- Persistir os tokens operacionais retornados pelo registro.
- Carregar a sessão autenticada e continuar o fluxo normal pós-login.
- Limpar token restrito, step-up token e senha em sucesso, falha ou cancelamento.

### Critérios de Aceite

- Login restrito leva ao fluxo de certificação.
- Certificação bem-sucedida termina em uma sessão operacional autenticada.
- Nenhum token operacional é persistido antes do sucesso do registro.
- Cancelamento retorna ao login e exige um novo login.
- Token restrito expirado retorna ao login e exige um novo login.

### Depende De

- Tarefa 6.
- Tarefa 7.
- Tarefa 8.

## Tarefa 10/12: Tratar limite de instalações e bloqueios de senha transacional

### Objetivo

Exibir estados de bloqueio aprovados para limite atingido e pré-condições de
senha transacional.

### Escopo

- Mapear `INSTALLATION_LIMIT_REACHED` para um estado ou erro tipado do app.
- Exibir a mensagem aprovada de limite atingido:
  `Esta conta já possui 3 instalações cadastradas. A instalação atual ainda não está autorizada.`
- Continuar a mensagem com:
  `Acesse sua conta por uma instalação já autorizada e remova uma instalação antiga para liberar espaço. Depois, tente entrar novamente neste app.`
- Usar o texto de botão `Entendi`.
- Retornar ao login após a confirmação.
- Não iniciar step-up quando o limite for atingido.
- Em `TRANSACTION_PASSWORD_NOT_SET` ou `TRANSACTION_PASSWORD_LOCKED`, não
  registrar a instalação.
- Descartar o estado temporário do fluxo restrito nesses bloqueios.

### Critérios de Aceite

- Limite atingido exibe título, corpo e botão aprovados.
- Limite atingido não solicita senha transacional.
- Senha transacional não configurada não chama `POST /security/installations`.
- Senha transacional bloqueada não chama `POST /security/installations`.
- Todos os resultados de bloqueio retornam o usuário ao login ou a um estado
  informativo que exige novo login.

### Depende De

- Tarefa 6.
- Tarefa 9.

## Tarefa 11/12: Adicionar telemetria segura e testes

### Objetivo

Tornar falhas observáveis sem vazar identificadores ou credenciais, e cobrir o
ciclo de vida crítico da identidade.

### Escopo

- Registrar falhas de armazenamento/leitura/escrita/validação sem o id completo
  da instalação.
- Não registrar tokens em log.
- Não registrar a senha transacional em log.
- Adicionar testes para criação da identidade na primeira execução.
- Adicionar testes para reutilização quando o marcador local estiver presente.
- Adicionar testes para substituição de um valor antigo do armazenamento seguro
  quando o marcador local estiver ausente.
- Adicionar testes para comportamento de bloqueio em falha.
- Adicionar testes para injeção do header pelo interceptor.
- Adicionar testes para injeção do header no refresh.
- Adicionar testes para login restrito e sucesso no registro.
- Adicionar testes para limite atingido e bloqueios de senha transacional.

### Critérios de Aceite

- Logs contêm contexto do evento, mas nenhum UUID completo, token ou senha.
- Testes do ciclo de vida da identidade passam.
- Testes de headers HTTP passam.
- Testes do repositório/API de autenticação cobrem resultados operacionais,
  restritos e bloqueados.

### Depende De

- Tarefa 2.
- Tarefa 4.
- Tarefa 5.
- Tarefa 9.
- Tarefa 10.

## Tarefa 12/12: Executar verificação mobile

### Objetivo

Confirmar que o app mobile permanece saudável após a mudança de identidade de
instalação.

### Escopo

- Executar `dart format` nos arquivos Dart alterados.
- Executar testes focados em identidade, interceptors, parsing de autenticação e
  registro.
- Executar `flutter analyze`.
- Executar `flutter test` se viável.
- Corrigir regressões introduzidas pela implementação.

### Critérios de Aceite

- Arquivos Dart alterados estão formatados.
- Testes focados passam.
- `flutter analyze` passa.
- `flutter test` passa ou qualquer execução pulada/indisponível é documentada.
- Comportamentos existentes de login, refresh, logout e step-up de transferência
  permanecem inalterados.

### Depende De

- Tarefa 11.

## Adiado

- `GET /security/installations`.
- `DELETE /security/installations/{installation_resource_id}`.
- Tela de gerenciamento de instalações.
- Revogar uma instalação existente pelo mobile.
- Identificar e desabilitar a remoção da instalação atual no gerenciamento.
