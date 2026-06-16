# Tasks do Backlog API 016

Backlog pai:

- `016 - installation-identity-auth-refresh.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/4: Exigir e validar `X-Installation-Id` no refresh

Status: Backlog

### Objetivo

Garantir que `POST /auth/refresh` exija `X-Installation-Id` canônico.

### Escopo

- Ler o header no handler de refresh.
- Rejeitar ausência ou valor malformado com `400 INVALID_INSTALLATION_ID`.
- Manter a validação desde o primeiro release da feature.

### Critérios de aceite

- Refresh sem header falha.
- Refresh com UUID inválido falha com erro de contrato estável.

### Depende de

- Nenhuma dependência.

## Task 2/4: Validar vínculo entre refresh, sessão e instalação

Status: Backlog

### Objetivo

Impedir que um refresh token seja usado a partir de outra instalação.

### Escopo

- Validar correspondência entre `X-Installation-Id` e a sessão persistida.
- Rejeitar divergência com `403 INSTALLATION_MISMATCH`.
- Preservar o vínculo do usuário com a instalação original.

### Critérios de aceite

- Refresh com instalação divergente é negado.
- Refresh com instalação correta segue o fluxo normal.

### Depende de

- Task 1.
- Backlog 018.

## Task 3/4: Negar refresh para instalação revogada

Status: Backlog

### Objetivo

Impedir renovação de tokens quando a instalação da sessão já estiver revogada.

### Escopo

- Verificar o estado da instalação no refresh.
- Bloquear renovação de instalação `revoked`.
- Cobrir integração com invalidação de sessão por revogação.

### Critérios de aceite

- Sessão de instalação revogada não recebe novos tokens.
- O refresh falha mesmo que o refresh token ainda exista materialmente.

### Depende de

- Task 2.
- Backlog 015.

## Task 4/4: Cobrir o refresh com testes

Status: Backlog

### Objetivo

Proteger o contrato e o vínculo da instalação no fluxo de refresh.

### Escopo

- Testar ausência e malformação do header.
- Testar divergência entre instalação e sessão.
- Testar instalação revogada.
- Testar sucesso com rotação do refresh token.

### Critérios de aceite

- Os principais ramos do refresh ficam cobertos por testes.
- A renovação mantém o vínculo correto da instalação.

### Depende de

- Tasks 1 a 3.
