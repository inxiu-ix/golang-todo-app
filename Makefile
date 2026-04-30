include .env
export

export PROJECT_ROOT := $(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres 

env-down:
	@docker compose down todoapp-postgres 

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных! (y/n): " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwrder && \
		rm -rf out/pgdata && \
		echo "Очистка окружения завершена"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwrder

env-port-close:
	@docker compose down port-forwrder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq is required"; \
		exit 1; \
	fi;
	
	docker compose run --rm todoapp-postgres-migrate \
		 create \
	 	 -ext sql \
	 -dir /migrations \
	 -seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "version is required, e.g: make migrate-force version=0"; \
		exit 1; \
	fi;
	docker compose run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	force $(version)

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "action is required"; \
		exit 1; \
	fi;

	docker compose run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	$(action)

run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go