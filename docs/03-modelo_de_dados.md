# Modelo de Dados — Bank API

## 1. Visão Geral

O modelo de dados foi projetado para suportar:

* consistência transacional
* rastreabilidade completa
* operações financeiras seguras

O banco de dados é a **fonte de verdade do sistema**.


## 2. Princípios de Modelagem

### 2.1 Valores monetários

* armazenados em **centavos**
* tipo: `BIGINT`
* evita problemas de precisão


### 2.2 Identificadores

* tipo: `UUID`
* gerados pela aplicação


### 2.3 Normalização

* normalização básica
* sem duplicação desnecessária
* foco em clareza e consistência


### 2.4 Imutabilidade parcial

* `transactions` são imutáveis
* `accounts.balance` é mutável, porém controlado


## 3. Tabelas


## 3.1 customers

Representa o titular do sistema.

```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    email VARCHAR(120) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Índices

* UNIQUE(cpf)
* UNIQUE(email)


## 3.2 accounts

Representa contas bancárias.

```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    number VARCHAR(20) NOT NULL UNIQUE,
    branch VARCHAR(10) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Índices

* UNIQUE(number)
* INDEX(customer_id)


### Observações

* `balance` é armazenado para leitura rápida
* alterações devem sempre ocorrer via transação controlada


## 3.3 transactions

Representa movimentações financeiras.

```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES accounts(id),
    type VARCHAR(30) NOT NULL,
    amount BIGINT NOT NULL,
    description TEXT,
    related_account_id UUID NULL,
    idempotency_key VARCHAR(100) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```


### Índices

```sql
CREATE INDEX idx_transactions_account_id
ON transactions(account_id);

CREATE INDEX idx_transactions_created_at
ON transactions(created_at DESC);
```


### Índice de idempotência

```sql
CREATE UNIQUE INDEX ux_transactions_idempotency
ON transactions(account_id, idempotency_key)
WHERE idempotency_key IS NOT NULL;
```


## 4. Relacionamentos

```text
customers 1 ─── N accounts
accounts  1 ─── N transactions
```


## 5. Integridade de Dados

### 5.1 Foreign Keys

* accounts.customer_id → customers.id
* transactions.account_id → accounts.id


### 5.2 Regras implícitas (nível aplicação)

O banco NÃO garante sozinho:

* saldo suficiente
* atomicidade de transferência entre contas

Essas garantias são feitas na aplicação com suporte do banco.


## 6. Estratégia de Consistência

### 6.1 Atualização de saldo

Sempre dentro de transação:

```text
1. lock da conta
2. validação de saldo
3. update balance
4. insert transaction
```


### 6.2 Transferência

Dentro de uma única transação:

```text
1. lock conta origem
2. lock conta destino
3. valida saldo
4. update origem
5. update destino
6. insert transfer_out
7. insert transfer_in
```


## 7. Tipos de Dados Relevantes

| Campo      | Tipo      | Justificativa              |
| ---------- | --------- | -------------------------- |
| amount     | BIGINT    | precisão financeira        |
| balance    | BIGINT    | consistência e performance |
| id         | UUID      | unicidade global           |
| created_at | TIMESTAMP | rastreabilidade            |


## 8. Estratégia de Ordenação

* extratos ordenados por `created_at DESC`
* índice dedicado para performance


## 9. Idempotência

* garantida por índice único
* escopo: `(account_id, idempotency_key)`

Evita:

* duplicação de depósitos
* duplicação de transferências


## 10. Decisões Intencionais

### 10.1 Sem tabela de ledger separada

As `transactions` cumprem esse papel.


### 10.2 Sem soft delete

* registros não são removidos
* mantém rastreabilidade


### 10.3 Sem versionamento de saldo

* saldo atual é suficiente neste estágio


## 11. Limitações

* não há particionamento de tabela
* não há histórico de alteração de contas
* não há suporte a múltiplas moedas


## 12. Evolução futura

Possíveis melhorias:

* particionamento de `transactions`
* indexação avançada
* auditoria formal
* separação entre ledger e snapshot de saldo


## 13. Conclusão

O modelo de dados:

* é simples, porém consistente
* suporta operações financeiras seguras
* mantém rastreabilidade completa

Ele está alinhado com:

* o domínio
* os fluxos de uso
* a arquitetura definida
