# Tasks do schema de banco da identidade de instalacao

Backlog pai:

- `013 - installation-identity-database-schema.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Area: API
- Tipo: Banco/Migrations/Seguranca

## Task 1/7: Criar migration base de instalacoes

Status: Concluída

### Objetivo

Criar a tabela `app_installations` para persistir associacoes entre usuario e
instalacao do app.

### Escopo

- Criar migration `up` e `down` com a proxima numeracao disponivel.
- Criar tabela `app_installations`.
- Incluir campos minimos:
  - `id`;
  - `resource_id`;
  - `user_id`;
  - `installation_id`;
  - `status`;
  - `platform`;
  - `app_version`;
  - `app_build`;
  - `first_seen_at`;
  - `last_seen_at`;
  - `revoked_at`;
  - `created_at`;
  - `updated_at`.
- Definir `resource_id` como identificador publico de gerenciamento separado
  do `installation_id`.
- Referenciar `users(id)` para `user_id`.

### Criterios de aceite

- Migration `up` cria `app_installations`.
- Migration `down` remove `app_installations`.
- `resource_id` nao reutiliza o `installation_id`.
- A tabela nao inclui dados de device fingerprinting, geolocalizacao ou
  attestation.

### Depende de

- Backlog 012.

## Task 2/7: Definir constraints de integridade de instalacoes

Status: Concluída

### Objetivo

Garantir no banco os invariantes persistidos da entidade de instalacao.

### Escopo

- Restringir `status` aos valores:
  - `known`;
  - `revoked`.
- Garantir unicidade de `resource_id`.
- Garantir unicidade de `(user_id, installation_id)`.
- Garantir que `revoked_at` seja nulo para `known`.
- Garantir que `revoked_at` seja obrigatorio para `revoked`.
- Garantir timestamps obrigatorios e consistentes quando possivel.

### Criterios de aceite

- Banco rejeita status fora do MVP, como `trusted`.
- Banco impede duplicidade da mesma instalacao para o mesmo usuario.
- Banco impede revogacao persistida sem `revoked_at`.
- Banco impede instalacao `known` com `revoked_at`.

### Depende de

- Task 1.

## Task 3/7: Criar indices de consulta de instalacoes

Status: Concluída

### Objetivo

Preparar as consultas necessarias para classificacao, listagem e revogacao sem
implementar repositorios ainda.

### Escopo

- Criar indice para busca por `(user_id, installation_id)`.
- Criar indice para busca por `(user_id, resource_id)`.
- Criar indice para listagem por `user_id`.
- Criar indice para consultas por `(user_id, status)`.
- Criar indice util para contagem de instalacoes `known`.

### Criterios de aceite

- Consultas previstas no backlog 014 podem usar indices diretos.
- Indices nao introduzem uma regra de negocio diferente das constraints.
- Nomes dos indices seguem o padrao local de migrations.

### Depende de

- Task 2.

## Task 4/7: Preparar suporte ao limite de tres instalacoes `known`

Status: Concluída

### Objetivo

Dar ao banco as garantias necessarias para que os repositorios consigam aplicar
atomicamente o limite de tres instalacoes `known` por usuario.

### Escopo

- Avaliar a estrategia relacional para suportar reserva atomica de vaga.
- Adicionar constraints, indices parciais ou estrutura auxiliar se necessario.
- Documentar no comentario da migration quando a garantia final depender de
  transacao no repositorio.
- Preservar instalacoes `revoked` no historico sem ocupar vaga.

### Criterios de aceite

- O schema permite contar e bloquear instalacoes `known` por usuario de forma
  segura no backlog de repositorios.
- Instalacoes `revoked` nao ocupam vaga do limite.
- O schema nao remove historico para liberar vaga.
- A decisao de concorrencia fica clara para quem implementar o backlog 014.

### Depende de

- Task 3.

## Task 5/7: Criar migration de autorizacoes restritas

Status: Concluída

### Objetivo

Criar a tabela `installation_registration_authorizations` para suportar o
futuro `restricted_access_token`.

### Escopo

- Criar tabela `installation_registration_authorizations`.
- Incluir campos minimos:
  - `id`;
  - `jti`;
  - `user_id`;
  - `installation_id`;
  - `scope`;
  - `status`;
  - `expires_at`;
  - `consumed_at`;
  - `created_at`.
- Referenciar `users(id)` para `user_id`.
- Definir `installation_id` como o valor apresentado no login, nao como
  `resource_id`.

### Criterios de aceite

- Migration `up` cria a tabela de autorizacoes.
- Migration `down` remove a tabela de autorizacoes.
- A tabela nao reutiliza `user_sessions`.
- A tabela nao persiste access token ou refresh token em texto claro.

### Depende de

- Backlog 012.

## Task 6/7: Definir constraints e indices de autorizacoes restritas

Status: Concluída

### Objetivo

Garantir unicidade, estado e consultas futuras das autorizacoes restritas.

### Escopo

- Garantir unicidade de `jti`.
- Restringir `scope` inicialmente a `installation.register`.
- Restringir `status` aos valores:
  - `active`;
  - `consumed`;
  - `revoked`.
- Garantir no maximo uma autorizacao `active` por
  `(user_id, installation_id, scope)`.
- Tratar expiracao como derivada de `expires_at`, sem status persistido
  `expired`.
- Criar indices para consulta por `jti`, usuario/instalacao/escopo/status e
  expiracao.
- Garantir que `consumed_at` seja nulo para `active`.
- Garantir que `consumed_at` exista para `consumed`.

### Criterios de aceite

- Banco impede `jti` duplicado.
- Banco impede duas autorizacoes `active` para o mesmo usuario, instalacao e
  escopo.
- Banco rejeita status `expired`.
- Consultas futuras por `jti` e limpeza por expiracao ficam indexadas.

### Depende de

- Task 5.

## Task 7/7: Validar migrations e rollback

Status: Concluída

### Objetivo

Garantir que o schema sobe e desce corretamente antes da implementacao de
repositorios.

### Escopo

- Rodar migrations em banco local de desenvolvimento/teste.
- Validar estrutura das tabelas.
- Validar constraints principais com inserts simples ou comandos SQL diretos.
- Validar rollback das migrations novas.
- Registrar qualquer decisao importante no backlog ou na propria migration.

### Criterios de aceite

- Migrations `up` executam sem erro.
- Migrations `down` removem as tabelas novas sem quebrar tabelas existentes.
- Constraints de status e unicidade sao verificadas.
- Nao ha implementacao de repositorio nesta task.

### Depende de

- Task 6.
