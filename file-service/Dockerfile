FROM golang:1.24

WORKDIR /usr/src/file-service

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/file-service

CMD ["file-service"]
