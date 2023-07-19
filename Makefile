
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

test:
	go test -v ./...
