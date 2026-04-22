-- MIGRATION: 000001_init_schema.down.sql
-- Tears down the full baseline schema in reverse dependency order.

-- ============================================================
-- TRIGGERS
-- ============================================================

DROP TRIGGER IF EXISTS trg_transactions_no_mutation ON transactions;

-- ============================================================
-- INDEXES (dropped automatically with tables, listed for clarity)
-- ============================================================

DROP INDEX IF EXISTS ux_transactions_transfer_pair;
DROP INDEX IF EXISTS ux_transactions_idempotency;
DROP INDEX IF EXISTS idx_transactions_reference_type;
DROP INDEX IF EXISTS idx_transactions_account_created;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_reference_id;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP INDEX IF EXISTS idx_user_sessions_expires_at;
DROP INDEX IF EXISTS idx_user_sessions_user_id;
DROP INDEX IF EXISTS idx_accounts_customer_id;

-- ============================================================
-- TABLES (in reverse FK dependency order)
-- ============================================================

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS customers;

-- ============================================================
-- FUNCTIONS
-- ============================================================

DROP FUNCTION IF EXISTS prevent_transactions_mutation();

-- ============================================================
-- SEQUENCES
-- ============================================================

DROP SEQUENCE IF EXISTS account_number_seq;

-- ============================================================
-- ENUMS
-- ============================================================

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS account_status;
