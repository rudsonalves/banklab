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
- No MVP, o único `endpoint_key` aceito é `internal_transfer.create`.
- A validação de `endpoint_key` será uma whitelist simples no módulo de
  segurança.
- A whitelist deve ser isolada o suficiente para evoluir depois para um policy
  registry ou policy engine.
- O estado `expired` não precisa ser persistido como status. Um token deve ser
  considerado expirado por regra quando `status=active` e `expires_at < now`.

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
8. Em caso de sucesso, zera falhas e prepara a emissão do step-up token.
9. Backend persiste `jti` com expiração de 2 minutos.
10. Backend emite JWT de step-up contendo o `jti` persistido.
11. Backend responde com step_up_token usando o envelope padrão.
```

Observações:

- A validação da senha transacional deve reutilizar o domínio, repositório e
  hasher criados no backlog `006a`.
- Em caso de senha correta, o fluxo deve chamar a regra de sucesso da senha
  transacional para zerar tentativas inválidas.
- Em caso de senha incorreta, o fluxo deve registrar falha, salvar o novo estado
  e aplicar bloqueio temporário se a terceira tentativa for atingida.
- Se a persistência do `jti` falhar, o JWT de step-up não deve ser entregue ao
  mobile.

## 5. Modelo conceitual

```text
step_up_tokens
- id
- jti
- user_id
- endpoint_key
- status
- expires_at
- consumed_at
- created_at
```

Regras do modelo:

- `jti` é obrigatório e deve ser único.
- `user_id` identifica o usuário autenticado para quem o token foi emitido.
- `endpoint_key` limita o uso do token ao endpoint lógico autorizado.
- `expires_at` define a validade curta de 2 minutos.
- `consumed_at` deve ser preenchido apenas quando o token for usado com sucesso.

Status persistidos iniciais:

```text
active
consumed
```

Estado derivado:

```text
expired = status active com expires_at anterior ao horário atual
```

## 6. Erros cobertos

- `TRANSACTION_PASSWORD_NOT_SET`
- `TRANSACTION_PASSWORD_INVALID`
- `TRANSACTION_PASSWORD_LOCKED`
- `STEP_UP_ENDPOINT_NOT_ALLOWED`

Observação:

- `TRANSACTION_PASSWORD_REQUIRED` pertence ao enforcement de endpoints
  sensíveis, tratado no backlog `006c`.
- `STEP_UP_TOKEN_INVALID`, `STEP_UP_TOKEN_EXPIRED` e
  `STEP_UP_TOKEN_CONSUMED` pertencem principalmente à validação/consumo do token
  no backlog `006c`.

O contrato completo de erros fica no backlog
[006d - zta-contracts-and-docs.md](<006d - zta-contracts-and-docs.md>).
