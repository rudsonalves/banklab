-- ENUMS
CREATE TYPE account_status AS ENUM ('active', 'inactive', 'blocked');
CREATE TYPE transaction_type AS ENUM ('deposit', 'withdraw', 'transfer_in', 'transfer_out');

-- CUSTOMERS
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    email VARCHAR(120) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_cpf_format
    CHECK (cpf ~ '^\d{11}$')
);

-- ACCOUNTS
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    number VARCHAR(20) NOT NULL UNIQUE,
    branch VARCHAR(10) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    status account_status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_customer_id
ON accounts(customer_id);

-- TRANSACTIONS (LOGICAL)
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    type transaction_type NOT NULL,
    amount BIGINT NOT NULL,
    description TEXT,
    related_account_id UUID NULL,
    reference_id UUID NULL,
    idempotency_key VARCHAR(100) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_amount_positive
    CHECK (amount > 0)
);

CREATE INDEX idx_transactions_account_id
ON transactions(account_id);

CREATE INDEX idx_transactions_created_at
ON transactions(created_at DESC);

CREATE INDEX idx_transactions_reference_id
ON transactions(reference_id);

CREATE UNIQUE INDEX ux_transactions_idempotency
ON transactions(account_id, idempotency_key)
WHERE idempotency_key IS NOT NULL;