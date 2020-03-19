# compiler image
FROM golang:1.13-alpine3.11 AS build-env

# set workdir where the app will work from.
WORKDIR /app

# add all files to image.
ADD . ./

# download dependencies.
RUN go mod download

# build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o goapp ./cmd/server/main.go

# create final application image.
FROM alpine:3.10
WORKDIR /app
COPY --from=build-env /app/goapp .
ENTRYPOINT ./goapp
