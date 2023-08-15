DB_NAME=simple_bank
DB_CONTAINER_NAME=postgres14-db
DB_CONN_STRING="postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"

docker-clean:
	docker rm ${DB_CONTAINER_NAME}

postgres:
	docker run --name ${DB_CONTAINER_NAME} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.1-alpine

create-db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=root --owner=root ${DB_NAME}

drop-db:
	docker exec -it ${DB_CONTAINER_NAME} dropdb ${DB_NAME}

migrate-up:
	migrate -path db/migration --database ${DB_CONN_STRING} -verbose up

migrate-down:
	migrate -path db/migration --database ${DB_CONN_STRING} -verbose down

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

.PHONY: postgres create-db drop-db server
