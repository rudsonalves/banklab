# Tasks do Backlog API 018

Backlog pai:

- `018 - installation-identity-shared-infrastructure.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/6: Criar migration de `app_installations`

Status: Backlog

### Objetivo

Preparar a persistência principal das instalações cadastradas.

### Escopo

- Criar tabela `app_installations`.
- Incluir campos do modelo do MVP.
- Representar estados `known` e `revoked`.
- Criar índices e constraints necessários.
- Criar rollback correspondente.

### Critérios de aceite

- A migration executa e faz rollback corretamente.
- O schema suporta o limite, o histórico e os metadados mínimos do MVP.

### Depende de

- Nenhuma dependência.

## Task 2/6: Criar migration de `installation_registration_authorizations`

Status: Backlog

### Objetivo

Persistir a autorização restrita usada antes do registro explícito da
instalação.

### Escopo

- Criar tabela `installation_registration_authorizations`.
- Incluir `jti`, `user_id`, `installation_id`, `scope`, `status`, `expires_at`,
  `consumed_at` e `created_at`.
- Garantir unicidade de `jti`.
- Garantir no máximo uma autorização `active` por
  `(user_id, installation_id, scope)`.
- Criar rollback correspondente.

### Critérios de aceite

- A tabela suporta o ciclo de vida do grant restrito.
- Constraints impedem duplicidade indevida de grant ativo.

### Depende de

- Nenhuma dependência.

## Task 3/6: Implementar domínio e contratos de repositório das instalações

Status: Backlog

### Objetivo

Definir as portas e os tipos compartilhados para instalações e grants
restritos.

### Escopo

- Criar tipos de domínio para instalação e autorização restrita.
- Definir estados `known` e `revoked` para instalação.
- Definir estados `active`, `consumed` e `revoked` para o grant.
- Definir interfaces de repositório necessárias para login, registro,
  listagem, revogação e refresh.

### Critérios de aceite

- Application consegue operar usando apenas contratos.
- Os estados e transições do MVP ficam explicitados no domínio.

### Depende de

- Tasks 1 e 2.

## Task 4/6: Implementar persistência Postgres de instalações e grants

Status: Backlog

### Objetivo

Criar a infraestrutura Postgres compartilhada para instalações e autorizações
restritas.

### Escopo

- Implementar repositório Postgres de `app_installations`.
- Implementar repositório Postgres de
  `installation_registration_authorizations`.
- Cobrir criação, consulta, listagem, revogação e consumo do grant.

### Critérios de aceite

- Os repositórios persistem corretamente instalações e grants.
- Casos de conflito e constraints são mapeados de forma estável.

### Depende de

- Task 3.

## Task 5/6: Implementar suporte de JWT e contexto para sessão e acesso restrito

Status: Backlog

### Objetivo

Preparar a infraestrutura de token e contexto usada pelos demais endpoints do
MVP.

### Escopo

- Emitir `access_token` operacional com `installation_id`.
- Emitir e validar `restricted_access_token`.
- Criar contexto autenticado para fluxo restrito.
- Separar middleware de sessão operacional e acesso restrito.

### Critérios de aceite

- Os dois tipos de token ficam representados e parseáveis.
- O contexto autenticado distingue sessão operacional de grant restrito.

### Depende de

- Tasks 2 a 4.

## Task 6/6: Definir retenção, auditoria e minimização dos metadados

Status: Backlog

### Objetivo

Fechar a política de retenção e auditoria dos metadados persistidos das
instalações.

### Escopo

- Definir prazo de retenção para instalações revogadas.
- Definir quais metadados são mantidos para auditoria.
- Definir o que não deve ser persistido além do mínimo do MVP.
- Registrar a decisão de forma explícita nos backlogs.

### Critérios de aceite

- A política de retenção fica explícita e implementável.
- Auditoria e minimização ficam alinhadas ao conjunto mínimo do MVP.

### Depende de

- Task 3.
