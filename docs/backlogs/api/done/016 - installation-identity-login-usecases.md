# Backlog API 016: Installation Identity Login Use Cases

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Implementar as regras de aplicação que decidem o resultado do login conforme a
instalação apresentada.

## 3. Escopo

- Classificar a instalação no `POST /auth/login`.
- Permitir sessão operacional para instalação `known`.
- Executar bootstrap automático da primeira instalação.
- Bloquear instalação `revoked`.
- Retornar estado de limite atingido quando houver três instalações `known`.
- Emitir autorização restrita quando a instalação for nova e ainda houver vaga.
- Garantir atomicidade da primeira instalação e da decisão de limite.
- Manter senha transacional fora do login.
- Não emitir refresh token em fluxo restrito.

## 4. Fora de escopo

- Criar tabelas ou repositorios.
- Criar endpoint `POST /security/installations`.
- Criar handlers novos.
- Exigir `X-Installation-Id` em refresh ou rotas autenticadas.

## 5. Dependencias

- Backlog 014 para repositórios e operações atômicas.
- Backlog 015 para sessão operacional e token restrito.

## 6. Preparacao para tasks

As tasks devem separar classificação, bootstrap, limite, token restrito e
respostas de aplicação. Testes devem cobrir cada classificação.

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
