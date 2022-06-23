phony: all

BINARY_NAME = bin/app
GQLGEN_DIR = server/gql

all: gqlgen build lint test

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: test
test: build
	go test ./...

build: ${BINARY_NAME}

${BINARY_NAME}: $(shell find . -name "*.go")
	go build -o ${BINARY_NAME} ./

gqlgen:
	(cd ${GQLGEN_DIR} && gqlgen)

clean:
	go clean
	rm ${BINARY_NAME}
