include .env
export

CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATION_DIR = internal/db/migrations

# import database
importdb:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql

# export database
exportdb:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql

# Run server
server:
	go run .

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

.PHONY: importdb exportdb server migrate-create migrate-up migrate-dow migrate-force migrate-goto migrate-down-n