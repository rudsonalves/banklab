# Backlog: pre-onboarding cadastral

## 1. Contexto

Antes de implementar o novo fluxo de onboarding, o BankLab precisa reorganizar a base cadastral atual.

Hoje, o modelo mistura dados de acesso, dados de cliente e documento principal em estruturas simples:

- `users` concentra autenticação, papel, status e vínculo com customer;
- `customers` concentra nome e CPF;
- ainda não existem tabelas próprias para documentos do customer;
- ainda não existem tabelas próprias para endereço do customer.

Essa estrutura foi suficiente para a primeira versão do projeto, mas o novo onboarding exige uma separação mais clara:

- `user` representa a identidade de acesso;
- `customer` representa a identidade de negócio;
- documentos pertencem ao customer;
- endereços pertencem ao customer;
- telefone e verificação de contato pertencem ao user;
- criação de conta bancária só deve acontecer depois da aprovação administrativa.

Este backlog define uma fase anterior ao onboarding, focada em preparar o modelo de dados e os contratos internos para o fluxo novo.

## 2. Objetivo

Reestruturar o modelo cadastral para separar dados de acesso, dados de cliente, documentos e endereço antes da implementação do onboarding por checkpoints.

Esta fase deve preparar o banco e o código para o onboarding, mas não deve implementar o onboarding completo.

## 3. Escopo

Incluído:

- adicionar telefone e campos de verificação em `users`;
- simplificar `customers` para identidade de negócio;
- criar `customer_documents`;
- migrar o CPF atual para o novo modelo de documentos;
- criar `customer_addresses` para uso futuro;
- ajustar domínio, repositórios, handlers, queries e testes afetados;
- manter compatibilidade funcional com o fluxo atual sempre que possível.

Não incluído:

- implementação de checkpoints de onboarding;
- upload de KYC;
- consulta de CEP;
- integração com serviço externo de email/SMS;
- integração com serviço externo de KYC;
- UI mobile de onboarding;
- painel administrativo;
- alteração de email/telefone após cadastro;
- múltiplos endereços em uso funcional.

## 4. Modelo alvo

### users

`users` representa a identidade de acesso.

Modelo alvo:

```text
users
- id
- email
- password
- phone
- email_verified_at
- phone_verified_at
- role
- customer_id
- status
- created_at
- updated_at
```

Observação: se o código/banco já usar `password_hash`, esse nome deve ser preferido a `password`.

### customers

`customers` representa a identidade de negócio.

Modelo alvo:

```text
customers
- id
- name
- birth_date
- created_at
```

O campo `cpf` não deve permanecer como responsabilidade direta de `customers` no modelo final desta fase.

### customer_documents

`customer_documents` representa documentos associados ao customer.

Modelo alvo:

```text
customer_documents
- id
- customer_id
- type
- value
- issuer
- issuer_state
- country
- is_primary
- created_at
- updated_at
```

Uso inicial:

- CPF deve ser migrado de `customers.cpf` para `customer_documents`;
- CPF deve ser registrado com `type = 'cpf'`;
- CPF deve ser o documento primário no contexto brasileiro inicial;
- RG/CNH serão preparados pelo modelo, mas o uso completo entra no onboarding.

### customer_addresses

`customer_addresses` representa endereços associados ao customer.

Modelo alvo:

```text
customer_addresses
- id
- customer_id
- postal_code
- street
- number
- complement
- neighborhood
- city
- state
- country
- is_primary
- created_at
- updated_at
```

Uso inicial:

- a tabela pode ser criada nesta fase;
- ela ainda não precisa ser usada por fluxos funcionais;
- o primeiro uso real será no onboarding.

## 5. Decisões iniciais

