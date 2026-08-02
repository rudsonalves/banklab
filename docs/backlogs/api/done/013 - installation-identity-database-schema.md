# Backlog API 013: Installation Identity Database Schema

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Criar a base relacional necessária para instalações e autorizações restritas
antes de implementar repositórios ou casos de uso.

## 3. Escopo

- Criar tabela `app_installations`.
- Criar tabela `installation_registration_authorizations`.
- Definir identificador público de gerenciamento separado do
  `installation_id` enviado pelo cliente.
- Persistir estados `known` e `revoked` para instalações.
- Persistir estados `active`, `consumed` e `revoked` para autorizações
  restritas, tratando expiração como derivada de `expires_at`.
- Criar constraints de unicidade e integridade.
- Criar índices para consultas por usuário, instalação, status, `jti`,
  expiração e identificador público.
- Preparar o limite de três instalações `known` por usuário de forma que os
  repositórios possam garanti-lo atomicamente.
- Criar migrations `up` e `down`.

## 4. Fora de escopo

- Implementar repositórios.
- Implementar use cases.
- Ligar login, refresh ou delivery ao novo schema.

## 5. Dependencias

- Backlog 012 com domínio e portas estabilizados.

## 6. Preparacao para tasks

As tasks futuras devem separar tabela de instalações, tabela de autorizações,
índices/constraints e testes ou validação de migration.

Tasks:

- [013 - installation-identity-database-schema_tasks.md](<013 - installation-identity-database-schema_tasks.md>)

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
