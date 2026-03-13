#!/bin/bash

# PairWise API Test Script
# This script tests all the basic API endpoints

BASE_URL="http://localhost:8080"

echo "🚀 Testing PairWise API endpoints..."
echo "================================="

# Test health check
echo "1. Testing health check..."
curl -s "$BASE_URL/health" | jq '.'
echo ""

# Test API info
echo "2. Testing API info..."
curl -s "$BASE_URL/" | jq '.'
echo ""

# Test create project
echo "3. Creating a new project..."
PROJECT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Website Redesign",
    "description": "A test project for API validation"
  }')

echo $PROJECT_RESPONSE | jq '.'
PROJECT_ID=$(echo $PROJECT_RESPONSE | jq -r '.id')
echo "Created project with ID: $PROJECT_ID"
echo ""

# Test get project
echo "4. Getting project details..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID" | jq '.'
echo ""

# Test list projects
echo "5. Listing all projects..."
curl -s "$BASE_URL/api/projects" | jq '.'
echo ""

# Test add attendee
echo "6. Adding attendee to project..."
ATTENDEE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/projects/$PROJECT_ID/attendees" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "role": "Product Manager",
    "is_facilitator": true
  }')

echo $ATTENDEE_RESPONSE | jq '.'
ATTENDEE_ID=$(echo $ATTENDEE_RESPONSE | jq -r '.id')
echo "Created attendee with ID: $ATTENDEE_ID"
echo ""

# Test add another attendee
echo "7. Adding another attendee..."
curl -s -X POST "$BASE_URL/api/projects/$PROJECT_ID/attendees" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "role": "Developer",
    "is_facilitator": false
  }' | jq '.'
echo ""

# Test get project attendees
echo "8. Getting project attendees..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/attendees" | jq '.'
echo ""

# Test update project
echo "9. Updating project..."
curl -s -X PUT "$BASE_URL/api/projects/$PROJECT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Website Redesign",
    "description": "Updated description for the test project",
    "status": "active"
  }' | jq '.'
echo ""

# Test delete attendee
echo "10. Deleting attendee..."
curl -s -X DELETE "$BASE_URL/api/projects/$PROJECT_ID/attendees/$ATTENDEE_ID"
echo "Deleted attendee with ID: $ATTENDEE_ID"
echo ""

# Test get attendees after deletion
echo "11. Getting attendees after deletion..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/attendees" | jq '.'
echo ""

# Test delete project
echo "12. Deleting project..."
curl -s -X DELETE "$BASE_URL/api/projects/$PROJECT_ID"
echo "Deleted project with ID: $PROJECT_ID"
echo ""

# Test get project after deletion (should return 404)
echo "13. Trying to get deleted project (should return 404)..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID" | jq '.'
echo ""

echo "✅ API testing completed!"
echo "================================="