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
            WHERE jobname = 'cleanup-installation-registration-authorizations'
        LOOP
            PERFORM cron.unschedule(scheduled_job.jobid);
        END LOOP;
    END IF;
END;
$$;

DROP FUNCTION IF EXISTS cleanup_installation_registration_authorizations();

DROP INDEX IF EXISTS idx_installation_registration_authorizations_revoked_created_at;

DROP INDEX IF EXISTS idx_installation_registration_authorizations_consumed_at;

DROP INDEX IF EXISTS idx_installation_registration_authorizations_active_expires_at;
