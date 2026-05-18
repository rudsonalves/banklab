INSERT INTO customer_documents (
    id,
    customer_id,
    type,
    value,
    country,
    is_primary,
    created_at,
    updated_at
)
SELECT
    gen_random_uuid(),
    c.id,
    'cpf',
    c.cpf,
    'BR',
    true,
    c.created_at,
    NOW()
FROM customers c
WHERE c.cpf IS NOT NULL
  AND NOT EXISTS (
        SELECT 1
        FROM customer_documents cd
        WHERE cd.type = 'cpf'
            AND cd.value = c.cpf
            AND cd.country = 'BR'
    );