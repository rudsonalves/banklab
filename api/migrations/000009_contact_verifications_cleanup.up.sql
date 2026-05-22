CREATE EXTENSION IF NOT EXISTS pg_cron;

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

DROP INDEX IF EXISTS idx_contact_verifications_target_channel;

CREATE INDEX IF NOT EXISTS idx_contact_verifications_unverified_expires_at
ON contact_verifications (expires_at)
WHERE verified_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_contact_verifications_verified_at
ON contact_verifications (verified_at)
WHERE verified_at IS NOT NULL;

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

DO $$
DECLARE
    scheduled_job record;
BEGIN
    FOR scheduled_job IN
        SELECT jobid
        FROM cron.job
        WHERE jobname = 'cleanup-contact-verifications'
    LOOP
        PERFORM cron.unschedule(scheduled_job.jobid);
    END LOOP;
END;
$$;

SELECT cron.schedule(
    'cleanup-contact-verifications',
    '0 3 * * *',
    $$SELECT cleanup_contact_verifications();$$
);

SELECT cleanup_contact_verifications();
