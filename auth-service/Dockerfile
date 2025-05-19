# Build stage
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o auth-service


# Clean final image
FROM debian:bookworm-slim
WORKDIR /usr/src/auth-service

COPY --from=builder /app/auth-service .

RUN apt update && apt install -y curl

CMD ["./auth-service"]
