## SSL certificates for local testing

```
$ openssl genrsa 2048 > local.miniboard.app.key
$ chmod 400 local.miniboard.app.key
$ openssl req -new -x509 -nodes -sha256 -days 365 -key local.miniboard.app.key -out local.miniboard.app.cert
```
