ALTER TABLE customers ADD COLUMN email VARCHAR(120);

ALTER TABLE customers ADD CONSTRAINT customers_email_key UNIQUE (email);
