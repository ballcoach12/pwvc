#!/bin/bash

echo "🗄️ Setting up PairWise PostgreSQL database..."

# Start PostgreSQL container
docker run --name pairwise-postgres \
    -e POSTGRES_PASSWORD=password \
    -e POSTGRES_DB=pairwise \
    -e POSTGRES_USER=pairwise \
    -p 5432:5432 \
    -d postgres:15-alpine

echo "⏳ Waiting for database to start..."
sleep 5

# Check if database is ready
while ! docker exec pairwise-postgres pg_isready -U pairwise; do
    echo "Waiting for database connection..."
    sleep 2
done

echo "✅ Database is ready!"
echo "📍 Connection string: postgres://pairwise:password@localhost:5432/pairwise?sslmode=disable"