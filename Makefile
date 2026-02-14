.PHONY: build test install clean fmt import-tool

# Default Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

# Binary names
PROVIDER_BINARY=terraform-provider-opnsense
IMPORT_TOOL_BINARY=opnsense-import

# Build directory
BUILD_DIR=./bin

# Version information
VERSION?=dev

all: build

build: fmt
	@echo "Building provider..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(PROVIDER_BINARY) -ldflags="-X main.version=$(VERSION)" .

import-tool: fmt
	@echo "Building import tool..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(IMPORT_TOOL_BINARY) ./tools/import

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

install: build
	@echo "Installing provider..."
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/yknx4/opnsense/$(VERSION)/linux_amd64/
	@cp $(BUILD_DIR)/$(PROVIDER_BINARY) ~/.terraform.d/plugins/registry.terraform.io/yknx4/opnsense/$(VERSION)/linux_amd64/

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f terraform.tfstate*
	@rm -f .terraform.lock.hcl
	@rm -rf .terraform/

deps:
	@echo "Downloading dependencies..."
	$(GOCMD) mod download
	$(GOCMD) mod tidy
