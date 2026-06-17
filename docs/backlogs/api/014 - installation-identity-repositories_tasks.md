# Tasks dos repositorios da identidade de instalacao

Backlog pai:

- `014 - installation-identity-repositories.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Area: API
- Tipo: Repositorio/Postgres/Seguranca

## Task 1/8: Criar estrutura do pacote de infraestrutura

Status: Concluída

### Objetivo

Preparar o pacote Postgres da identidade de instalacao sem ligar o repositório
a use cases, handlers ou wiring.

### Escopo

- Criar pacote de infraestrutura para `internal/installation`.
- Definir struct de repositório com `pgxpool.Pool` ou executor compatível com
  os padrões existentes.
- Criar helpers de scan para:
  - `app_installations`;
  - `installation_registration_authorizations`.
- Mapear linhas SQL para entidades do domínio do backlog 012.
- Mapear `pgx.ErrNoRows` para erros de domínio adequados.

### Criterios de aceite

- O pacote compila.
- Nenhum handler ou use case passa a depender do novo pacote.
- Helpers de scan validam entidades usando os construtores/restauradores de
  domínio.
- Erros de ausência não vazam detalhes de Postgres.

### Depende de

- Backlog 012.
- Backlog 013.

## Task 2/8: Implementar consultas de leitura de instalacoes

Status: Concluída

### Objetivo

Implementar as portas de leitura necessárias para classificação, listagem e
revogação futura.

### Escopo

- Implementar busca por `(user_id, installation_id)`.
- Implementar busca por `(user_id, resource_id)`.
- Implementar `CountKnownByUserID`.
- Implementar `HasAnyByUserID`.
- Implementar `ListByUserID`.
- Garantir ordenação estável para listagem.

### Criterios de aceite

- Instalação existente é restaurada corretamente.
- Instalação ausente retorna erro/resultado conforme contrato de domínio.
- Contagem considera apenas `status = 'known'`.
- `HasAnyByUserID` considera histórico, incluindo instalações revogadas.
- Listagem não expõe `installation_id` como identificador público de
  gerenciamento.

### Depende de

- Task 1.

## Task 3/8: Implementar bootstrap atomico da primeira instalacao

Status: Concluída

### Objetivo

Garantir que somente uma primeira instalação seja criada para usuário que nunca
teve associação histórica.

### Escopo

- Implementar `BootstrapFirstInstallation`.
- Executar a operação dentro de transação.
- Bloquear o usuário ou o conjunto relevante de linhas para evitar corrida.
- Confirmar que não existe qualquer instalação histórica para o usuário.
- Criar instalação `known` com `known_slot = 1`.
- Retornar erro de conflito quando outro processo vencer a corrida.

### Criterios de aceite

- Dois bootstraps concorrentes não criam duas primeiras instalações.
- Instalações revogadas contam como histórico e impedem novo bootstrap.
- A operação cria instalação `known` válida.
- Em conflito, o erro é mapeável pelo domínio/aplicação.

### Depende de

- Task 2.

## Task 4/8: Implementar reserva atomica de vaga para instalacao conhecida

Status: Concluída

### Objetivo

Criar instalação `known` para uma nova instalação respeitando o limite de três
instalações conhecidas.

### Escopo

- Implementar `ReserveKnownInstallation`.
- Executar a operação dentro de transação.
- Bloquear estado necessário para o usuário.
- Encontrar slot livre entre `1..3`.
- Criar instalação `known` com `known_slot` reservado.
- Retornar `ErrInstallationLimitReached` quando não houver slot.
- Preservar instalações revogadas no histórico sem ocupar slot.

### Criterios de aceite

- Repositório não permite mais de três instalações `known` por usuário.
- Registros concorrentes não ultrapassam o limite.
- Revogadas liberam slot.
- Duplicidade de `(user_id, installation_id)` é tratada sem erro SQL cru.

### Depende de

- Task 3.

## Task 5/8: Implementar revogacao logica de instalacao

Status: Concluída

### Objetivo

Revogar uma instalação por identificador público de gerenciamento preservando o
histórico e liberando vaga.

### Escopo

- Implementar `RevokeByResourceID`.
- Buscar instalação pelo `resource_id` e usuário.
- Alterar `status` para `revoked`.
- Preencher `revoked_at`.
- Limpar `known_slot`.
- Atualizar `updated_at`.
- Tratar revogação repetida como erro de domínio.

### Criterios de aceite

- A linha não é removida fisicamente.
- Instalação revogada permanece listável como histórico.
- Slot é liberado para reserva futura.
- Repositório não permite reativar instalação revogada.

### Depende de

- Task 4.

## Task 6/8: Implementar repositorio de autorizacoes restritas

Status: Concluída

### Objetivo

Implementar persistência de autorizações restritas usadas no futuro fluxo de
registro de instalação.

### Escopo

- Implementar `Create`.
- Implementar `FindByJTI`.
- Implementar `ConsumeByJTI`.
- Implementar `RevokeByJTI`.
- Implementar `RevokeActiveByUserIDAndInstallationID`.
- Mapear estados `active`, `consumed` e `revoked`.
- Tratar expiração como derivada de `expires_at`.

### Criterios de aceite

- `jti` duplicado retorna erro de domínio adequado.
- Consumo é atômico e de uso único.
- Autorização expirada não é consumida como válida.
- Revogação de autorização ativa altera estado sem remover linha.
- Busca por `jti` restaura entidade de domínio válida.

### Depende de

- Task 1.

## Task 7/8: Cobrir testes de repositorios e concorrencia

Status: Concluída

### Objetivo

Garantir que os repositórios cumprem contratos de domínio e invariantes de
concorrência.

### Escopo

- Criar testes de integração Postgres para repositório de instalações.
- Criar testes de integração Postgres para repositório de autorizações
  restritas.
- Testar ausência, sucesso e estados inválidos.
- Testar bootstrap concorrente.
- Testar reserva concorrente no limite de três instalações `known`.
- Testar consumo concorrente de autorização restrita.

### Criterios de aceite

- Testes falham se duas primeiras instalações forem criadas.
- Testes falham se mais de três instalações `known` forem criadas.
- Testes falham se autorização restrita for consumida duas vezes.
- Testes não dependem de delivery HTTP.

### Depende de

- Task 6.

## Task 8/8: Validar integração de build sem wiring

Status: Concluída

### Objetivo

Confirmar que os repositórios estão prontos para os próximos backlogs sem
serem ligados prematuramente ao fluxo de login ou delivery.

### Escopo

- Rodar `go test ./...`.
- Confirmar que o pacote novo não é importado por handlers.
- Confirmar que o login atual não classifica instalações por causa deste
  backlog.
- Confirmar que migrations do backlog 013 continuam aplicadas nos testes de
  repositório.

### Criterios de aceite

- Suíte Go passa.
- Nenhuma rota nova é criada.
- Nenhum token novo é emitido.
- Nenhum middleware novo é ligado.

### Depende de

- Task 7.
