#!/bin/bash

echo "Setting up test authentication for attendees..."

# First, run the migration to add auth fields
echo "Running migration..."
cd /workspaces/pwvc
docker-compose exec backend sh -c "cd /home/pairwise && sqlite3 data/pairwise.db < migrations/008_add_attendee_auth.up.sql" 2>/dev/null || echo "Migration may have already run"

# Add simple PINs for testing (1234 for User A, 5678 for User B)
echo "Adding test PINs..."

# Hash for PIN "1234" 
PIN_1234="03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4"
# Hash for PIN "5678"
PIN_5678="ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f"

# Update attendees with PINs and emails
curl -X POST http://localhost:3001/api/projects/8/attendees -H "Content-Type: application/json" -d '{
  "name": "Test Update User A",
  "role": "Tester", 
  "is_facilitator": false
}' > /dev/null 2>&1

# Use SQL to directly update the existing attendees for testing
docker-compose exec backend sh -c "
cd /home/pairwise && sqlite3 data/pairwise.db \"
UPDATE attendees SET email='usera@test.com', pin='$PIN_1234' WHERE id=10;
UPDATE attendees SET email='userb@test.com', pin='$PIN_5678' WHERE id=11;
SELECT 'Updated attendees:' as message;
SELECT id, name, role, email FROM attendees WHERE project_id=8;
\""

echo ""
echo "Test credentials created:"
echo "User A - PIN: 1234 (email: usera@test.com)"  
echo "User B - PIN: 5678 (email: userb@test.com)"
echo ""
echo "Now rebuild the backend to apply the auth API changes."