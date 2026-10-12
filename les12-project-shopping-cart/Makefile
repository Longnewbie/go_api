include .env
export

CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATION_DIR = internal/db/migrations

ENV_FILE=.env

PROD_COMPOSE=docker-compose.prod.yml
NOAPP_COMPOSE=docker-compose.noapp.yml
DEV_COMPOSE=docker-compose.dev.yml

# import database
importdb:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql

# export database
exportdb:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql

# build a binary file
build:
	go build -o bin/myapp ./cmd/api

run-binary:
	./bin/myapp

# Run server
server:
	go run ./cmd/api

# Generate sqlc
sqlc:
	sqlc generate

# create a new migration (ex: make migrate-create NAME=profiles)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

# run all pending migration (ex: make migrate-up)
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(CONN_STRING)" up

# rollback the last migration (ex: make migrate-down)
migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(CONN_STRING)" down 1

# rollback N migration (ex: make migrate-down-n N=2)
migrate-down-n:
	migrate -path $(MIGRATION_DIR) -database "$(CONN_STRING)" down $(N)

# force migration version (use with caution) || (ex: make migrate-force VERSION=1)
migrate-force:
	migrate -path $(MIGRATION_DIR) -database "$(CONN_STRING)" force $(VERSION)

# drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATION_DIR) -database "$(CONN_STRING)" drop

# apply specific migration by version (ex: make migrate-goto VERSION=1)
migrate-goto:
	migrate -path $(MIGRATION_DIR) -database "$(CONN_STRING)" goto $(VERSION)

# Build container for noapp
noapp:
	docker compose -f $(NOAPP_COMPOSE) down
	docker compose -f $(NOAPP_COMPOSE) --env-file $(ENV_FILE) up -d --build

# Stop all containers noapp
stop-noapp:
	docker compose -f $(NOAPP_COMPOSE) down

# Build container for dev
dev:
	docker compose -f $(DEV_COMPOSE) down
	docker compose -f $(DEV_COMPOSE) --env-file $(ENV_FILE) up --build

# Stop all containers dev
stop-dev:
	docker compose -f $(DEV_COMPOSE) down

# Build container for production
prod:
	docker compose -f $(PROD_COMPOSE) down
	docker compose -f $(PROD_COMPOSE) --env-file $(ENV_FILE) up -d --build

# Stop all containers production
stop-prod:
	docker compose -f $(PROD_COMPOSE) down

# Show logs container production
logs-prod:
	docker compose -f $(PROD_COMPOSE) logs -f --tail=100

# Go to container golang-api
bash:
	docker exec -it golang-api /bin/sh

.PHONY: importdb exportdb server migrate-create migrate-up migrate-dow migrate-force migrate-goto migrate-down-n sqlc build run-binary prod stop-prod logs-prod bash noapp stop-noapp dev stop-dev