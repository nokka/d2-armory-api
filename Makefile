docker/local:
	docker build -t d2-armory-api:local -f Dockerfile .

test:
	go test -v ./...

lint:
	golangci-lint run
