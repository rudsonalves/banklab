CREATE EXTENSION IF NOT EXISTS pg_cron;

CREATE INDEX IF NOT EXISTS idx_installation_registration_authorizations_active_expires_at
ON installation_registration_authorizations (expires_at)
WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_installation_registration_authorizations_consumed_at
ON installation_registration_authorizations (consumed_at)
WHERE status = 'consumed';

CREATE INDEX IF NOT EXISTS idx_installation_registration_authorizations_revoked_created_at
ON installation_registration_authorizations (created_at)
WHERE status = 'revoked';

CREATE OR REPLACE FUNCTION cleanup_installation_registration_authorizations()
RETURNS integer AS $$
DECLARE
    deleted_count integer;
BEGIN
    DELETE FROM installation_registration_authorizations
    WHERE (status = 'active' AND expires_at < NOW() - INTERVAL '24 hours')
       OR (status = 'consumed' AND consumed_at < NOW() - INTERVAL '24 hours')
       OR (status = 'revoked' AND created_at < NOW() - INTERVAL '24 hours');

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
        WHERE jobname = 'cleanup-installation-registration-authorizations'
    LOOP
        PERFORM cron.unschedule(scheduled_job.jobid);
    END LOOP;
END;
$$;

SELECT cron.schedule(
    'cleanup-installation-registration-authorizations',
    '30 3 * * *',
    $$SELECT cleanup_installation_registration_authorizations();$$
);

SELECT cleanup_installation_registration_authorizations();
