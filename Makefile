ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	@go build -o bin/out ./cmd/app/main.go

run: build
	@./bin/out

test:
	@echo "Testing..."
	@go test ./tests -v

db-status:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose status

db-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose up

db-reset:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose reset

swag:
	@go run github.com/swaggo/swag/cmd/swag@latest fmt
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/app/main.go -o docs

dev:
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