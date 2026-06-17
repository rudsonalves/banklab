# Tasks dos use cases de registro e gerenciamento da identidade de instalacao

Backlog pai:

- `017 - installation-identity-management-usecases.md`

Campos sugeridos para todas as tasks:

- Status: Concluída
- Prioridade: Alta
- Area: API
- Tipo: Aplicacao/Seguranca/Instalacoes

## Task 1/9: Definir contratos de aplicacao para registro de instalacao

Status: Concluída

### Objetivo

Modelar a entrada e a saida do caso de uso que registra explicitamente uma nova
instalacao a partir de contexto restrito.

### Escopo

- Definir input com:
  - `user_id` do contexto restrito;
  - `installation_id` do contexto restrito;
  - `jti` da autorizacao restrita;
  - step-up validado para a operacao publica esperada;
  - `installation_id` apresentado no header ou contrato equivalente.
- Definir output com sessao operacional:
  - access token;
  - refresh token;
  - dados seguros da instalacao criada, quando necessario.
- Separar contrato de aplicacao de DTO HTTP.
- Nao criar handler ou middleware nesta task.

### Criterios de aceite

- Use case futuro consegue receber contexto restrito sem depender de HTTP.
- Output deixa claro que registro bem-sucedido cria sessao operacional.
- Contrato permite comparar o `installation_id` apresentado no login com o
  valor apresentado no registro.
- Nenhum DTO de delivery e introduzido.

### Depende de

- Backlog 015.
- Backlog 016.

## Task 2/9: Implementar validacao de contexto restrito e step-up

Status: Concluída

### Objetivo

Garantir que o registro de instalacao so avance com autorizacao restrita ativa e
step-up valido para `POST /security/installations`.

### Escopo

- Ler `user_id`, `installation_id`, `jti` e escopo do contexto restrito.
- Validar escopo `installation.register`.
- Receber ou consultar resultado de step-up ja validado pela camada de
  seguranca.
- Confirmar que o step-up corresponde a operacao publica de registro de
  instalacao.
- Retornar erro de aplicacao quando contexto restrito ou step-up forem
  invalidos.

### Criterios de aceite

- Contexto restrito ausente ou incompleto bloqueia o registro.
- Escopo diferente de `installation.register` e rejeitado.
- Step-up de outra operacao e rejeitado.
- Testes cobrem contexto valido, contexto ausente, escopo errado e step-up
  divergente.

### Depende de

- Task 1.
- Fluxo de step-up existente.

## Task 3/9: Confirmar correspondencia do `installation_id`

Status: Concluída

### Objetivo

Impedir que uma autorizacao restrita emitida para uma instalacao seja usada para
registrar outra.

### Escopo

- Comparar `installation_id` do contexto restrito com o `installation_id`
  apresentado na chamada de registro.
- Retornar erro de instalacao divergente quando os valores nao coincidirem.
- Garantir que a comparacao use o value object canonico do dominio.
- Nao aceitar `installation_id` vazio, malformado ou diferente.

### Criterios de aceite

- Mesmo `installation_id` permite avancar.
- `installation_id` divergente retorna erro especifico.
- Erro nao consome autorizacao restrita.
- Testes cobrem igualdade, divergencia e valor invalido.

### Depende de

- Task 2.

## Task 4/9: Reservar vaga e criar instalacao `known`

Status: Concluída

### Objetivo

Criar a associacao `known` para a nova instalacao respeitando o limite de tres
instalacoes conhecidas.

### Escopo

- Usar `ReserveKnownInstallation` do repositorio.
- Gerar `resource_id` publico de gerenciamento.
- Executar reserva dentro de uma operacao atomica.
- Propagar `ErrInstallationLimitReached` quando nao houver vaga.
- Preservar historico de instalacoes revogadas sem ocupar slot.

### Criterios de aceite

- Registro cria instalacao `known`.
- Registro nao permite ultrapassar tres instalacoes `known`.
- Instalacao revogada libera vaga para nova instalacao.
- Erros de concorrencia nao vazam detalhe SQL.
- Testes cobrem sucesso e limite atingido.

### Depende de

- Backlog 014.
- Task 3.

## Task 5/9: Consumir autorizacao restrita de uso unico

Status: Concluída

### Objetivo

Garantir que a autorizacao restrita usada no registro seja consumida uma unica
vez e nao possa registrar multiplas instalacoes.

