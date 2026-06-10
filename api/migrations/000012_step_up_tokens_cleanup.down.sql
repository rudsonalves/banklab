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
            WHERE jobname = 'cleanup-step-up-tokens'
        LOOP
            PERFORM cron.unschedule(scheduled_job.jobid);
        END LOOP;
    END IF;
END;
$$;

DROP FUNCTION IF EXISTS cleanup_step_up_tokens();

DROP INDEX IF EXISTS idx_step_up_tokens_consumed_at;

DROP INDEX IF EXISTS idx_step_up_tokens_active_expires_at;

CREATE INDEX IF NOT EXISTS idx_step_up_tokens_expires_at
ON step_up_tokens (expires_at);
