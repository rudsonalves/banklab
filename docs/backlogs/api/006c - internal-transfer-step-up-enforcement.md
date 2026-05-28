# Backlog: enforcement na transferência interna

## 1. Objetivo

Proteger a transferência interna com o step-up token emitido pelo módulo de
segurança.

No primeiro corte do MVP, apenas a transferência interna exigirá senha
transacional.

## 2. Endpoint protegido

```http
POST /accounts/internal-transfers
Authorization: Bearer <access_token>
X-Step-Up-Token: <step_up_token>
```

Endpoint lógico:

```text
internal_transfer.create
```

## 3. Regras de enforcement

Antes de chamar o use case de transferência, o backend deve validar:

- `X-Step-Up-Token` está presente;
- o token é válido e assinado;
- o token não expirou;
- o `jti` existe no backend;
- o `jti` ainda não foi consumido;
- o usuário do JWT normal é o mesmo usuário do step-up token;
- o `endpoint_key` do token é `internal_transfer.create`.

Se a validação passar:

- o `jti` deve ser consumido de forma atômica;
- o use case de transferência interna pode ser chamado.

Se a validação falhar:

- o use case de transferência interna não deve ser chamado;
- a API deve responder usando o envelope padrão de erro.

## 4. Posicionamento arquitetural

No MVP, a delivery do módulo de transferência chama explicitamente o enforcement
do módulo de segurança antes do use case sensível.

```text
internal/account/transaction/delivery
  -> internal/security/application
  -> internal/account/transaction/application
```

O use case de transferência não deve receber senha transacional.

## 5. Fluxo

```text
1. Mobile chama POST /accounts/internal-transfers com JWT e X-Step-Up-Token.
2. Delivery de transferência identifica endpoint_key=internal_transfer.create.
3. Delivery chama o enforcement de security/application.
4. Security valida token, usuário, endpoint, expiração e consumo.
5. Security consome o jti.
6. Delivery chama o use case de transferência.
```

## 6. Erros cobertos

- `STEP_UP_TOKEN_REQUIRED`
- `STEP_UP_TOKEN_INVALID`
- `STEP_UP_TOKEN_EXPIRED`
- `STEP_UP_TOKEN_CONSUMED`
- `STEP_UP_ENDPOINT_MISMATCH`

O contrato completo de erros fica no backlog
[006d - zta-contracts-and-docs.md](<006d - zta-contracts-and-docs.md>).
