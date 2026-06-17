# Backlog API 010: Split por Dependência da Identidade de Instalação

## 1. Objetivo

Separar o MVP de identidade de instalação em backlogs menores seguindo a ordem
mais natural de construção da implementação.

Esta separação substitui a divisão por fluxo ou endpoint. A implementação deve
avançar por dependência técnica:

1. contrato mínimo que não depende de persistência;
2. domínio e portas;
3. banco de dados;
4. repositórios e operações atômicas;
5. sessão, tokens e contexto autenticado;
6. casos de uso;
7. delivery HTTP, middleware e documentação pública;
8. auditoria, retenção e acabamento operacional.

## 2. Princípio de separação

Cada backlog deve deixar o próximo tecnicamente possível sem criar contrato
fictício.

Regras:

- domínio não deve depender de Postgres, HTTP ou JWT;
- interfaces devem nascer antes das implementações concretas;
- migrations devem existir antes dos repositórios Postgres;
- casos de uso devem depender de interfaces estáveis, não de detalhes de
  delivery;
- delivery deve ser a última camada a conectar o fluxo ao mundo externo;
- classificação, bootstrap, registro, refresh e enforcement não devem ser
  ligados antes da base compartilhada existir.

## 3. Exceção permitida: contrato mínimo de entrada

O único corte que pode acontecer antes da base de domínio/persistência é o
contrato mínimo de `X-Installation-Id` no login:

- constante compartilhada de header;
- validação de UUID v4 canônico;
- erro `INVALID_INSTALLATION_ID`;
- propagação do valor validado até a camada de aplicação.

Esse corte não consulta instalações, não classifica estado operacional, não
cria associação, não altera sessão e não emite token restrito.

## 4. Backlogs resultantes

### 011 - Entry Contract

Contrato mínimo de entrada, sem dependência de persistência.

Arquivo:

- [011 - installation-identity-entry-contract.md](<011 - installation-identity-entry-contract.md>)

### 012 - Domain Contracts

Entidades, estados, value objects, erros e portas de domínio/aplicação.

Arquivo:

- [012 - installation-identity-domain-contracts.md](<012 - installation-identity-domain-contracts.md>)

### 013 - Database Schema

Migrations, tabelas, constraints, índices e modelo relacional mínimo.

Arquivo:

- [013 - installation-identity-database-schema.md](<013 - installation-identity-database-schema.md>)

### 014 - Repositories

Implementações Postgres das portas e operações atômicas compartilhadas.

Arquivo:

- [014 - installation-identity-repositories.md](<014 - installation-identity-repositories.md>)

### 015 - Session, Tokens and Context

Vínculo de sessão, claims, token restrito e contexto autenticado/restrito.

Arquivo:

- [015 - installation-identity-session-tokens-context.md](<015 - installation-identity-session-tokens-context.md>)

### 016 - Login and Restricted Authorization Use Cases

Classificação da instalação no login, bootstrap da primeira instalação,
limite, bloqueios e emissão de autorização restrita.

Arquivo:

- [016 - installation-identity-login-usecases.md](<016 - installation-identity-login-usecases.md>)

### 017 - Registration and Management Use Cases

Registro explícito, listagem, revogação e efeitos sobre sessões.

Arquivo:

- [017 - installation-identity-management-usecases.md](<017 - installation-identity-management-usecases.md>)

### 018 - Delivery and Enforcement

Handlers, DTOs, rotas, middlewares operacionais/restritos e contrato REST.

Arquivo:

- [018 - installation-identity-delivery-enforcement.md](<018 - installation-identity-delivery-enforcement.md>)

### 019 - Audit, Retention and Operational Docs

Retenção, auditoria, minimização, documentação final e validação operacional.

Arquivo:

- [019 - installation-identity-audit-retention.md](<019 - installation-identity-audit-retention.md>)

## 5. Sequência esperada

```text
011 entry contract
  -> 012 domain contracts
  -> 013 database schema
  -> 014 repositories
  -> 015 session/tokens/context
  -> 016 login use cases
  -> 017 registration/management use cases
  -> 018 delivery/enforcement
  -> 019 audit/retention/docs
```

Backlogs posteriores podem ser refinados em tasks menores somente depois que o
backlog anterior tiver seus contratos principais estabilizados.
