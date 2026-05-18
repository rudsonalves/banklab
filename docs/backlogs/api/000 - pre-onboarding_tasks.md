# Tasks do Pre-Onboarding Cadastral

Backlog pai:

- `000 - pre-onboarding.md`

Campos sugeridos para todas as tasks:

- Status: Backlog
- Prioridade: Alta
- Área: API
- Tipo: Arquitetura/Modelo de Dados

## Task 1/12: Criar a migration de preparação cadastral

Status: Backlog

### Objetivo

Preparar o schema do banco para o novo modelo cadastral sem implementar o fluxo
completo de checkpoints do onboarding.

### Escopo

- Adicionar `phone`, `email_verified_at` e `phone_verified_at` em `users`.
- Manter `users.email` como identificador de login e campo único.
- Tornar `users.phone` único.
- Adicionar `birth_date` em `customers`.
- Criar `customer_documents`.
- Criar `customer_addresses`.
- Adicionar constraint de documento único em `(type, value, country)`.
- Adicionar índice único parcial permitindo no máximo um documento primário por
  customer.
- Adicionar índice único parcial permitindo no máximo um endereço primário por
  customer.
- Usar `country CHAR(2) NOT NULL DEFAULT 'BR'` em documentos e endereços.
- Manter `issuer` e `issuer_state` como campos opcionais no banco.
- Não remover `customers.cpf` antes de migrar os valores existentes.

### Critérios de aceite

- A nova migration executa contra o schema atual.
- `customer_documents` suporta CPF, RG, CNH e tipos futuros de documento.
- `customer_addresses` existe, mas ainda não é exigida por nenhum fluxo
  funcional.
- A unicidade de documento primário e endereço primário é protegida por índices
  no banco.
- A ordem da migration mantém `customers.cpf` disponível para a cópia dos dados.

### Depende de

- Nenhuma dependência.

## Task 2/12: Migrar customers.cpf para customer_documents

Status: Backlog

### Objetivo

Mover os dados existentes de CPF de `customers.cpf` para o novo modelo de
documentos.

### Escopo

- Inserir uma linha em `customer_documents` para cada `customers.cpf` existente
  e não nulo.
- Usar `type = 'cpf'`.
- Usar `country = 'BR'`.
- Usar `is_primary = true`.
- Preservar `customers.created_at` como `created_at` do documento quando
  possível.
- Definir `updated_at` com o timestamp de execução da migration.
- Tornar a migration suficientemente idempotente para retentativas locais quando
  viável.
- Não remover `customers.cpf` nesta task.

### Critérios de aceite

- Customers existentes passam a ter o CPF representado em
  `customer_documents`.
- A unicidade de CPF é protegida por `(type, value, country)`.
- Cada CPF migrado é primário para o respectivo customer.
- Nenhum dado de CPF é perdido antes de o código parar de depender de
  `customers.cpf`.

### Depende de

- Task 1.

## Task 3/12: Introduzir contratos de domínio para documentos do customer

Status: Backlog

### Objetivo

Representar documentos do customer explicitamente na fronteira de
domínio/aplicação.

### Escopo

- Adicionar um modelo de domínio `CustomerDocument`.
- Definir o tratamento de tipos de documento para o caso inicial de `cpf`.
- Manter o comportamento de validação e normalização de CPF hoje coberto por
  `api/internal/customer/domain/cpf.go`.
- Mover o uso da validação de CPF da criação de `Customer` para o caminho de
  criação de documento.
- Definir contratos de repositório para criar documentos e ler documentos
  primários ou CPF por customer.
- Manter `issuer` e `issuer_state` opcionais no banco.
- Permitir que regras de aplicação exijam campos de órgão emissor futuramente
  para tipos específicos de documento.

### Critérios de aceite

- CPF pode ser validado e normalizado sem ser um campo direto de `Customer`.
- A aplicação consegue criar um documento do customer na mesma transação da
  criação do customer.
