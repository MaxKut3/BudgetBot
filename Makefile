.DEFAULT_GOAL := run

run:
	gofmt -w .
	goimports -w .
	go run cmd/BudgetBot/main.go

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

mig-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=postgres dbname=postgres sslmode=disable host=localhost port=5434" goose -dir migrations up

docker:
	docker-compose up