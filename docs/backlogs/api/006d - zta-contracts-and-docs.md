# Backlog: contratos e documentação do ZTA MVP

## 1. Objetivo

Consolidar os contratos HTTP, códigos de erro e documentação do MVP de ZTA com
senha transacional e step-up token.

Este backlog não introduz nova regra de negócio. Ele garante que API, mobile e
documentação falem a mesma língua.

## 2. Padrão de resposta

Todas as respostas devem seguir o padrão definido em
`api/docs/05-error_and_response.md`.

Sucesso:

```json
{
  "data": {},
  "error": null
}
```

Erro:

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {}
  }
}
```

Clientes devem depender de `error.code`, não de `error.message`.

## 3. Nomes definidos

```text
Endpoint de criação da senha transacional:
POST /security/transaction-password

Endpoint de autorização de step-up:
POST /security/step-up/authorize

Endpoint sensível protegido no primeiro corte:
POST /accounts/internal-transfers

Header do token de step-up:
X-Step-Up-Token

Endpoint lógico:
internal_transfer.create
```

Campos JSON:

```text
endpoint_key
transaction_password
transaction_password_confirmation
step_up_token
expires_in
```

## 4. JWT de step-up

O step-up token é um JWT curto, assinado pelo backend com `HS256`, com validade
de 120 segundos. O backend também persiste o `jti` para permitir consumo único
no enforcement.

Claims mínimos:

```text
user_id
endpoint_key
scope=step_up
exp
iat
jti
```

O token não deve conter senha transacional, hash, dados sensíveis ou payload da
operação.

## 5. Resposta de autorização de step-up

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 120
  },
  "error": null
}
```

## 6. Contrato inicial de erros

| Código                             | HTTP | Cenário                                                |
| ---------------------------------- | ---: | ------------------------------------------------------ |
| `TRANSACTION_PASSWORD_ALREADY_SET` |  409 | usuário tenta criar senha transacional já ativa        |
| `TRANSACTION_PASSWORD_REQUIRED`    |  403 | endpoint sensível exige step-up                        |
| `TRANSACTION_PASSWORD_NOT_SET`     |  409 | usuário ainda não possui senha transacional cadastrada |
| `TRANSACTION_PASSWORD_INVALID`     |  401 | PIN transacional inválido                              |
| `TRANSACTION_PASSWORD_LOCKED`      |  403 | PIN transacional bloqueado temporariamente             |
| `STEP_UP_ENDPOINT_NOT_ALLOWED`     |  403 | endpoint lógico não autorizado para emissão step-up    |
| `STEP_UP_TOKEN_REQUIRED`           |  401 | endpoint sensível foi chamado sem `X-Step-Up-Token`    |
| `STEP_UP_TOKEN_INVALID`            |  401 | token de step-up inválido ou malformado                |
| `STEP_UP_TOKEN_EXPIRED`            |  401 | token de step-up expirado                              |
| `STEP_UP_TOKEN_CONSUMED`           |  401 | token de step-up já utilizado                          |
| `STEP_UP_ENDPOINT_MISMATCH`        |  403 | token válido, mas emitido para outro endpoint lógico   |

Os códigos acima são contrato do MVP. As mensagens podem seguir o padrão textual
da API, desde que o `error.code` permaneça estável.

Separação entre autorização e enforcement:

- `STEP_UP_ENDPOINT_NOT_ALLOWED` pertence à emissão/autorização de step-up:
  o cliente pediu autorização para um `endpoint_key` fora da whitelist.
- `TRANSACTION_PASSWORD_REQUIRED`, `STEP_UP_TOKEN_REQUIRED`,
  `STEP_UP_TOKEN_INVALID`, `STEP_UP_TOKEN_EXPIRED`,
  `STEP_UP_TOKEN_CONSUMED` e `STEP_UP_ENDPOINT_MISMATCH` pertencem ao
  enforcement do endpoint sensível.
- `STEP_UP_TOKEN_INVALID` cobre token malformado, assinatura inválida,
  `scope` diferente de `step_up`, `jti` inexistente ou claims obrigatórios
  ausentes.

## 7. Documentos a atualizar

Quando a implementação avançar, revisar:

- `api/docs/05-error_and_response.md`;
- `api/docs/07-api-rest.md`;
- `api/docs/implementations/03-zta-step-up-transaction-password.md`;
- documentação mobile, se houver impacto direto no consumo.

## 8. Status de alinhamento (MVP enforcement)

Contrato confirmado com a implementação atual:

- Endpoint sensível protegido: `POST /accounts/internal-transfers`.
- Header obrigatório de enforcement: `X-Step-Up-Token`.
- Endpoint lógico exigido no token: `internal_transfer.create`.
- Ponto arquitetural do enforcement:
  `internal/account/transaction/delivery -> internal/security/application`
  antes de `internal/account/transaction/application`.
- O use case de transferência não recebe step-up token, senha transacional nem
  detalhes de política ZTA.
- O step-up token autoriza uma chamada ao endpoint lógico e não é vinculado ao
  payload da transferência no MVP.
- Step-up token é de uso único com consumo atômico por `jti` antes do use case
  sensível.
- Se a transferência falhar depois do enforcement, o step-up token permanece
  consumido.
- Retry com o mesmo `X-Step-Up-Token` após consumo retorna
  `STEP_UP_TOKEN_CONSUMED`.
- Retry com mesmo `idempotency_key` pode exigir novo step-up token.

Verificação final:

- `go test ./...` passa na API.
- Não há divergência conhecida entre constantes, registry de erros,
  signer/verifier JWT, enforcement, handler de transferência interna e
  documentação revisada.
- Este backlog pode ser considerado fechado para o escopo do MVP.

Fora do escopo do MVP:

- vincular o step-up token ao payload detalhado da operação;
- ampliar a policy para outros endpoints sensíveis;
- implementar dispositivo confiável, biometria local, prova de vida ou sinais
  de risco;
- exigir coleção Postman como artefato deste backlog.
