# 🧱 1. Step 1 — Introduzir `UserStatus` no domínio

## 🎯 Objetivo

Adicionar controle de lifecycle **sem afetar fluxos atuais**

---

## ✔️ Definição

No domínio de auth:

```go
type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
)
```

---

## ✔️ Atualizar entidade `User`

Adicionar:

```go
Status UserStatus
```

---

## ✔️ Regra de criação

No `NewUser`:

```go
status: UserStatusPending,
```

👉 Todo usuário nasce como `pending`

---

# 🧠 2. Invariante novo (IMPORTANTE)

Você introduziu uma regra implícita:

```text
User MUST be active to operate in the system
```

Isso não existia antes.

---

# ⚙️ 3. Step 2 — Persistência

## ✔️ Migration

```sql
ALTER TABLE users
ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending';
```

---

## ✔️ Repository

Garantir:

* insert inclui status
* update de status possível

Adicionar método:

```go
UpdateStatus(ctx context.Context, userID uuid.UUID, status UserStatus) error
```

---

# 🧩 4. Step 3 — Use Case: Approve User

## 🎯 Nome

```text
ApproveUser
```

---

## ✔️ Input

```go
type ApproveUserInput struct {
	UserID uuid.UUID
}
```

---

## ✔️ Fluxo (crítico)

```text
1. begin transaction
2. load user (FOR UPDATE)
3. validate:
   - user exists
   - status == pending
4. update status → active
5. create account (reuse domain)
6. commit
```

---

## ✔️ Regra essencial

Isso precisa ser **atômico**.

👉 status + account devem acontecer na mesma transação

---

## ✔️ Erros

* USER_NOT_FOUND
* USER_ALREADY_ACTIVE
* INTERNAL_ERROR

---

# 🔒 5. Step 4 — Bloqueio no CreateAccount

Você precisa proteger isso:

```text
POST /accounts
```

---

## ✔️ Regra nova

Antes de criar conta:

```go
if user.Status != UserStatusActive {
	return ErrForbidden
}
```

---

👉 Isso garante:

* ninguém cria conta antes de aprovação
* fluxo fica consistente

---

# 🌐 6. Step 5 — Endpoint Admin

## ✔️ Definição

```text
POST /admin/users/{id}/approve
```

---

## ✔️ Header

```http
Authorization: Bearer <admin_token>
```

---

## ✔️ Response

```json
{
  "data": {
    "user_id": "...",
    "status": "active",
    "account_id": "..."
  },
  "error": null
}
```

---

# ⚠️ 7. Autorização (não ignore isso)

Você já tem roles:

* admin
* customer

👉 Aqui precisa validar:

```go
if user.Role != admin {
	return ErrForbidden
}
```

---

# 🔁 8. Fluxo final (agora correto)

```text
Register
  ↓
User (pending)
  ↓
Admin approve
  ↓
User (active) + Account created
  ↓
Operações financeiras
```

---

# 🧠 9. Impacto no modelo (importante)

Você NÃO alterou:

* Customer
* Account
* Transaction

✔ Isso é ótimo — isolamento correto do domínio 

---

# 🚀 10. Ordem de implementação (prática)

Siga exatamente isso:

```text
1. Domain (UserStatus)
2. Migration
3. Repository
4. Use case (ApproveUser)
5. Endpoint admin
6. Guard no CreateAccount
```

---

# 📌 11. Decisão arquitetural (muito boa)

Você acabou de introduzir:

> separação entre identidade e capacidade operacional

Isso é base de sistemas reais:

* KYC
* aprovação manual
* antifraude

---

# 🔚 Conclusão

Você não está apenas “adicionando um campo”.

Você está introduzindo:

* lifecycle
* controle de acesso real
* ponto de governança

E fez isso sem quebrar:

* consistência
* arquitetura
* domínio
