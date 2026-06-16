# Backlog API 013: Installation Identity Domain and Repositories

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Definir o domínio e os repositórios compartilhados que todos os fluxos de
identidade de instalação usarão.

Este backlog deve ser concluído antes de ligar classificação, bootstrap,
registro, listagem, revogação ou refresh ao fluxo real.

## 3. Escopo

- Criar tipos de domínio para instalação.
- Criar tipos de domínio para autorização restrita de registro.
- Representar estados de instalação: `known` e `revoked`.
- Representar estados de autorização: `active`, `consumed` e `revoked`.
- Implementar repositório Postgres de instalações.
- Implementar repositório Postgres de autorizações restritas.
- Mapear conflitos e violações de constraint para erros estáveis.
- Expor operações atômicas necessárias para:
  - bootstrap da primeira instalação;
  - reserva de vaga respeitando o limite de três instalações `known`;
  - consumo de autorização restrita;
  - revogação lógica de instalação.

## 4. Fora de escopo

- Implementar handlers HTTP.
- Emitir JWT operacional ou restrito.
- Alterar middleware de autenticação.
- Retornar respostas específicas do login.

## 5. Dependências

- Backlog 012: base de dados criada e validada.

## 6. Orientação para tasks

As tasks deste backlog devem nascer em torno de contratos de domínio,
repositórios e operações transacionais. A quebra deve evitar duplicar lógica de
classificação dentro dos handlers.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
