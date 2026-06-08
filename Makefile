.PHONY: run build test lint docker-build

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

docker-build:
	docker build -t be-panganlink-data-handler .
