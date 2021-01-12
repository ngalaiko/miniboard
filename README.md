[![CI Status](https://github.com/ngalaiko/miniboard/workflows/CI/badge.svg)](https://github.com/ngalaiko/miniboard/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngalaiko/miniboard)](https://goreportcard.com/report/github.com/ngalaiko/miniboard)

### Command line arguments

| Command line            | Default                  | Description                                      |
| ----------------------- | ------------------------ | ------------------------------------------------ |
| config                  |                          | Path to the configuration file, required.        |

### Configuration file

See [example](./server/config.dev.yaml).

It is also possible to define any configuration value by setting an environment value, for example:
* `MINIBOARD_HTTP_ADDR` will override `http.addr`
* `MINIBOARD_DB_DRIVER` will override `db.driver`

## Development

1. Run server: 

```bash
$ cd ./server && go run \ 
    cmd/miniboard/main.go \
        --config config.dev.yaml
```

2. Open browser

```bash
$ open http://localhost:8080
```
