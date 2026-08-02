CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE app_installations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id     UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    installation_id UUID NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'known',
    known_slot      SMALLINT,
    platform        VARCHAR(40),
    app_version     VARCHAR(40),
    app_build       VARCHAR(40),
    first_seen_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_seen_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    revoked_at      TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_app_installations_status
        CHECK (status IN ('known', 'revoked')),

    CONSTRAINT chk_app_installations_timestamps
        CHECK (
            last_seen_at >= first_seen_at
            AND updated_at >= created_at
        ),

    CONSTRAINT chk_app_installations_revoked_at_consistency
        CHECK (
            (status = 'known' AND revoked_at IS NULL)
            OR
            (status = 'revoked' AND revoked_at IS NOT NULL)
        ),

    CONSTRAINT chk_app_installations_known_slot_consistency
        CHECK (
            (status = 'known' AND known_slot BETWEEN 1 AND 3)
            OR
            (status = 'revoked' AND known_slot IS NULL)
        )
);

COMMENT ON COLUMN app_installations.resource_id IS
    'Public management identifier. It must not expose the client-provided installation_id.';

COMMENT ON COLUMN app_installations.known_slot IS
    'Auxiliary slot used to enforce at most three known installations per user. Repositories must reserve it transactionally.';

CREATE UNIQUE INDEX ux_app_installations_resource_id
ON app_installations (resource_id);

CREATE UNIQUE INDEX ux_app_installations_user_installation_id
ON app_installations (user_id, installation_id);

CREATE UNIQUE INDEX ux_app_installations_known_slot
ON app_installations (user_id, known_slot)
WHERE status = 'known';

CREATE INDEX idx_app_installations_user_id
ON app_installations (user_id);

CREATE INDEX idx_app_installations_user_resource_id
ON app_installations (user_id, resource_id);

CREATE INDEX idx_app_installations_user_status
ON app_installations (user_id, status);

CREATE INDEX idx_app_installations_known_count
ON app_installations (user_id)
WHERE status = 'known';

CREATE TABLE installation_registration_authorizations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jti             VARCHAR(120) NOT NULL,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    installation_id UUID NOT NULL,
    scope           VARCHAR(80) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    expires_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    consumed_at     TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_installation_registration_authorizations_jti_not_blank
        CHECK (length(trim(jti)) > 0),

    CONSTRAINT chk_installation_registration_authorizations_scope
        CHECK (scope = 'installation.register'),

    CONSTRAINT chk_installation_registration_authorizations_status
        CHECK (status IN ('active', 'consumed', 'revoked')),

    CONSTRAINT chk_installation_registration_authorizations_expires_after_created
        CHECK (expires_at > created_at),

    CONSTRAINT chk_installation_registration_authorizations_consumed_at_consistency
        CHECK (
            (status = 'active' AND consumed_at IS NULL)
            OR
            (status = 'consumed' AND consumed_at IS NOT NULL)
            OR
            (status = 'revoked')
        )
);

CREATE UNIQUE INDEX ux_installation_registration_authorizations_jti
ON installation_registration_authorizations (jti);

CREATE UNIQUE INDEX ux_installation_registration_authorizations_active
ON installation_registration_authorizations (user_id, installation_id, scope)
WHERE status = 'active';

CREATE INDEX idx_installation_registration_authorizations_user_installation_scope_status
ON installation_registration_authorizations (user_id, installation_id, scope, status);

CREATE INDEX idx_installation_registration_authorizations_expires_at
ON installation_registration_authorizations (expires_at);
