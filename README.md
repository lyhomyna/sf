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
1. Download `template-files` folder and cd into it
2. Create dir named `certs`
3. Generate SSL certificates into `certs/`:
```
$ mkcert <IP> localhost 127.0.0.1
```
4. Rename certificate to `sf-cert.pem` and key to `sf-key.pem`
5. Change variables in all `.env` files (recommended)
6. Start the stack:
```
$ docker compose up -d
```
7. Register new user:
```
$ curl -H "Content-Type: application/json" -d '{"email":"mail@example.com", "password":"strong"}' https://tmpl:strong@localhost/api/auth/register -v -k
```
>[!NOTE] `tmpl:strong` are credentials that located in `.auth.env file`

## Usage
Open in your browser: `https://<IP>` and log in using the credentials created in step 7.

https://github.com/user-attachments/assets/c8704fdc-2fa3-431a-bdfc-a0fb5f559527
