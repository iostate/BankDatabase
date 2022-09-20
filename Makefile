# postgres:
# 		docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
# createdb:
# 		docker exec -it postgres12 createdb --username=root --owner=root simple_bank

# dropdb:
# 		docker exec -it postgres12 dropdb simple_bank

# migrateup:
# 		migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

# migratedown:
# 		migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

# sqlc:
# 		sqlc generate

# .PHONY: postgres createdb dropdb migrateup migratedown

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:6432/simple_bank?sslmode=disable" -verbose up 


migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:6432/simple_bank?sslmode=disable" -verbose down 

sqlc: 
	sqlc generate
	sed -i -e 's/sql.NullString/string/g' db/sqlc/*.sql.go db/sqlc/models.go
	sed -i -e 's/sql.NullInt64/int64/g' db/sqlc/*.sql.go db/sqlc/models.go
	sed -i -e 's/sql.NullTime/time.Time/g' db/sqlc/*.sql.go db/sqlc/models.go
	goimports -w db/sqlc/

test: 
	go test -v -cover ./...


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
