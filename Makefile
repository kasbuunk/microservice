phony: all

.PHONY: lint
lint: build
	golangci-lint run -v ./...

.PHONY: test
test: build
	go test ./...

build: bin/app

bin/app:
	go build -o bin/app ./

gqlgen:
	gqlgen