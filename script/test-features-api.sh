#!/bin/bash

# P-WVC Feature Management API Test Script
# This script tests all the feature management endpoints

BASE_URL="http://localhost:8080"

echo "🚀 Testing P-WVC Feature Management API endpoints..."
echo "=================================================="

# Test health check
echo "1. Testing health check..."
curl -s "$BASE_URL/health" | jq '.'
echo ""

# Create a test project first
echo "2. Creating a test project for features..."
PROJECT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Feature Test Project",
    "description": "A test project for feature API validation"
  }')

echo $PROJECT_RESPONSE | jq '.'
PROJECT_ID=$(echo $PROJECT_RESPONSE | jq -r '.id')
echo "Created project with ID: $PROJECT_ID"
echo ""

# Test create feature
echo "3. Creating a new feature..."
FEATURE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/projects/$PROJECT_ID/features" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "User Authentication",
    "description": "Implement secure user login and registration functionality with email verification and password reset capabilities.",
    "acceptance_criteria": "Given a user with valid credentials, when they login, then they should be authenticated and redirected to dashboard. Given invalid credentials, then appropriate error message should be displayed."
  }')

echo $FEATURE_RESPONSE | jq '.'
FEATURE_ID=$(echo $FEATURE_RESPONSE | jq -r '.id')
echo "Created feature with ID: $FEATURE_ID"
echo ""

# Test create another feature
echo "4. Creating another feature..."
curl -s -X POST "$BASE_URL/api/projects/$PROJECT_ID/features" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Dashboard Analytics",
    "description": "Create a comprehensive dashboard showing key performance metrics and user activity charts.",
    "acceptance_criteria": "Dashboard should load within 2 seconds and display real-time data updates."
  }' | jq '.'
echo ""

# Test get specific feature
echo "5. Getting specific feature details..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/features/$FEATURE_ID" | jq '.'
echo ""

# Test list all features in project
echo "6. Listing all features in project..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/features" | jq '.'
echo ""

# Test update feature
echo "7. Updating feature..."
curl -s -X PUT "$BASE_URL/api/projects/$PROJECT_ID/features/$FEATURE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Enhanced User Authentication",
    "description": "Implement secure user login and registration functionality with email verification, password reset capabilities, and two-factor authentication.",
    "acceptance_criteria": "Given a user with valid credentials, when they login, then they should be authenticated and redirected to dashboard. Given invalid credentials, then appropriate error message should be displayed. Users should be able to enable 2FA for additional security."
  }' | jq '.'
echo ""

# Test CSV export
echo "8. Testing CSV export..."
echo "Exporting features to CSV file..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/features/export" -o "features_export.csv"
echo "Features exported to features_export.csv"
cat features_export.csv
echo ""

# Create a sample CSV for import testing
echo "9. Creating sample CSV for import testing..."
cat > features_import.csv << 'EOF'
title,description,acceptance_criteria
"Search Functionality","Implement full-text search across all content with filters and sorting options","Search results should be relevant and load within 1 second"
"User Profile Management","Allow users to update their profile information including avatar upload","Profile changes should be saved immediately and reflected across the application"
"Notification System","Real-time notification system for important events and updates","Notifications should appear within 3 seconds of trigger event"
EOF

echo "Sample CSV created:"
cat features_import.csv
echo ""

# Test CSV import
echo "10. Testing CSV import..."
IMPORT_RESULT=$(curl -s -X POST "$BASE_URL/api/projects/$PROJECT_ID/features/import" \
  -F "file=@features_import.csv")
echo $IMPORT_RESULT | jq '.'
echo ""

# List features after import
echo "11. Listing features after CSV import..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/features" | jq '.'
echo ""

# Test delete feature
echo "12. Deleting a feature..."
curl -s -X DELETE "$BASE_URL/api/projects/$PROJECT_ID/features/$FEATURE_ID"
echo "Deleted feature with ID: $FEATURE_ID"
echo ""

# List features after deletion
echo "13. Listing features after deletion..."
curl -s "$BASE_URL/api/projects/$PROJECT_ID/features" | jq '.'
echo ""

# Test error handling - invalid project ID
echo "14. Testing error handling with invalid project ID..."
curl -s "$BASE_URL/api/projects/999999/features" | jq '.'
echo ""

# Test error handling - invalid feature data
echo "15. Testing validation with invalid feature data..."
curl -s -X POST "$BASE_URL/api/projects/$PROJECT_ID/features" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "",
    "description": "Missing title should cause validation error"
  }' | jq '.'
echo ""

# Clean up - delete test project
echo "16. Cleaning up - deleting test project..."
curl -s -X DELETE "$BASE_URL/api/projects/$PROJECT_ID"
echo "Deleted test project"
echo ""

# Clean up files
rm -f features_export.csv features_import.csv

echo "✅ Feature Management API testing completed!"
echo "=================================================="