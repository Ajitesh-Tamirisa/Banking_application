postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root banking_app

dropdb:
	docker exec -it postgres14 dropdb banking_app

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/banking_app?sslmode=disable" -verbose up
	
migrateup1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/banking_app?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/banking_app?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/banking_app?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go banking_application/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock