# Tasks do Backlog API 013

Backlog pai:

- `013 - installation-identity-register-installation.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/7: Criar o caso de uso de registro explícito da instalação

Status: Backlog

### Objetivo

Implementar a orquestração principal de `POST /security/installations`.

### Escopo

- Criar caso de uso específico para o registro da instalação.
- Receber contexto restrito, `step_up_token` e `X-Installation-Id`.
- Encadear validação, criação da instalação, invalidação do grant e criação da
  sessão operacional.

### Critérios de aceite

- O caso de uso representa o fluxo sem depender do delivery layer.
- As etapas críticas do fluxo ficam centralizadas e testáveis.

### Depende de

- Backlogs 012 e 018.

## Task 2/7: Validar grant restrito e correspondência do `X-Installation-Id`

Status: Backlog

### Objetivo

Garantir que só a instalação autorizada no login possa concluir o registro.

### Escopo

- Validar `restricted_access_token`.
- Validar grant persistido e seu estado.
- Confirmar que o `X-Installation-Id` atual corresponde ao autorizado no login.
- Bloquear divergência com erro estável.

### Critérios de aceite

- Instalação divergente não conclui o registro.
- Grant expirado, consumido ou revogado falha corretamente.

### Depende de

- Task 1.
- Backlog 018.

## Task 3/7: Validar e consumir `step_up_token`

Status: Backlog

### Objetivo

Exigir prova de intenção antes de cadastrar a nova instalação.

### Escopo

- Enforçar `step_up_token` para `POST /security/installations`.
- Consumir o token de forma atômica.
- Garantir que ele pertença ao mesmo usuário e endpoint.

### Critérios de aceite

- Token ausente, inválido, expirado ou consumido é rejeitado.
- Token válido é consumido uma única vez.

### Depende de

- Task 1.
- Backlog 012.

## Task 4/7: Registrar a instalação com controle atômico de limite

Status: Backlog

### Objetivo

Criar a nova instalação `known` sem ultrapassar o limite de três instalações.

### Escopo

- Verificar vaga disponível no momento do registro.
- Criar a instalação como `known`.
- Reservar a vaga de forma atômica.
- Preservar instalações revogadas apenas como histórico.

### Critérios de aceite

- Registros concorrentes não ultrapassam o limite.
- A nova instalação nasce `known`.
- Instalações revogadas não ocupam vaga.

### Depende de

- Task 2.
- Backlog 018.

## Task 5/7: Criar sessão operacional e emitir tokens normais

Status: Backlog

### Objetivo

Encerrar o fluxo restrito já devolvendo a sessão definitiva vinculada à nova
instalação.

### Escopo

- Invalidar a autorização restrita após sucesso.
- Criar sessão operacional vinculada ao par usuário + `installation_id`.
- Emitir `access_token` com claim `installation_id`.
- Emitir `refresh_token` operacional.

### Critérios de aceite

- O endpoint retorna `access_token` e `refresh_token` normais.
- O `restricted_access_token` não continua utilizável após sucesso.
- O `access_token` carrega `installation_id`.

### Depende de

- Tasks 2 a 4.

## Task 6/7: Implementar o endpoint HTTP de registro de instalação

Status: Backlog

### Objetivo

Expor `POST /security/installations` com o contrato final do backlog.

### Escopo

- Criar handler e wiring da rota.
- Ler `Authorization`, `X-Step-Up-Token` e `X-Installation-Id`.
- Aplicar envelope padrão da API.
- Mapear erros de grant, `step_up`, header e limite.

### Critérios de aceite

- A rota responde com o contrato esperado.
- Erros relevantes são mapeados de forma estável.
- O endpoint não exige payload adicional no MVP.

### Depende de

- Tasks 1 a 5.

## Task 7/7: Cobrir o registro explícito da instalação com testes

Status: Backlog

### Objetivo

Proteger o fluxo fim a fim do endpoint de registro.

### Escopo

- Testar grant restrito válido.
- Testar `step_up_token` inválido ou consumido.
- Testar divergência de `X-Installation-Id`.
- Testar limite atingido.
- Testar sucesso com emissão de sessão operacional.

### Critérios de aceite

- O fluxo de sucesso e os principais ramos de recusa ficam cobertos.
- O contrato de resposta com tokens operacionais fica protegido por testes.

### Depende de

- Tasks 1 a 6.
