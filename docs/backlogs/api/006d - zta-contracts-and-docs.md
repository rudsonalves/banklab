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

## 4. Resposta de autorização de step-up

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 120
  },
  "error": null
}
```

## 5. Contrato inicial de erros

| Código                             | HTTP | Cenário                                                |
| ---------------------------------- | ---: | ------------------------------------------------------ |
| `TRANSACTION_PASSWORD_ALREADY_SET` |  409 | usuário tenta criar senha transacional já ativa        |
| `TRANSACTION_PASSWORD_REQUIRED`    |  403 | endpoint sensível exige step-up                        |
| `TRANSACTION_PASSWORD_NOT_SET`     |  409 | usuário ainda não possui senha transacional cadastrada |
| `TRANSACTION_PASSWORD_INVALID`     |  401 | PIN transacional inválido                              |
| `TRANSACTION_PASSWORD_LOCKED`      |  403 | PIN transacional bloqueado temporariamente             |
| `STEP_UP_TOKEN_REQUIRED`           |  401 | endpoint sensível foi chamado sem `X-Step-Up-Token`    |
| `STEP_UP_TOKEN_INVALID`            |  401 | token de step-up inválido ou malformado                |
| `STEP_UP_TOKEN_EXPIRED`            |  401 | token de step-up expirado                              |
| `STEP_UP_TOKEN_CONSUMED`           |  401 | token de step-up já utilizado                          |
| `STEP_UP_ENDPOINT_MISMATCH`        |  403 | token válido, mas emitido para outro endpoint lógico   |

Os códigos acima são contrato do MVP. As mensagens podem seguir o padrão textual
da API, desde que o `error.code` permaneça estável.

## 6. Documentos a atualizar

Quando a implementação avançar, revisar:

- `api/docs/05-error_and_response.md`;
- `api/docs/07-api-rest.md`;
- `api/docs/implementations/03-zta-step-up-transaction-password.md`;
- coleção Postman, se o fluxo estiver coberto nela;
- documentação mobile, se houver impacto direto no consumo.
