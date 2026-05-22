DELETE FROM contact_verifications cv
USING (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY target, channel
            ORDER BY created_at DESC, id DESC
        ) AS row_number
    FROM contact_verifications
) ranked
WHERE cv.id = ranked.id
  AND ranked.row_number > 1;

CREATE UNIQUE INDEX IF NOT EXISTS contact_verifications_unique_target_channel
ON contact_verifications (target, channel);

CREATE INDEX IF NOT EXISTS idx_contact_verifications_unverified_expires_at
ON contact_verifications (expires_at)
WHERE verified_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_contact_verifications_verified_at
ON contact_verifications (verified_at)
WHERE verified_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS contact_verification_cleanup_runs (
    name TEXT PRIMARY KEY,
    last_run_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE OR REPLACE FUNCTION cleanup_contact_verifications()
RETURNS integer AS $$
DECLARE
    deleted_count integer;
BEGIN
    DELETE FROM contact_verifications
    WHERE (verified_at IS NULL AND expires_at < NOW() - INTERVAL '24 hours')
       OR (verified_at IS NOT NULL AND verified_at < NOW() - INTERVAL '7 days');

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION cleanup_contact_verifications_if_due()
RETURNS trigger AS $$
DECLARE
    last_run_at TIMESTAMP WITH TIME ZONE;
BEGIN
    IF NOT pg_try_advisory_xact_lock(hashtext('contact_verifications_cleanup')) THEN
        RETURN NULL;
    END IF;

    SELECT c.last_run_at
    INTO last_run_at
    FROM contact_verification_cleanup_runs c
    WHERE c.name = 'contact_verifications';

    IF last_run_at IS NOT NULL AND last_run_at > NOW() - INTERVAL '1 day' THEN
        RETURN NULL;
    END IF;

    PERFORM cleanup_contact_verifications();

    INSERT INTO contact_verification_cleanup_runs (name, last_run_at)
    VALUES ('contact_verifications', NOW())
    ON CONFLICT (name)
    DO UPDATE SET last_run_at = EXCLUDED.last_run_at;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_cleanup_contact_verifications
AFTER INSERT ON contact_verifications
FOR EACH STATEMENT
EXECUTE FUNCTION cleanup_contact_verifications_if_due();

SELECT cleanup_contact_verifications();

INSERT INTO contact_verification_cleanup_runs (name, last_run_at)
VALUES ('contact_verifications', NOW())
ON CONFLICT (name)
DO UPDATE SET last_run_at = EXCLUDED.last_run_at;
