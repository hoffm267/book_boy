IMAGE_NAME=book_boy
DEV_COMPOSE_FILE=docker-compose.dev.yml
PROD_COMPOSE_FILE=docker-compose.yml

.PHONY: build-dev build-prod run-dev run-prod

build-dev:
	docker-compose -f $(DEV_COMPOSE_FILE) build

build-prod:
	docker-compose -f $(PROD_COMPOSE_FILE) build

run-dev:
	@if [ -z "$$(docker images -q $(IMAGE_NAME))" ]; then \
		echo "Dev image '$(IMAGE_NAME)' not found. Building..."; \
		make build-dev; \
	fi
	docker-compose -f $(DEV_COMPOSE_FILE) run --rm book_boy $(ARGS)

run-prod:
	@if [ -z "$$(docker images -q $(IMAGE_NAME))" ]; then \
		echo "Prod image '$(IMAGE_NAME)' not found. Building..."; \
		make build-prod; \
	fi
	docker-compose -f $(PROD_COMPOSE_FILE) run --rm book_boy $(ARGS)

rebuild:
	docker-compose -f docker-compose.yml build --no-cache

inspect:
	docker image inspect book_boy | jq '.[0].Config.Entrypoint'

