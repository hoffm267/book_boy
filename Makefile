DEV_IMAGE_NAME=book_boy:dev
DEV_COMPOSE_FILE=docker-compose.dev.yml

PROD_IMAGE_NAME=book_boy:prod
PROD_COMPOSE_FILE=docker-compose.yml

.PHONY: build-dev build-prod run-dev run-prod clean

build-dev:
	docker compose -f $(DEV_COMPOSE_FILE) build --no-cache

build-prod:
	docker compose -f $(PROD_COMPOSE_FILE) build --no-cache

run-dev:
	@if [ -z "$$(docker images -q $(DEV_IMAGE_NAME))" ]; then \
		echo "Dev image '$(DEV_IMAGE_NAME)' not found. Building..."; \
		make build-dev; \
	fi
	docker compose -f $(DEV_COMPOSE_FILE) run --rm book_boy $(ARGS)

run-prod:
	@if [ -z "$$(docker images -q $(PROD_IMAGE_NAME))" ]; then \
		echo "Prod image '$(PROD_IMAGE_NAME)' not found. Building..."; \
		make build-prod; \
	fi
	docker compose -f $(PROD_COMPOSE_FILE) run --rm book_boy $(ARGS)

clean:
	@if [ -n "$$(docker images -aq)" ]; then \
		echo "Removing images..."; \
		docker rmi -f $$(docker images -aq); \
	else \
		echo "No images to remove."; \
	fi