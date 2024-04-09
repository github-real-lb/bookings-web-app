postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

start:
	docker start postgres16

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root bookings
	docker exec -it postgres16 createdb --username=root --owner=root bookings_test

dropdb:
	docker exec -it postgres16 dropdb bookings
	docker exec -it postgres16 dropdb bookings_test

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bookings?sslmode=disable" -verbose up
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bookings_test?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bookings?sslmode=disable" -verbose down
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/bookings_test?sslmode=disable" -verbose down

sqlc:
	sqlc generate

mock:
	mockery --output db/mocks --filename store.go --outpkg mocks  --dir db --name DatabaseStore --structname MockStore

test:
	go test -v ./... -cover

server:
	go run ./cmd/web

.PHONY: postgrs start createdb  dropdb migrateup migratedown sqlc mock test server
