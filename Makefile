
.PHONY: run dev tidy migrate seed clean help


dev:
	@echo "Checking for Air..."
	@if ! command -v air >/dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@echo ">>> Starting Server..."
	@air

run:
	@echo ">>> Starting Server..."
	@go run cmd/web/main.go

tidy:
	@echo "Tidying up Go modules..."
	@go mod tidy

migrate:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go

seed:
	@echo "Seeding database..."
	@go run cmd/seed/main.go

db-reset:
	@echo "Resetting Database..."
	@
	@
	@make migrate
	@make seed


clean:
	@echo "Cleaning tmp files..."
	@go clean
	@rm -rf tmp/

help:
	@echo "FoodCourt Makefile Commands:"
	@echo "  make dev       - Jalankan server dengan Hot Reload (Air)"
	@echo "  make run       - Jalankan server biasa"
	@echo "  make tidy      - Rapikan go.mod dan download library"
	@echo "  make migrate   - Update struktur database"
	@echo "  make seed      - Masukkan data awal (Admin user)"
	@echo "  make clean     - Hapus file temporary"