BASE_COMPOSE_FILE=docker-compose.yml

.PHONY: cli-dev cli-prod clean

cli-dev:
	ARGS="$(ARGS)" docker compose --profile cli_dev up

cli-prod:
	ARGS="$(ARGS)" docker compose --profile cli_prod up --build

clean:
	@if [ -n "$$(docker ps -aq)" ]; then \
		echo "Removing containers..."; \
		docker rm -vf $$(docker ps -aq) > /dev/null 2>&1; \
	else \
		echo "No containers to remove."; \
	fi
	@if [ -n "$$(docker images -aq)" ]; then \
		echo "Removing images..."; \
		docker rmi -f $$(docker images -aq) > /dev/null 2>&1; \
	else \
		echo "No images to remove."; \
	fi