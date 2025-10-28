#!/bin/bash

echo "Testing vote submission API..."

# First, get the current session and comparisons
echo "1. Getting current pairwise session..."
curl -s http://localhost:3001/api/projects/8/pairwise | jq '.'

echo -e "\n2. Getting comparisons..."
COMPARISONS=$(curl -s http://localhost:3001/api/projects/8/pairwise/comparisons)
echo $COMPARISONS | jq '.'

# Extract the first comparison ID
COMPARISON_ID=$(echo $COMPARISONS | jq -r '.comparisons[0].comparison.id')
echo -e "\n3. Using comparison ID: $COMPARISON_ID"

# Get attendees
echo -e "\n4. Getting attendees..."
ATTENDEES=$(curl -s http://localhost:3001/api/projects/8/attendees)
echo $ATTENDEES | jq '.'

# Extract the first attendee ID
ATTENDEE_ID=$(echo $ATTENDEES | jq -r '.attendees[0].id')
echo -e "\n5. Using attendee ID: $ATTENDEE_ID"

# Submit a test vote (choosing Feature A)
echo -e "\n6. Submitting vote..."
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{\"comparison_id\": $COMPARISON_ID, \"attendee_id\": $ATTENDEE_ID, \"preferred_feature_id\": 20, \"is_tie_vote\": false}" \
  "http://localhost:3001/api/projects/8/pairwise/votes" \
  | jq '.'

echo -e "\n7. Getting updated comparison to verify vote..."
curl -s "http://localhost:3001/api/projects/8/pairwise/comparisons" | jq '.'