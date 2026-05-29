# Tasks do Enforcement de Step-Up na Transferência Interna

Backlog pai:

- `006c - internal-transfer-step-up-enforcement.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/7: Definir contrato de validação do step-up token

Status: Backlog

### Objetivo

Criar as portas e tipos necessários para validar o JWT de step-up no enforcement
sem acoplar application à implementação JWT concreta.

### Escopo

- Criar contrato `StepUpTokenVerifier` no módulo de segurança.
- Definir modelo de claims validados contendo:
  - `user_id`;
  - `endpoint_key`;
  - `scope`;
  - `jti`;
  - `expires_at`;
  - `issued_at`.
- Prever erro para token ausente no nível de enforcement:
  - `STEP_UP_TOKEN_REQUIRED`.
- Reutilizar erros já definidos para token inválido, expirado e consumido
  quando fizer sentido.
- Garantir que o contrato permita testes com fake/mock sem JWT real.
- Não alterar o signer de emissão exceto se houver necessidade de compartilhar
  tipos internos sem vazar detalhes de infraestrutura.

### Critérios de aceite

- Application consegue depender de uma interface para validar step-up token.
- Claims obrigatórios têm representação tipada.
- O contrato exige `scope=step_up`.
- O contrato não depende de HTTP, delivery ou Postgres.
- Testes unitários conseguem simular token válido, inválido, expirado e ausente.

### Depende de

- `006b - step-up-token`.

## Task 2/7: Implementar verifier JWT do step-up token

Status: Backlog

### Objetivo

Validar e extrair claims do JWT de step-up emitido no backlog `006b`.

### Escopo

- Criar implementação em `internal/security/infrastructure`.
- Validar assinatura do JWT com o segredo/configuração usada pelo signer.
- Validar algoritmo esperado.
- Validar expiração do JWT.
- Exigir claims obrigatórios:
  - `user_id`;
  - `endpoint_key`;
  - `scope=step_up`;
  - `exp`;
  - `iat`;
  - `jti`.
- Rejeitar token malformado, assinatura inválida, claims ausentes ou `scope`
  diferente de `step_up`.
- Mapear falhas públicas para erro de token inválido.
- Não validar consumo nem existência do `jti` nesta camada; isso fica no
  enforcement application/repositório.

### Critérios de aceite

- Token emitido pelo signer de `006b` é aceito pelo verifier.
- Token com segredo diferente é rejeitado.
- Token com `scope` diferente de `step_up` é rejeitado.
- Token sem `jti`, `user_id` ou `endpoint_key` é rejeitado.
- Token expirado é rejeitado de forma controlada.
- Testes cobrem assinatura, claims obrigatórios, expiração e algoritmo inválido.

### Depende de

- Task 1.

## Task 3/7: Implementar caso de uso de enforcement de step-up

Status: Backlog

### Objetivo

Criar `EnforceStepUpUseCase` para validar o step-up token antes da transferência
interna e consumir o `jti` de forma atômica.

### Escopo

- Criar `EnforceStepUpUseCase` em `internal/security/application`.
- Receber:
  - usuário autenticado;
  - `endpoint_key` esperado;
  - step-up token bruto;
  - horário atual.
- Retornar `STEP_UP_TOKEN_REQUIRED` se o token bruto estiver ausente.
- Validar JWT usando `StepUpTokenVerifier`.
- Buscar/consumir `jti` usando `StepUpTokenRepository.ConsumeByJTI`.
- Comparar `user_id` do JWT com usuário autenticado.
- Comparar `user_id` do JWT com registro persistido.
- Comparar `endpoint_key` do JWT com registro persistido.
- Comparar `endpoint_key` do JWT com endpoint esperado.
- Validar expiração do JWT e expiração derivada do registro persistido.
- Consumir o `jti` antes de liberar a chamada ao use case sensível.
- Não chamar nem conhecer o use case de transferência.

### Critérios de aceite

- Token válido para `internal_transfer.create` libera o enforcement.
- Header ausente retorna `STEP_UP_TOKEN_REQUIRED`.
- Token malformado, assinatura inválida, `scope` errado, claims ausentes ou
  `jti` inexistente retornam `STEP_UP_TOKEN_INVALID`.
- Token expirado retorna `STEP_UP_TOKEN_EXPIRED`.
- Token já consumido retorna `STEP_UP_TOKEN_CONSUMED`.
- Token emitido para outro endpoint retorna `STEP_UP_ENDPOINT_MISMATCH`.
- Token de outro usuário retorna erro estável e não libera a transferência.
- Divergência entre claims do JWT e registro persistido retorna
  `STEP_UP_TOKEN_INVALID`.
- Consumo do `jti` acontece de forma atômica antes do retorno de sucesso.

### Depende de

- Task 1.
- Task 2.
- `006b - step-up-token`.

## Task 4/7: Registrar erros de enforcement do step-up

Status: Backlog

### Objetivo

Adicionar os erros específicos do enforcement ao contrato compartilhado de
erros da API.

### Escopo

- Definir erros de domínio/application necessários:
  - `TRANSACTION_PASSWORD_REQUIRED`;
  - `STEP_UP_TOKEN_REQUIRED`;
  - `STEP_UP_ENDPOINT_MISMATCH`.
- Confirmar reutilização de:
  - `STEP_UP_TOKEN_INVALID`;
  - `STEP_UP_TOKEN_EXPIRED`;
  - `STEP_UP_TOKEN_CONSUMED`.
- Adicionar códigos compartilhados quando ainda não existirem.
- Registrar mapeamentos do módulo `security`.
- Mapear HTTP conforme `006d`:
  - `TRANSACTION_PASSWORD_REQUIRED`: 403;
  - `STEP_UP_TOKEN_REQUIRED`: 401;
  - `STEP_UP_TOKEN_INVALID`: 401;
  - `STEP_UP_TOKEN_EXPIRED`: 401;
  - `STEP_UP_TOKEN_CONSUMED`: 401;
  - `STEP_UP_ENDPOINT_MISMATCH`: 403.
- Manter `STEP_UP_ENDPOINT_NOT_ALLOWED` separado da autorização de step-up.

### Critérios de aceite

- Todos os códigos de enforcement estão registrados e testados.
- Mobile pode depender de `error.code`.
- Mapeamentos usam o envelope padrão de erro.
- Nenhum erro de enforcement é confundido com erro de autorização de emissão.

### Depende de

- Task 3.

## Task 5/7: Integrar enforcement no handler de transferência interna

Status: Backlog

### Objetivo

Proteger `POST /accounts/internal-transfers` exigindo o header
`X-Step-Up-Token` antes de chamar o use case de transferência.

### Escopo

- Injetar dependência de enforcement no handler de transferência.
- No handler de `Transfer`, extrair `X-Step-Up-Token`.
- Definir endpoint lógico esperado como `internal_transfer.create`.
- Chamar `EnforceStepUpUseCase` antes de `Transfer.Execute`.
- Não passar step-up token, senha transacional ou detalhes de policy para o use
  case de transferência.
- Se o enforcement falhar, responder erro e não chamar a transferência.
- Se o enforcement passar, chamar a transferência normalmente.
- Conectar wiring no startup da API.

### Critérios de aceite

- Transferência sem `X-Step-Up-Token` não chama o use case de transferência.
- Transferência com token inválido não chama o use case de transferência.
- Transferência com token válido chama o use case de transferência.
- O endpoint continua exigindo JWT normal via `Authorization: Bearer`.
- Nenhum log contém token completo ou material sensível.
- O use case de transferência permanece sem dependência de security/ZTA.

### Depende de

- Task 3.
- Task 4.

## Task 6/7: Cobrir enforcement com testes

Status: Backlog

### Objetivo

Garantir cobertura automatizada do enforcement antes de considerar a
transferência interna protegida.

### Escopo

- Adicionar testes do verifier JWT para:
  - token válido;
  - assinatura inválida;
  - token expirado;
  - `scope` inválido;
  - claims ausentes;
  - algoritmo inválido.
- Adicionar testes do `EnforceStepUpUseCase` para:
  - sucesso;
  - token ausente;
  - token inválido;
  - token expirado;
  - token consumido;
  - mismatch de endpoint;
  - mismatch de usuário;
  - divergência entre JWT e registro persistido;
  - consumo atômico antes do sucesso.
- Adicionar testes de delivery para:
  - ausência de `X-Step-Up-Token`;
  - enforcement falhando;
  - enforcement passando;
  - garantia de que transferência não é chamada quando enforcement falha.
- Adicionar ou ajustar testes de integração quando necessário para cobrir o
  fluxo com repositório real.

### Critérios de aceite

- Regras principais de enforcement têm cobertura automatizada.
- Testes provam que o use case de transferência não roda sem step-up válido.
- Testes provam que `jti` não pode ser reutilizado.
- Testes provam que retry com mesmo token consumido retorna
  `STEP_UP_TOKEN_CONSUMED`.
- Testes afetados da API passam.

### Depende de

- Task 2.
- Task 3.
- Task 5.

## Task 7/7: Atualizar documentação do enforcement

Status: Backlog

### Objetivo

Manter a documentação REST e ZTA alinhada com o enforcement implementado.

### Escopo

- Atualizar `api/docs/07-api-rest.md` se o comportamento final divergir do
  contrato atual.
- Atualizar `api/docs/implementations/03-zta-step-up-transaction-password.md`.
- Confirmar contrato de erros em `006d - zta-contracts-and-docs.md`.
- Documentar que o step-up token é de uso único.
- Documentar que retry com mesmo `X-Step-Up-Token` após consumo retorna
  `STEP_UP_TOKEN_CONSUMED`.
- Documentar que um novo step-up token pode ser necessário mesmo quando o
  `idempotency_key` da transferência é o mesmo.
- Confirmar que o endpoint protegido continua sendo:
  - `POST /accounts/internal-transfers`;
  - header `X-Step-Up-Token`;
  - endpoint lógico `internal_transfer.create`.

### Critérios de aceite

- Documentação REST mostra header obrigatório e erros esperados.
- Documentação ZTA explica consumo único e retry/idempotência.
- `006d` continua alinhado com os códigos implementados.
- Não há divergência entre docs e comportamento testado.

### Depende de

- Task 5.
- Task 6.