- Erros de CPF duplicado são mapeados a partir da constraint de unicidade de
  `customer_documents`, e não de `customers_cpf_key`.
- Testes existentes de validação de CPF são preservados ou movidos com cobertura
  equivalente.

### Depende de

- Task 1.

## Task 4/12: Simplificar o modelo de domínio Customer

Status: Backlog

### Objetivo

Fazer `Customer` representar a identidade de negócio sem possuir CPF
diretamente.

### Escopo

- Remover `CPF` de `api/internal/customer/domain/entity.go`.
- Alterar `NewCustomer` para receber apenas campos do customer: nome e
  `birth_date`.
- Validar `name` obrigatório.
- Validar `birth_date` obrigatório conforme a regra de domínio escolhida.
- Manter geração de ID do customer e comportamento de `created_at`.
- Atualizar erros de domínio de customer quando necessário.
- Não remover a validação de CPF em si; ela pertence à criação de documento.

### Critérios de aceite

- `Customer` não tem mais campo direto `CPF`.
- `NewCustomer` não recebe mais CPF.
- Criação de customer não falha por validação específica de CPF.
- Testes de customer refletem a nova fronteira de responsabilidade.

### Depende de

- Task 3.

## Task 5/12: Atualizar persistência e leitura do repositório de customer

Status: Backlog

### Objetivo

Persistir e ler customers usando a separação entre customer e documento.

### Escopo

- Atualizar `api/internal/customer/infrastructure/postgres.go`.
- Remover `cpf` do `INSERT INTO customers`.
- Incluir `birth_date` na persistência e leitura de customer.
- Substituir `SELECT c.cpf` por join ou query auxiliar contra
  `customer_documents` quando a resposta ainda precisar de CPF.
- Ler CPF usando `type = 'cpf'` e `country = 'BR'`.
- Preferir o documento CPF primário quando a consulta precisar semanticamente do
  documento principal.
- Mapear violações de unicidade de documento para o erro de CPF duplicado
  existente ou novo.

### Critérios de aceite

- O repositório de customer não escreve mais em `customers.cpf`.
- Leituras de perfil do customer continuam retornando CPF onde o contrato atual
  da API exigir.
- O código do repositório não depende funcionalmente de `customers.cpf`.
- Testes do repositório de customer passam contra o novo schema.

### Depende de

- Task 2.
- Task 3.
- Task 4.

## Task 6/12: Atualizar registro para criar customer e documento CPF

Status: Backlog

### Objetivo

Manter CPF no payload inicial de registro, mas armazená-lo como documento do
customer.

### Escopo

- Atualizar `api/internal/auth/application/register_user.go`.
- Criar o customer sem CPF direto.
- Criar o documento CPF em `customer_documents` na mesma transação.
- Usar `type = 'cpf'`, `country = 'BR'` e `is_primary = true`.
- Manter validação e normalização de CPF antes da persistência.
- Mapear CPF duplicado a partir de `customer_documents_unique_document`.
- Manter criação de conta bancária fora desta fase.
- Preparar o caso de uso para retornar sessão com escopo de onboarding quando
  as tasks de autenticação/sessão forem implementadas.

### Critérios de aceite

- Registrar usuário cria uma linha em `users`, uma linha em `customers` e uma
  linha de documento CPF.
- CPF duplicado falha com o erro esperado de CPF duplicado.
- Se a criação do documento CPF falhar, a criação de user e customer é
  revertida.
- Registro não cria conta bancária.

### Depende de

- Task 3.
- Task 4.
- Task 5.

## Task 7/12: Adicionar schema e fluxo de verificação de contato

Status: Backlog

### Objetivo

Suportar verificação simulada de email e telefone durante o registro.

### Escopo

- Adicionar armazenamento para solicitações de verificação de contato, caso não
  exista estrutura adequada.
