-- MIGRATION: 000007_consolidate_ledger.up.sql
-- Consolidates the `transactions` table into `account_transactions`.
-- `transactions` was used solely as an idempotency store for transfers;
-- those fields are now part of the ledger itself.

-- 1. Add new columns to the ledger
ALTER TABLE account_transactions
    ADD COLUMN related_account_id UUID NULL,
    ADD COLUMN idempotency_key VARCHAR(100) NULL;

-- 2. Backfill idempotency data from transactions into the ledger
--    Only transfer_out rows carry the idempotency_key (origin side).
UPDATE account_transactions at
SET
    idempotency_key    = t.idempotency_key,
    related_account_id = t.related_account_id
FROM transactions t
WHERE at.account_id   = t.account_id
  AND at.reference_id = t.reference_id
  AND at.type         = 'transfer_out'
  AND t.idempotency_key IS NOT NULL;

-- 3. Promote the type column to use the enum (was VARCHAR(20))
ALTER TABLE account_transactions
    ALTER COLUMN type TYPE transaction_type
    USING type::transaction_type;

-- 4. Unique index for idempotency (replaces ux_transactions_idempotency)
CREATE UNIQUE INDEX ux_account_transactions_idempotency
    ON account_transactions(account_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

-- 5. Drop the now-redundant transactions table and its indexes
DROP INDEX IF EXISTS ux_transactions_idempotency;
DROP INDEX IF EXISTS idx_transactions_reference_id;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP TABLE transactions;
