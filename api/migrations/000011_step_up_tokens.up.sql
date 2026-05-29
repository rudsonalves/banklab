CREATE TABLE step_up_tokens (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jti          VARCHAR(120) NOT NULL,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    endpoint_key VARCHAR(120) NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'active',
    expires_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    consumed_at  TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_step_up_tokens_status
        CHECK (status IN ('active', 'consumed')),

    CONSTRAINT chk_step_up_tokens_jti_not_blank
        CHECK (length(trim(jti)) > 0),

    CONSTRAINT chk_step_up_tokens_endpoint_key_not_blank
        CHECK (length(trim(endpoint_key)) > 0),

    CONSTRAINT chk_step_up_tokens_expires_after_created
        CHECK (expires_at > created_at),

    CONSTRAINT chk_step_up_tokens_consumed_at_consistency
        CHECK (
            (status = 'active' AND consumed_at IS NULL)
            OR
            (status = 'consumed' AND consumed_at IS NOT NULL)
        )
);

CREATE UNIQUE INDEX ux_step_up_tokens_jti
ON step_up_tokens (jti);

CREATE INDEX idx_step_up_tokens_user_id
ON step_up_tokens (user_id);

CREATE INDEX idx_step_up_tokens_endpoint_key
ON step_up_tokens (endpoint_key);

CREATE INDEX idx_step_up_tokens_expires_at
ON step_up_tokens (expires_at);