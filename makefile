# Define variables
GOCMD := go
GOINSTALL := $(GOCMD) install
GOTEST := $(GOCMD) test
GOBUILD := $(GOCMD) build
GOLINT := golangci-lint run

# Define default target
.PHONY: all
all: install lint test

# Lint target
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	$(GOLINT)

# Test target
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Install target
.PHONY: install
install:
	@echo "Installing the project..."
	$(GOINSTALL) ./...

# Build target
.PHONY: build
build:
	@echo "Building the project..."
	$(GOBUILD) ./...

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning build cache..."
	$(GOCMD) clean -cache -testcache -modcache

# Dev target
.PHONY: dev
dev: lint test
	@echo "Linting and testing done."