[![CI Status](https://github.com/ngalaiko/miniboard/workflows/CI/badge.svg)](https://github.com/ngalaiko/miniboard/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngalaiko/miniboard)](https://goreportcard.com/report/github.com/ngalaiko/miniboard)

### Command line arguments

| Command line            | Default                  | Description                                      |
| ----------------------- | ------------------------ | ------------------------------------------------ |
| addr                    | :8080                    | Address to listen for connections.               |
| bolt-path               | ./bolt.db                | Path to the bolt storage.                        |
| redis-uri               |                          | Redis URI to connect to.                         |
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

1. Watch front changes:

```
$ cd ./web && yarn watch
```

2. Run server: 

```
$ cd ./server && go run \ 
    cmd/miniboard/main.go \
    --static-path=../web/dist
```

3. Open browser

```
$ open http://localhost:8080
```
