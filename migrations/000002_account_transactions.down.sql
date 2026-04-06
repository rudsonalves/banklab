-- DROP TRIGGER AND FUNCTION

DROP TRIGGER IF EXISTS trg_account_transactions_no_mutation ON account_transactions;
DROP FUNCTION IF EXISTS prevent_account_transactions_mutation;

-- DROP INDEXES

DROP INDEX IF EXISTS idx_account_transactions_account_created;
DROP INDEX IF EXISTS idx_account_transactions_reference_id;
DROP INDEX IF EXISTS idx_account_transactions_account_id;

-- DROP TABLE

DROP TABLE IF EXISTS account_transactions;