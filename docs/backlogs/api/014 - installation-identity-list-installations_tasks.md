# Tasks do Backlog API 014

Backlog pai:

- `014 - installation-identity-list-installations.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Média
- Área: API
- Tipo: Segurança/ZTA

## Task 1/4: Definir e implementar o caso de uso de listagem

Status: Backlog

### Objetivo

Permitir que o usuário autenticado liste suas instalações cadastradas e
revogadas.

### Escopo

- Criar caso de uso de listagem de instalações por usuário.
- Consultar somente instalações do usuário autenticado.
- Ordenar ou estruturar a lista de forma consistente para o cliente.

### Critérios de aceite

- O caso de uso retorna apenas instalações do usuário autenticado.
- Instalações `known` e `revoked` aparecem corretamente.

### Depende de

- Backlog 018.

## Task 2/4: Modelar a resposta pública da listagem

Status: Backlog

### Objetivo

Definir a resposta HTTP de `GET /security/installations` sem expor
`installation_id` como identificador de gerenciamento.

### Escopo

- Expor `installation_resource_id` público.
- Incluir estado da instalação.
- Incluir metadados mínimos do MVP:
  `platform`, `app_version`, `app_build`, `first_seen_at`, `last_seen_at`,
  `revoked_at`, `created_at` e `updated_at`.

### Critérios de aceite

- A resposta não usa `installation_id` como identificador público.
- Os metadados mínimos aparecem de forma consistente.

### Depende de

- Task 1.

## Task 3/4: Implementar o endpoint HTTP de listagem

Status: Backlog

### Objetivo

Expor `GET /security/installations` no contrato autenticado final.

### Escopo

- Criar handler e wiring da rota.
- Exigir autenticação operacional e `X-Installation-Id`.
- Aplicar envelope padrão da API.

### Critérios de aceite

- A rota lista as instalações do usuário autenticado.
- A resposta segue o envelope padrão.
- O vínculo da instalação atual é respeitado pelo middleware.

### Depende de

- Tasks 1 e 2.
- Backlog 017.

## Task 4/4: Cobrir listagem de instalações com testes

Status: Backlog

### Objetivo

Proteger a resposta de listagem e o escopo por usuário.

### Escopo

- Testar usuário sem instalações revogadas.
- Testar usuário com histórico revogado.
- Testar isolamento entre usuários.

### Critérios de aceite

- A rota não vaza instalações de outro usuário.
- Estados e metadados mínimos ficam cobertos por testes.

### Depende de

- Tasks 1 a 3.
