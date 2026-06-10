CREATE EXTENSION IF NOT EXISTS pg_cron;

DROP INDEX IF EXISTS idx_step_up_tokens_expires_at;

CREATE INDEX IF NOT EXISTS idx_step_up_tokens_active_expires_at
ON step_up_tokens (expires_at)
WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_step_up_tokens_consumed_at
ON step_up_tokens (consumed_at)
WHERE status = 'consumed';

CREATE OR REPLACE FUNCTION cleanup_step_up_tokens()
RETURNS integer AS $$
DECLARE
    deleted_count integer;
BEGIN
    DELETE FROM step_up_tokens
    WHERE (status = 'active' AND expires_at < NOW() - INTERVAL '24 hours')
       OR (status = 'consumed' AND consumed_at < NOW() - INTERVAL '24 hours');

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
        WHERE jobname = 'cleanup-step-up-tokens'
    LOOP
        PERFORM cron.unschedule(scheduled_job.jobid);
    END LOOP;
END;
$$;

SELECT cron.schedule(
    'cleanup-step-up-tokens',
    '15 3 * * *',
    $$SELECT cleanup_step_up_tokens();$$
);

SELECT cleanup_step_up_tokens();
