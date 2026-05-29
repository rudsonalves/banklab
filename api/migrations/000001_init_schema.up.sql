-- MIGRATION: 000001_init_schema.up.sql
-- Baseline migration representing the full current database schema.
-- Replaces all prior incremental migrations.

-- ============================================================
-- EXTENSIONS
-- ============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================================
-- ENUMS
-- ============================================================

CREATE TYPE account_status AS ENUM ('active', 'inactive', 'blocked');

CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'withdraw',
    'transfer_in',
    'transfer_out'
);

-- ============================================================
-- SEQUENCES
-- ============================================================

CREATE SEQUENCE account_number_seq
    START WITH 10000000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- ============================================================
-- FUNCTIONS
-- ============================================================

CREATE OR REPLACE FUNCTION prevent_transactions_mutation()
RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'transactions is immutable';
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- TABLES
-- ============================================================

CREATE TABLE customers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(120) NOT NULL,
    cpf        VARCHAR(11)  NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_cpf_format CHECK (cpf ~ '^\d{11}$')
);

CREATE TABLE accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID        NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    number      VARCHAR(20) NOT NULL UNIQUE,
    branch      VARCHAR(10) NOT NULL,
    balance     BIGINT      NOT NULL DEFAULT 0,
    status      account_status NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(120) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL,
    customer_id   UUID         UNIQUE REFERENCES customers(id) ON DELETE SET NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    status        VARCHAR(20)  NOT NULL DEFAULT 'pending',

    CONSTRAINT chk_users_customer_role_consistency
        CHECK ((role = 'customer' AND customer_id IS NOT NULL) OR role <> 'customer')
);

CREATE TABLE user_sessions (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash CHAR(64)    NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Immutable financial ledger. All balance changes are recorded here.
CREATE TABLE transactions (
    id                UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id        UUID         NOT NULL REFERENCES accounts(id),
    type              transaction_type NOT NULL,
    amount            BIGINT       NOT NULL,
    balance_after     BIGINT       NOT NULL,
    reference_id      UUID,
    related_account_id UUID,
    idempotency_key   VARCHAR(100),
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================
-- INDEXES
-- ============================================================

CREATE INDEX idx_accounts_customer_id
    ON accounts(customer_id);

CREATE INDEX idx_transactions_account_id
    ON transactions(account_id);

CREATE INDEX idx_transactions_reference_id
    ON transactions(reference_id);

CREATE INDEX idx_transactions_created_at
    ON transactions(created_at DESC);

CREATE INDEX idx_transactions_account_created
    ON transactions(account_id, created_at DESC);

CREATE INDEX idx_transactions_reference_type
    ON transactions(reference_id, type);

CREATE UNIQUE INDEX ux_transactions_idempotency
    ON transactions(account_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

CREATE UNIQUE INDEX ux_transactions_transfer_pair
    ON transactions(reference_id, type)
    WHERE reference_id IS NOT NULL
      AND type IN ('transfer_in', 'transfer_out');

CREATE INDEX idx_user_sessions_user_id
    ON user_sessions(user_id);

CREATE INDEX idx_user_sessions_expires_at
    ON user_sessions(expires_at);

-- ============================================================
-- TRIGGERS
-- ============================================================

CREATE TRIGGER trg_transactions_no_mutation
BEFORE UPDATE OR DELETE ON transactions
FOR EACH ROW
EXECUTE FUNCTION prevent_transactions_mutation();
