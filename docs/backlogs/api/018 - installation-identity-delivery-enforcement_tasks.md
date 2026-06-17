# Tasks de delivery e enforcement da identidade de instalacao

Backlog pai:

- `018 - installation-identity-delivery-enforcement.md`

Campos sugeridos para todas as tasks:

- Status: Concluída
- Prioridade: Alta
- Area: API
- Tipo: Delivery/Middleware/Contrato/Seguranca

## Task 1/10: Consolidar resposta HTTP do login com instalacao

Status: Concluída

### Objetivo

Expor no `POST /auth/login` os resultados de autenticacao definidos no backlog
016 de forma estavel para o cliente.

### Escopo

- Confirmar resposta operacional com `access_token` e `refresh_token`.
- Confirmar resposta restrita com `restricted_access_token`, tipo, escopo e
  expiracao.
- Garantir que fluxo restrito nao serialize `refresh_token` vazio como contrato
  obrigatorio.
- Mapear instalacao revogada e limite atingido para erros HTTP publicos.
- Preservar o envelope `{ data, error }`.

### Criterios de aceite

- Instalacao `known` recebe resposta operacional.
- Instalacao nova com vaga recebe resposta restrita sem refresh token.
- Instalacao revogada retorna erro publico estavel.
- Limite atingido retorna erro publico estavel.
- Testes de handler cobrem todos os formatos.

### Depende de

- Backlog 016.

## Task 2/10: Exigir `X-Installation-Id` no refresh

Status: Concluída

### Objetivo

Garantir que `POST /auth/refresh` valide a instalacao apresentada antes de
rotacionar a sessao.

### Escopo

- Ler e validar `X-Installation-Id` no refresh usando a mesma regra canonica do
  login.
- Propagar o valor validado para o use case de refresh, ou comparar no delivery
  se o contrato escolhido assim exigir.
- Comparar header com a `installation_id` persistida na refresh session.
- Retornar `INVALID_INSTALLATION_ID` para ausencia ou formato invalido.
- Retornar `INSTALLATION_MISMATCH` quando header e sessao divergirem.

### Criterios de aceite

- Refresh sem header retorna 400.
- Refresh com UUID invalido retorna 400.
- Refresh com instalacao divergente nao rotaciona token.
- Refresh com instalacao correta rotaciona access e refresh token.
- Testes cobrem sessao legada sem instalacao conforme regra de transicao
  escolhida.

### Depende de

- Backlog 015.
- Backlog 016.

## Task 3/10: Criar middleware operacional de instalacao

Status: Concluída

### Objetivo

Validar rotas autenticadas por access token comparando header, claim e contexto
operacional da instalacao.

### Escopo

- Exigir `X-Installation-Id` em rotas operacionais autenticadas.
- Ler `installation_id` do access token.
- Comparar header com claim.
- Rejeitar token sem claim quando enforcement estiver ativo.
- Popular `authctx.OperationalSession` com usuario, role, customer e
  instalacao.
- Manter erro explicito para mismatch.

### Criterios de aceite

- Header ausente ou invalido retorna `INVALID_INSTALLATION_ID`.
- Claim ausente ou invalida retorna erro de token ou mismatch conforme contrato.
- Header diferente da claim retorna `INSTALLATION_MISMATCH`.
- Contexto operacional contem `installation_id`.
- Testes cobrem armazenamento por valor e uso em handler protegido.

### Depende de

- Backlog 015.

## Task 4/10: Criar middleware restrito

Status: Concluída

### Objetivo

Permitir rotas do fluxo de registro com `restricted_access_token`, sem aceitar
esse token como sessao operacional.

### Escopo

- Ler token restrito do contrato HTTP escolhido.
- Validar assinatura, expiracao, tipo, escopo, `jti` e autorizacao persistida.
- Comparar `X-Installation-Id` com a claim restrita.
- Popular `authctx.RestrictedSession`.
- Garantir que contexto restrito nao popula `AuthenticatedUser` operacional.
- Rejeitar token operacional em rota restrita.

### Criterios de aceite

- Token restrito valido popula contexto restrito.
- Token ausente, expirado, consumido ou revogado e rejeitado.
- Token com escopo diferente e rejeitado.
- Header divergente retorna `INSTALLATION_MISMATCH`.
- Testes cobrem que rotas operacionais nao aceitam token restrito.

### Depende de

- Backlog 015.
- Backlog 016.

## Task 5/10: Permitir step-up com contexto restrito para registro

Status: Concluída

### Objetivo

Habilitar `POST /security/step-up/authorize` durante o fluxo restrito de
registro de instalacao.

### Escopo

- Adaptar o handler de step-up para aceitar contexto restrito quando a operacao
  for `POST /security/installations`.
- Manter operacoes operacionais exigindo autenticacao operacional.
- Usar `installation.register` como endpoint key para a operacao publica.
- Garantir que o step-up emitido pertença ao usuario do contexto restrito.
- Nao permitir outras operacoes sensiveis com contexto restrito.

