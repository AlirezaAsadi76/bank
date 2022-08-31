postgres:
	docker run --name mybank -e POSTGRES_USER=root -p 5432:5432 -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine
createdb:
	docker exec -it mybank createdb --username=root --owner=root bank
dropdb:
	docker exec -it mybank dropdb bank
migrateup:
	 migrate -path db/migration -database "postgresql://root:123456@127.0.0.1:5432/bank?sslmode=disable" -verbose up
migratedown:
	 migrate -path db/migration -database "postgresql://root:123456@127.0.0.1:5432/bank?sslmode=disable" -verbose down
test:
	go test ./... -v -cover
.PHONY:createdb postgres dropdb migratedown migrateup test