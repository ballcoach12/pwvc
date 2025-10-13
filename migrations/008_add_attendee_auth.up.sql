-- Add authentication fields to attendees table
ALTER TABLE attendees ADD COLUMN email VARCHAR(255);
ALTER TABLE attendees ADD COLUMN pin VARCHAR(10);

-- Create index for email lookup
CREATE INDEX idx_attendees_email ON attendees(email);