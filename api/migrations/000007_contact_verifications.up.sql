CREATE TABLE contact_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel VARCHAR(20) NOT NULL,
    target VARCHAR(160) NOT NULL,
    token VARCHAR(20) NOT NULL,
    verification_token VARCHAR(100),
    verified_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_contact_verifications_channel
        CHECK (channel IN ('email', 'phone'))
);

CREATE INDEX idx_contact_verifications_target_channel
ON contact_verifications (target, channel);

CREATE UNIQUE INDEX contact_verifications_unique_verification_token
ON contact_verifications (verification_token)
WHERE verification_token IS NOT NULL;
