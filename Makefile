.PHONY: clean

clean:
	docker compose --profile docker-test --profile db down -v --remove-orphans
	@docker compose rm -f 2>/dev/null || true
	@docker volume rm book_boy_db_data 2>/dev/null || true
	@docker volume prune -f
	@docker system prune -af --filter "label=com.docker.compose.project=book_boy" 2>/dev/null || true
