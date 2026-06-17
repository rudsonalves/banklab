# Backlog API 011: Installation Identity Entry Contract

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Estabelecer o contrato mínimo de entrada para identidade de instalação antes de
qualquer dependência de domínio persistente, tabela de instalações ou sessão
vinculada.

Este backlog é a exceção inicial da ordem em camadas porque o contrato HTTP do
header pode ser validado sem consultar estado.

## 3. Escopo

- Centralizar constantes de headers compartilhados.
- Exigir `X-Installation-Id` no `POST /auth/login`.
- Validar UUID v4 canônico.
- Retornar `400 INVALID_INSTALLATION_ID` para ausência ou formato inválido.
- Propagar o `installation_id` validado até a camada de aplicação.

## 4. Fora de escopo

- Consultar instalações cadastradas.
- Classificar instalação como `known`, `new`, `revoked` ou `limit_reached`.
- Criar bootstrap automático da primeira instalação.
- Emitir `restricted_access_token`.
- Alterar `access_token`, `refresh_token` ou sessão.
- Exigir `X-Installation-Id` em refresh ou rotas autenticadas.

## 5. Dependências

- Nenhuma dependência interna da identidade de instalação.

## 6. Preparação para tasks

As tasks futuras devem ser pequenas e verificáveis por contrato HTTP e testes de
handler/use case. Elas não devem criar repositórios fictícios nem antecipar a
classificação operacional.

Tasks:

- [011 - installation-identity-entry-contract_tasks.md](<011 - installation-identity-entry-contract_tasks.md>)

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
