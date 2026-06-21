ifneq (,$(wildcard .env))
    include .env
    export
endif

postgres:
	docker run --name postgres_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d docker.io/library/postgres
createdb:
	docker exec -it postgres_db createdb --username=root --owner=root game_db
dropdb:
	docker exec -it postgres_db dropdb game_db
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/game_db?sslmode=disable" --verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/game_db?sslmode=disable" --verbose down
sqlc:
	sqlc generate
seed:
	docker exec -i postgres_db psql -U $(DB_USER) -d $(DB_NAME) -f - < db/seed.sql
.PHONY: postgres createdb dropdb migrateup migratedown sqlc seed
