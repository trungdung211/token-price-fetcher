# Set up tools.
install:
	go install github.com/cosmtrek/air@v1.27.3
	go install github.com/swaggo/swag/cmd/swag@latest

# Start dev server.
start:
	air

swag:
	swag init -g internal/app/app.go -o ./gen/docs
.PHONY: swag

test: 
	go test -race -cover -coverprofile=coverage.txt -covermode=atomic -v ./... -timeout 600s
.PHONY: test