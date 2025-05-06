BINARY_NAME := app

APP_PATH := ./cmd/main.go

# flags go build/test
GOFLAGS ?=

.PHONY: all build run test fmt vet clean

all: build test

build:
	@echo "==> Build $(BINARY_NAME)"
	go build $(GOFLAGS) -o ./bin/$(BINARY_NAME) $(APP_PATH)

run: build
	@echo "==> Run $(BINARY_NAME)"
	./bin/$(BINARY_NAME) -config sunny_5_skiers/config.json -events sunny_5_skiers/events

test:
	@echo "==> Testing"
	go test $(GOFLAGS) ./...

fmt:
	@echo "==> Formatting"
	go fmt ./...

clean:
	@echo "==> Cleaning"
	rm -rf bin/

