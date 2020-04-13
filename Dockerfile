FROM golang:1.14.2-alpine as go_builder

RUN apk add --no-cache gcc musl-dev

COPY /server /server
WORKDIR /server

RUN go build -o miniboard ./cmd/miniboard/main.go


FROM node:13.8.0-alpine as node_builder

ARG VERSION=development
ENV VERSION=$VERSION

COPY /web web
WORKDIR /web

RUN npm install --global rollup
RUN yarn install && yarn build


FROM alpine:3.11.3

COPY --from=go_builder /server/miniboard /app/miniboard
COPY --from=node_builder /web/dist /app/dist
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
WORKDIR /app

ENTRYPOINT ["/app/entrypoint.sh"]
