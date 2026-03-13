#!/bin/bash

echo "Testing if a vote was created..."

# Check if the vote we created earlier is still there
curl -s http://localhost:3001/api/projects/8/pairwise/comparisons | jq '.comparisons[0].votes'

echo -e "\n\nTesting vote creation from frontend..."
echo "The frontend should now work with:"
echo "1. Auto-selected attendee: User A (ID: 10)"
echo "2. Available comparison: Feature X vs Feature Y (ID: 8)" 
echo "3. Vote buttons should call the API correctly"
echo ""
echo "Debug logs should show in browser console when clicking vote buttons"