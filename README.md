# SF
This is my grad work. 
File server built in Go using microservice architecture, running in Docker.
Uses PostgreSQL for data storage and Nginx as a reverse proxy.
HTTPS is enabled via a self-signed SSL certificate, which you should generate on your own :) (e.g. using mkcert).
Authentication and file handling are split into separate services.

## Requirements
- Docker
- Docker Compose
- mkcert (or any other tool for generating local SSL certificates)

## How to Run
1. Copy nginx.conf, docker-compose.yml, .env (I'll create these files sometime later)
2. Create dir named `certs` in repo's root directory
2. Generate SSL certificates into certs/:
```
$ mkcert <IP> localhost 127.0.0.1
```
3. Rename certificate to `sf-cert.pem` and key to `sf-key.pem`
4. Start the stack:
```
$ docker compose up -d
```

## Usage
Open in your browser: `https://<IP>`
