# Backlog API 018: Shared infrastructure for Installation Identity

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Concentrar os elementos compartilhados entre os endpoints do MVP de identidade
de instalação.

## 3. Modelo de dados

### `app_installations`

```text
- id
- user_id
- installation_id
- status
- platform
- app_version
- app_build
- first_seen_at
- last_seen_at
- revoked_at
- created_at
- updated_at
```

Estados consolidados:

```text
known
revoked
```

### `installation_registration_authorizations`

```text
- id
- jti
- user_id
- installation_id
- scope
- status
- expires_at
- consumed_at
- created_at
```

Status mínimos:

```text
active
consumed
revoked
```

## 4. Regras consolidadas

- a sessão operacional nasce vinculada ao par usuário + `installation_id`;
- o `access_token` operacional carrega `installation_id`;
- o `restricted_access_token` é JWT curto com:
  `sub`, `jti`, `token_type=restricted_access`,
  `scope=installation.register`, `installation_id`, `iat` e `exp`;
- `jti` deve ser único;
- deve existir no máximo uma autorização `active` por
  `(user_id, installation_id, scope)`;
- o header é obrigatório desde o primeiro release da feature, sem fase de
  compatibilidade;
- não existe fluxo de recuperação no MVP quando todas as instalações anteriores
  estiverem indisponíveis ou revogadas.

## 5. Decisão ainda aberta

- definir retenção, auditoria e minimização dos metadados persistidos.

## 6. Tasks derivadas

- criar tabela `app_installations`;
- criar tabela `installation_registration_authorizations`;
- adaptar emissão de JWT operacional com `installation_id`;
- criar middleware específico para acesso restrito;
- adaptar contexto autenticado para vínculo com instalação.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
