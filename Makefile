BINARY_NAME=s3backup
BIN_DIR=bin
BINARY_PATH=$(BIN_DIR)/$(BINARY_NAME)
GO_FILES=$(shell find . -name "*.go")

.PHONY: all build clean test fmt install uninstall

all: build

build: $(BINARY_PATH)

VERSION?=0.0.0
LDFLAGS=-X github.com/niklas/s3backup/cmd.Version=$(VERSION)

$(BINARY_PATH): $(GO_FILES)
	mkdir -p $(BIN_DIR)
	go build -v -ldflags "$(LDFLAGS)" -o $(BINARY_PATH) main.go

clean:
	rm -rf $(BIN_DIR)
	go clean

test:
	go test -v ./...

fmt:
	go fmt ./...

install: build
	install -m 755 $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)

uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the binary"
	@echo "  clean       Remove build artifacts"
	@echo "  test        Run tests"
	@echo "  fmt         Format source code"
	@echo "  install     Install the binary to /usr/local/bin"
	@echo "  uninstall   Remove the binary from /usr/local/bin"
	@echo "  help        Show this help message"
