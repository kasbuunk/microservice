phony: all

.PHONY: lint
lint: build
	golangci-lint run -v ./...

.PHONY: test
test: build
	go test ./...

.PHONY: build
build: bin/app

# TODO: build should not be phony, fix the pattern matching
bin/app: go.mod
	go build -o bin/app ./

gqlgen:
	gqlgen