-- Remove invite token fields
DROP INDEX IF EXISTS idx_attendees_invite_token;
ALTER TABLE attendees DROP COLUMN IF EXISTS invite_token_expires_at;
ALTER TABLE attendees DROP COLUMN IF EXISTS invite_token;