FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xzf - "migrate" && mv migrate /usr/bin

COPY . ./migrations
