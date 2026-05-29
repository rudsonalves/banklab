# Tasks do Step-Up Token

Backlog pai:

- `006b - step-up-token.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/9: Criar migration do step-up token

Status: Backlog

### Objetivo

Preparar o banco para persistir os `jti` emitidos para step-up token e permitir
controle de uso único.

### Escopo

- Criar tabela `step_up_tokens`.
- Incluir `id`, `jti`, `user_id`, `endpoint_key`, `status`, `expires_at`,
  `consumed_at` e `created_at`.
- Garantir vínculo com `users`.
- Garantir unicidade de `jti`.
- Definir status persistidos iniciais: `active` e `consumed`.
- Não persistir `expired` como status.
- Criar índices para consulta por `jti`, `user_id`, `endpoint_key` e expiração.
- Criar migration de rollback compatível.

### Critérios de aceite

- A migration executa contra o schema atual.
- O banco impede duplicidade de `jti`.
- A tabela permite representar token ativo, consumido e expiração por regra.
- Rollback remove os objetos criados por esta migration.

### Depende de

- Nenhuma dependência.

## Task 2/9: Introduzir domínio do step-up token

Status: Backlog

### Objetivo

Criar os tipos e regras de domínio para representar o step-up token sem expor
detalhes HTTP, JWT ou infraestrutura.

### Escopo

- Criar entidade ou modelo `StepUpToken` em `internal/security/domain`.
- Representar status persistidos: `active` e `consumed`.
- Representar expiração como regra derivada de `expires_at`.
- Definir duração padrão de 2 minutos.
- Validar `jti`, `user_id`, `endpoint_key`, `expires_at` e `created_at`.
- Criar comportamento para identificar token expirado.
- Criar comportamento para consumir token ativo.
- Definir erros de domínio necessários para emissão e uso futuro.

### Critérios de aceite

- Tokens expirados são identificados por `expires_at < now`.
- Token consumido não pode ser consumido novamente.
- Token expirado não pode ser consumido.
- O domínio não depende de HTTP, banco, JWT ou delivery.
- Os erros necessários para token inválido, expirado e consumido estão
  definidos.

### Depende de

- Nenhuma dependência.

## Task 3/9: Definir contratos de repositório, signer e policy de endpoint

Status: Backlog

### Objetivo

Definir as portas necessárias para emitir step-up token sem acoplar application
à infraestrutura de banco, JWT ou política de endpoint.

### Escopo

- Criar interface de repositório para step-up token.
- Prever criação/persistência de token emitido.
- Prever busca por `jti`.
- Prever consumo atômico de `jti` para uso futuro no backlog `006c`.
- Criar contrato para assinatura de JWT de step-up.
- Criar contrato para validação de `endpoint_key`.
- Definir whitelist inicial com `internal_transfer.create`.
- Isolar a whitelist para evolução futura para policy registry ou policy
  engine.

### Critérios de aceite

- Application consegue emitir step-up token usando apenas interfaces.
- O contrato suporta uso único por consumo atômico.
- A validação de endpoint não fica espalhada por handlers.
- O contrato permite testes com mocks/fakes sem banco real ou JWT real.

### Depende de

- Task 2.

## Task 4/9: Implementar persistência Postgres do step-up token

Status: Backlog

### Objetivo

Implementar a infraestrutura Postgres para persistência e consumo do step-up
token.

### Escopo

- Criar implementação em `internal/security/infrastructure`.
- Implementar criação de registro com `jti`, `user_id`, `endpoint_key`,
  `status=active`, `expires_at` e `created_at`.
- Implementar busca por `jti`.
- Implementar consumo atômico de `jti` quando `status=active` e não expirado.
- Não alterar status para `expired`; expiração deve ser regra de leitura.
- Mapear token inexistente, expirado ou consumido para erros de domínio
  apropriados.

### Critérios de aceite

- Repositório persiste e recupera step-up token corretamente.
- `jti` duplicado falha de forma controlada.
- Consumo atômico impede reutilização do mesmo token.
- Token expirado não é consumido.
- Testes de infraestrutura cobrem criação, busca, duplicidade, expiração e
  consumo.

### Depende de

- Task 1.
- Task 3.

## Task 5/9: Implementar signer JWT do step-up token

Status: Backlog

### Objetivo

Implementar a assinatura do JWT curto de step-up seguindo o modelo híbrido
definido no backlog.

### Escopo

- Criar implementação em `internal/security/infrastructure`.
- Gerar JWT assinado com os claims mínimos:
  - `user_id`;
  - `endpoint_key`;
  - `scope=step_up`;
  - `exp`;
  - `iat`;
  - `jti`.
- Usar duração de 2 minutos.
- Reaproveitar configuração/segredo JWT existente quando fizer sentido, ou
  isolar configuração específica se o padrão atual pedir.
- Não incluir senha transacional, hash ou payload da operação no JWT.

### Critérios de aceite

- JWT gerado contém todos os claims obrigatórios.
- `exp` corresponde à validade curta definida.
- `scope` é sempre `step_up`.
- O token assinado pode ser validado posteriormente pelo backend.
- Testes cobrem claims obrigatórios e expiração.

### Depende de

- Task 2.
- Task 3.

## Task 6/9: Implementar caso de uso de autorização de step-up

Status: Backlog

### Objetivo

Permitir que um usuário autenticado autorize um endpoint sensível usando senha
transacional e receba um step-up token curto.

### Escopo

- Criar `AuthorizeStepUpUseCase` em `internal/security/application`.
- Receber usuário autenticado, `endpoint_key` e senha transacional.
- Validar usuário autenticado e ativo conforme padrão já usado em `006a`.
- Validar `endpoint_key` pela policy/whitelist.
- Buscar senha transacional do usuário.
- Retornar erro se a senha transacional não existir.
- Normalizar bloqueio expirado antes de validar.
- Retornar erro se a senha transacional estiver bloqueada.
- Comparar PIN recebido com hash armazenado.
- Em caso de falha, registrar tentativa, persistir estado e retornar erro
  apropriado.
- Em caso de sucesso, zerar falhas, persistir estado, criar `jti`, persistir
  step-up token e gerar JWT.
- Não entregar JWT se a persistência do `jti` falhar.
- Retornar `step_up_token` e `expires_in=120`.

### Critérios de aceite

- Senha correta emite step-up token para `internal_transfer.create`.
- Senha inexistente retorna `TRANSACTION_PASSWORD_NOT_SET`.
- Senha inválida incrementa falhas e retorna `TRANSACTION_PASSWORD_INVALID`.
- Terceira falha bloqueia e retorna `TRANSACTION_PASSWORD_LOCKED`.
- Senha bloqueada não emite token.
- `endpoint_key` fora da whitelist retorna erro estável.
- O `jti` é persistido antes de o token ser retornado.
- Nenhum material sensível é retornado ou logado.

### Depende de

- Task 3.
- Task 4.
- Task 5.
- `006a - transaction-password`.

## Task 7/9: Implementar endpoint de autorização de step-up

Status: Backlog

### Objetivo

Expor a autorização de step-up pelo contrato HTTP definido no backlog.

### Escopo

- Adicionar rota `POST /security/step-up/authorize`.
- Proteger rota com JWT existente.
- Aceitar payload:

```json
{
  "endpoint_key": "internal_transfer.create",
  "transaction_password": "123456"
}
```

- Usar `DisallowUnknownFields` se seguir o padrão dos handlers atuais.
- Chamar `AuthorizeStepUpUseCase`.
- Responder sucesso com envelope padrão:

```json
{
  "data": {
    "step_up_token": "<token>",
    "expires_in": 120
  },
  "error": null
}
```

- Não logar senha transacional nem token completo.
- Conectar wiring no startup da API.

### Critérios de aceite

- A rota exige `Authorization: Bearer <access_token>`.
- Payload inválido retorna erro no envelope padrão.
- Sucesso retorna `step_up_token` e `expires_in`.
- Erros de domínio são mapeados para `error.code` estável.
- Nenhum log contém senha transacional ou material sensível.

### Depende de

- Task 6.

## Task 8/9: Registrar erros do step-up authorization

Status: Backlog

### Objetivo

Adicionar os erros específicos da autorização de step-up ao contrato
compartilhado de erros da API.

### Escopo

- Adicionar código compartilhado para endpoint não permitido:
  - `STEP_UP_ENDPOINT_NOT_ALLOWED`.
- Confirmar reutilização dos códigos já criados no `006a`:
  - `TRANSACTION_PASSWORD_NOT_SET`;
  - `TRANSACTION_PASSWORD_INVALID`;
  - `TRANSACTION_PASSWORD_LOCKED`.
- Registrar mapeamento dos erros do módulo `security`.
- Mapear respostas usando o envelope definido em
  `api/docs/05-error_and_response.md`.
- Separar erros de autorização de step-up dos erros de enforcement do backlog
  `006c`.

### Critérios de aceite

- `TRANSACTION_PASSWORD_NOT_SET` continua mapeando para HTTP 409.
- `TRANSACTION_PASSWORD_INVALID` continua mapeando para HTTP 401.
- `TRANSACTION_PASSWORD_LOCKED` continua mapeando para HTTP 403.
- `STEP_UP_ENDPOINT_NOT_ALLOWED` possui HTTP status definido e testado.
- Mobile pode depender de `error.code`.

### Depende de

- Task 2.
- Task 6.

## Task 9/9: Cobrir step-up token com testes e documentação mínima

Status: Backlog

### Objetivo

Garantir que a autorização de step-up e emissão do token estejam cobertas antes
de avançar para o enforcement da transferência interna.

### Escopo

- Adicionar testes de domínio para expiração, consumo e status.
- Adicionar testes de policy/whitelist de `endpoint_key`.
- Adicionar testes de signer JWT para claims obrigatórios.
- Adicionar testes de aplicação para senha correta, senha ausente, senha
  inválida, bloqueio e endpoint não permitido.
- Adicionar testes de infraestrutura para persistência e consumo atômico.
- Adicionar testes de delivery para payload inválido, sucesso e erros
  principais.
- Adicionar testes de mapeamento de erro.
- Atualizar documentação REST mínima da nova rota.
- Atualizar documentação ZTA quando necessário.
- Executar testes afetados da API.

### Critérios de aceite

- Regras principais do step-up token têm cobertura automatizada.
- Endpoint responde no envelope padrão.
- JWT contém claims obrigatórios.
- `jti` é persistido antes da resposta.
- Códigos de erro esperados estão cobertos por testes.
- Documentação mínima da nova rota está disponível.
- Testes afetados da API passam.

### Depende de

- Task 7.
- Task 8.
