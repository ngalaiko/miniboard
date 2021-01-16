[![CI Status](https://github.com/ngalaiko/miniboard/workflows/CI/badge.svg)](https://github.com/ngalaiko/miniboard/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngalaiko/miniboard)](https://goreportcard.com/report/github.com/ngalaiko/miniboard)

## API

You can access Swagger UI [here](https://docs.miniboard.app/) if you want to expole and try out the api.

Swagger description itsef is available [here](https://docs.miniboard.app/api.swagger.yaml).

## Configuration

### Command line arguments

| Command line            | Default                  | Description                     |
| ----------------------- | ------------------------ | ------------------------------- |
| config                  |                          | Path to the configuration file. |

### Configuration file

```yaml
db:
  driver: "sqlite3"     # available values: "sqlite3", "postgres"
  addr: "./db.sqilite3" # db address
http:
  addr: ":8080" # address to listen on
  tls:
    key_path: "key.pem" # path to tls key
    cert_path "crt.pem" # path to tls certificate
operations:
  workers: 10 # number of workers that execute longrunning operations
users:
  bcrypt_cose: 14 # bcrypt cost
```

### Environment variables

It is also possible to define any configuration value by setting an environment value, for example:
* `MINIBOARD_HTTP_ADDR` will override `http.addr`
* `MINIBOARD_DB_DRIVER` will override `db.driver`

## Development

1. Run server: 

```bash
$ cd ./backend && go run \
    cmd/miniboard/main.go \
        --config config.dev.yaml
```

2. Open browser

```bash
$ open http://localhost:8080
```
