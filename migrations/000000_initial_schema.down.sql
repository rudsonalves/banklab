-- DROP INDEXES (optional, PostgreSQL drops with table, but explicit is cleaner)

DROP INDEX IF EXISTS ux_transactions_idempotency;
DROP INDEX IF EXISTS idx_transactions_reference_id;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_account_id;

DROP INDEX IF EXISTS idx_account_transactions_reference_id;
DROP INDEX IF EXISTS idx_account_transactions_account_created;

DROP INDEX IF EXISTS idx_accounts_customer_id;

-- DROP TRIGGER AND FUNCTION (before table)

DROP TRIGGER IF EXISTS trg_account_transactions_no_mutation ON account_transactions;
DROP FUNCTION IF EXISTS prevent_account_transactions_mutation;

-- DROP TABLES (children first)

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS account_transactions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS customers;

-- DROP ENUMS (after tables)

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS account_status;