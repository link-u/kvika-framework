.PHPNY: all
all: kvika ;

kvika: $(shell find . -type f -name '*.go')
	go build -o ./kvika ./cmd/kvika

.PHONY: musl-static
musl-static: $(shell find . -type f -name '*.go')
	CGO_ENABLED=1 \
		go build --ldflags '-linkmode external -extldflags "-static -lnghttp2 -lssl -lcrypto -lz"' \
		-o ./kvika ./cmd/kvika

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -Rfv kvika