# Tasks do Backlog API 012

Backlog pai:

- `012 - installation-identity-step-up-authorize.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/5: Registrar `POST /security/installations` na política de step-up

Status: Backlog

### Objetivo

Adicionar `POST /security/installations` à política/lista de operações públicas
autorizáveis por senha transacional.

### Escopo

- Registrar o endpoint na resolução de operações públicas do `step-up`.
- Garantir consistência entre método HTTP, path e endpoint key.
- Cobrir testes de resolução e validação do endpoint.

### Critérios de aceite

- `POST /security/installations` é reconhecido como operação elegível.
- Método ou path divergente continuam sendo rejeitados.

### Depende de

- Nenhuma dependência.

## Task 2/5: Permitir `restricted_access_token` em `POST /security/step-up/authorize`

Status: Backlog

### Objetivo

Fazer com que o endpoint de `step-up` aceite autorização restrita no fluxo de
registro de instalação.

### Escopo

- Adaptar middleware/handler para aceitar contexto restrito.
- Garantir que o fluxo não exija sessão operacional completa.
- Preservar comportamento atual das demais rotas protegidas.

### Critérios de aceite

- O endpoint aceita `restricted_access_token` válido para este fluxo.
- Tokens operacionais normais continuam funcionando conforme esperado.

### Depende de

- Task 1.
- Backlog 018.

## Task 3/5: Bloquear emissão do `step_up_token` quando a senha transacional não estiver apta

Status: Backlog

### Objetivo

Impedir autorização de registro de instalação quando a senha transacional
estiver `not_set` ou `locked`.

### Escopo

- Integrar consulta do estado da senha transacional.
- Bloquear emissão do `step_up_token` para `not_set`.
- Bloquear emissão do `step_up_token` para `locked`.
- Mapear erro de forma estável no contrato do endpoint.

### Critérios de aceite

- Usuário sem senha transacional ativa não recebe `step_up_token`.
- Usuário com senha bloqueada não recebe `step_up_token`.
- O contrato de erro fica estável e testável.

### Depende de

- Task 2.

## Task 4/5: Emitir `step_up_token` escopado para o registro de instalação

Status: Backlog

### Objetivo

Garantir que o token emitido para esse fluxo permaneça curto, de uso único e
escopado a `POST /security/installations`.

### Escopo

- Emitir token com endpoint key correspondente ao registro de instalação.
- Persistir o token para consumo posterior.
- Manter TTL curto e sem payload sensível.
- Reusar o modelo de emissão/consumo já existente no `step-up`.

### Critérios de aceite

- O token emitido só serve para `POST /security/installations`.
- O token é persistido para consumo único.
- O payload não expõe material sensível.

### Depende de

- Tasks 1 a 3.

## Task 5/5: Cobrir o fluxo de `step-up` para registro de instalação com testes

Status: Backlog

### Objetivo

Proteger o contrato HTTP e a integração do `step-up` no novo fluxo.

### Escopo

- Testar autorização com `restricted_access_token` válido.
- Testar recusa para senha `not_set`.
- Testar recusa para senha `locked`.
- Testar emissão com endpoint correto e rejeição para endpoint divergente.

### Critérios de aceite

- O fluxo principal e os ramos de recusa ficam cobertos por testes.
- O contrato do endpoint permanece protegido contra regressões.

### Depende de

- Tasks 1 a 4.
