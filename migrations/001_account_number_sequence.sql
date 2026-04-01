CREATE SEQUENCE account_number_seq START 10000000;

ALTER TABLE accounts
ADD CONSTRAINT accounts_number_key UNIQUE (number);