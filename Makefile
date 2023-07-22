
ping:
	# @echo "Yo, i'm alive"

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

.PHONY: server
