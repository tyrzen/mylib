THIS_FILE := $(lastword $(MAKEFILE_LIST))
DOCKER_COMPOSE_FILE := "docker-comose.yml"
ENV := ".env"
LOG_LEVEL := "DEBUG"
.PHONY: build up start down remove stop restart logs logs-save

build:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV)  --log-level $(LOG_LEVEL) build
up:
	sudo docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL) up
start:
	sudo docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL) start
down:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) down
remove:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) down --remove-orphans
stop:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) stop
restart:
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --log-level $(LOG_LEVEL) stop
	docker-compose --verbose --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL) run --detach
logs:
	docker-compose --file $(DOCKER_COMPOSE_FILE) logs --follow
logs-save:
	docker-compose --file $(DOCKER_COMPOSE_FILE) logs >>docker.log