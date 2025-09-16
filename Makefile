APP_NAME=myapp

.PHONY: run build clean

# Run the application
run:
	go run cmd/web/main.go

# Build the application binary
build:
	go build -o bin/$(APP_NAME) cmd/web/main.go

# Clean build artifacts
clean:
	rm -rf bin
