# Backlog API 017: Installation Identity Registration and Management Use Cases

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Implementar os casos de uso de registro explícito, listagem e revogação de
instalações depois que login, token restrito e repositórios estiverem prontos.

## 3. Escopo

- Autorizar registro de instalação com contexto restrito.
- Exigir step-up válido para `POST /security/installations`.
- Confirmar que o `X-Installation-Id` corresponde ao valor apresentado no
  login.
- Confirmar vaga de forma atômica.
- Criar associação `known`.
- Consumir autorização restrita.
- Criar sessão operacional vinculada após registro bem-sucedido.
- Listar instalações do usuário autenticado.
- Revogar instalação por identificador público de gerenciamento.
- Impedir revogação da instalação da sessão atual.
- Invalidar refresh tokens e access tokens da instalação revogada.

## 4. Fora de escopo

- Criar schema ou repositórios.
- Definir DTOs HTTP finais.
- Implementar painel administrativo.
- Exigir step-up para revogação neste MVP.

## 5. Dependencias

- Backlog 014 para repositórios.
- Backlog 015 para contexto e sessão.
- Backlog 016 para autorização restrita emitida no login.
- Fluxo de step-up já existente.

## 6. Preparacao para tasks

As tasks devem separar registro, listagem, revogação e invalidação de sessão.
Testes devem cobrir consumo único de autorização e efeitos da revogação.

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
