.PHONY: test

test:
	go test ./... -cover

lint:
	golangci-lint run ./...