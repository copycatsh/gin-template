.PHONY: run build test test-v fmt lint docker-build docker-up docker-down clean

run:
	go run ./cmd/api

build:
	go build -o server ./cmd/api

test:
	go test ./...

test-v:
	go test ./... -v

fmt:
	docker run --rm -v $(PWD):/app -w /app golang:1.22-alpine gofmt -w .

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run

docker-build:
	docker build -t gin-template .

up:
	docker compose up -d

down:
	docker compose down

clean:
	rm -f server
