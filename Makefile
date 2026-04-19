include .env
export

export PROJECT_ROOT=$(shell pwd)

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Dont have seq"; \
		exit 1; \
	fi; \
		docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Dont have action"; \
		exit 1; \
	fi; \
	docker compose run --rm \
		-path /migrations/ \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@data-base:5432/${POSTGRES_DB}?sslmode=disable \
		- "$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

docker-start:
	@docker compose up -d

docker-down:
	@docker compose down