# Backlog API 012: Step-up for `POST /security/installations`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Permitir que a senha transacional autorize especificamente o registro de uma
nova instalação, usando o fluxo já existente de `step-up`.

## 3. Contrato

```http
POST /security/step-up/authorize
Authorization: Bearer <restricted_access_token>
```

```json
{
  "method": "POST",
  "path": "/security/installations",
  "transaction_password": "123456"
}
```

Resposta esperada:

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 300
  },
  "error": null
}
```

## 4. Regras consolidadas

- `POST /security/installations` é operação elegível para `step-up`;
- o fluxo usa `restricted_access_token`, não sessão operacional completa;
- senha transacional ativa é pré-requisito;
- se a senha estiver `not_set` ou `locked`, o token não é emitido;
- o `step_up_token` continua curto, de uso único e escopado ao endpoint.

## 5. Tasks derivadas

- incluir `POST /security/installations` na política/lista de operações do
  `step-up`;
- aceitar contexto de autorização restrita;
- bloquear `not_set` e `locked` no fluxo;
- cobrir testes de emissão e recusa do `step_up_token`.

## 6. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [Installation registration](<013 - installation-identity-register-installation.md>)
