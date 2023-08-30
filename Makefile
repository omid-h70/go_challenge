DB_NAME=simple_bank
DB_CONTAINER_NAME=postgres14-db
DB_CONN_STRING="postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"

docker-clean:
	docker rm ${DB_CONTAINER_NAME}

postgres:
	docker run --name ${DB_CONTAINER_NAME} -v $(CURDIR)/data:/data/postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.1-alpine

psql:
	docker exec -it ${DB_CONTAINER_NAME} psql -U root

create-db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=root --owner=root ${DB_NAME}

drop-db:
	docker exec -it ${DB_CONTAINER_NAME} dropdb ${DB_NAME}

migrate-up:
	migrate -path db/migration --database ${DB_CONN_STRING} -verbose up

migrate-down:
	migrate -path db/migration --database ${DB_CONN_STRING} -verbose down

#docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate is Windows Specific Command
#$(CURDIR) is gnu makefile variable and works every where (?)
sqlc:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc generate

#generate dbdocs from db.dbml
db_docs:
	dbdocs build doc/db.dbml

#Convert dbml file to postgres sql
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

# rm -f pb/*.go
# swagger creates documentation for u, that machines understand and can share
# statik helps you to generate static files in binary format, it will be so much faster and docker image will be unchanged
# it is used for swagger-ui project => get files from swagger-ui-master.zip\swagger-ui-master\dist
#    statik -src=./doc/swagger -dest=./doc
proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb  --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

# a tool to test grpc
evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

ping:
	@echo "Yo, i'm alive"

tidy:
	go mod tidy
	go mod vendor

clean-build:
	# @echo "Clean and Build"
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up --force-recreate
	# @echo "Clean Build Done!"

#Calling Server with Our Desired Emv
#SERVER_ADDRESS=0.0.0.0:8085 make server
server:
	go run main.go

test:
	go test -v ./...

# meaning that the target name , doesn't represent  an existing file

.PHONY: postgres create-db drop-db server proto evans redis
