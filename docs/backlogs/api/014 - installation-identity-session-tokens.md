# Backlog API 014: Installation Identity Session and Tokens

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Preparar sessão, JWT e contexto autenticado para carregar o vínculo entre
usuário, sessão e instalação.

## 3. Escopo

- Vincular sessão operacional ao par usuário + `installation_id`.
- Emitir `access_token` operacional com claim `installation_id`.
- Preparar `POST /auth/refresh` para validar `X-Installation-Id`.
- Emitir e validar `restricted_access_token`.
- Criar middleware específico para acesso restrito.
- Carregar `user_id`, `installation_id`, `jti` e `scope` no contexto restrito.
- Preparar middleware operacional para comparar header, token e sessão.

## 4. Fora de escopo

- Classificar instalação no login.
- Registrar nova instalação.
- Consumir step-up.
- Listar ou revogar instalações.
- Aplicar bloqueio final em todas as rotas autenticadas antes dos fluxos
  dependentes estarem prontos.

## 5. Dependências

- Backlog 013: domínio e repositórios compartilhados.

## 6. Orientação para tasks

As tasks deste backlog devem separar token operacional, token restrito, sessão
persistida e contexto de middleware. O objetivo é preparar a autenticação sem
misturar regra de registro de instalação no middleware.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
