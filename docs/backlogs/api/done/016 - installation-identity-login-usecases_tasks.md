# Tasks dos use cases de login da identidade de instalacao

Backlog pai:

- `016 - installation-identity-login-usecases.md`

Campos sugeridos para todas as tasks:

- Status: Concluída
- Prioridade: Alta
- Area: API
- Tipo: Aplicacao/Auth/Seguranca

## Task 1/9: Definir contrato de resultado do login com instalacao

Status: Concluída

### Objetivo

Representar na camada de aplicacao os possiveis resultados do login conforme a
instalacao apresentada, sem acoplar a decisao ao delivery HTTP.

### Escopo

- Definir saida de login para sessao operacional.
- Definir saida de login para autorizacao restrita de registro.
- Definir saida de login para instalacao revogada.
- Definir saida de login para limite de instalacoes atingido.
- Garantir que a saida explicite quando ha ou nao refresh token.
- Evitar expor detalhes internos como `known_slot`, `jti` bruto quando nao for
  necessario ou identificador publico de gerenciamento de outra instalacao.

### Criterios de aceite

- Resultado operacional contem access token e refresh token.
- Resultado restrito contem somente `restricted_access_token` e metadados
  seguros para o cliente.
- Resultado de limite atingido nao cria sessao operacional.
- Resultado de revogada e distinguivel dos demais erros de autenticacao.
- O contrato nao importa pacotes de HTTP.

### Depende de

- Backlog 012.
- Backlog 015.

## Task 2/9: Implementar classificador de instalacao para login

Status: Concluída

### Objetivo

Centralizar a decisao que classifica a instalacao apresentada no login depois
que as credenciais do usuario forem validadas.

### Escopo

- Consultar instalacao por usuario e `installation_id`.
- Retornar `known` quando a instalacao existir e estiver conhecida.
- Retornar `revoked` quando a instalacao existir e estiver revogada.
- Retornar `first` quando o usuario nunca teve instalacao historica.
- Retornar `new` quando a instalacao nao existir e ainda houver vaga.
- Retornar `limit_reached` quando a instalacao nao existir e ja houver tres
  instalacoes `known`.
- Tratar erro de repositorio sem mascarar falhas operacionais como ausencia.

### Criterios de aceite

- Classificacao nao altera estado persistido.
- Instalacoes revogadas contam como historico para impedir novo bootstrap.
- Contagem de limite considera apenas instalacoes `known`.
- Testes cobrem todas as classificacoes.

### Depende de

- Backlog 014.
- Task 1.

## Task 3/9: Permitir login operacional para instalacao conhecida

Status: Concluída

### Objetivo

Emitir sessao operacional normal quando o login apresentar uma instalacao
`known`.

### Escopo

- Reutilizar a validacao atual de credenciais e elegibilidade da conta.
- Criar refresh session vinculada ao `installation_id`.
- Emitir access token operacional com claim `installation_id`.
- Preservar os dados de contexto ja emitidos no login atual.
- Nao exigir senha transacional no login.

### Criterios de aceite

- Instalacao `known` recebe access token e refresh token.
- Sessao persistida fica vinculada ao usuario e instalacao.
- Access token contem `installation_id`.
- Nenhum `restricted_access_token` e emitido para instalacao conhecida.
- Testes garantem que a senha transacional nao participa do fluxo.

### Depende de

- Backlog 015.
- Task 2.

## Task 4/9: Implementar bootstrap atomico da primeira instalacao

Status: Concluída

### Objetivo

Transformar a primeira instalacao historica do usuario em `known` durante o
login, de forma atomica e idempotente sob concorrencia.

### Escopo

- Usar a operacao atomica de bootstrap do repositorio.
- Criar a primeira instalacao como `known`.
- Vincular a sessao operacional criada a essa instalacao.
- Emitir tokens operacionais iguais ao fluxo de instalacao conhecida.
- Tratar conflito concorrente como decisao de aplicacao mapeavel.

### Criterios de aceite

- Primeiro login com primeira instalacao cria associacao `known`.
- Dois logins concorrentes nao criam duas primeiras instalacoes.
- Conflito nao retorna erro SQL cru.
- Login bem-sucedido apos bootstrap recebe refresh token.
- Testes cobrem sucesso e corrida simulada ou exercitada no repositorio.

### Depende de

- Backlog 014.
- Task 3.

## Task 5/9: Bloquear login de instalacao revogada

Status: Concluída

### Objetivo

Impedir que uma instalacao revogada volte a operar por novo login.

### Escopo

