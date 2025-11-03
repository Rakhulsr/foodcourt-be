# Makefile — FoodCourt Clean Architecture
# Jalankan dari mana saja: make run

.PHONY: run dev db-reset test clean help

# === DEFAULT: Run Server ===
run:
	@echo "Starting server..."
	@go run cmd/web/main.go

# === Hot Reload (Windows compatible) ===
dev:
	@echo "Starting with hot reload..."
	@if ! command -v air >/dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@air


# === Test ===
test:
	@echo "Running tests..."
	@go test -v ./...

# === Clean ===
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf tmp/

# === Help ===
help:
	@echo "FoodCourt Makefile"
	@echo ""
	@echo "Commands:"
	@echo "  make run        - Start server (AutoMigrate if first time)"
	@echo "  make dev        - Hot reload (AutoMigrate if first time)"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build"
	@echo "  make help       - Show this help"