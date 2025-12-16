.PHONY: clean

clean:
	docker compose --profile docker-test --profile db down -v --remove-orphans
	@docker compose rm -f 2>/dev/null || true
	@docker volume prune -f
	@docker image prune -af
	@docker system prune -af
