.PHONY: build run serve start_db migrate_up migrate_down install_deps


build:
	echo "Building project..."
	go build

run:
	go run .

test:
	go test ./...

serve:
	go build -o out && ./out

start_db:
	echo "Starting postgresql engine..."
	brew services start postgresql@17

migrate_up:
	@echo "Running migrations up..."
	@set -a; \
	source .env; \
	goose -dir sql/schema postgres "$$DB_URL" up

migrate_down:
	@echo "Running migrations down..."
	@set -a; \
	source .env; \
	goose -dir sql/schema postgres "$$DB_URL" down

install_deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/joho/godotenv/cmd/dotenv@latest

queries:
	sqlc generate