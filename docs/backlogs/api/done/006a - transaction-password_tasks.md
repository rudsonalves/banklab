# Tasks da Senha Transacional

Backlog pai:

- `006a - transaction-password.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Segurança/ZTA

## Task 1/8: Criar migration da senha transacional

Status: Backlog

### Objetivo

Preparar o banco para armazenar a credencial de senha transacional com controle
de tentativas e bloqueio temporário.

### Escopo

- Criar tabela `transaction_passwords`.
- Incluir `id`, `user_id`, `password_hash`, `status`, `failed_attempts`,
  `locked_until`, `created_at`, `updated_at` e `changed_at`.
- Garantir vínculo com `users`.
- Garantir no máximo uma senha transacional ativa por usuário.
- Definir defaults seguros para `status` e `failed_attempts`.
- Criar índices necessários para busca por `user_id`.
- Criar migration de rollback compatível.

### Critérios de aceite

- A migration executa contra o schema atual.
- A tabela permite representar senha ativa e bloqueio temporário.
- O banco impede duplicidade indevida de senha transacional ativa por usuário.
- Rollback remove os objetos criados por esta migration.

### Depende de

- Nenhuma dependência.

## Task 2/8: Introduzir domínio da senha transacional

Status: Backlog

### Objetivo

Criar os tipos e regras de domínio para representar a senha transacional sem
expor detalhes HTTP ou infraestrutura.

### Escopo

- Criar módulo `internal/security/domain`, caso ainda não exista.
- Criar entidade ou modelo `TransactionPassword`.
- Representar status iniciais: `active` e `blocked`.
- Validar PIN numérico de 6 dígitos.
- Definir regra de bloqueio após 3 falhas consecutivas.
- Definir duração de bloqueio de 30 minutos.
- Definir comportamento de desbloqueio automático por leitura.
- Definir reset do contador quando o bloqueio expirar.
- Definir erros de domínio necessários para os cenários do backlog.

### Critérios de aceite

- PINs fora do formato de 6 dígitos numéricos são rejeitados.
- A regra de 3 falhas e bloqueio de 30 minutos está representada no domínio ou
  em serviço de aplicação claramente testável.
- O domínio não depende de HTTP, banco, JWT ou bibliotecas de delivery.
- Os erros necessários para criação, validação e bloqueio estão definidos.

### Depende de

- Nenhuma dependência.

## Task 3/8: Definir contratos de repositório e hasher

Status: Backlog

### Objetivo

Definir as portas necessárias para persistir e validar senha transacional sem
acoplar application à infraestrutura.

### Escopo

- Criar interface de repositório para senha transacional.
- Prever busca por `user_id`.
- Prever criação de senha transacional.
- Prever atualização de falhas, bloqueio, desbloqueio e `changed_at`.
- Criar contrato de hasher/comparador para a senha transacional.
- Reaproveitar o padrão de hashing existente quando fizer sentido, sem acoplar
  `security/domain` ao módulo `auth`.

### Critérios de aceite

- Application consegue criar e validar senha transacional usando apenas
  interfaces.
- O contrato suporta as regras de bloqueio temporário.
- O contrato permite testes com mocks/fakes sem banco real.

### Depende de

- Task 2.

## Task 4/8: Implementar persistência Postgres da senha transacional

Status: Backlog

### Objetivo

Implementar a infraestrutura Postgres para a senha transacional.

### Escopo

- Criar implementação em `internal/security/infrastructure`.
- Implementar criação da senha transacional.
- Implementar busca por usuário.
- Implementar atualização de tentativas inválidas.
- Implementar bloqueio temporário com `locked_until`.
- Implementar desbloqueio/reset de falhas quando aplicável.
- Mapear conflito de senha já existente para erro de domínio apropriado.

### Critérios de aceite

- Repositório persiste e recupera senha transacional corretamente.
- Conflito de criação duplicada é mapeado para erro esperado.
- Atualizações de falhas e bloqueio preservam consistência.
- Testes de infraestrutura cobrem criação, busca, conflito e bloqueio.

### Depende de

- Task 1.
- Task 3.

## Task 5/8: Implementar caso de uso de criação da senha transacional

Status: Backlog

### Objetivo

Permitir que qualquer usuário ativo autenticado crie sua senha transacional
inicial.

### Escopo

- Criar `CreateTransactionPasswordUseCase` em `internal/security/application`.
- Receber usuário autenticado, PIN e confirmação.
- Validar usuário ativo.
- Validar PIN numérico de 6 dígitos.
- Validar confirmação do PIN.
- Bloquear criação se já existir senha transacional ativa.
- Gerar hash antes de persistir.
- Não retornar senha, hash ou material sensível.

### Critérios de aceite

- Usuário ativo autenticado consegue criar senha transacional.
- Usuário não ativo não consegue criar senha transacional.
- PIN inválido falha antes de persistir.
- Confirmação divergente falha antes de persistir.
- Segunda criação com senha ativa falha com erro esperado.
- Senha é persistida apenas como hash.

### Depende de

- Task 3.
- Task 4.

## Task 6/8: Implementar endpoint de criação da senha transacional

Status: Backlog

### Objetivo

Expor a criação da senha transacional pelo contrato HTTP definido no backlog.

### Escopo

- Criar handler em `internal/security/delivery`.
- Registrar rota `POST /security/transaction-password`.
- Proteger rota com JWT existente.
- Aceitar payload:

```json
{
  "transaction_password": "123456",
  "transaction_password_confirmation": "123456"
}
```

- Usar `DisallowUnknownFields` se seguir o padrão dos handlers atuais.
- Responder usando o envelope padrão da API.
- Não logar PIN, confirmação ou hash.
- Conectar wiring no startup da API.

### Critérios de aceite

- A rota exige `Authorization: Bearer <access_token>`.
- Payload inválido retorna erro no envelope padrão.
- Criação com sucesso retorna resposta sem material sensível.
- Erros de domínio são mapeados para `error.code` estável.
- Nenhum log contém senha transacional.

### Depende de

- Task 5.

## Task 7/8: Registrar erros da senha transacional

Status: Backlog

### Objetivo

Adicionar os erros da senha transacional ao contrato compartilhado de erros da
API.

### Escopo

- Adicionar códigos compartilhados necessários:
  - `TRANSACTION_PASSWORD_ALREADY_SET`;
  - `TRANSACTION_PASSWORD_NOT_SET`;
  - `TRANSACTION_PASSWORD_INVALID`;
  - `TRANSACTION_PASSWORD_LOCKED`.
- Registrar mapeamento dos erros do módulo `security`.
- Mapear respostas usando o envelope definido em
  `api/docs/05-error_and_response.md`.
- Usar HTTP status definidos no backlog de contratos.

### Critérios de aceite

- `TRANSACTION_PASSWORD_ALREADY_SET` mapeia para HTTP 409.
- `TRANSACTION_PASSWORD_NOT_SET` mapeia para HTTP 409.
- `TRANSACTION_PASSWORD_INVALID` mapeia para HTTP 401.
- `TRANSACTION_PASSWORD_LOCKED` mapeia para HTTP 403.
- Mobile pode depender de `error.code`.

### Depende de

- Task 2.
- Task 5.

## Task 8/8: Cobrir senha transacional com testes e documentação mínima

Status: Backlog

### Objetivo

Garantir que a criação da senha transacional e suas regras principais estejam
cobertas antes de avançar para o step-up token.

### Escopo

- Adicionar testes de domínio/aplicação para PIN válido e inválido.
- Adicionar testes para confirmação divergente.
- Adicionar testes para usuário não ativo.
- Adicionar testes para criação duplicada.
- Adicionar testes de delivery para payload inválido e sucesso.
- Adicionar testes de mapeamento de erro.
- Atualizar documentação REST mínima, se a rota passar a existir.
- Executar testes afetados da API.

### Critérios de aceite

- Regras principais da senha transacional têm cobertura automatizada.
- Endpoint responde no envelope padrão.
- Códigos de erro esperados estão cobertos por testes.
- Documentação mínima da nova rota está disponível.
- Testes afetados da API passam.

### Depende de

- Task 6.
- Task 7.
