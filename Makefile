include .env
export

export PROJECT_ROOT=$(shell pwd)


env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down data-base && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

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
	docker compose run --rm migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@data-base:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

clean_migrate:
	@docker compose run --rm migrate \
		-path /migrations \
     	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@data-base:5432/${POSTGRES_DB}?sslmode=disable force 0

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

docker-start:
	@docker compose up -d

docker-down:
	@docker compose down

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/trackerapp/main.go