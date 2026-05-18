ALTER TABLE users
    ADD COLUMN phone VARCHAR(20) UNIQUE,
    ADD COLUMN email_verified_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN phone_verified_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE customers
    ADD COLUMN birth_date DATE;

CREATE TABLE customer_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    value VARCHAR(80) NOT NULL,
    issuer VARCHAR(80),
    issuer_state VARCHAR(30),
    country CHAR(2) NOT NULL DEFAULT 'BR',
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_documents_unique_document
        UNIQUE (type, value, country)
);

CREATE UNIQUE INDEX customer_documents_one_primary_per_customer
ON customer_documents (customer_id)
WHERE is_primary = true;

CREATE TABLE customer_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    postal_code VARCHAR(20) NOT NULL,
    number VARCHAR(30) NOT NULL,
    neighborhood VARCHAR(120),
    city VARCHAR(120) NOT NULL,
    state VARCHAR(60) NOT NULL,
    street VARCHAR(160) NOT NULL,
    complement VARCHAR(120),
    country CHAR(2) NOT NULL DEFAULT 'BR',
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX customer_addresses_one_primary_per_customer
ON customer_addresses (customer_id)
WHERE is_primary = true;