- Telefone pertence ao `user`.
- Email pertence ao `user`.
- Verificação de email e telefone pertence ao `user`.
- Nome completo pertence ao `customer`.
- Data de nascimento pertence ao `customer`.
- CPF, RG e CNH pertencem ao modelo de documentos do customer.
- Endereço pertence ao customer, por meio de tabela própria.
- `customers` não deve receber `address_id` direto.
- `customer_addresses` deve permitir evolução futura para múltiplos endereços.
- A conta bancária não deve ser criada nesta fase.
- Como o projeto ainda está em desenvolvimento e não há dados reais a preservar, `customers.cpf` deve ser removido nesta fase.
- A consistência do novo modelo é mais importante do que manter compatibilidade temporária com dados locais antigos.
- A unicidade de documentos deve considerar `(type, value, country)`.
- `issuer` deve ser opcional no banco.
- `issuer_state` deve ser opcional no banco.
- A obrigatoriedade de `issuer` e `issuer_state` pode ser aplicada pela camada de aplicação conforme o tipo de documento.
- `country` deve ter default `BR`.
- `email` deve ser único em `users`, pois é usado como identificador de login.
- `phone` deve ser único em `users`.
- `email_verified_at` e `phone_verified_at` devem ser exigidos para login.
- No fluxo mobile, email e telefone devem ser verificados durante o registro, antes da criação efetiva do usuário.
- O registro inicial deve criar o usuário já com `email_verified_at` e `phone_verified_at` preenchidos.
- A verificação por link de email/SMS pode ser discutida futuramente para app web ou outros canais.
- A fase inicial deve simular a verificação de email e telefone, não apenas preparar os campos.
- Em ambiente local/dev, os tokens de verificação podem ser retornados pela API para permitir validação do fluxo sem provedor externo.
- `customer_documents` deve permitir no máximo um documento primário por customer.
- `customer_addresses` deve permitir no máximo um endereço primário por customer.
- As regras de primário devem ser protegidas por constraint ou índice único parcial no banco.

## 6. Impactos esperados no código

Áreas provavelmente afetadas:

- migrations;
- domínio de customer;
- domínio de auth/register;
- repositórios de customer;
- repositórios de auth/user;
- handlers de registro;
- DTOs de request/response;
- consultas que hoje leem `customers.cpf`;
- consultas de conta que retornam dados do customer;
- testes de auth, customer, account e transaction.

Pontos de atenção já identificados:

- `customers.cpf` aparece em queries de customer e account;
- testes de integração criam `customers` com CPF direto;
- validação de CPF hoje pertence ao domínio de customer;
- erro de CPF duplicado hoje depende de constraint em `customers.cpf`;
- o registro de usuário hoje recebe CPF diretamente.

## 7. Estratégia de migração

Estratégia sugerida:

1. Criar `customer_documents`.
2. Migrar dados existentes de `customers.cpf` para `customer_documents`.
3. Criar constraint de unicidade para documento por tipo/valor/país.
4. Ajustar código para ler CPF a partir de `customer_documents`.
5. Ajustar registro para criar customer e documento CPF.
6. Remover dependência funcional de `customers.cpf`.
7. Remover `customers.cpf` nesta fase.
8. Criar `customer_addresses` sem uso funcional imediato.
9. Adicionar `phone`, `email_verified_at` e `phone_verified_at` em `users`.

A remoção física de `customers.cpf` deve acontecer nesta fase, junto com a migração do código para `customer_documents`.

## 8. Pontos para discutir

Nenhum ponto em aberto nesta seção no momento.

## 9. Proposta de implementação

### 9.1 Migrations

Criar uma nova migration para preparar o modelo cadastral.

Alterações em `users`:

```sql
ALTER TABLE users
    ADD COLUMN phone VARCHAR(20) UNIQUE,
    ADD COLUMN email_verified_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN phone_verified_at TIMESTAMP WITH TIME ZONE;
```

Observações:

- `email` já é único e deve continuar sendo o identificador de login.
- `phone` deve ser único.
- `email_verified_at` e `phone_verified_at` serão preenchidos durante o fluxo de registro mobile.
- Login deve exigir email e telefone verificados.

