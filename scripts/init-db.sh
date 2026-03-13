#!/bin/bash
set -e

# Database initialization script for Docker

# Wait for PostgreSQL to be ready
until pg_isready -h localhost -p 5432 -U "$POSTGRES_USER"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "PostgreSQL is ready!"

# Create additional databases if needed
# psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
#     CREATE DATABASE pwvc_test;
#     GRANT ALL PRIVILEGES ON DATABASE pwvc_test TO $POSTGRES_USER;
# EOSQL

echo "Database initialization completed!"