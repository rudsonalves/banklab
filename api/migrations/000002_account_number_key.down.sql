ALTER TABLE accounts
    DROP CONSTRAINT IF EXISTS accounts_branch_number_key;

ALTER TABLE accounts
    ADD CONSTRAINT accounts_number_key UNIQUE (number);
