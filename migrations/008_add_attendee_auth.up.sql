-- Add authentication fields to attendees table
ALTER TABLE attendees ADD COLUMN email VARCHAR(255);
ALTER TABLE attendees ADD COLUMN pin_hash VARCHAR(64);  -- For SHA256 hex string

-- Create index for email lookup
CREATE INDEX idx_attendees_email ON attendees(email);

-- Set test data for existing attendees
-- Attendee ID 10: PIN 1234 (hash: 03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4)
UPDATE attendees SET 
    email = 'usera@test.com',
    pin_hash = '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4'
WHERE id = 10;

-- Attendee ID 11: PIN 5678 (hash: 4a4133c42ba991494840c336e768653654c7bb60fbbc0936a46b0b6d668553)
UPDATE attendees SET 
    email = 'userb@test.com',
    pin_hash = '4a4133c42ba991494840c336e768653654c7bb60fbbc0936a46b0b6d668553'
WHERE id = 11;