- Implementar `POST /auth/contact-verifications`.
- Implementar `POST /auth/contact-verifications/confirm`.
- Suportar `channel = 'email'` e `channel = 'phone'`.
- Retornar o token bruto em respostas de ambiente local/dev.
- Retornar um `verification_token` assinado ou aleatório após a confirmação.
- Garantir que o endpoint de registro valide que cada token de verificação
  pertence ao email ou telefone enviado.
- Não integrar provedores externos de email/SMS nesta fase.

### Critérios de aceite

- Um cliente consegue solicitar token de verificação por email.
- Um cliente consegue solicitar token de verificação por telefone.
- Um cliente consegue confirmar cada token e receber um token de verificação
  para registro.
- Respostas em local/dev expõem o token bruto para testabilidade.
- Tokens de verificação não podem ser reutilizados para outro destino.

### Depende de

- Task 1.

## Task 8/12: Estender modelo de user, repositório e handler de registro

Status: Backlog

### Objetivo

Armazenar telefone e estado de verificação de contato em users e expor o novo
contrato de registro.

### Escopo

- Adicionar `Phone`, `EmailVerifiedAt` e `PhoneVerifiedAt` ao modelo de auth
  user.
- Atualizar inserts e leituras do repositório de auth/user.
- Atualizar o body de `POST /auth/register` para incluir:
  `email`, `phone`, `password`, `name`, `birth_date`, `cpf`,
  `email_verification_token` e `phone_verification_token`.
- Validar telefone e data de nascimento obrigatórios.
- Validar os dois tokens de verificação antes de criar o user.
- Criar users com `email_verified_at` e `phone_verified_at` preenchidos.
- Preservar o envelope de resposta.
- Retornar sessão com `scope = "onboarding"` quando o suporte de
  token/sessão for atualizado.

### Critérios de aceite

- Registro exige telefone.
- Registro exige tokens de verificação de email e telefone.
- Users criados possuem telefone, timestamp de verificação de email e timestamp
  de verificação de telefone.
- Unicidade de telefone é aplicada e mapeada para erro estável de API.
- Testes existentes de campos obrigatórios são atualizados.

### Depende de

- Task 6.
- Task 7.

## Task 9/12: Exigir contatos verificados no login

Status: Backlog

### Objetivo

Impedir login quando email ou telefone ainda não tiver sido verificado.

### Escopo

- Atualizar o caso de uso de login para verificar `email_verified_at` e
  `phone_verified_at`.
- Retornar erro padronizado `CONTACT_NOT_VERIFIED` quando qualquer campo estiver
  ausente.
- Incluir detalhes indicando `email_verified` e `phone_verified`.
- Manter comportamento de credenciais inválidas inalterado.
- Garantir que a falha de verificação de contato aconteça antes da criação de
  token e refresh session.

### Critérios de aceite

- Users com email não verificado não conseguem fazer login.
- Users com telefone não verificado não conseguem fazer login.
- Users com ambos os contatos verificados seguem para as regras atuais de
  elegibilidade de login.
- `CONTACT_NOT_VERIFIED` usa o envelope padrão `data`/`error`.
- Credenciais inválidas não revelam estado de verificação de contato.

### Depende de

- Task 8.

## Task 10/12: Atualizar queries de conta e destinatário de transferência

Status: Backlog

### Objetivo

Remover dependências de `customers.cpf` nas consultas de conta e destinatário de
transferência.

### Escopo

- Atualizar `api/internal/account/bankaccount/infrastructure/repository.go`.
- Substituir `SELECT c.cpf` por join com `customer_documents`.
- Substituir `WHERE c.cpf = $1` por busca de documento usando:
  `customer_documents.type = 'cpf'`,
  `customer_documents.value = $1`,
  `customer_documents.country = 'BR'`.
- Preservar o formato de resposta de destinatário de transferência para
  compatibilidade com o mobile.
- Garantir que busca de conta por agência/número continue retornando documento
  do titular.
