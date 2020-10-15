# compiler image
FROM golang:1.15.1-alpine AS build-env

# set workdir where the app will work from.
WORKDIR /app

# download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# add all files to image.
ADD . ./

# build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o armory ./cmd/server/main.go

# create final application image.
FROM alpine:3.12
WORKDIR /app
COPY --from=build-env /app/armory .
ENTRYPOINT ./armory
