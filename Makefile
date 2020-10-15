# Builds a local docker image for local use.
docker/local:
	docker build -f Dockerfile -t d2-armory-api:local .

# Builds docker image with the Github container registry commit tag.
docker/build:
	docker build -f Dockerfile -t $(GHCR_COMMIT_TAG) .

docker/tag:
	docker tag $(GHCR_COMMIT_TAG) $(GHCR_REPO):$(VALID_TAG)

# docker/login:
# 	echo $(CR_TOKEN)
# 	echo $(secrets.CR_TOKEN)
# 	echo $(secrets.CR_TOKEN) | docker login ghcr.io -u nokka --password-stdin

test:
	go test -v ./...

lint:
	golangci-lint run --disable-all -E gocyclo -E golint -E staticcheck -E structcheck -E unused -E gocritic -E gofmt -E interfacer -E misspell -E stylecheck -E unconvert -E unparam -E scopelint -E prealloc
