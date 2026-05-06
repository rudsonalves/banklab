ALTER TABLE accounts
    DROP CONSTRAINT IF EXISTS accounts_number_key;

ALTER TABLE accounts
    ADD CONSTRAINT accounts_branch_number_key UNIQUE (branch, number);