- Interromper o fluxo quando a classificacao for `revoked`.
- Nao recriar instalacao.
- Nao emitir access token, refresh token ou token restrito.
- Retornar erro de aplicacao especifico para instalacao revogada.
- Garantir que a senha do usuario nao altera a revogacao.

### Criterios de aceite

- Instalacao revogada nao recebe nenhum token.
- Instalacao revogada nao e reativada.
- O erro e distinguivel de credenciais invalidas no nivel de aplicacao.
- Testes cobrem instalacao revogada com credenciais validas.

### Depende de

- Task 2.

## Task 6/9: Retornar limite atingido para nova instalacao sem vaga

Status: Concluída

### Objetivo

Responder de forma segura quando uma instalacao nova tenta login e o usuario ja
possui tres instalacoes `known`.

### Escopo

- Detectar `limit_reached` depois da validacao de credenciais.
- Nao criar instalacao.
- Nao criar autorizacao restrita.
- Nao emitir access token, refresh token ou token restrito.
- Retornar erro ou resultado de aplicacao especifico para limite atingido.
- Preservar atomicidade da decisao quando a contagem estiver no limite.

### Criterios de aceite

- Usuario com tres instalacoes `known` nao recebe autorizacao para quarta.
- Instalacoes revogadas nao ocupam vaga.
- Corrida de logins novos nao permite ultrapassar o limite nos fluxos
  posteriores.
- Testes cobrem limite atingido e vaga liberada por revogacao.

### Depende de

- Backlog 014.
- Task 2.

## Task 7/9: Emitir autorizacao restrita para instalacao nova com vaga

Status: Concluída

### Objetivo

Permitir que uma instalacao nova avance para o registro explicito sem receber
sessao operacional completa.

### Escopo

- Criar autorizacao restrita persistida para usuario e `installation_id`.
- Emitir `restricted_access_token` com escopo `installation.register`.
- Definir expiracao curta conforme configuracao ou constante da aplicacao.
- Revogar autorizacoes restritas ativas anteriores para o mesmo usuario e
  `installation_id`, se o contrato escolhido exigir uma autorizacao ativa por
  instalacao.
- Nao criar refresh session.
- Nao emitir access token operacional.

### Criterios de aceite

- Instalacao nova com vaga recebe `restricted_access_token`.
- Token restrito referencia autorizacao persistida ativa.
- Resposta nao contem refresh token.
- Resposta nao cria instalacao `known`.
- Testes cobrem criacao, expiracao configurada e ausencia de refresh token.

### Depende de

- Backlog 014.
- Backlog 015.
- Task 6.

## Task 8/9: Orquestrar fluxo completo no use case de login

Status: Concluída

### Objetivo

Integrar validacao de credenciais, classificacao de instalacao e emissao do
resultado correto em um fluxo de aplicacao coeso.

### Escopo

- Manter a ordem segura do login atual:
  - validar entrada ja normalizada pelo handler;
  - validar credenciais;
  - validar elegibilidade da conta;
  - classificar instalacao;
  - executar decisao de sessao, bootstrap, bloqueio, limite ou restricao.
- Garantir que credenciais invalidas nao vazem informacao sobre instalacoes.
- Preservar comportamento de conta pendente, bloqueada ou inelegivel.
- Manter senha transacional fora do fluxo.
- Nao criar novo handler HTTP nesta task.

### Criterios de aceite

- Credenciais invalidas retornam o erro atual sem classificar instalacao.
- Contas inelegiveis preservam o comportamento esperado.
- Cada classificacao leva exatamente a um tipo de resultado.
- O use case permanece testavel sem HTTP.

### Depende de

- Task 3.
- Task 4.
- Task 5.
- Task 7.

## Task 9/9: Cobrir testes de aplicacao e validar build

Status: Concluída

### Objetivo

Garantir que as regras de login com identidade de instalacao estejam cobertas
por testes automatizados e prontas para o backlog de delivery/enforcement.

### Escopo

- Testar login com instalacao `known`.
- Testar bootstrap da primeira instalacao.
- Testar instalacao `revoked`.
- Testar instalacao nova com vaga.
- Testar instalacao nova com limite atingido.
- Testar credenciais invalidas sem classificacao de instalacao.
- Testar que fluxo restrito nao emite refresh token.
- Testar que senha transacional nao e exigida.
- Rodar `go test ./...`.

### Criterios de aceite

- Todas as classificacoes do backlog 012 sao exercitadas.
- Testes falham se token restrito criar refresh session.
- Testes falham se instalacao revogada receber token.
- Testes falham se primeira instalacao nao for criada de forma atomica.
- Suite Go passa.

### Depende de

- Task 8.
