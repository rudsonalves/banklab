# Backlog API 014: Installation Identity Repositories

## 1. Status

- Tipo: Planning
- Area: Security
- Prioridade: High
- Estado: Discussao

## 2. Objetivo

Implementar as portas de instalação e autorização restrita usando Postgres,
incluindo as operações atômicas que sustentam os use cases posteriores.

## 3. Escopo

- Implementar repositório de instalações.
- Implementar repositório de autorizações restritas.
- Implementar consultas para classificação de login.
- Implementar bootstrap atômico da primeira instalação.
- Implementar reserva atômica de vaga respeitando o limite de três
  instalações `known`.
- Implementar consumo e revogação de autorização restrita.
- Implementar revogação lógica de instalação.
- Garantir que concorrência não crie duas primeiras instalações nem ultrapasse
  o limite de instalações `known`.
- Cobrir erros de conflito, ausência, estado inválido e falhas de transação.

## 4. Fora de escopo

- Alterar contrato HTTP.
- Emitir tokens.
- Ligar repositórios ao login ou a handlers.

## 5. Dependencias

- Backlog 012 para portas e domínio.
- Backlog 013 para schema e constraints.

## 6. Preparacao para tasks

As tasks devem ser divididas por repositório e por operação atômica. Testes
devem cobrir os caminhos de concorrência mais importantes.

## 7. Referencias

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
- [Split por dependencia](<010 - split-installation-identity-by-dependency.md>)