Alterações em `customers`:

```sql
ALTER TABLE customers
    DROP CONSTRAINT IF EXISTS chk_cpf_format,
    DROP CONSTRAINT IF EXISTS customers_cpf_key,
    DROP COLUMN cpf;
```

Observação:

- Como o projeto ainda está em desenvolvimento e os dados locais podem ser reiniciados, a remoção de `customers.cpf` pode acontecer nesta fase.

Criar `customer_documents`:

```sql
CREATE TABLE customer_documents (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    value VARCHAR(80) NOT NULL,
    issuer VARCHAR(80),
    issuer_state VARCHAR(30),
    country CHAR(2) NOT NULL DEFAULT 'BR',
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_documents_unique_document
        UNIQUE (type, value, country)
);
```

Criar índice para um documento primário por customer:

```sql
CREATE UNIQUE INDEX customer_documents_one_primary_per_customer
ON customer_documents (customer_id)
WHERE is_primary = true;
```

Criar `customer_addresses`:

```sql
CREATE TABLE customer_addresses (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    postal_code VARCHAR(20) NOT NULL,
    street VARCHAR(160) NOT NULL,
    number VARCHAR(30) NOT NULL,
    complement VARCHAR(120),
    neighborhood VARCHAR(120),
    city VARCHAR(120) NOT NULL,
    state VARCHAR(60) NOT NULL,
    country CHAR(2) NOT NULL DEFAULT 'BR',
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

Criar índice para um endereço primário por customer:

```sql
CREATE UNIQUE INDEX customer_addresses_one_primary_per_customer
ON customer_addresses (customer_id)
WHERE is_primary = true;
```

### 9.2 Modelos de domínio

Atualizar `Customer` para deixar de carregar CPF diretamente:

```text
Customer
- ID
- Name
- CreatedAt
```

Criar modelo de documento do customer:

```text
CustomerDocument
- ID
- CustomerID
- Type
- Value
- Issuer
- IssuerState
- Country
- IsPrimary
- CreatedAt
- UpdatedAt
```

Criar modelo de endereço do customer, ainda sem uso funcional no primeiro fluxo:

```text
CustomerAddress
- ID
- CustomerID
- PostalCode
- Street
- Number
- Complement
- Neighborhood
- City
- State
- Country
- IsPrimary
- CreatedAt
- UpdatedAt
```

Atualizar `User` para incluir:

```text
User
- Phone
- EmailVerifiedAt
- PhoneVerifiedAt
```

### 9.3 Endpoints alterados

#### POST /auth/register

O endpoint atual deve ser alterado para receber telefone e tokens de verificação.

Ele continuará orquestrando a criação inicial de `user`, `customer` e CPF em uma única operação transacional. A separação entre identidade de acesso e identidade de negócio deve existir no modelo interno, não necessariamente em endpoints separados neste primeiro momento.

Payload atual:

```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "Maria Silva",
  "cpf": "12345678909"
}
```

Payload proposto:

```json
{
  "email": "user@example.com",
  "phone": "+5511999999999",
  "password": "abc12345",
  "name": "Maria Silva",
  "birth_date": "1990-01-15",
  "cpf": "12345678909",
  "email_verification_token": "123456",
  "phone_verification_token": "654321"
}
```

Regras:

- `email` é obrigatório e único.
- `phone` é obrigatório e único.
- `password` deve ter no mínimo 8 caracteres, com letras e números.
- `name` é obrigatório.
- `birth_date` é obrigatório.
- `cpf` é obrigatório no contexto brasileiro inicial.
- `email_verification_token` é obrigatório.
- `phone_verification_token` é obrigatório.
- O registro só cria usuário se email e telefone forem verificados.
- O customer é criado sem CPF direto.
- O CPF é criado em `customer_documents` com `type = 'cpf'`, `country = 'BR'` e `is_primary = true`.
- A conta bancária não é criada nesta fase.
- A resposta deve retornar uma sessão com escopo de onboarding, permitindo que o usuário continue o fluxo sem acesso à área principal do app.

Resposta:

```json
{
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "scope": "onboarding",
    "user_id": "user-id",
    "customer_id": "customer-id"
  },
  "error": null
}
```

#### POST /auth/login

O login deve passar a exigir email e telefone verificados.

Se email ou telefone não estiverem verificados, o login deve retornar erro padronizado.

Exemplo:

```json
{
  "data": null,
  "error": {
    "code": "CONTACT_NOT_VERIFIED",
    "message": "Contact verification required",
    "details": {
      "email_verified": false,
      "phone_verified": true
    }
  }
}
```

### 9.4 Endpoints novos

#### POST /auth/contact-verifications

Cria uma solicitação de verificação para email ou telefone.

Payload:

```json
{
  "channel": "email",
  "target": "user@example.com"
}
```

Ou:

```json
{
  "channel": "phone",
  "target": "+5511999999999"
}
```

Resposta em ambiente local/dev:

```json
{
  "data": {
    "verification_id": "verification-id",
    "channel": "email",
    "target": "user@example.com",
    "token": "123456",
    "expires_at": "2026-05-15T18:00:00Z"
  },
  "error": null
}
```

Observação:

- Em ambiente futuro com provedor externo, `token` não deve ser retornado na resposta.
- O backend deve enviar o token por email/SMS usando serviço externo.

#### POST /auth/contact-verifications/confirm

Confirma um token de verificação.

Payload:

```json
{
  "verification_id": "verification-id",
  "token": "123456"
}
```

Resposta:

```json
{
  "data": {
    "verification_token": "signed-or-random-verification-token",
    "channel": "email",
    "target": "user@example.com",
    "verified_at": "2026-05-15T17:45:00Z"
  },
  "error": null
}
```

Uso:

- O `verification_token` retornado deve ser enviado no `POST /auth/register`.
- O backend deve validar que o token pertence ao mesmo email ou telefone usado no registro.

#### PUT /onboarding/customer

Atualiza dados gerais do customer durante o onboarding ou em fluxo de revisão.

Payload:

```json
{
  "name": "Maria Silva",
  "birth_date": "1990-01-15"
}
```

Resposta:

```json
{
  "data": {
    "customer_id": "customer-id",
    "name": "Maria Silva",
    "birth_date": "1990-01-15"
  },
  "error": null
}
```

Observações:

- Este endpoint não altera email, telefone ou senha.
- Este endpoint não altera CPF.
- Email, telefone e senha pertencem a fluxos próprios fora do onboarding.

#### PUT /onboarding/documents

Atualiza documentos estruturados do customer durante o onboarding ou em fluxo de revisão.

Este endpoint deve ser usado para documentos como RG, CNH ou outros documentos adicionais. O CPF inicial é criado no `POST /auth/register`.

Payload:

```json
{
  "documents": [
    {
      "type": "rg",
      "value": "123456789",
      "issuer": "SSP",
      "issuer_state": "SP",
      "country": "BR",
      "is_primary": false
    },
    {
      "type": "cnh",
      "value": "12345678900",
      "issuer": "DETRAN",
      "issuer_state": "SP",
      "country": "BR",
      "is_primary": false
    }
  ]
}
```

Resposta:

```json
{
  "data": {
    "customer_id": "customer-id",
    "documents": [
      {
        "id": "document-id",
        "type": "rg",
        "value": "123456789",
        "issuer": "SSP",
        "issuer_state": "SP",
        "country": "BR",
        "is_primary": false
      }
    ]
  },
  "error": null
}
```

Observações:

- O endpoint deve receber o cadastro completo da etapa de documentos.
- Imagens de documentos pertencem ao KYC documental, não a este endpoint.
- CPF permanece como documento primário criado no registro inicial.

### 9.5 Repositórios e contratos internos

Criar ou ajustar contratos para:

- criar customer sem CPF direto;
- criar documento do customer;
- atualizar dados gerais do customer;
- criar ou substituir documentos adicionais do customer no onboarding;
- buscar documento principal do customer quando necessário;
- buscar CPF do customer quando fluxos existentes ainda precisarem exibir/consultar CPF;
- criar usuário com telefone e verificações preenchidas;
- verificar unicidade de telefone;
- validar tokens de verificação de contato.

### 9.6 Ajustes em consultas existentes

Consultas que hoje retornam `customers.cpf` devem ser ajustadas para buscar o documento principal ou CPF em `customer_documents`.

Exemplos de áreas afetadas:

- perfil do customer;
- listagem/dados de conta;
- busca de destinatários por CPF;
- testes de transação e conta que criam customers diretamente.

## 10. Critérios de aceite

- `users` possui telefone e campos de verificação.
- `customers` deixa de depender funcionalmente de CPF direto.
- CPF existente é representado em `customer_documents`.
- Registro de usuário cria customer e documento CPF.
- Queries existentes continuam funcionando após ajuste.
- Testes afetados são atualizados.
- `customer_addresses` existe, mas não precisa ser consumida ainda.
- Documentação de onboarding referencia esta fase como preparação.
- `POST /auth/register` recebe telefone e tokens de verificação.
- `POST /auth/register` cria CPF em `customer_documents`.
- `POST /auth/register` retorna sessão com `scope = onboarding`.
- Endpoints de verificação de contato existem e usam envelope padrão.
- Em ambiente local/dev, tokens de verificação podem ser retornados pela API.
- `PUT /onboarding/customer` atualiza dados gerais do customer.
- `PUT /onboarding/documents` atualiza documentos adicionais do customer.

## 11. Ordem sugerida de implementação

1. Criar migration de estrutura.
2. Migrar CPF para `customer_documents`.
3. Atualizar domínio e repositórios de customer.
4. Atualizar fluxo de registro.
5. Atualizar consultas de account que dependem de CPF.
6. Atualizar testes.
7. Criar `customer_addresses`.
8. Atualizar documentação técnica.

## 12. Relação com o onboarding

Este backlog prepara o terreno para o backlog de onboarding.

Depois desta fase, o onboarding poderá assumir que:

- customer já existe;
- user já existe;
- email e telefone pertencem ao user;
- documentos são extensíveis;
- endereço possui tabela própria;
- conta bancária ainda depende de aprovação administrativa.

## 13. Documentação e histórico de decisões

Este backlog deve funcionar também como registro das decisões tomadas antes da implementação.

A intenção é que novos colaboradores consigam entender:

- por que a reestruturação cadastral foi feita antes do onboarding;
- por que `user` e `customer` foram separados conceitualmente;
- por que CPF saiu de `customers` e passou para `customer_documents`;
- por que email e telefone ficam em `users`;
- por que email e telefone são verificados durante o registro mobile;
- por que o registro ainda cria `user`, `customer` e CPF em uma única operação transacional;
- por que `customer_addresses` é criada antes de ser usada funcionalmente.

Ao implementar esta fase, a documentação técnica deve ser atualizada para refletir as decisões consolidadas.

Documentos que provavelmente precisam ser revisados:

- [api/docs/09-database.md](../../../api/docs/09-database.md)
- [api/docs/01-domain_model.md](../../../api/docs/01-domain_model.md)
- [api/docs/07-api-rest.md](../../../api/docs/07-api-rest.md)
- [api/docs/08-auth_implementation.md](../../../api/docs/08-auth_implementation.md)
- [api/docs/visao_geral/chapters/15-guia_pratico_de_contribuicao.md](../../../api/docs/visao_geral/chapters/15-guia_pratico_de_contribuicao.md)

Também é recomendável manter este backlog no repositório mesmo depois da implementação, como histórico de deliberação do projeto.
