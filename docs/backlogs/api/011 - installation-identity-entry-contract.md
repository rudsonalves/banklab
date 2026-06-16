# Backlog API 011: Installation Identity Entry Contract

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Estabelecer o contrato mínimo de entrada para identidade de instalação antes de
qualquer dependência de persistência.

Este backlog cobre apenas o que pode ser implementado sem tabela de
instalações, sem sessão vinculada e sem classificação operacional.

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

## 6. Orientação para tasks

As tasks deste backlog devem ser pequenas e verificáveis por contrato HTTP e
testes de handler/use case. Não devem criar dependências fictícias com
repositórios de instalação ainda inexistentes.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
