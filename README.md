[![CI Status](https://github.com/ngalaiko/miniboard/workflows/CI/badge.svg)](https://github.com/ngalaiko/miniboard/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngalaiko/miniboard)](https://goreportcard.com/report/github.com/ngalaiko/miniboard)

### Command line arguments

| Command line            | Default                  | Description                                      |
| ----------------------- | ------------------------ | ------------------------------------------------ |
| addr                    | :8080                    | Address to listen for connections.               |
| db-addr                 | db.sqlite                | Database URI to connect to.                      |
| db-type                 | sqlite3                  | Database type (sqlite3, postgres).               |
| domain                  | http://localhost:8080    | Service domain.                                  |
| smtp-host               |                          | SMTP server host.                                |
| smtp-port               |                          | SMTP server port.                                |
| smtp-sender             |                          | SMTP sender.                                     |
| ssl-cert                |                          | Path to ssl certificate.                         |
| ssl-key                 |                          | Path to ssl key.                                 |
| static-path             |                          | Path to static files.

### Environment variables

| Name                       | Description                                      |
| -------------------------- | ------------------------------------------------ |
| SMTP_USERNAME              | Username for SMTP server authentication          |
| SMTP_PASSWORD              | Password for SMTP server authentication          |

## Development

1. Run server: 

```bash
$ cd ./server && go run \ 
    cmd/miniboard/main.go \
    --static-path=../web/src
```

2. Open browser

```bash
$ open http://localhost:8080
```
