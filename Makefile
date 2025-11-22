.PHONY: run dev dev-backend dev-frontend db-reset test clean help install


install:
	@echo "Installing Go dependencies..."
	@go mod tidy
	@echo "Installing NPM dependencies..."
	@npm install


run:
	@echo "Starting server..."
	@go run cmd/web/main.go


dev:
	@echo "Starting Development Environment..."
	@make -j 2 dev-backend 


dev-backend:
	@echo "   >>> Starting Backend (Air)..."
	@if ! command -v air >/dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@air

migrate:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go

seed:
	@echo "Running seeder..."
	@go run cmd/seed/main.go


test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf tmp/
	@rm -rf public/css/style.css 


help:
	@echo "FoodCourt Makefile"
	@echo ""
	@echo "Commands:"
	@echo "  make install     - Install Go & NPM dependencies"
	@echo "  make dev         - Start Air (Go) & Tailwind (CSS) concurrently"
	@echo "  make run         - Start server only (no hot reload)"
	@echo "  make migrate     - Run migrations"
	@echo "  make seed        - Run seeder"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build & temp files"