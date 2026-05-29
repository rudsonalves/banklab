# Tasks do contrato público de step-up por operação HTTP

Backlog pai:

- `007 - public-step-up-endpoint-contract.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA/Contrato

## Task 1/6: Definir modelo de operação HTTP pública

Status: Concluída

### Objetivo

Introduzir um modelo interno para representar a operação HTTP pública recebida
no endpoint de autorização de step-up.

### Escopo

- Criar tipo/estrutura para operação pública com:
  - `method`;
  - `path`.
- Normalizar `method` para maiúsculas.
- Validar `path` como path público:
  - começa com `/`;
  - não contém scheme;
  - não contém host;
  - não contém query string;
  - não contém fragmento.
- Definir que rotas parametrizadas futuras usam template público, como
  `/accounts/{id}/withdraw`.
- Não expor `endpoint_key` como input público.

### Critérios de aceite

- Operação válida `POST /accounts/internal-transfers` é aceita.
- `method` vazio é rejeitado.
- `path` vazio ou sem `/` inicial é rejeitado.
- `path` com `http://`, `https://`, host, query string ou fragmento é rejeitado.
- O modelo não depende de HTTP delivery, JWT ou banco.

### Depende de

- Nenhuma dependência.

## Task 2/6: Criar resolvedor de operação pública para policy interna

Status: Concluída

### Objetivo

Traduzir `method + path` público para a chave interna de policy usada pelo
step-up token e pelo enforcement.

### Escopo

- Criar resolvedor/whitelist no módulo `internal/security`.
- Mapear inicialmente:
  - `POST /accounts/internal-transfers` -> `internal_transfer.create`.
- Retornar `ErrStepUpEndpointNotAllowed` quando a operação não estiver
  permitida.
- Preservar `internal_transfer.create` como chave interna de policy.
- Preparar o resolvedor para suportar templates públicos em evolução futura.

### Critérios de aceite

- `POST /accounts/internal-transfers` resolve para
  `internal_transfer.create`.
- Método diferente para o mesmo path é rejeitado.
- Path não configurado é rejeitado.
- Testes cobrem sucesso e rejeições.

### Depende de

- Task 1.

## Task 3/6: Atualizar autorização de step-up para receber `method` e `path`

Status: Concluída

### Objetivo

Alterar `POST /security/step-up/authorize` para receber a operação pública em
vez de `endpoint_key`.

### Escopo

- Atualizar DTO/request do delivery:
  - remover `endpoint_key`;
  - adicionar `method`;
  - adicionar `path`.
- Rejeitar payload antigo com `endpoint_key` como campo inesperado.
- Resolver `method + path` para `endpoint_key` interno antes de emitir o token.
- Manter o use case de autorização emitindo token vinculado à policy interna.
- Manter o JWT com claim interno `endpoint_key`, se essa continuar sendo a
  melhor representação para o enforcement.
- Manter resposta atual:
  - `step_up_token`;
  - `expires_in`.

### Critérios de aceite

- Payload novo com `method=POST` e `path=/accounts/internal-transfers` emite
  token.
- Payload antigo com `endpoint_key` é rejeitado.
- Operação pública fora da whitelist retorna `STEP_UP_ENDPOINT_NOT_ALLOWED`.
- Campos ausentes ou malformados retornam erro controlado.
- Token emitido continua aceito pelo verifier e enforcement existentes.

### Depende de

- Task 2.

## Task 4/6: Atualizar testes do contrato público de step-up

Status: Concluída

### Objetivo

Garantir cobertura automatizada para o novo contrato público e para a transição
fora de `endpoint_key`.

### Escopo

- Atualizar testes de delivery de `POST /security/step-up/authorize`.
- Atualizar testes de application/use case quando necessário.
- Adicionar testes do resolvedor de operação pública.
- Cobrir rejeição de:
  - payload antigo com `endpoint_key`;
  - método inválido;
  - path inválido;
  - path com host;
  - path com query string;
  - path com fragmento;
  - operação fora da whitelist.
- Manter cobertura de sucesso para emissão de token e `expires_in`.
- Garantir que testes de enforcement continuem passando sem alterar o contrato
  do header `X-Step-Up-Token`.

### Critérios de aceite

- Testes cobrem `method + path -> endpoint_key`.
- Testes provam que o mobile não precisa enviar `internal_transfer.create`.
- Testes de regressão do signer/verifier/enforcement seguem passando.

### Depende de

- Task 3.

## Task 5/6: Atualizar documentação API e mobile para `method` + `path`

Status: Concluída

### Objetivo

Remover das documentações de consumo a orientação para clientes enviarem
`endpoint_key`.

### Escopo

- Atualizar `api/docs/07-api-rest.md`.
- Atualizar `api/docs/implementations/03-zta-step-up-transaction-password.md`.
- Atualizar docs de erro se houver mudança de descrição.
- Atualizar READMEs quando citarem o contrato de step-up.
- Atualizar `docs/backlogs/mobile/011 - senha-transacional-e-step-up.md`.
- Documentar que `endpoint_key` pode existir internamente, mas não é input
  público.
- Documentar que rotas parametrizadas devem usar template público, como
  `/accounts/{id}/withdraw`, quando forem suportadas.

### Critérios de aceite

- Documentação REST usa `method` e `path`.
- Documentação mobile usa `method` e `path`.
- Nenhuma documentação de consumo orienta cliente a enviar `endpoint_key`.
- O contrato continua deixando claro que `X-Step-Up-Token` é exigido na
  transferência interna.

### Depende de

- Task 3.
- Task 4.

## Task 6/6: Verificar alinhamento final do contrato público de step-up

Status: Concluída

### Objetivo

Fechar o backlog garantindo que implementação, testes e documentação usam o
novo contrato público.

### Escopo

- Conferir:
  - delivery de autorização de step-up;
  - resolvedor `method + path`;
  - use case de autorização;
  - signer/verifier JWT;
  - enforcement da transferência interna;
  - documentação REST;
  - documentação mobile.
- Rodar `go test ./...`.
- Registrar no backlog o status de alinhamento final.

### Critérios de aceite

- `go test ./...` passa.
- `POST /security/step-up/authorize` não exige `endpoint_key` público.
- `POST /security/step-up/authorize` aceita `method`, `path` e
  `transaction_password`.
- Token emitido para `POST /accounts/internal-transfers` é aceito na
  transferência interna.
- Não há divergência conhecida entre código, testes e documentação.

### Depende de

- Task 5.

### Registro de alinhamento final

Data de fechamento: 2026-05-29

Verificações executadas:

- `go test ./...` executado na API com sucesso.
- `POST /security/step-up/authorize` atualizado para contrato público com:
  - `method`;
  - `path`;
  - `transaction_password`.
- Payload legado com `endpoint_key` rejeitado como campo inesperado.
- Resolvedor `method + path -> endpoint_key` implementado com whitelist inicial:
  - `POST /accounts/internal-transfers` -> `internal_transfer.create`.
- JWT de step-up mantém claim interno `endpoint_key` para enforcement.
- Enforcement de `POST /accounts/internal-transfers` segue exigindo
  `X-Step-Up-Token` válido e de uso único.
- Testes de delivery/application/domain cobrem:
  - emissão com contrato novo;
  - rejeição de payload legado;
  - rejeição de método/path inválidos;
  - rejeição de operação fora da whitelist;
  - resolução de operação pública para policy interna.
- Documentação REST e mobile atualizada para `method` + `path`, sem orientar
  clientes a enviar `endpoint_key`.

Conclusão:

- Não há divergência conhecida entre implementação, testes e documentação para
  o contrato público de step-up por operação HTTP.
