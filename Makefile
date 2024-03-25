postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

start:
	docker start postgres16

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root bookings

dropdb:
	docker exec -it postgres16 dropdb bookings

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bookings?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bookings?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./... -cover

server:
	go run ./cmd/web

.PHONY: postgrs start createdb  dropdb migrateup migratedown sqlc test server
