DROP TRIGGER IF EXISTS trg_cleanup_contact_verifications ON contact_verifications;

DROP FUNCTION IF EXISTS cleanup_contact_verifications_if_due();

DROP FUNCTION IF EXISTS cleanup_contact_verifications();

DROP TABLE IF EXISTS contact_verification_cleanup_runs;

DROP INDEX IF EXISTS idx_contact_verifications_verified_at;

DROP INDEX IF EXISTS idx_contact_verifications_unverified_expires_at;

DROP INDEX IF EXISTS contact_verifications_unique_target_channel;
