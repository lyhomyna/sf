FROM golang:latest

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . ./migrations

CMD ["bash"]
