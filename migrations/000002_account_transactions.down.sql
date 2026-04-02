DROP TRIGGER IF EXISTS trg_account_transactions_no_mutation ON account_transactions;
DROP FUNCTION IF EXISTS prevent_account_transactions_mutation();
DROP TABLE IF EXISTS account_transactions;
