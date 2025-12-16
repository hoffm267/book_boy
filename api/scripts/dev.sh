#!/bin/sh

cd "$(dirname "$0")/.."

make dev

until docker exec $(docker compose ps -q db) pg_isready -U postgres > /dev/null 2>&1; do
  sleep 1
done

go run ./cmd/server/main.go

