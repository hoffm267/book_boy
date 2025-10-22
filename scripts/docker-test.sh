#!/bin/sh

cd "$(dirname "$0")/.."

echo "Running unit tests..."
make test

if [ $? -ne 0 ]; then
    echo "Tests failed!"
    exit 1
fi

echo ""
echo "Building Docker image and starting services..."
make docker-test

if [ $? -ne 0 ]; then
    echo "Docker build failed!"
    exit 1
fi

echo ""
echo "Waiting for services to be ready..."
sleep 10

echo ""
echo "Testing API health..."

HEALTH_CHECK=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null || echo "000")

if [ "$HEALTH_CHECK" = "200" ]; then
    echo "API is healthy (HTTP $HEALTH_CHECK)"
else
    echo "API not responding (HTTP $HEALTH_CHECK)"
    echo ""
    echo "Checking logs..."
    docker compose logs backend-docker | tail -20
    exit 1
fi

echo ""
echo "Docker containers running:"
docker compose ps
echo ""
echo "All tests passed and API is running!"
echo ""
echo "API available at: http://localhost:8080"
echo ""
echo "Commands:"
echo "  make logs  - View backend logs"
echo "  make down  - Stop containers"
