CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE transaction_passwords (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash   TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until    TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    changed_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_transaction_passwords_status
        CHECK (status IN ('active', 'blocked')),

    CONSTRAINT chk_transaction_passwords_failed_attempts
        CHECK (failed_attempts >= 0),

    CONSTRAINT chk_transaction_passwords_blocked_locked_until
        CHECK (status <> 'blocked' OR locked_until IS NOT NULL)
);

CREATE UNIQUE INDEX ux_transaction_passwords_user_id
ON transaction_passwords (user_id);

CREATE INDEX idx_transaction_passwords_locked_until
ON transaction_passwords (locked_until)
WHERE locked_until IS NOT NULL;