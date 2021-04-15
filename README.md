# Miniboard 

[![CI Status](https://github.com/ngalaiko/miniboard/workflows/CI/badge.svg)](https://github.com/ngalaiko/miniboard/actions)[![Coverage Status](https://coveralls.io/repos/github/ngalaiko/miniboard/badge.svg?branch=master)](https://coveralls.io/github/ngalaiko/miniboard?branch=master)[![Go Report Card](https://goreportcard.com/badge/github.com/ngalaiko/miniboard)](https://goreportcard.com/report/github.com/ngalaiko/miniboard)

## API

You can access Swagger UI [here](https://docs.miniboard.app/) if you want to explore and try out the api.

Swagger description itsef is available [here](https://docs.miniboard.app/api.swagger.yaml).

## Configuration

### Backend

#### Command line arguments

| Command line            | Default                  | Description                    |
| ----------------------- | ------------------------ | ------------------------------ |
| config                  |                          | Path to the configuration file |
| verbose                 | false                    | Enable verbose logging         |

#### Configuration file

```yaml
authorizations:
  domain: "example.com" # domain to set cookie to
  secure: false         # if cookie should be Secure
  cookie_lifetime: 720h # lifetime of auth cookie
db:
  driver: "sqlite3"       # available values: "sqlite3", "postgres"
  addr: "./db.sqilite3"   # db address
  max_open_connections: 0 # max open connections to db
http:
  addr: ":8080" # address to listen on
  tls:
    enabled: true          # if false, plaintext http will be used
    key_path: "./key.pem"  # path to tls key
    cert_path: "./crt.pem" # path to tls certificate
operations:
  workers: 10 # number of workers that execute longrunning operations
subsciptions:
  updates:
    workers: 10 # number of workers that update subscriptions in background
    interval: 5m # interval between feed updates
users:
  bcrypt_cose: 14 # bcrypt cost
web:
  fs: true # if true, files will be served from the filesystem
```

#### Environment variables

It is also possible to define any configuration value by setting an environment value, for example:

* `MINIBOARD_HTTP_ADDR` will override `http.addr`
* `MINIBOARD_DB_DRIVER` will override `db.driver`
* etc.

## Development

1. Run server:

```bash
$ cd ./backend \
    && go run cmd/miniboard/main.go \
        --verbose \
        --config config.dev.yaml
```

2. Open browser:

```
$ open http://localhost/
```

## Acknowledgments 

* [Miniflux](https://miniflux.app)
* [Pinboard](https://pinboard.in)
