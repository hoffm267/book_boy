BASE_COMPOSE_FILE=docker-compose.yml

.PHONY: cli-dev cli-prod backend_prod clean

#docker build --tag docker-gs-ping .
#docker run -p 8080:8080 docker-gs-ping

cli-dev:
	ARGS="$(ARGS)" docker compose --profile cli_dev up

cli-prod:
	ARGS="$(ARGS)" docker compose --profile cli_prod up

backend-dev:
	docker compose --profile backend_dev up -d

backend-prod:
	docker compose --profile backend_prod up -d

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