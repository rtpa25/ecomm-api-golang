postgres:
	docker run --name ecomm-go -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it ecomm-go createdb --username=root --owner=root ecomm

dropdb:
	docker exec -it ecomm-go dropdb ecomm

migrateup:
	migrate -path db/migration -database postgresql://root:secret@localhost:5432/ecomm?sslmode=disable -verbose up

migratedown:
	migrate -path db/migration -database postgresql://root:secret@localhost:5432/ecomm?sslmode=disable -verbose down

pullSQLC:
	docker pull kjconroy/sqlc

startupSQLC:
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate

server:
	go run main.go