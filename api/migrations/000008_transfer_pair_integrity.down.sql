-- MIGRATION: 000008_transfer_pair_integrity.down.sql

DROP INDEX IF EXISTS ux_account_transactions_transfer_pair;
DROP INDEX IF EXISTS idx_account_transactions_reference_type;
