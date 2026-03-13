-- Remove authentication fields from attendees table
ALTER TABLE attendees DROP COLUMN email;
ALTER TABLE attendees DROP COLUMN pin;

-- Drop index
DROP INDEX idx_attendees_email;