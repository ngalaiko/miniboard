FROM golang:1.15.0-alpine as go_builder

RUN apk add --no-cache gcc musl-dev

COPY /server /server
WORKDIR /server

RUN go build -o miniboard ./cmd/miniboard/main.go


FROM alpine:3.11.6

COPY --from=go_builder /server/miniboard /app/miniboard
COPY ./web/src /app/dist
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
WORKDIR /app

ENTRYPOINT ["/app/entrypoint.sh"]
