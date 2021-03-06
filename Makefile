all: format generate test build

format:
	go fmt ./...

generate:
	go generate ./...

test:
	go test ./...

build:
	go build ./cmd/clustertest
	go build ./cmd/clustertestd

install:
	go install ./cmd/clustertest
	go install ./cmd/clustertestd
