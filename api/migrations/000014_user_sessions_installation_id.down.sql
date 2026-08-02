DROP INDEX IF EXISTS idx_user_sessions_user_installation;

ALTER TABLE user_sessions
    DROP COLUMN IF EXISTS installation_id;
