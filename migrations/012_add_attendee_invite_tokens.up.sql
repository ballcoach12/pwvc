-- Add invite token fields for PIN setup workflow
ALTER TABLE attendees ADD COLUMN invite_token VARCHAR(32);
ALTER TABLE attendees ADD COLUMN invite_token_expires_at TIMESTAMP;

-- Create index for invite token lookup
CREATE INDEX idx_attendees_invite_token ON attendees(invite_token);