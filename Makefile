.PHONY: build up start down remove stop restart logs logs-save migrate-create

ENV := .env
include $(ENV)

# App
run:
	go run ./cmd/main.go

# Migrations
DB_DSN := "host=$(DB_HOST) port=$(DB_HOST_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSL_MODE)"
MIG_DIR := ./mig

ifeq (migrate-create,$(firstword $(MAKECMDGOALS)))
  SQL_FILE := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(SQL_FILE):;@:)
endif

migrate-create:
	goose -dir $(MIG_DIR) create $(SQL_FILE) sql
migrate-status:
	goose -dir $(MIG_DIR) -v $(DB_DIALECT) $(DB_DSN) status
migrate-up:
	goose -dir $(MIG_DIR) -v $(DB_DIALECT) $(DB_DSN) up
migrate-down:
	goose -dir $(MIG_DIR) -v $(DB_DIALECT) $(DB_DSN) down
migrate-fix:
	goose -dir $(MIG_DIR) -v fix

# Docker
build:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV)  --log-level $(LOG_LEVEL) build
rebuild: clean
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV)  --log-level $(LOG_LEVEL) build --no-cache
up:
	sudo docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL) up
start:
	sudo docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL) start
down:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) down
clean:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) down --remove-orphans
stop:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) stop
restart:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) stop
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL) run --detach
logs:
	docker-compose --file $(DOCKER_COMPOSE_FILE) logs --follow
logs-save:
	docker-compose --file $(DOCKER_COMPOSE_FILE) logs >> docker.log