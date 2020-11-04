.PHPNY: all
all: kvika ;

kvika: $(shell find . -type f -name '*.go')
	go build -o ./kvika ./cmd/kvika

.PHONY: test
test:
	go test ./...
