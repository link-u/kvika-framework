FROM golang:1.15-alpine as builder

WORKDIR /go/src/github.com/link-u/kvika-framework
COPY . .

RUN apk add --no-cache gcc musl-dev make \
            libcurl curl-dev curl-static \
            openssl openssl-dev libssl1.1 libcrypto1.1 openssl-libs-static \
            nghttp2 nghttp2-libs nghttp2-dev nghttp2-static \
            zlib zlib-dev zlib-static

ENV GOOS=linux
ENV GOARCH=amd64

RUN make clean && make musl-static

FROM alpine:3.12

WORKDIR /

RUN apk add --no-cache ca-certificates && update-ca-certificates

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
ENV SSL_CERT_DIR=/etc/ssl/certs

COPY --chown=nobody:nogroup --from=builder /go/src/github.com/link-u/kvika-framework/kvika kvika
RUN ["chmod", "a+x", "/kvika"]
ENTRYPOINT ["/kvika"]
