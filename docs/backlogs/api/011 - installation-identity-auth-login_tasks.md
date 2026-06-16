# Tasks do Backlog API 011

Backlog pai:

- `011 - installation-identity-auth-login.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/8: Validar `X-Installation-Id` no login

Status: Backlog

### Objetivo

Garantir que `POST /auth/login` exija `X-Installation-Id` em formato UUID v4
canônico desde o primeiro release da feature.

### Escopo

- Ler `X-Installation-Id` no handler de login.
- Rejeitar ausência do header.
- Rejeitar valor malformado ou fora do formato canônico.
- Mapear erro para `400 INVALID_INSTALLATION_ID`.
- Não aceitar fallback silencioso sem header.

### Critérios de aceite

- Login sem `X-Installation-Id` falha com erro de contrato estável.
- Login com UUID inválido falha com `400 INVALID_INSTALLATION_ID`.
- Login com UUID v4 canônico segue o fluxo normal.

### Depende de

- Nenhuma dependência.

## Task 2/8: Introduzir serviço de classificação da instalação no login

Status: Backlog

### Objetivo

Classificar a instalação recebida no login para decidir entre instalação
conhecida, primeira instalação, nova instalação, instalação revogada ou limite
atingido.

### Escopo

- Definir contrato de consulta de instalações por usuário e `installation_id`.
- Representar os estados necessários para a decisão no login.
- Detectar:
  - instalação `known`;
  - primeira instalação do usuário;
  - nova instalação com vagas disponíveis;
  - instalação `revoked`;
  - limite de três instalações `known` atingido.
- Considerar instalações revogadas como histórico prévio que impede novo
  bootstrap automático.

### Critérios de aceite

- O login consegue distinguir todos os estados do fluxo sem depender do
  delivery layer.
- Instalações revogadas contam como histórico anterior.
- O limite considera apenas instalações `known`.

### Depende de

- Task 1.
- Backlog 018.

## Task 3/8: Implementar bootstrap atômico da primeira instalação

Status: Backlog

### Objetivo

Cadastrar automaticamente a primeira instalação de uma conta durante o login,
sem permitir corrida entre logins concorrentes.

### Escopo

- Executar verificação e criação da primeira instalação em bloco atômico.
- Criar associação `known` para a primeira instalação.
- Vincular a sessão operacional à instalação recém-criada.
- Preservar a regra de que instalações revogadas impedem novo bootstrap.

### Critérios de aceite

- Dois logins concorrentes não criam duas “primeiras instalações”.
- A primeira instalação criada já nasce vinculada à sessão.
- Histórico revogado não reabre elegibilidade para bootstrap automático.

### Depende de

- Task 2.
- Backlog 018.

## Task 4/8: Emitir sessão operacional normal para instalação conhecida

Status: Backlog

### Objetivo

Fazer com que uma instalação já conhecida siga o login normal, com tokens
operacionais vinculados à instalação.

### Escopo

- Reutilizar o fluxo atual de validação de credenciais.
- Criar sessão operacional vinculada ao par usuário + `installation_id`.
- Emitir `access_token` com claim `installation_id`.
- Emitir `refresh_token` normal vinculado à sessão.

### Critérios de aceite

- Instalação `known` recebe `access_token` e `refresh_token` normais.
- O `access_token` inclui `installation_id`.
- A sessão criada pode ser validada depois por refresh e middleware.

### Depende de

- Task 2.
- Backlog 018.

## Task 5/8: Emitir autorização restrita para nova instalação com vaga

Status: Backlog

### Objetivo

Retornar `restricted_access_token` quando as credenciais forem válidas, mas a
instalação ainda não estiver cadastrada e houver vaga disponível.

### Escopo

- Integrar o fluxo de login com a persistência de
  `installation_registration_authorizations`.
- Emitir `restricted_access_token` com TTL curto.
- Retornar `authentication_status = installation_registration_required`.
- Não emitir `refresh_token` nesse cenário.
- Garantir consistência com o modelo definido no backlog 018.

### Critérios de aceite

- Nova instalação com vaga não recebe sessão operacional completa.
- A resposta contém `restricted_access_token` e `expires_in`.
- O grant fica persistido para consumo posterior.

### Depende de

- Task 2.
- Backlog 018.

## Task 6/8: Negar login para instalação revogada

Status: Backlog

### Objetivo

Impedir que uma instalação revogada seja reativada silenciosamente por novo
login.

### Escopo

- Detectar instalação `revoked` durante a classificação do login.
- Bloquear o login sem emitir sessão operacional nem autorização restrita.
- Mapear erro de negócio estável para o contrato HTTP do login.
- Preservar ausência de fluxo de recuperação no MVP.

### Critérios de aceite

- Instalação revogada não recebe tokens operacionais.
- Instalação revogada não recebe `restricted_access_token`.
- O login falha com erro consistente e testável.

### Depende de

- Task 2.

## Task 7/8: Retornar `installation_limit_reached`

Status: Backlog

### Objetivo

Encerrar o login com estado estruturado quando o usuário já tiver três
instalações `known` e tentar autenticar em uma nova instalação.

### Escopo

- Detectar limite de três instalações `known`.
- Retornar HTTP 200 com:
  - `authentication_status = installation_limit_reached`
  - `max_installations`
  - `known_installations_count`
  - `installation_registered = false`
  - `next_action = revoke_existing_installation`
- Não emitir `restricted_access_token`.
- Não criar sessão operacional.

### Critérios de aceite

- Nova instalação sem vaga recebe `installation_limit_reached`.
- A resposta segue exatamente o contrato final do backlog 011.
- O login não abre fluxo de step-up nem registro nesse cenário.

### Depende de

- Task 2.

## Task 8/8: Cobrir o login com testes de contrato e fluxo

Status: Backlog

### Objetivo

Garantir cobertura automatizada para os novos comportamentos do login com
identidade de instalação.

### Escopo

- Adicionar testes de aplicação e/ou integração para:
  - header ausente;
  - header malformado;
  - instalação conhecida;
  - primeira instalação;
  - nova instalação com vaga;
  - instalação revogada;
  - limite atingido;
  - corrida no bootstrap da primeira instalação, quando viável no nível
    escolhido.
- Validar claims e payloads de resposta relevantes.

### Critérios de aceite

- Os principais ramos de decisão do login ficam cobertos por testes.
- Os contratos `installation_registration_required` e
  `installation_limit_reached` ficam protegidos por testes.
- O vínculo da sessão com `installation_id` aparece nas asserções
  correspondentes.

### Depende de

- Tasks 1 a 7.
