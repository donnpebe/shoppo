.PHONY: test lint build

build: 
	go build -o build/shoppo cmd/main.go

test:
	go test ./... -cover

lint:
	golangci-lint run ./...