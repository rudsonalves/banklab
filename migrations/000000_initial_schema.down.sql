-- DROP INDEXES

DROP INDEX IF EXISTS ux_transactions_idempotency;
DROP INDEX IF EXISTS idx_transactions_reference_id;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_account_id;

DROP INDEX IF EXISTS idx_accounts_customer_id;

-- DROP TABLES

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS customers;

-- DROP ENUMS

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS account_status;