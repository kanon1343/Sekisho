.PHONY: build test run lint clean certs

APP_NAME := sekisho
BUILD_DIR := bin
CERTS_DIR := certs

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/sekisho

test:
	go test -v -race ./...

run: build certs
	./$(BUILD_DIR)/$(APP_NAME)

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR) $(CERTS_DIR)

certs: $(CERTS_DIR)/localhost.pem $(CERTS_DIR)/localhost-key.pem

$(CERTS_DIR)/localhost.pem $(CERTS_DIR)/localhost-key.pem:
	mkdir -p $(CERTS_DIR)
	mkcert -install
	mkcert -cert-file $(CERTS_DIR)/localhost.pem -key-file $(CERTS_DIR)/localhost-key.pem localhost 127.0.0.1 ::1
