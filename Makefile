GHCR_REPO=ghcr.io/nokka/d2-armory-api
GHCR_COMMIT_TAG=$(GHCR_REPO):commit-$(GITHUB_SHA)
VALID_TAG=$(shell echo $(TAG_NAME) | sed 's/[^a-z0-9_\.-]/-/g')
COVER_PKG=$(shell go list ./... | grep -v /cmd/ | grep -v /integrationtest/ | paste -sd "," -)

# Builds a local docker image for local use.
docker/local:
	docker build -f Dockerfile -t d2-armory-api:local .

# Builds docker image with the Github container registry commit tag.
docker/build:
	docker build -f Dockerfile -t $(GHCR_COMMIT_TAG) .

docker/tag:
	docker tag $(GHCR_COMMIT_TAG) $(GHCR_REPO):$(VALID_TAG)

docker/push:
	docker push $(GHCR_REPO)

test:
	go test -v ./...

# test integration will setup a database using docker and test related functionality.
test/integration: test/integration/prepare test/integration/run test/integration/teardown

test/integration/prepare:
	docker-compose -f integrationtest/docker-compose.yml up -d

test/integration/run:
	go test -v ./... -tags=integration

# teardown removes all resources created for the integration test.
test/integration/teardown:
	docker-compose -f integrationtest/docker-compose.yml down -v

lint:
	golangci-lint run --disable-all -E gocyclo -E golint -E staticcheck -E structcheck -E unused -E gocritic -E gofmt -E interfacer -E misspell -E stylecheck -E unconvert -E unparam -E scopelint -E prealloc
