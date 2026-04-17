-- MIGRATION: 000007_consolidate_ledger.down.sql
-- Restores the separate transactions table and reverts account_transactions.

-- 1. Recreate transactions table
CREATE TABLE transactions (
    id                UUID PRIMARY KEY,
    account_id        UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    type              transaction_type NOT NULL,
    amount            BIGINT NOT NULL,
    description       TEXT,
    related_account_id UUID NULL,
    reference_id      UUID NULL,
    idempotency_key   VARCHAR(100) NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_amount_positive CHECK (amount > 0)
);

-- 2. Restore idempotency records from the ledger
INSERT INTO transactions (id, account_id, type, amount, reference_id, related_account_id, idempotency_key, created_at)
SELECT id, account_id, type, amount, reference_id, related_account_id, idempotency_key, created_at
FROM account_transactions
WHERE idempotency_key IS NOT NULL;

-- 3. Recreate indexes on transactions
CREATE INDEX idx_transactions_account_id
    ON transactions(account_id);

CREATE INDEX idx_transactions_created_at
    ON transactions(created_at DESC);

CREATE INDEX idx_transactions_reference_id
    ON transactions(reference_id);

CREATE UNIQUE INDEX ux_transactions_idempotency
    ON transactions(account_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

-- 4. Revert account_transactions
DROP INDEX  IF EXISTS ux_account_transactions_idempotency;

ALTER TABLE account_transactions
    ALTER COLUMN type TYPE VARCHAR(20);

ALTER TABLE account_transactions
    DROP COLUMN related_account_id,
    DROP COLUMN idempotency_key;
