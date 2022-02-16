# test image
FROM golang:1.17.1-alpine3.14

# set workdir where the app will work from.
WORKDIR /app

# download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# add all files to image.
ADD . ./

# Run only the integration tests
CMD go test -v ./... -tags=integration
