ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS cpf VARCHAR(11);

UPDATE customers c
SET cpf = cd.value
FROM customer_documents cd
WHERE cd.customer_id = c.id
    AND cd.type = 'cpf'
    AND cd.country = 'BR'
    AND cd.is_primary = true
    AND c.cpf IS NULL;

ALTER TABLE customers
    DROP CONSTRAINT IF EXISTS chk_cpf_format,
    ADD CONSTRAINT chk_cpf_format CHECK (cpf ~ '^\\d{11}$');

ALTER TABLE customers
    DROP CONSTRAINT IF EXISTS customers_cpf_key,
    ADD CONSTRAINT customers_cpf_key UNIQUE (cpf);