- Garantir que busca por CPF/documento continue funcionando após remover
  `customers.cpf`.

### Critérios de aceite

- Busca de destinatário por agência e número retorna CPF a partir de
  `customer_documents`.
- Busca de destinatário por CPF/documento filtra por `customer_documents`.
- Nenhuma query do repositório de account referencia `customers.cpf`.
- Testes de account e transfer passam com o novo schema.

### Depende de

- Task 2.
- Task 5.

## Task 11/12: Remover customers.cpf após a migração do código

Status: Backlog

### Objetivo

Finalizar a transição do modelo removendo fisicamente CPF de `customers`.

### Escopo

- Remover a coluna `customers.cpf` apenas depois que a aplicação não a utilizar
  mais.
- Remover `chk_cpf_format`.
- Remover `customers_cpf_key`.
- Atualizar o schema inicial ou a cadeia de migrations conforme a política de
  migrations do projeto.
- Atualizar SQL de setup de testes de integração que criem tabelas `customers`
  manualmente.
- Substituir inserts diretos em customer por inserts correspondentes em
  `customer_documents`.
- Buscar no código da API referências remanescentes a `customers.cpf`, `c.cpf` e
  `customers_cpf_key`.

### Critérios de aceite

- O banco de teste da API pode ser criado sem `customers.cpf`.
- Nenhuma query de produção lê ou escreve `customers.cpf`.
- Nenhum fixture de teste exige `customers.cpf`.
- Comportamento de CPF duplicado vem de `customer_documents`.
- `rg "customers\\.cpf|c\\.cpf|customers_cpf_key" api` não encontra dependência
  ativa de código.

### Depende de

- Task 5.
- Task 6.
- Task 10.

## Task 12/12: Atualizar docs da API, coleção Postman e testes afetados

Status: Backlog

### Objetivo

Fazer documentação e verificação do projeto refletirem o novo modelo cadastral.

### Escopo

- Atualizar `api/docs/09-database.md`.
- Atualizar `api/docs/01-domain_model.md`.
- Atualizar `api/docs/07-api-rest.md`.
- Atualizar `api/docs/08-auth_implementation.md`, se existir e for relevante.
- Atualizar docs de contribuição ou visão geral que ainda descrevam CPF em
  `customers`.
- Atualizar `tools/postman/Banklab_API.postman_collection.json`.
- Atualizar exemplos voltados ao mobile apenas onde o contrato da API mudar.
- Rodar testes de domínio de customer.
- Rodar testes de aplicação e delivery de auth.
- Rodar testes de repositório de customer.
- Rodar testes de account e transaction afetados por lookup de destinatário ou
  fixtures.
- Rodar a suíte completa da API se for viável.

### Critérios de aceite

- Documentação descreve CPF como documento do customer.
- Documentação de registro inclui telefone, data de nascimento e tokens de
  verificação.
- Documentação de login menciona `CONTACT_NOT_VERIFIED`.
- Exemplos da coleção Postman refletem o novo fluxo de registro e verificação.
- Testes afetados passam.
- Testes não executados são listados com o motivo.

### Depende de

- Task 7.
- Task 8.
- Task 9.
- Task 11.

## Ordem sugerida no GitHub Project

1. Criar a migration de preparação cadastral.
2. Migrar `customers.cpf` para `customer_documents`.
3. Introduzir contratos de domínio para documentos do customer.
4. Simplificar o modelo de domínio `Customer`.
5. Atualizar persistência e leitura do repositório de customer.
6. Atualizar registro para criar customer e documento CPF.
7. Adicionar schema e fluxo de verificação de contato.
8. Estender modelo de user, repositório e handler de registro.
9. Exigir contatos verificados no login.
10. Atualizar queries de conta e destinatário de transferência.
11. Remover `customers.cpf` após a migração do código.
12. Atualizar docs da API, coleção Postman e testes afetados.
