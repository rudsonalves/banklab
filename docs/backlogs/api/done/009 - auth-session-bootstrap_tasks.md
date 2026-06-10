# Tasks do endpoint de sessão pós-login

> Decisão superveniente: `can_access_home` foi removido do contrato de
> `GET /auth/session`. As menções abaixo registram a implementação histórica;
> o contrato vigente expõe somente os estados objetivos de readiness.

Backlog pai:

- `009 - auth-session-bootstrap.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Contrato/Sessão/Bootstrap

## Task 1/7: Definir DTO de resposta de sessão

Status: Concluída

### Objetivo

Definir o contrato de saída de `GET /auth/session` seguindo o envelope
padronizado da aplicação.

### Escopo

- Criar DTOs de resposta no delivery de auth, ou em local equivalente:
  - `sessionData`;
  - `sessionUserData`;
  - `sessionCustomerData`;
  - `sessionReadinessData`.
- Incluir em `user`:
  - `id`;
  - `email`;
  - `phone`;
  - `role`.
- Incluir em `customer`:
  - `id`;
  - `name`;
  - `cpf`;
  - `birth_date`;
  - `created_at`.
- Incluir em `readiness`:
  - `onboarding_completed`;
  - `approved`;
  - `has_operational_account`;
  - `transaction_password_status`;
  - `can_access_home`.
- Não incluir `user.customer_id`.
- Não incluir `customer.email`.
- Não incluir `next_required_step` no primeiro corte.

### Critérios de aceite

- O formato de resposta segue `{ data, error }`.
- A resposta contém `user`, `customer` e `readiness`.
- Campos redundantes não aparecem no contrato.
- O contrato contém todos os campos necessários para o mobile substituir a
  composição atual de `/auth/me` + `/customers/me`.

### Depende de

- Nenhuma dependência.

## Task 2/7: Criar use case de sessão autenticada

Status: Concluída

### Objetivo

Criar a orquestração de aplicação responsável por montar o snapshot de sessão
pós-login.

### Escopo

- Criar use case no módulo de auth, ou módulo transversal apropriado.
- Ler o usuário autenticado a partir do contexto.
- Buscar dados completos do usuário pelo repositório de auth.
- Para usuário customer, buscar o customer vinculado.
- Retornar erro de estado inválido quando usuário customer não possuir customer
  associado.
- Retornar erro apropriado quando o customer vinculado não for encontrado.
- Montar saída de aplicação com:
  - dados do usuário;
  - dados do customer;
  - readiness inicial.

### Critérios de aceite

- Use case retorna sessão para usuário autenticado válido.
- Use case falha para sessão ausente ou inválida.
- Use case falha para usuário customer sem customer associado.
- Use case falha para customer inexistente quando o usuário customer exigir
  vínculo cadastral.
- Use case não depende de HTTP delivery.

### Depende de

- Task 1.

## Task 3/7: Consultar senha transacional e calcular readiness

Status: Concluída

### Objetivo

Adicionar ao snapshot de sessão os sinais de readiness necessários para o
bootstrap pós-login.

### Escopo

- Consultar o repositório de senha transacional pelo usuário autenticado.
- Mapear o estado da senha transacional para:
  - `active`;
  - `not_set`;
  - `locked`;
  - `unknown`.
- Expor `onboarding_completed = true` no primeiro corte.
- Expor `approved` conforme o estado atual de aprovação usado pela API.
- Expor `has_operational_account` conforme contas operacionais disponíveis para
  o customer.
- Calcular `can_access_home` na API.
- Não retornar senha transacional, hash, pepper ou qualquer material sensível.

### Critérios de aceite

- Usuário com senha transacional ativa retorna
  `transaction_password_status = active`.
- Usuário sem senha transacional retorna
  `transaction_password_status = not_set`.
- `onboarding_completed` retorna `true`.
- `approved` é preenchido.
- `has_operational_account` é preenchido.
- `can_access_home` é `false` quando a senha transacional obrigatória não
  existir.
- `can_access_home` é calculado pela API, não pelo cliente.

### Depende de

- Task 2.

## Task 4/7: Implementar handler e rota `GET /auth/session`

Status: Concluída

### Objetivo

Expor o novo contrato HTTP autenticado.

### Escopo

- Adicionar método `Session` no handler de auth, ou handler equivalente.
- Registrar `GET /auth/session` no router de auth.
- Proteger a rota com middleware JWT.
- Converter a saída do use case para DTO de resposta.
- Usar `sharedhttp.WriteJSON` para respeitar o envelope padronizado.
- Mapear erros pelo registry/padrão existente.

### Critérios de aceite

- `GET /auth/session` exige JWT válido.
- Sessão válida retorna HTTP 200.
- Sessão ausente ou inválida retorna erro padronizado.
- Resposta não inclui dados sensíveis.
- `/auth/me` e `/customers/me` permanecem inalterados.

### Depende de

- Task 3.

## Task 5/7: Cobrir testes de application e delivery

Status: Concluída

### Objetivo

Garantir cobertura automatizada do contrato e das regras de sessão.

### Escopo

- Adicionar testes do use case para:
  - usuário com senha transacional ativa;
  - usuário sem senha transacional;
  - usuário customer sem customer associado;
  - customer vinculado não encontrado;
  - sessão ausente ou inválida.
- Adicionar testes do handler para:
  - HTTP 200 com envelope esperado;
  - rota exigindo autenticação;
  - erro padronizado;
  - ausência de campos redundantes.
- Garantir que `/auth/me` e `/customers/me` continuem passando nos testes
  existentes.

### Critérios de aceite

- Testes cobrem os cenários principais do backlog.
- Testes validam `user.phone`.
- Testes validam ausência de `user.customer_id`.
- Testes validam ausência de `customer.email`.
- Testes validam `readiness.transaction_password_status`.
- Testes validam `readiness.can_access_home`.

### Depende de

- Task 4.

## Task 6/7: Atualizar documentação REST e referências mobile

Status: Concluída

### Objetivo

Documentar o novo endpoint e sua relação com o bootstrap pós-login do mobile.

### Escopo

- Atualizar `api/docs/07-api-rest.md`.
- Documentar:
  - método e path;
  - autenticação;
  - resposta de sucesso;
  - erros relevantes;
  - campos de readiness;
  - ausência de dados sensíveis.
- Atualizar referências em backlogs mobile quando necessário.
- Atualizar collection Postman se ela for mantida como artefato de contrato.

### Critérios de aceite

- Documentação descreve `GET /auth/session`.
- Documentação deixa claro que o endpoint substitui a composição pós-login de
  `/auth/me` + `/customers/me` no mobile.
- Documentação mantém `/auth/me` e `/customers/me` como endpoints existentes.
- Documentação usa envelope padronizado.

### Depende de

- Task 4.

## Task 7/7: Validar alinhamento final da API

Status: Concluída

### Objetivo

Fechar a implementação garantindo alinhamento entre contrato, código, testes e
documentação.

### Escopo

- Rodar `go test ./...` na API.
- Revisar:
  - rota;
  - handler;
  - use case;
  - DTOs;
  - readiness;
  - documentação REST;
  - backlogs mobile relacionados.
- Confirmar que o mobile pode migrar `AuthApi.getProfile` para uma única
  chamada.

### Critérios de aceite

- `go test ./...` passa.
- `GET /auth/session` retorna todos os campos necessários ao bootstrap
  pós-login.
- `/auth/me` e `/customers/me` continuam disponíveis.
- Não há divergência conhecida entre backlog, documentação e implementação.

### Depende de

- Task 5.
- Task 6.

### Registro de alinhamento final

Data de fechamento: 2026-06-02

Verificações executadas:

- `go test ./...` executado na API com sucesso.
- `GET /auth/session` registrado com autenticação JWT.
- Handler `Session` retorna envelope padronizado `{ data, error }`.
- Resposta inclui:
  - `user.id`;
  - `user.email`;
  - `user.phone`;
  - `user.role`;
  - `customer.id`;
  - `customer.name`;
  - `customer.cpf`;
  - `customer.birth_date`;
  - `customer.created_at`;
  - `readiness.onboarding_completed`;
  - `readiness.approved`;
  - `readiness.has_operational_account`;
  - `readiness.transaction_password_status`;
  - `readiness.can_access_home`.
- Resposta não inclui:
  - `user.customer_id`;
  - `customer.email`;
  - `readiness.next_required_step`.
- `/auth/me` e `/customers/me` permanecem disponíveis.
- Documentação REST atualizada.
- Collection Postman atualizada.
- Backlog mobile de cadastro de senha transacional alinhado para usar
  `GET /auth/session` no gate pós-login.

Conclusão:

- Implementação da API está alinhada com contrato, testes e documentação.
- O mobile pode migrar o bootstrap pós-login para `GET /auth/session` em etapa
  posterior.
