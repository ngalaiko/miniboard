FROM golang:1.12.4 as builder

WORKDIR /app

COPY . .

RUN go build \
    -mod vendor \
    -o miniboard \
    ./cmd/miniboard/main.go 

FROM alpine:3.9

COPY --from=builder /app/miniboard miniboard

RUN ["/miniboard"]
