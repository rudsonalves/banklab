# Tasks do Backlog API 015

Backlog pai:

- `015 - installation-identity-revoke-installation.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/5: Implementar o caso de uso de revogação lógica

Status: Backlog

### Objetivo

Revogar uma instalação preservando histórico e `revoked_at`.

### Escopo

- Criar caso de uso de revogação por `installation_resource_id`.
- Atualizar estado para `revoked`.
- Preencher `revoked_at`.
- Preservar o registro para auditoria e histórico.

### Critérios de aceite

- A instalação não é removida fisicamente.
- O estado passa a `revoked` com carimbo temporal correspondente.

### Depende de

- Backlog 018.

## Task 2/5: Impedir revogação da instalação em uso

Status: Backlog

### Objetivo

Bloquear tentativa de revogar a instalação vinculada à sessão atual.

### Escopo

- Comparar o alvo da revogação com a instalação da sessão atual.
- Retornar erro estável quando a instalação em uso for alvo da operação.

### Critérios de aceite

- O usuário não consegue revogar a instalação em uso.
- O erro deixa claro que a instalação atual não pode ser removida por esse
  fluxo.

### Depende de

- Task 1.
- Backlog 017.

## Task 3/5: Invalidar imediatamente sessões e tokens da instalação revogada

Status: Backlog

### Objetivo

Paralisar o acesso da instalação revogada no mesmo momento da revogação.

### Escopo

- Invalidar `refresh_token`/sessões associadas à instalação revogada.
- Impedir aceitação posterior de `access_token` da instalação revogada.
- Garantir efeito imediato do corte de acesso.

### Critérios de aceite

- Sessões da instalação revogada deixam de renovar tokens.
- Requisições autenticadas da instalação revogada passam a ser negadas.

### Depende de

- Task 1.
- Backlogs 016 e 017.

## Task 4/5: Implementar o endpoint HTTP de revogação

Status: Backlog

### Objetivo

Expor `DELETE /security/installations/{installation_resource_id}` no contrato
final do MVP.

### Escopo

- Criar handler e wiring da rota.
- Exigir autenticação operacional e `X-Installation-Id`.
- Não exigir `step-up`.
- Aplicar envelope padrão da API.

### Critérios de aceite

- A rota revoga logicamente outra instalação do usuário.
- A rota falha ao tentar revogar a instalação em uso.

### Depende de

- Tasks 1 a 3.

## Task 5/5: Cobrir a revogação com testes

Status: Backlog

### Objetivo

Proteger a operação de revogação e o corte imediato de acesso.

### Escopo

- Testar revogação com sucesso.
- Testar tentativa de autorrevogação.
- Testar invalidação imediata de refresh/sessão.

### Critérios de aceite

- O fluxo de sucesso e recusa ficam cobertos por testes.
- O corte imediato de acesso fica protegido por testes.

### Depende de

- Tasks 1 a 4.
