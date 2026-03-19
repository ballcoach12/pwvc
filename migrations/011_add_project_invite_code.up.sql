-- Add invite_code column for project invite links (T016 - US1)
ALTER TABLE projects ADD COLUMN invite_code VARCHAR(12) UNIQUE;