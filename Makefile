.PHONY: build test run lint clean

APP_NAME := sekisho
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/sekisho

test:
	go test -v -race ./...

run: build
	./$(BUILD_DIR)/$(APP_NAME)

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)
