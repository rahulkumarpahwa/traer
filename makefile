build:
	@go build -o bin/api ./cmd/api
run: build
	@./bin/api
tui:
	@go run ./cmd/tui/main.go
continue:
	@air
migrate-up:
	@go run ./cmd/migrate up
migrate-down:
	@go run ./cmd/migrate down
seed:
	@go run ./cmd/seed