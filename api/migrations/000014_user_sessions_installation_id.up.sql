ALTER TABLE user_sessions
    ADD COLUMN installation_id UUID;

CREATE INDEX idx_user_sessions_user_installation
    ON user_sessions(user_id, installation_id)
    WHERE installation_id IS NOT NULL;
