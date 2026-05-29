# Tasks dos Contratos e Documentação do ZTA MVP

Backlog pai:

- `006d - zta-contracts-and-docs.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA/Documentação

## Task 1/6: Consolidar padrão de resposta e nomes públicos

Status: Done

### Objetivo

Garantir que os nomes públicos do MVP ZTA estejam documentados de forma única e
consistente para API, mobile e documentação.

### Escopo

- Confirmar o envelope padrão:
  - sucesso com `data` preenchido e `error: null`;
  - erro com `data: null` e `error.code`, `error.message`, `error.details`.
- Documentar que clientes devem depender de `error.code`, não de
  `error.message`.
- Consolidar nomes públicos:
  - `POST /security/transaction-password`;
  - `POST /security/step-up/authorize`;
  - `POST /accounts/internal-transfers`;
  - `X-Step-Up-Token`;
  - `internal_transfer.create`.
- Consolidar campos JSON:
  - `endpoint_key`;
  - `transaction_password`;
  - `transaction_password_confirmation`;
  - `step_up_token`;
  - `expires_in`.

### Critérios de aceite

- Os nomes públicos aparecem iguais em todos os documentos revisados.
- O envelope padrão não diverge de `api/docs/05-error_and_response.md`.
- Não há dependência documentada em mensagens textuais de erro.

### Depende de

- `006a - transaction-password`.
- `006b - step-up-token`.
- `006c - internal-transfer-step-up-enforcement`.

## Task 2/6: Consolidar contrato do JWT de step-up

Status: Done

### Objetivo

Documentar o contrato estável do step-up token emitido pelo backend.

### Escopo

- Documentar validade de 120 segundos.
- Documentar persistência do `jti` para consumo único.
- Documentar claims mínimos:
  - `user_id`;
  - `endpoint_key`;
  - `scope=step_up`;
  - `exp`;
  - `iat`;
  - `jti`.
- Documentar que o token não contém:
  - senha transacional;
  - hash;
  - dados sensíveis;
  - payload da operação.
- Confirmar que o token autoriza uma chamada ao endpoint lógico e não o payload
  detalhado da transferência no MVP.

### Critérios de aceite

- O contrato documentado corresponde ao signer e verifier implementados.
- A documentação deixa claro que o `jti` persistido é parte da regra de uso
  único.
- Não há orientação para clientes inspecionarem ou dependerem do payload do JWT.

### Depende de

- `006b - step-up-token`.
- `006c - internal-transfer-step-up-enforcement`.

## Task 3/6: Consolidar contrato de erros ZTA

Status: Done

### Objetivo

Garantir que os códigos de erro e status HTTP do MVP ZTA estejam alinhados com
o contrato compartilhado da API.

### Escopo

- Confirmar os erros de senha transacional:
  - `TRANSACTION_PASSWORD_ALREADY_SET`: 409;
  - `TRANSACTION_PASSWORD_REQUIRED`: 403;
  - `TRANSACTION_PASSWORD_NOT_SET`: 409;
  - `TRANSACTION_PASSWORD_INVALID`: 401;
  - `TRANSACTION_PASSWORD_LOCKED`: 403.
- Confirmar erro de autorização de emissão:
  - `STEP_UP_ENDPOINT_NOT_ALLOWED`: 403.
- Confirmar erros de enforcement:
  - `STEP_UP_TOKEN_REQUIRED`: 401;
  - `STEP_UP_TOKEN_INVALID`: 401;
  - `STEP_UP_TOKEN_EXPIRED`: 401;
  - `STEP_UP_TOKEN_CONSUMED`: 401;
  - `STEP_UP_ENDPOINT_MISMATCH`: 403.
- Documentar a separação entre:
  - autorização/emissão de step-up;
  - enforcement do endpoint sensível.
- Documentar que `STEP_UP_TOKEN_INVALID` cobre token malformado, assinatura
  inválida, `scope` diferente de `step_up`, `jti` inexistente ou claims
  obrigatórios ausentes.

### Critérios de aceite

- Todos os códigos existem no contrato compartilhado da API.
- Todos os mapeamentos HTTP estão registrados e testados.
- `STEP_UP_ENDPOINT_NOT_ALLOWED` não é confundido com
  `STEP_UP_ENDPOINT_MISMATCH`.
- Mobile pode depender de `error.code` para todos os cenários listados.

### Depende de

- `006a - transaction-password`.
- `006b - step-up-token`.
- `006c - internal-transfer-step-up-enforcement`.

## Task 4/6: Atualizar documentação REST do fluxo ZTA

Status: Done

### Objetivo

Manter a documentação REST alinhada com os endpoints, payloads, headers,
respostas e erros implementados no MVP ZTA.

### Escopo

- Atualizar `api/docs/07-api-rest.md`.
- Confirmar documentação de:
  - criação de senha transacional;
  - autorização de step-up;
  - transferência interna protegida.
- Documentar `X-Step-Up-Token` como obrigatório em
  `POST /accounts/internal-transfers`.
- Documentar que o endpoint lógico esperado é `internal_transfer.create`.
- Documentar consumo atômico antes da execução da transferência.
- Documentar comportamento de retry:
  - mesmo `X-Step-Up-Token` após consumo retorna `STEP_UP_TOKEN_CONSUMED`;
  - mesmo `idempotency_key` pode exigir novo step-up token.
- Confirmar exemplos de erro para os cenários de enforcement.

### Critérios de aceite

- REST docs mostram header obrigatório e erros esperados.
- REST docs não sugerem uso de branch/account-number para o payload novo de
  transferência interna.
- REST docs deixam claro que o step-up token é single-use.
- Documentação REST não diverge dos testes de delivery e application.

### Depende de

- Task 1.
- Task 2.
- Task 3.

## Task 5/6: Atualizar documentação ZTA, READMEs e referências de consumo

Status: Done

### Objetivo

Garantir que a documentação de arquitetura/implementação e as referências
rápidas de consumo reflitam o fluxo final do MVP ZTA.

### Escopo

- Atualizar `api/docs/implementations/03-zta-step-up-transaction-password.md`.
- Atualizar `api/docs/05-error_and_response.md`, se houver divergência de
  códigos ou status.
- Atualizar `README.md`, `README_en.md` e `api/README.md` com a exigência de
  `X-Step-Up-Token` na transferência interna.
- Avaliar documentação mobile quando houver impacto direto no consumo.
- Não incluir coleção Postman como requisito deste backlog.

### Critérios de aceite

- Docs ZTA explicam consumo único, retry e relação com idempotência.
- READMEs deixam claro que `POST /accounts/internal-transfers` exige
  `X-Step-Up-Token`.
- Documentos de erro listam todos os códigos ZTA do MVP.
- Não há menção a Postman como pendência do backlog `006d`.

### Depende de

- Task 1.
- Task 2.
- Task 3.
- Task 4.

## Task 6/6: Verificar alinhamento final do contrato ZTA

Status: Done

### Objetivo

Fechar o backlog com uma verificação cruzada entre implementação, testes e
documentação.

### Escopo

- Conferir que os contratos documentados batem com:
  - constantes de erro compartilhadas;
  - registry de erros do módulo `security`;
  - signer/verifier JWT;
  - `EnforceStepUpUseCase`;
  - handler de transferência interna;
  - documentação REST e ZTA.
- Rodar os testes afetados da API.
- Registrar no backlog o status de alinhamento do MVP enforcement:
  - endpoint protegido;
  - header obrigatório;
  - endpoint lógico;
  - ponto arquitetural;
  - consumo atômico;
  - retry e idempotência.

### Critérios de aceite

- `go test ./...` passa.
- Não há divergência conhecida entre código e documentação.
- O backlog `006d` pode ser considerado fechado ou movido para done.
- Pendências futuras ficam explicitamente fora do escopo do MVP.

### Depende de

- Task 4.
- Task 5.
