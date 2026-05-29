# Backlog: enforcement na transferência interna

## 1. Objetivo

Proteger a transferência interna com o step-up token emitido pelo módulo de
segurança.

No primeiro corte do MVP, apenas a transferência interna exigirá senha
transacional.

## 2. Endpoint protegido

```http
POST /accounts/internal-transfers
Authorization: Bearer <access_token>
X-Step-Up-Token: <step_up_token>
```

Endpoint lógico:

```text
internal_transfer.create
```

## 3. Regras de enforcement

Antes de chamar o use case de transferência, o backend deve validar:

- `X-Step-Up-Token` está presente;
- o token é válido e assinado;
- o token contém `scope=step_up`;
- o token não expirou;
- o `jti` existe no backend;
- o `jti` ainda não foi consumido;
- o usuário do JWT normal é o mesmo usuário do step-up token;
- o `endpoint_key` do token é `internal_transfer.create`.

Se a validação passar:

- o `jti` deve ser consumido de forma atômica;
- o use case de transferência interna pode ser chamado.

Se a validação falhar:

- o use case de transferência interna não deve ser chamado;
- a API deve responder usando o envelope padrão de erro.

## 4. Posicionamento arquitetural

No MVP, a delivery do módulo de transferência chama explicitamente o enforcement
do módulo de segurança antes do use case sensível.

```text
internal/account/transaction/delivery
  -> internal/security/application
  -> internal/account/transaction/application
```

O use case de transferência não deve receber senha transacional.

## 5. Fluxo

```text
1. Mobile chama POST /accounts/internal-transfers com JWT e X-Step-Up-Token.
2. Delivery de transferência identifica endpoint_key=internal_transfer.create.
3. Delivery chama o enforcement de security/application.
4. Security valida assinatura, `scope`, usuário, endpoint, expiração e consumo.
5. Security consome o `jti` persistido de forma atômica.
6. Delivery chama o use case de transferência.
```

## 6. Erros cobertos

- `TRANSACTION_PASSWORD_REQUIRED`
- `STEP_UP_TOKEN_REQUIRED`
- `STEP_UP_TOKEN_INVALID`
- `STEP_UP_TOKEN_EXPIRED`
- `STEP_UP_TOKEN_CONSUMED`
- `STEP_UP_ENDPOINT_MISMATCH`

Separação inicial:

- `TRANSACTION_PASSWORD_REQUIRED`: o endpoint exige step-up e a requisição não
  apresenta autorização suficiente para seguir.
- `STEP_UP_TOKEN_REQUIRED`: o header `X-Step-Up-Token` obrigatório não foi
  enviado no request protegido.
- `STEP_UP_TOKEN_INVALID`: token malformado, assinatura inválida, `scope`
  diferente de `step_up`, `jti` inexistente ou claims obrigatórios ausentes.
- `STEP_UP_TOKEN_EXPIRED`: token válido e persistido, mas expirado pela regra
  `status=active` com `expires_at < now`.
- `STEP_UP_TOKEN_CONSUMED`: `jti` já consumido anteriormente.
- `STEP_UP_ENDPOINT_MISMATCH`: token válido, mas emitido para outro
  `endpoint_key`.

O contrato completo de erros fica no backlog
[006d - zta-contracts-and-docs.md](<006d - zta-contracts-and-docs.md>).

## 7. Propostas para fechar antes das tasks

### 7.1 Contrato de validação do JWT de step-up

Criar uma porta explícita no módulo de segurança para validar e extrair os
claims do JWT de step-up, separada da assinatura criada no backlog `006b`.

Nome sugerido:

```text
StepUpTokenVerifier
```

Responsabilidade:

- validar assinatura;
- validar expiração do JWT;
- exigir `scope=step_up`;
- extrair `user_id`;
- extrair `endpoint_key`;
- extrair `jti`;
- rejeitar claims obrigatórios ausentes ou malformados.

Resultado esperado da validação:

```text
user_id
endpoint_key
jti
expires_at
issued_at
```

### 7.2 Comparação entre JWT e registro persistido

Após validar o JWT, o enforcement deve buscar o `jti` persistido e comparar:

- `user_id` do JWT com `user_id` persistido;
- `endpoint_key` do JWT com `endpoint_key` persistido;
- `endpoint_key` esperado pelo endpoint protegido;
- expiração do JWT;
- expiração derivada do registro persistido.

Proposta de erro:

- divergência entre JWT e registro persistido: `STEP_UP_TOKEN_INVALID`;
- token emitido para outro endpoint lógico: `STEP_UP_ENDPOINT_MISMATCH`;
- `jti` inexistente: `STEP_UP_TOKEN_INVALID`.

### 7.3 Momento do consumo do `jti`

Proposta: consumir o `jti` de forma atômica antes de chamar o use case de
transferência.

Consequência aceita:

- se a transferência falhar depois do enforcement por dados inválidos, saldo
  insuficiente, conta inativa ou outra regra de negócio, o step-up token já terá
  sido consumido;
- o cliente precisará obter novo step-up token para uma nova tentativa.

Justificativa:

- mantém a garantia de uso único forte;
- evita reutilização do mesmo step-up token em múltiplas tentativas de operação;
- simplifica o MVP sem vincular o token ao payload da operação.

### 7.4 Idempotência e retry

A transferência já usa `idempotency_key`, mas o step-up token continua sendo de
uso único.

Proposta para o MVP:

- retry com o mesmo `X-Step-Up-Token` após consumo retorna
  `STEP_UP_TOKEN_CONSUMED`;
- retry com novo step-up token e o mesmo `idempotency_key` pode alcançar a regra
  de idempotência da transferência;
- o mobile deve tratar `STEP_UP_TOKEN_CONSUMED` pedindo nova autorização de
  step-up quando necessário.

Esta decisão pode ser revisitada em evolução futura, especialmente se o step-up
token passar a ser vinculado à intenção detalhada da operação.

### 7.5 Separação entre erros required

Proposta de semântica:

- `TRANSACTION_PASSWORD_REQUIRED`: erro de challenge indicando que o endpoint
  sensível exige autorização de step-up;
- `STEP_UP_TOKEN_REQUIRED`: erro específico para requisição protegida sem o
  header obrigatório `X-Step-Up-Token`.

No endpoint de transferência, a ausência do header tende a retornar
`STEP_UP_TOKEN_REQUIRED`. `TRANSACTION_PASSWORD_REQUIRED` pode ser usado por
camadas de policy/challenge quando a API precisar informar previamente que uma
operação exige step-up.

### 7.6 Nome do caso de uso de enforcement

Criar um caso de uso no módulo de segurança, chamado pela delivery de
transferência antes do use case sensível.

Nome sugerido:

```text
EnforceStepUpUseCase
```

Entrada sugerida:

```text
authenticated_user
endpoint_key
step_up_token
now
```

Saída sugerida:

```text
allow / error
```

O use case de transferência interna não deve receber step-up token, senha
transacional ou detalhes de policy.
