# Backlog API 018: Installation Identity Delivery and Enforcement

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Conectar os casos de uso a HTTP, middleware e contrato REST depois que domínio,
persistência, repositórios, tokens e application estiverem prontos.

## 3. Escopo

- Atualizar delivery de `POST /auth/login` para expor os novos resultados de
  autenticação.
- Exigir `X-Installation-Id` em `POST /auth/refresh`.
- Exigir `X-Installation-Id` nas rotas autenticadas por access token.
- Criar middleware operacional para comparar header, claim e sessão.
- Criar middleware restrito para validar `restricted_access_token`.
- Permitir `POST /security/step-up/authorize` com contexto restrito para
  `POST /security/installations`.
- Implementar handlers:
  - `POST /security/installations`;
  - `GET /security/installations`;
  - `DELETE /security/installations/{installation_resource_id}`.
- Padronizar respostas e erros HTTP, incluindo `INSTALLATION_MISMATCH`.
- Atualizar documentação REST e coleções de teste, quando existirem.

## 4. Fora de escopo

- Criar regras de domínio novas.
- Criar schema.
- Alterar limite de instalacoes.
- Adicionar biometria, attestation ou device fingerprinting.

## 5. Dependencias

- Backlog 016 para login.
- Backlog 017 para registro e gerenciamento.
- Backlog 015 para middlewares e contexto.

## 6. Preparacao para tasks

As tasks devem ser divididas por endpoint/middleware apenas depois que os use
cases estiverem estáveis. Testes devem cobrir contrato HTTP e integração de
rotas.

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
