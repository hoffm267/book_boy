#!/bin/sh

cd "$(dirname "$0")/.."

make test

if [ $? -ne 0 ]; then
    exit 1
fi

docker compose --profile docker-test up -d --build

if [ $? -ne 0 ]; then
    exit 1
fi

sleep 10

HEALTH_CHECK=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null || echo "000")

if [ "$HEALTH_CHECK" = "200" ]; then
    docker compose ps
else
    docker compose logs api-docker | tail -20
    exit 1
fi
