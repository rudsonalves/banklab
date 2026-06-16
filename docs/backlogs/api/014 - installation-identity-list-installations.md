# Backlog API 014: List installations on `GET /security/installations`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: Medium
- Estado: Discussão

## 2. Objetivo

Listar as instalações associadas ao usuário autenticado para gerenciamento e
auditoria básica.

## 3. Contrato

```http
GET /security/installations
Authorization: Bearer <access_token>
X-Installation-Id: <UUID v4>
```

## 4. Regras consolidadas

- listar apenas instalações do usuário autenticado;
- expor um identificador público de gerenciamento distinto de
  `installation_id`;
- refletir os estados `known` e `revoked`;
- usar como metadados mínimos do MVP:
  `platform`, `app_version`, `app_build`, `first_seen_at`, `last_seen_at`,
  `revoked_at`, `created_at` e `updated_at`;
- exigir correspondência entre `X-Installation-Id`, sessão e claim do token.

## 5. Decisão ainda aberta

- definir retenção, auditoria e minimização desses metadados persistidos.

## 6. Tasks derivadas

- modelar resposta de listagem;
- mapear `installation_resource_id` público;
- retornar estados e metadados do MVP;
- cobrir usuário sem instalações revogadas e com histórico revogado.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [Revoke installation](<015 - installation-identity-revoke-installation.md>)
