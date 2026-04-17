-- MIGRATION: 000008_transfer_pair_integrity.up.sql
-- Hardens transfer pair integrity and lookup performance in account_transactions.

-- Fast lookup for replay by (reference_id, type).
CREATE INDEX IF NOT EXISTS idx_account_transactions_reference_type
ON account_transactions(reference_id, type);

-- Guardrail: for transfer pairs, each reference_id may have at most one leg per type.
-- We scope this to transfer_in/transfer_out rows only.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM account_transactions
        WHERE type IN ('transfer_in', 'transfer_out')
          AND reference_id IS NOT NULL
        GROUP BY reference_id, type
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot add transfer pair uniqueness: duplicate (reference_id, type) rows exist';
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS ux_account_transactions_transfer_pair
ON account_transactions(reference_id, type)
WHERE reference_id IS NOT NULL
  AND type IN ('transfer_in', 'transfer_out');
