# Set up tools.
install:
	go install github.com/cosmtrek/air@v1.27.3
	go install github.com/swaggo/swag/cmd/swag@latest

# Start dev server.
start:
	air

# Set up database.
setup_db:
	./bin/init_db.sh

# Migrate scheme to database.
migrate_schema:
	go run ./cmd/migration/main.go

.PHONY: install setup_db start migrate_schema

swag:
	swag init -g internal/app/app.go -o ./gen/docs
.PHONY: swag