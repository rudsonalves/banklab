-- MIGRATION: 000002_account_transactions.up.sql
-- This migration creates the account_transactions table, which serves as an immutable ledger of all transactions affecting
CREATE TABLE account_transactions (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES accounts(id),
    type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reference_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION prevent_account_transactions_mutation()
RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'account_transactions is immutable';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_account_transactions_no_mutation
BEFORE UPDATE OR DELETE ON account_transactions
FOR EACH ROW
EXECUTE FUNCTION prevent_account_transactions_mutation();

CREATE INDEX idx_account_transactions_account_id
ON account_transactions(account_id);

CREATE INDEX idx_account_transactions_reference_id
ON account_transactions(reference_id);

CREATE INDEX idx_account_transactions_created_at
ON account_transactions(created_at DESC);

CREATE INDEX idx_account_transactions_account_created
ON account_transactions(account_id, created_at DESC);