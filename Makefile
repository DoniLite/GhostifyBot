APP_NAME=ghostify-bot
BIN_OUTPUT=./bin/$(APP_NAME)
SRC=main.go

build:
	go build -o $(BIN_OUTPUT) $(SRC)

run: build
	$(BIN_OUTPUT)

clean:
	rm -f $(BIN_OUTPUT)

lint:
	golangci-lint run ./...

test:
	go test ./...

fmt:
	go fmt ./...

install-deps:
	go mod tidy

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run --rm --env-file .env $(APP_NAME):latest

.PHONY: build run clean test lint fmt install-deps
