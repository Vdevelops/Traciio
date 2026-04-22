.PHONY: build run test clean docker-build docker-run docker-up docker-down

# Build the application
build:
	go build -o server ./cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f server
	go clean

# Build Docker image
docker-build:
	docker build -t crm-healthcare-api .

# Run Docker container
docker-run:
	docker run -p 8080:8080 crm-healthcare-api

# Start with docker-compose
docker-up:
	docker-compose up --build

# Stop docker-compose
docker-down:
	docker-compose down

# Download dependencies
deps:
	go mod download
	go mod tidy

