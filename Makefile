ENV := .env
.PHONY: build up start down remove stop restart logs logs-save

include $(ENV)
# App
run:
	go run ./cmd/main.go

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