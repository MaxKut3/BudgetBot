.DEFAULT_GOAL := run

run:
	gofmt -w .
	goimports -w .
	go run .

tests:
	gofmt -w .
	goimports -w .
	go test

build:
	go build .

fmt:
	gofmt -w .

imports:
	goimports -w .

lint:
	golangci-lint run main.go



