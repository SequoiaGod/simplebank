postgres:
	docker run --name yan-postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12125412 -d postgres
createdb:
	docker exec -it yan-postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it yan-postgres dropdb simple_bank
migrateup1:
	migrate -path db/migration -database "postgresql://root:12125412@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migrateup:
	migrate -path db/migration -database "postgresql://root:12125412@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown1:
	migrate -path db/migration -database "postgresql://root:12125412@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
migratedown:
	migrate -path db/migration -database "postgresql://root:12125412@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server migrateup1 migratedown1