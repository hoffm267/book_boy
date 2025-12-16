.PHONY: dev test db-shell clean

dev:
	docker compose --profile db up -d

test:
	JWT_SECRET=test-secret-for-unit-tests go test ./...

db-shell:
	docker exec -it $$(docker compose ps -q db) psql -U postgres

clean:
	docker compose --profile docker-test --profile db down -v --remove-orphans
	@docker compose rm -f 2>/dev/null || true
	@docker volume rm book_boy_db 2>/dev/null || true
	@docker system prune -af --filter "label=com.docker.compose.project=book_boy" 2>/dev/null || true