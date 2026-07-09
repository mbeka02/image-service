
# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/server/main.go

# Run the application
run:
	@go run cmd/server/main.go
# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Database migrations (requires goose: go install github.com/pressly/goose/v3/cmd/goose@v3.22.1)
migrate-up:
	@goose -dir sql/schema postgres "$(DB_URI)" up

migrate-down:
	@goose -dir sql/schema postgres "$(DB_URI)" down

migrate-status:
	@goose -dir sql/schema postgres "$(DB_URI)" status

migrate-create:
	@read -p "Migration name: " name; \
	goose -dir sql/schema create $$name sql

.PHONY: all build run test clean watch migrate-up migrate-down migrate-status migrate-create
