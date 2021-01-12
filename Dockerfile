ARG GO_VERSION=1.15.2		
FROM golang:${GO_VERSION}-alpine as go_builder		

RUN apk add --no-cache gcc musl-dev		

COPY /server /server		
COPY /web /web		
WORKDIR /server		

RUN go test ./... -v

RUN go build -o miniboard ./cmd/miniboard/main.go		

FROM alpine:3.12.3

COPY --from=go_builder /server/miniboard /bin/miniboard		

ENTRYPOINT ["/bin/miniboard"]
CMD ["/bin/miniboard"]
