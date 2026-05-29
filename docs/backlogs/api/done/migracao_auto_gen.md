Eu faria essa migração em etapas pequenas, porque mudar geração de ID mexe em domínio, aplicação, repositórios e testes.

**Plano De Migração**

**Etapa 1: Banco Como Fonte De ID**

Objetivo: garantir que todas as tabelas já conseguem gerar UUID sozinhas.

Escopo:

- `customers`
- `accounts`
- `users`
- `user_sessions`
- `transactions`
- `customer_documents`
- `customer_addresses`
- `contact_verifications`
- `transaction_passwords`

Critério:

```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

Status: praticamente feito nas migrations.

**Etapa 2: Customer**

Objetivo: remover geração de `Customer.ID` no domínio.

Mudanças:

- `NewCustomer` não chama mais `uuid.New()`;
- `Customer.ID` continua `uuid.UUID`;
- `customerRepo.Create` deixa de inserir `id`;
- `customerRepo.Create` usa `RETURNING id`;
- fluxo de registro cria documento CPF depois de persistir customer.

Ordem:

```text
Customer domain
Customer repository
RegisterUserUseCase
Testes de customer/auth
```

**Etapa 3: CustomerDocument**

Objetivo: remover geração de `CustomerDocument.ID` no domínio.

Mudanças:

- `NewCPFDocument` não chama mais `uuid.New()`;
- `CreateDocument` deixa de inserir `id`;
- `CreateDocument` usa `RETURNING id`;
- entidade recebe ID após persistência.

Depende da etapa 2, porque documento precisa do `customer.ID` já persistido.

**Etapa 4: User**

Objetivo: remover geração de `User.ID` na aplicação.

Mudanças:

- ajustar `domain.NewUser` para não exigir `id` como argumento;
- `PostgresUserRepository.Create` deixa de inserir `id`;
- usa `RETURNING id`;
- `RegisterUserUseCase` usa `user.ID` retornado pelo repositório.

Depende da etapa 2, porque user customer precisa de `customer.ID`.

**Etapa 5: Account**

Objetivo: remover geração de `Account.ID` no domínio.

Mudanças:

- `NewAccount` não chama `uuid.New()`;
- `accountRepo.Create` deixa de inserir `id`;
- usa `RETURNING id`;
- fluxos que precisam do ID após criação usam entidade preenchida pelo repo.

**Etapa 6: UserSession**

Objetivo: remover `uuid.New()` direto no repositório.

Mudanças:

- `PostgresSessionRepository.Create` deixa de inserir `id`;
- como o método não retorna ID hoje, pode apenas inserir sem `RETURNING`;
- banco gera o ID.

Essa é bem simples.

**Etapa 7: ContactVerification**

Objetivo: remover geração de `ContactVerification.ID`.

Atenção: esse caso é mais delicado, porque a API retorna `verification_id`.

Mudanças:

- `NewContactVerification` não chama `uuid.New()`;
- repository usa `INSERT ... RETURNING id`;
- em caso de `ON CONFLICT DO UPDATE`, também precisa retornar o novo `id`;
- use case só retorna depois do repositório preencher `verification.ID`.

**Etapa 8: Transactions**

Objetivo: remover geração de `Transaction.ID` no domínio.

Mudanças:

- `NewTransaction` e `NewTransactionWithIdempotency` não geram ID;
- `CreateTransaction` deixa de inserir `id`;
- se precisar do ID após insert, usar `RETURNING id`;
- preservar `referenceID`, porque ele não é PK e continua sendo ID de correlação criado na aplicação.

Atenção: `referenceID := uuid.New()` em transferência deve continuar.

**Etapa 9: Security / Transaction Password**

Objetivo: já nascer no padrão novo.

Mudanças:

- não gerar ID no domínio/application;
- `transaction_passwords.id` gerado no banco;
- repository usa `RETURNING id`.

**Ordem recomendada**

```text
1. Banco
2. UserSession
3. Customer
4. CustomerDocument
5. User
6. Account
7. ContactVerification
8. Transactions
9. Security novo já no padrão
```

Eu começaria por `UserSession`, porque é o menor e valida o padrão sem mexer em fluxo de domínio. Depois iria para `Customer`, que é o primeiro caso realmente importante para o registro.