### Escopo

- Consumir autorizacao por `jti`.
- Validar que a autorizacao consumida pertence ao mesmo usuario e
  `installation_id`.
- Rejeitar autorizacao expirada, consumida ou revogada.
- Ordenar consumo e reserva de forma atomica para evitar estado parcial.
- Definir rollback esperado quando criacao de instalacao ou sessao falhar.

### Criterios de aceite

- Autorizacao ativa e consumida no registro bem-sucedido.
- Segunda tentativa com mesmo `jti` falha.
- Autorizacao de outro usuario ou instalacao falha.
- Falha durante o registro nao deixa autorizacao consumida sem instalacao,
  conforme transacao definida.
- Testes cobrem consumo unico e falhas de estado.

### Depende de

- Backlog 014.
- Task 4.

## Task 6/9: Criar sessao operacional apos registro

Status: Concluída

### Objetivo

Transformar o registro bem-sucedido em sessao operacional vinculada a nova
instalacao.

### Escopo

- Emitir access token operacional com claim `installation_id`.
- Emitir refresh token.
- Persistir refresh session vinculada ao usuario e instalacao.
- Reutilizar TTL e regras de sessao do login.
- Nao emitir novo `restricted_access_token`.

### Criterios de aceite

- Registro bem-sucedido retorna access token e refresh token.
- Sessao persistida contem `installation_id`.
- Access token contem a instalacao registrada.
- Nenhuma autorizacao restrita nova e criada.
- Testes cobrem tokens, hash de refresh e ausencia de token restrito.

### Depende de

- Backlog 015.
- Task 5.

## Task 7/9: Implementar listagem de instalacoes do usuario

Status: Concluída

### Objetivo

Permitir que o usuario autenticado consulte suas instalacoes registradas e o
historico necessario para gerenciamento.

### Escopo

- Criar use case de listagem por `user_id` autenticado.
- Usar `ListByUserID`.
- Retornar `resource_id`, status e timestamps seguros.
- Nao retornar o `installation_id` bruto apresentado pelo app.
- Definir ordenacao estavel conforme repositorio.

### Criterios de aceite

- Usuario recebe apenas suas instalacoes.
- Resposta nao expõe `installation_id`.
- Instalacoes revogadas aparecem como historico quando listadas pelo
  repositorio.
- Testes cobrem lista vazia, lista com conhecidas e lista com revogadas.

### Depende de

- Backlog 014.
- Backlog 015.

## Task 8/9: Implementar revogacao de instalacao

Status: Concluída

### Objetivo

Permitir que o usuario revogue uma instalacao propria por identificador publico
de gerenciamento.

### Escopo

- Criar use case de revogacao por `resource_id`.
- Validar que a instalacao pertence ao usuario autenticado.
- Impedir revogacao da instalacao da sessao atual.
- Usar `RevokeByResourceID` para revogacao logica.
- Nao exigir step-up para revogacao neste MVP.

### Criterios de aceite

- Instalacao propria pode ser revogada.
- Instalacao ausente retorna erro especifico.
- Instalacao de outro usuario nao e revogada.
- Instalacao da sessao atual nao pode ser revogada.
- Revogacao repetida retorna erro de dominio.
- Testes cobrem sucesso, ausente, outro usuario, atual e repetida.

### Depende de

- Backlog 014.
- Backlog 015.

## Task 9/9: Invalidar sessoes da instalacao revogada e validar build

Status: Concluída

### Objetivo

Garantir que a revogacao corte sessoes operacionais da instalacao revogada e
que os use cases fiquem prontos para o delivery do backlog 018.

### Escopo

- Chamar invalidacao de sessoes por usuario + `installation_id` apos revogacao.
- Usar o timestamp de revogacao como referencia de corte.
- Garantir que sessoes de outras instalacoes nao sejam afetadas.
- Cobrir caso sem sessoes ativas como operacao segura.
- Rodar `go test ./...`.
- Confirmar que nenhum handler HTTP novo foi criado neste backlog.

### Criterios de aceite

- Refresh sessions da instalacao revogada sao invalidadas.
- Sessoes de outras instalacoes permanecem ativas.
- Access tokens futuros da instalacao revogada serao bloqueaveis pelo
  enforcement do backlog 018.
- Suite Go passa.
- Use cases permanecem testaveis sem delivery HTTP.

### Depende de

- Task 8.
