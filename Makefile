postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d docker.io/library/postgres
createdb:
	docker exec -it postgres createdb --username=root --owner=root game_db
dropdb:
	docker exec -it postgres dropdb game_db
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/game_db?sslmode=disable" --verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/game_db?sslmode=disable" --verbose down

.PHONY: postgres createdb dropdb migrateup migratedown
