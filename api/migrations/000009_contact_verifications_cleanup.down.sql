DO $$
DECLARE
    scheduled_job record;
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_extension
        WHERE extname = 'pg_cron'
    ) THEN
        FOR scheduled_job IN
            SELECT jobid
            FROM cron.job
            WHERE jobname = 'cleanup-contact-verifications'
        LOOP
            PERFORM cron.unschedule(scheduled_job.jobid);
        END LOOP;
    END IF;
END;
$$;

DROP FUNCTION IF EXISTS cleanup_contact_verifications();

DROP INDEX IF EXISTS idx_contact_verifications_verified_at;

DROP INDEX IF EXISTS idx_contact_verifications_unverified_expires_at;

DROP INDEX IF EXISTS contact_verifications_unique_target_channel;

CREATE INDEX IF NOT EXISTS idx_contact_verifications_target_channel
ON contact_verifications (target, channel);
