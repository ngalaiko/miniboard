[![CI Status](https://github.com/ngalaiko/miniboard/workflows/CI/badge.svg)](https://github.com/ngalaiko/miniboard/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngalaiko/miniboard)](https://goreportcard.com/report/github.com/ngalaiko/miniboard)

### Environment variables

| Name                       | Description                                      |
| -------------------------- | ------------------------------------------------ |
| SMTP_USERNAME              | Username for SMTP server authentication          |
| SMTP_PASSWORD              | Password for SMTP server authentication          |

## Development

Requirements: 

* [bazel](https://bazel.build)

Run: 
```
$ bazel run :miniboard
```

Run tests:
```
$ bazel test //...
```
