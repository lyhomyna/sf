# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o file-service


# Clean final image
FROM debian:bookworm-slim
WORKDIR /usr/src/file-service

COPY --from=builder /app/file-service .

RUN apt update && apt install -y curl

CMD ["./file-service"]
