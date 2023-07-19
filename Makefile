SHELL=cmd.exe
APP_BINARY=test.exe

ping:
	@echo "Yo, i'm alive"

tidy:
	go mod tidy
	go mod vendor

rebuild:
	docker-compose down -v
	docker-compose up

test:
	go test -v ./...
