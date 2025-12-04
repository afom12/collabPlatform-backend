.PHONY: build run docker-up docker-down docker-build test clean

build:
	go build -o bin/server ./cmd/server

run:
	go run cmd/server/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

docker-logs:
	docker-compose logs -f

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

deps:
	go mod download
	go mod tidy

