#!/bin/bash

echo "🗄️ Setting up P-WVC PostgreSQL database..."

# Start PostgreSQL container
docker run --name pwvc-postgres \
    -e POSTGRES_PASSWORD=password \
    -e POSTGRES_DB=pwvc \
    -e POSTGRES_USER=pwvc \
    -p 5432:5432 \
    -d postgres:15-alpine

echo "⏳ Waiting for database to start..."
sleep 5

# Check if database is ready
while ! docker exec pwvc-postgres pg_isready -U pwvc; do
    echo "Waiting for database connection..."
    sleep 2
done

echo "✅ Database is ready!"
echo "📍 Connection string: postgres://pwvc:password@localhost:5432/pwvc?sslmode=disable"