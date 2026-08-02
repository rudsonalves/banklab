# Tasks de sessao, tokens e contexto da identidade de instalacao

Backlog pai:

- `015 - installation-identity-session-tokens-context.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Area: API
- Tipo: Sessao/JWT/Contexto/Seguranca

## Task 1/8: Preparar schema de sessao operacional vinculada

Status: Backlog

### Objetivo

Permitir que sessoes operacionais sejam vinculadas a uma instalacao sem alterar
handlers ou fluxo de login neste backlog.

### Escopo

- Adicionar suporte persistido para `installation_id` em `user_sessions`.
- Criar indice para consultas por `(user_id, installation_id)`.
- Preservar compatibilidade de dados existentes quando aplicavel.
- Documentar se a obrigatoriedade do campo sera ativada apenas quando o fluxo
  de login com instalacao estiver completo.

### Criterios de aceite

- Schema permite associar refresh session a uma instalacao.
- Schema permite invalidar sessoes por usuario + instalacao.
- Migration tem rollback.
- Nenhum handler passa a exigir o vinculo nesta task.

### Depende de

- Backlog 013.

## Task 2/8: Atualizar contrato de repositorio de sessoes

Status: Backlog

### Objetivo

Preparar a camada de persistencia de sessoes para criar, consultar e invalidar
sessoes vinculadas a instalacao.

### Escopo

- Definir entrada para criacao de sessao com `installation_id` opcional ou
  explicito conforme a transicao escolhida.
- Definir consulta de sessao que retorne `installation_id` quando persistido.
- Definir metodo para invalidar refresh tokens por usuario + instalacao.
- Atualizar implementacao Postgres de `user_sessions`.
- Manter o refresh atual funcionando enquanto o fluxo novo nao estiver ligado.

### Criterios de aceite

- Sessoes antigas continuam suportadas durante a transicao, se houver dados
  legados.
- Nova sessao pode persistir `installation_id`.
- Repositorio consegue revogar sessoes por instalacao.
- Testes cobrem criacao, consulta e revogacao por instalacao.

### Depende de

- Task 1.

## Task 3/8: Adicionar claim operacional `installation_id` ao access token

Status: Backlog

### Objetivo

Preparar tokens operacionais para carregar o vinculo de instalacao da sessao.

### Escopo

- Estender `TokenClaims` com `InstallationID`.
- Atualizar `GenerateAccessToken` para emitir claim `installation_id` quando
  disponivel.
- Atualizar `ParseAccessToken` para ler e validar `installation_id`.
- Manter compatibilidade com tokens antigos enquanto o enforcement ainda nao
  estiver ligado.
- Cobrir claim ausente, valida e invalida nos testes.

### Criterios de aceite

- Access token novo pode carregar `installation_id`.
- Parser retorna `installation_id` quando presente.
- Claim malformada torna o token invalido.
- Rotas atuais continuam funcionando antes do enforcement do backlog 018.

### Depende de

- Backlog 012.
- Task 2.

## Task 4/8: Criar servico de `restricted_access_token`

Status: Backlog

### Objetivo

Emitir e validar token restrito para registro futuro de instalacao sem criar
sessao operacional completa.

### Escopo

- Criar signer/verifier para `restricted_access_token`.
- Definir claims:
  - `sub`;
  - `jti`;
  - `token_type = restricted_access`;
  - `scope = installation.register`;
  - `installation_id`;
  - `iat`;
  - `exp`.
- Validar assinatura, expiracao, tipo, escopo, `jti` e `installation_id`.
- Usar o repositorio de autorizacoes restritas do backlog 014 para validar
  existencia e estado persistido.
- Nao emitir refresh token.

### Criterios de aceite

- Token restrito valido retorna claims verificadas.
- Token com tipo incorreto e rejeitado.
- Token com escopo incorreto e rejeitado.
- Token expirado e rejeitado.
- Token sem autorizacao persistida ativa e rejeitado.

### Depende de

- Backlog 014.

## Task 5/8: Definir contexto autenticado operacional com instalacao

Status: Backlog

### Objetivo

Preparar o contexto usado por use cases operacionais para carregar a
instalacao autenticada.

### Escopo

- Estender ou complementar `shared/authctx` com contexto operacional contendo:
  - `user_id`;
  - `role`;
  - `customer_id`;
  - `installation_id`.
- Manter helpers de leitura obrigatoria e opcional.
- Preservar compatibilidade com o contexto autenticado atual ate o middleware
  operacional ser atualizado no backlog 018.

### Criterios de aceite

- Use cases futuros conseguem ler `installation_id` do contexto.
- Ausencia de contexto retorna erro explicito.
- Testes cobrem armazenamento por valor e ponteiro quando aplicavel.
- Nenhum middleware passa a exigir `installation_id` nesta task.

### Depende de

- Task 3.

## Task 6/8: Definir contexto restrito

Status: Backlog

### Objetivo

Preparar o contexto usado por rotas permitidas durante o fluxo de registro de
nova instalacao.

### Escopo

- Criar helpers de contexto restrito contendo:
  - `user_id`;
  - `installation_id`;
  - `jti`;
  - `scope`.
- Expor leitura obrigatoria e opcional.
- Manter esse contexto separado da sessao operacional.
- Nao criar middleware HTTP neste backlog se isso ligar delivery
  prematuramente.

### Criterios de aceite

- Contexto restrito nao e confundido com contexto operacional.
- Helpers retornam erro explicito quando o contexto falta.
- Testes cobrem dados validos e ausencia de contexto.

### Depende de

- Task 4.

## Task 7/8: Preparar invalidacao de sessoes por instalacao revogada

Status: Backlog

### Objetivo

Disponibilizar infraestrutura para cortar refresh sessions de uma instalacao
revogada.

### Escopo

- Implementar ou adaptar porta `InstallationSessionInvalidator`.
- Invalidar sessoes por usuario + `installation_id`.
- Garantir que a operacao use timestamp de revogacao.
- Cobrir caso sem sessoes ativas como operacao segura.

### Criterios de aceite

- Refresh tokens da instalacao podem ser revogados em lote.
- Operacao nao revoga sessoes de outra instalacao.
- Operacao nao depende de handler HTTP.
- Testes cobrem uma, varias e nenhuma sessao ativa.

### Depende de

- Task 2.

## Task 8/8: Validar build sem enforcement

Status: Backlog

### Objetivo

Confirmar que a infraestrutura de sessao, tokens e contexto ficou pronta para
os use cases posteriores sem ativar enforcement prematuro.

### Escopo

- Rodar `go test ./...`.
- Confirmar que `POST /auth/refresh` ainda nao exige `X-Installation-Id`.
- Confirmar que rotas autenticadas ainda nao exigem `X-Installation-Id`.
- Confirmar que login ainda nao classifica ou vincula instalacao por causa
  deste backlog.

### Criterios de aceite

- Suíte Go passa.
- Nenhum handler novo e criado.
- Nenhum middleware novo e ligado no router.
- Nenhuma rota muda contrato antes do backlog de delivery/enforcement.

### Depende de

- Task 7.
