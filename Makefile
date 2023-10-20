DB_NAME=simple_bank
DB_CONTAINER_NAME=postgres14-db
DB_CONN_STRING="postgres://root:secret@localhost:5432/simple-bank?sslmode=disable"
#AWS_DB_CONN_STRING="postgres://root:a}1K&{LuGe}jw*BCasFdhxp83y7f@simple-bank.c1lri2mnc8os.us-east-2.rds.amazonaws.com:5432/awsdude"
#AWS_DB_CONN_STRING="postgres%3A%2F%2Froot%3Aawsdude%23123456%40simple-bank.c1lri2mnc8os.us-east-2.rds.amazonaws.com%3A5432%2Fawsdude"
AWS_DB_CONN_STRING="postgres://root:a%7D1K%26%7BLuGe%7Djw%2ABCasFdhxp83y7f@simple-bank.c1lri2mnc8os.us-east-2.rds.amazonaws.com:5432/awsdude"

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

aws-migrate-up:
	migrate -path "./db/migration" -database ${AWS_DB_CONN_STRING} -verbose up

migrate-down:
	migrate -path db/migration --database ${DB_CONN_STRING} -verbose down

# make new-migration name="add_verify_emails"
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

#docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate is Windows Specific Command
#$(CURDIR) is gnu makefile variable and works every where (?)
sqlc:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc generate

#generate dbdocs from db.dbml
db_docs:
	dbdocs build doc/db.dbml

#Convert dbml file to postgres sql
db_schema:
	dbml2sql --postgres -o db/doc/schema.sql db/doc/db.dbml

# Create MockDb
# Store, TaskDistributor are interfaces
mock:
	mockgen -package mockdb -destination db/mock/store.go go_challenge/db/sqlc Store
	mockgen -package mockwr -destination worker/mock/distributor.go go_challenge/worker TaskDistributor

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

test-short:
	go test -v -cover -short ./...

test:
	go test -v -cover  ./...

# meaning that the target name , doesn't represent  an existing file

.PHONY: postgres create-db drop-db server proto evans redis
