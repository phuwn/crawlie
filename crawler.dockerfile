FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --update supervisor && rm  -rf /tmp/* /var/cache/apk/*

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY ./cmd/crawler/*.go ./

COPY ./src ./src

COPY ./config/config.crawler.json /config/config.json

RUN go build -o /main

WORKDIR /

RUN go install github.com/rubenv/sql-migrate/...@latest

COPY ./data/migration ./migration

COPY ./dbconfig.yml ./

COPY ./entrypoint-crawler.sh ./

RUN chmod +x /entrypoint-crawler.sh

COPY ./crontab.txt /etc/crontabs/root

ENTRYPOINT ["/entrypoint-crawler.sh"]