# Simple Makefile for a Go project

build:
	@echo "Building..."

	@./tailwindcss -i internal/assets/input.css -o public/assets/css/output.css  -m -c internal/assets/tailwind.config.js
	@go build -o bin/main cmd/main.go

# Run the application
run:
	@go run cmd/main.go

tailwind:
	@if [ ! -f tailwindcss ]; then curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-x64 -o tailwindcss; fi
	@chmod +x tailwindcss

# Create DB container
docker-run:
	@if docker compose up -d --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air -c ./air.toml; \
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

.PHONY: build run watch tailwind
