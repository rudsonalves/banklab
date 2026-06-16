# Backlog API 012: Installation Identity Persistence Foundation

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Criar a base de dados que sustenta o MVP antes dos fluxos de login, registro,
listagem, revogação e refresh passarem a depender da identidade de instalação.

## 3. Escopo

### `app_installations`

- Criar tabela de instalações.
- Incluir os metadados mínimos do MVP:
  `platform`, `app_version`, `app_build`, `first_seen_at`, `last_seen_at`,
  `revoked_at`, `created_at` e `updated_at`.
- Representar apenas os estados `known` e `revoked`.
- Garantir identificador público de gerenciamento separado do
  `installation_id` enviado pelo cliente.
- Criar índices e constraints para consulta por usuário, instalação e status.
- Preparar suporte ao limite de três instalações `known` por usuário.

### `installation_registration_authorizations`

- Criar tabela de autorizações restritas.
- Persistir `jti`, `user_id`, `installation_id`, `scope`, `status`,
  `expires_at`, `consumed_at` e `created_at`.
- Garantir unicidade de `jti`.
- Garantir no máximo uma autorização `active` por
  `(user_id, installation_id, scope)`.
- Tratar expiração como estado derivado de `expires_at`.

## 4. Fora de escopo

- Implementar handlers HTTP.
- Emitir tokens.
- Consumir autorização restrita.
- Executar bootstrap da primeira instalação.
- Aplicar middleware de sessão.

## 5. Dependências

- Backlog 011 apenas para manter o contrato de header consistente.

## 6. Orientação para tasks

As tasks deste backlog devem focar migrations, constraints, índices e rollback.
Elas precisam deixar claro como a base protege limite, histórico e unicidade
antes da camada de aplicação usar esses dados.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
