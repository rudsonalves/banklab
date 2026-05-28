# Backlog: step-up token

## 1. Objetivo

Implementar a autorização de step-up que valida a senha transacional e emite um
token curto, de uso único, vinculado a um endpoint lógico.

## 2. Decisões fechadas

- O step-up token tem duração de 2 minutos.
- O step-up token é de uso único.
- O modelo é híbrido: JWT assinado com `jti` persistido no backend.
- O JWT de step-up deve conter, no mínimo:
  - `user_id`;
  - `endpoint_key`;
  - `scope=step_up`;
  - `exp`;
  - `iat`;
  - `jti`.
- O `jti` deve ser persistido no momento da emissão.
- O uso único deve ser garantido por consumo atômico do `jti` persistido.
- O token é escopado para um endpoint lógico específico.
- No MVP, o token não é vinculado ao payload da operação.

## 3. Autorização de step-up

Endpoint definido:

```http
POST /security/step-up/authorize
Authorization: Bearer <access_token>
```

Payload definido:

```json
{
  "endpoint_key": "internal_transfer.create",
  "transaction_password": "123456"
}
```

Resposta definida:

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 120
  },
  "error": null
}
```

## 4. Fluxo

```text
1. Mobile chama POST /security/step-up/authorize.
2. Backend identifica o usuário autenticado.
3. Backend valida endpoint_key.
4. Backend busca a senha transacional do usuário.
5. Backend verifica se a senha existe e não está bloqueada.
6. Backend compara PIN recebido com hash armazenado.
7. Em caso de falha, incrementa tentativas e aplica bloqueio se necessário.
8. Em caso de sucesso, zera falhas e emite JWT de step-up.
9. Backend persiste jti com expiração de 2 minutos.
10. Backend responde com step_up_token usando o envelope padrão.
```

## 5. Modelo conceitual

```text
step_up_tokens
- id
- user_id
- endpoint_key
- status
- expires_at
- consumed_at
- created_at
```

Status iniciais:

```text
active
consumed
expired
```

## 6. Erros cobertos

- `TRANSACTION_PASSWORD_REQUIRED`
- `TRANSACTION_PASSWORD_NOT_SET`
- `TRANSACTION_PASSWORD_INVALID`
- `TRANSACTION_PASSWORD_LOCKED`
- `STEP_UP_TOKEN_INVALID`

O contrato completo de erros fica no backlog
[006d - zta-contracts-and-docs.md](<006d - zta-contracts-and-docs.md>).
