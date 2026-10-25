include .env
export
export PROJECT_ROOT=$(shell pwd)

DB_SERVICES   := chatping-postgres chatping-redis
PORT_SERVICES := postgres-port-forwarder redis-port-forwarder

env-up:
	@docker compose up -d $(DB_SERVICES)

env-down:
	@docker compose down $(DB_SERVICES) $(PORT_SERVICES)

env-cleanup:
	@read -p "Clear all environment volume files? Warning: Data will be lost! [y/N]: " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down $(DB_SERVICES) $(PORT_SERVICES) && \
		rm -rf "${PROJECT_ROOT}/out/pgdata" && \
		rm -rf "${PROJECT_ROOT}/out/redisdata" && \
		echo "Environment volume files cleared"; \
	else \
		echo "Cleanup aborted."; \
	fi

env-port-forward:
	@docker compose up -d $(PORT_SERVICES)

env-port-close:
	@docker compose down $(PORT_SERVICES)

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing required parameter 'seq'. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm chatping-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter 'action'. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm chatping-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@chatping-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

logs-cleanup:
	@read -p "Clear all log files? Warning: Data will be lost! [y/N]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf "${PROJECT_ROOT}/out/logs" && \
		echo "Log files cleared"; \
	else \
		echo "Cleanup aborted."; \
	fi

run: fmt vet staticcheck
	@export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs" && \
	export POSTGRES_HOST=localhost && \
	export REDIS_HOST=localhost && \
	go mod tidy && \
	go run "${PROJECT_ROOT}/cmd/chatping/main.go"

fmt:
	@go fmt ./...

vet:
	@go vet ./...

staticcheck:
	@staticcheck ./...