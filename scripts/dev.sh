#!/bin/sh

cd "$(dirname "$0")/.."

make backend-dev

echo "Waiting for database to be ready..."
until docker exec $(docker compose ps -q db) pg_isready -U postgres > /dev/null 2>&1; do
  echo "Database is unavailable - sleeping"
  sleep 1
done

echo "Database is ready!"
echo "Starting backend server..."
go run ./cmd/server/main.go

