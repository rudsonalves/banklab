# Tasks do Backlog API 017

Backlog pai:

- `017 - installation-identity-session-enforcement.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/5: Estender o token operacional com `installation_id`

Status: Backlog

### Objetivo

Permitir que o `access_token` carregue o vínculo da instalação operacional.

### Escopo

- Adaptar claims do token operacional.
- Incluir `installation_id` na emissão do `access_token`.
- Atualizar parsing e representação do contexto autenticado.

### Critérios de aceite

- Tokens operacionais novos incluem `installation_id`.
- O parsing do token expõe esse claim de forma segura e consistente.

### Depende de

- Backlog 018.

## Task 2/5: Adaptar middleware de autenticação operacional

Status: Backlog

### Objetivo

Comparar `X-Installation-Id` com o `installation_id` do token em toda rota
autenticada.

### Escopo

- Ler `X-Installation-Id` no middleware de sessão operacional.
- Rejeitar ausência/malformação com `400 INVALID_INSTALLATION_ID`.
- Rejeitar divergência com `403 INSTALLATION_MISMATCH`.
- Propagar o vínculo da instalação para o contexto autenticado.

### Critérios de aceite

- O middleware bloqueia instalação divergente.
- O contexto autenticado passa a carregar o vínculo da instalação.

### Depende de

- Task 1.

## Task 3/5: Aplicar enforcement nas rotas autenticadas relevantes

Status: Backlog

### Objetivo

Garantir que as rotas autenticadas do MVP passem a respeitar o vínculo de
instalação.

### Escopo

- Aplicar o enforcement em logout, quando exposto.
- Aplicar o enforcement nas rotas autenticadas por `access_token`.
- Preservar `restricted_access_token` em middleware separado.

### Critérios de aceite

- Rotas operacionais autenticadas passam a exigir `X-Installation-Id`.
- O fluxo restrito continua separado do fluxo operacional.

### Depende de

- Task 2.

## Task 4/5: Impedir uso de `access_token` de instalação revogada

Status: Backlog

### Objetivo

Negar requisições autenticadas quando a instalação associada ao token tiver sido
revogada.

### Escopo

- Consultar o estado da instalação quando necessário para enforcement.
- Rejeitar `access_token` de instalação revogada.
- Integrar com o corte imediato de acesso definido na revogação.

### Critérios de aceite

- Requisições autenticadas de instalação revogada são negadas.
- O corte imediato permanece consistente com a revogação.

### Depende de

- Task 3.
- Backlog 015.

## Task 5/5: Cobrir middleware e enforcement com testes

Status: Backlog

### Objetivo

Proteger o enforcement de sessão com cobertura automatizada.

### Escopo

- Testar ausência e malformação do header.
- Testar divergência entre header e claim.
- Testar propagação do contexto autenticado com instalação.
- Testar negação de instalação revogada.

### Critérios de aceite

- O middleware e seus ramos principais ficam cobertos.
- O vínculo da instalação fica protegido contra regressões.

### Depende de

- Tasks 1 a 4.
