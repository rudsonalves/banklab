DROP INDEX IF EXISTS customer_addresses_one_primary_per_customer;
DROP TABLE IF EXISTS customer_addresses;

DROP INDEX IF EXISTS customer_documents_one_primary_per_customer;
DROP TABLE IF EXISTS customer_documents;

ALTER TABLE customers
    DROP COLUMN IF EXISTS birth_date;

ALTER TABLE users
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS email_verified_at,
    DROP COLUMN IF EXISTS phone_verified_at;