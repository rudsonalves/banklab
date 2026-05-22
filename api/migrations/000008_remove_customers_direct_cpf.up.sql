ALTER TABLE customers
    DROP CONSTRAINT IF EXISTS chk_cpf_format,
    DROP CONSTRAINT IF EXISTS customers_cpf_key,
    DROP COLUMN IF EXISTS cpf;