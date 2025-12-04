.PHONY: build run docker-up docker-down docker-build test clean swagger

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

swagger:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

swagger-serve: swagger
	@echo "Swagger docs generated. Start server and visit http://localhost:8080/swagger/index.html"
