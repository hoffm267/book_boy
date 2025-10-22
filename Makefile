.PHONY: backend-dev docker-test test db-shell db-logs logs down stop clean

backend-dev:
	docker compose --profile db up -d

test:
	JWT_SECRET=test-secret-for-unit-tests go test ./...

docker-test:
	docker compose --profile docker-test up -d --build

db-shell:
	docker exec -it $$(docker compose ps -q db) psql -U postgres

db-logs:
	docker compose logs -f db

logs:
	docker compose logs -f backend-docker

down:
	docker compose down

stop:
	docker compose down

clean:
	@echo "Cleaning up Docker resources..."
	@echo "Stopping and removing containers (all profiles)..."
	docker compose --profile docker-test --profile db down -v --remove-orphans
	@echo "Removing any remaining containers..."
	@docker compose rm -f 2>/dev/null || true
	@echo "Removing volumes..."
	@docker volume rm book_boy_db 2>/dev/null || true
	@echo "Removing images..."
	@docker rmi book_boy-backend-docker 2>/dev/null || true
	@docker rmi book_boy-postgres:14.1-alpine 2>/dev/null || true
	@docker rmi book_boy-prod-test 2>/dev/null || true
	@echo "Pruning system..."
	@docker system prune -af --filter "label=com.docker.compose.project=book_boy" 2>/dev/null || true
	@echo "Cleaned up all containers, volumes, and images"