FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY ./cmd/server/*.go ./

COPY ./src ./src

COPY ./config/config.server.json /config/config.json

RUN go build -o /main

WORKDIR /

RUN go install github.com/rubenv/sql-migrate/...@latest

COPY ./data/migration ./migration

COPY ./dbconfig.yml ./

COPY ./entrypoint.sh ./

RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]