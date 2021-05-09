ARG GO_VERSION=1.16.3
ARG ALPINE_VERSION=3.12

FROM golang:$GO_VERSION-alpine$ALPINE_VERSION as builder

COPY backend /src/miniboard

WORKDIR /src/miniboard

RUN apk add --no-cache gcc musl-dev

ARG LITESTREAM_VERSION=v0.3.4
ARG LITESTREAM_SHA256=1539ec40440e5167790d4b1702366fea7ecf3a6aab44c2a302e64f7f2dbb634d
ADD https://github.com/benbjohnson/litestream/releases/download/${LITESTREAM_VERSION}/litestream-${LITESTREAM_VERSION}-linux-amd64-static.tar.gz /tmp/
RUN echo "${LITESTREAM_SHA256}  /tmp/litestream-${LITESTREAM_VERSION}-linux-amd64-static.tar.gz" | sha256sum -c \
    && tar -C /usr/local/bin -xzf /tmp/litestream-${LITESTREAM_VERSION}-linux-amd64-static.tar.gz

ARG S6_OVERLAY_VERSION=v2.2.0.3
ARG S6_OVERLAY_SHASUM=7140eafc62720ecc43f81292c9bdd75bc7f4d0421518d42707e69f8e78e55088
ADD https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64-installer /tmp/
RUN echo "${S6_OVERLAY_SHASUM}  /tmp/s6-overlay-amd64-installer" | sha256sum -c \
    && chmod +x /tmp/s6-overlay-amd64-installer

RUN go build \
        -o /usr/local/bin/miniboard \
        -ldflags '-s -w -extldflags "-static"' \
        -tags osusergo,netgo,sqlite_omit_load_extension \
        ./cmd/miniboard/main.go \
    && chmod +x /usr/local/bin/miniboard


FROM alpine:$ALPINE_VERSION

COPY --from=builder /usr/local/bin/miniboard /usr/local/bin/miniboard
COPY --from=builder /usr/local/bin/litestream /usr/local/bin/litestream

COPY --from=builder /tmp/s6-overlay-amd64-installer /tmp/s6-overlay-amd64-installer
RUN /tmp/s6-overlay-amd64-installer /

RUN apk add bash

RUN mkdir -p /data

COPY etc/cont-init.d /etc/cont-init.d
COPY etc/services.d /etc/services.d

# The kill grace time is set to zero because our app handles shutdown through SIGTERM.
ENV S6_KILL_GRACETIME=0

# Sync disks is enabled so that data is properly flushed.
ENV S6_SYNC_DISKS=1

# Run the s6 init process on entry.
ENTRYPOINT [ "/init" ]
