# Backlog API 012: Installation Identity Domain Contracts

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Definir o domínio e as portas internas da identidade de instalação antes de
qualquer implementação de banco, use case ou delivery.

Este backlog deve criar a linguagem estável usada pelos próximos cortes.

## 3. Escopo

- Modelar instalação de app como associação entre usuário e
  `installation_id`.
- Definir estados persistidos `known` e `revoked`.
- Definir classificações derivadas para login: `known`, `first`, `new`,
  `revoked` e `limit_reached`.
- Definir value objects e validações compartilhadas para UUID v4 canônico,
  resource id público, status e escopos.
- Definir erros internos necessários, incluindo divergência de instalação e
  limite atingido.
- Definir interfaces/portas para:
  - consulta de instalação por usuário e `installation_id`;
  - contagem de instalacoes `known`;
  - verificacao de existencia historica;
  - bootstrap atômico da primeira instalação;
  - reserva atomica de vaga;
  - listagem e revogação lógica;
  - autorizacoes restritas.

## 4. Fora de escopo

- Criar migrations.
- Implementar repositórios Postgres.
- Alterar login, refresh ou middleware.
- Criar handlers HTTP.

## 5. Dependencias

- Backlog 011 apenas para o contrato de entrada do header.

## 6. Preparacao para tasks

As tasks devem separar modelo, erros e interfaces. Testes devem focar regras
puras de domínio e contratos de porta, sem banco real.

Tasks:

- [012 - installation-identity-domain-contracts_tasks.md](<012 - installation-identity-domain-contracts_tasks.md>)

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
