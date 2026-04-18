UPDATE customers c
SET email = COALESCE(u.email, CONCAT(c.id::text, '@rollback.local'))
FROM users u
WHERE c.email IS NULL
  AND u.customer_id = c.id;

UPDATE customers
SET email = CONCAT(id::text, '@rollback.local')
WHERE email IS NULL;

ALTER TABLE customers
ALTER COLUMN email SET NOT NULL;
