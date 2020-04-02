docker/local:
	docker build -t d2-armory-api:local -f Dockerfile .

test:
	go test -v ./...

lint:
	golangci-lint run --disable-all -E gocyclo -E golint -E staticcheck -E structcheck -E unused -E gocritic -E gofmt -E interfacer -E misspell -E stylecheck -E unconvert -E unparam -E scopelint -E prealloc
