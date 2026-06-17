# Backlog API 015: Installation Identity Session, Tokens and Context

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Preparar a infraestrutura de sessão, JWT, token restrito e contexto
autenticado antes dos use cases dependerem do vínculo de instalação.

## 3. Escopo

- Vincular sessão operacional ao par usuário + instalação.
- Emitir `access_token` operacional com claim `installation_id`.
- Preparar refresh token para revalidar a instalação vinculada.
- Definir e emitir `restricted_access_token` com:
  - `sub`;
  - `jti`;
  - `token_type = restricted_access`;
  - `scope = installation.register`;
  - `installation_id`;
  - `iat`;
  - `exp`.
- Validar token restrito contra autorização persistida.
- Definir contexto operacional com usuário, customer e instalação.
- Definir contexto restrito com `user_id`, `installation_id`, `jti` e `scope`.
- Preparar interfaces para invalidar sessões por instalação revogada.

## 4. Fora de escopo

- Implementar handlers HTTP.
- Classificar instalação no login.
- Registrar ou revogar instalações por endpoint.

## 5. Dependencias

- Backlog 012 para domínio e portas.
- Backlog 014 para repositórios de instalação/autorização.
- Fluxo de sessões existente da API.

## 6. Preparacao para tasks

As tasks devem separar mudanças de claims, serviço de token restrito, contexto
autenticado e persistência do vínculo de sessão.

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
