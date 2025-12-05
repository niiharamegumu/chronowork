.PHONY: run run-prod build clean test test-verbose test-coverage test-usecase

SRC_DIR := .

run:
	 DATABASE_NAME=sqlite_dev.db go run $(SRC_DIR)/main.go

run-prod:
	 DATABASE_NAME=sqlite.db go run $(SRC_DIR)/main.go

build:
	go build -o chronowork main.go

test:
	@go test ./... 2>&1 | grep -v "no test files" || true

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

test-usecase:
	go test -v ./internal/usecase/...

clean:
	@rm -rf $(SRC_DIR)
