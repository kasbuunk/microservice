phony: all

all: gqlgen build lint test

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: test
test: build
	go test ./...

build: bin/app

bin/app: $(shell find . -name "*.go")
	go build -o bin/app ./

gqlgen:
	(cd input/server && gqlgen)

clean:
	go clean
	rm ${BINARY_NAME}