### Criterios de aceite

- Contexto restrito autoriza step-up apenas para registro de instalacao.
- Contexto restrito nao autoriza step-up para transferencia interna.
- Contexto operacional continua funcionando para operacoes existentes.
- Testes cobrem permissao restrita e rejeicoes.

### Depende de

- Backlog 017.
- Task 4.

## Task 6/10: Implementar handler `POST /security/installations`

Status: Concluída

### Objetivo

Conectar o use case de registro de instalacao ao contrato REST.

### Escopo

- Criar rota `POST /security/installations`.
- Proteger rota com middleware restrito.
- Exigir `X-Installation-Id`.
- Exigir step-up token conforme contrato escolhido.
- Chamar use case de registro do backlog 017.
- Retornar sessao operacional vinculada apos sucesso.
- Consumir autorizacao restrita e step-up conforme use cases.

### Criterios de aceite

- Registro bem-sucedido retorna access token e refresh token.
- Registro consome autorizacao restrita.
- Registro com `installation_id` divergente retorna `INSTALLATION_MISMATCH`.
- Registro com step-up ausente ou invalido nao cria instalacao.
- Testes HTTP cobrem sucesso, mismatch, limite e replay de autorizacao.

### Depende de

- Backlog 017.
- Task 4.
- Task 5.

## Task 7/10: Implementar handler `GET /security/installations`

Status: Concluída

### Objetivo

Expor a listagem de instalacoes do usuario autenticado.

### Escopo

- Criar rota `GET /security/installations`.
- Proteger rota com middleware operacional.
- Chamar use case de listagem.
- Retornar `resource_id`, status e timestamps seguros.
- Nao retornar `installation_id` bruto.

### Criterios de aceite

- Usuario autenticado lista apenas suas instalacoes.
- Resposta nao contem `installation_id`.
- Instalacoes revogadas aparecem como historico quando retornadas pelo use
  case.
- Header/claim divergentes bloqueiam a chamada.
- Testes HTTP cobrem lista vazia e lista com itens.

### Depende de

- Backlog 017.
- Task 3.

## Task 8/10: Implementar handler `DELETE /security/installations/{installation_resource_id}`

Status: Concluída

### Objetivo

Expor a revogacao de instalacao por identificador publico de gerenciamento.

### Escopo

- Criar rota `DELETE /security/installations/{installation_resource_id}`.
- Proteger rota com middleware operacional.
- Validar UUID canonico do `installation_resource_id`.
- Chamar use case de revogacao.
- Impedir revogacao da instalacao da sessao atual.
- Retornar resposta segura de sucesso sem expor `installation_id`.

### Criterios de aceite

- Instalacao propria e revogada.
- Instalacao atual nao pode ser revogada.
- Identificador invalido retorna erro de request.
- Revogacao invalida nao corta sessoes indevidas.
- Testes HTTP cobrem sucesso, atual, ausente e id invalido.

### Depende de

- Backlog 017.
- Task 3.

## Task 9/10: Ligar wiring e rotas no `main`

Status: Concluída

### Objetivo

Conectar repositórios, token services, middlewares e handlers novos no processo
real da API.

### Escopo

- Instanciar use cases de registro, listagem e revogacao.
- Instanciar middleware operacional e restrito.
- Registrar rotas de seguranca de instalacao.
- Garantir ordem correta de middlewares com app token, auth operacional e auth
  restrito.
- Preservar rotas existentes e comportamento de transferencia interna.

### Criterios de aceite

- API compila com wiring completo.
- Rotas novas aparecem no router.
- Rotas antigas continuam registradas.
- Nenhuma rota autenticada operacional fica sem enforcement de instalacao quando
  o backlog estiver completo.
- Testes de `cmd/api` passam.

### Depende de

- Task 3.
- Task 4.
- Task 6.
- Task 7.
- Task 8.

## Task 10/10: Atualizar contrato REST, colecoes e validacao final

Status: Concluída

### Objetivo

Documentar e validar o contrato HTTP final do MVP de identidade de instalacao.

### Escopo

- Atualizar documentacao REST afetada.
- Documentar `X-Installation-Id` em login, refresh e rotas autenticadas.
- Documentar resposta restrita do login.
- Documentar endpoints de instalacao.
- Documentar erros:
  - `INVALID_INSTALLATION_ID`;
  - `INSTALLATION_MISMATCH`;
  - `INSTALLATION_REVOKED`;
  - `INSTALLATION_LIMIT_REACHED`.
- Atualizar colecoes Bruno ou equivalentes, se existirem.
- Rodar `go test ./...`.

### Criterios de aceite

- Documentacao reflete o contrato implementado.
- Colecoes de teste conseguem exercitar login conhecido, login restrito,
  registro, listagem e revogacao.
- Suite Go passa.
- Contrato deixa claro que senha transacional fica fora do login.

### Depende de

- Task 9.
