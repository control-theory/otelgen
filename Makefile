# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint
REPO=controltheory
# Binary name
BINARY_NAME=otelgen

# Main package path
MAIN_PACKAGE=./cmd/otelgen

# Build directory
BUILD_DIR=./build

# Source files
SRC:=$(shell find . -name "*.go")
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE:=$(GIT_BRANCH)::$(shell date +"%D-%T:%s")
COMMIT_ID:=$(shell git rev-parse --short=8 HEAD)
BUILD_VERSION=$(shell git describe --tags --always --dirty)

# Test coverage output
COVERAGE_OUTPUT=coverage.out

.PHONY: build clean test coverage lint deps tidy run help

build: $(BUILD_DIR)/$(BINARY_NAME) ## Build the application

$(BUILD_DIR)/$(BINARY_NAME): $(SRC)
	$(GOBUILD) -ldflags "-X main.version=${BUILD_VERSION} -X main.date=${BUILD_DATE} -X main.commit=${COMMIT_ID}" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_OUTPUT)

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=$(COVERAGE_OUTPUT) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUTPUT)

lint: ## Run linter
	$(GOLINT) run

deps: ## Download dependencies
	$(GOGET) -v -t -d ./...

tidy: ## Tidy and verify dependencies
	$(GOMOD) tidy
	$(GOMOD) verify


run: build ## Run the application
	$(BUILD_DIR)/$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

docker-build: ## Build Docker image
	docker buildx build --platform linux/arm64 \
	--build-arg BUILD_DATE=$(BUILD_DATE) \
	--build-arg COMMIT_ID=$(COMMIT_ID) \
	--build-arg BUILD_VERSION=$(BUILD_VERSION) \
	-t $(REPO)/$(BINARY_NAME) .

docker-run: ## Run Docker container
	docker run --rm $(REPO)/$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

push: ## build and Push Docker image
	docker buildx build --platform linux/amd64,linux/arm64 \
	--build-arg BUILD_DATE=$(BUILD_DATE) \
	--build-arg COMMIT_ID=$(COMMIT_ID) \
	--build-arg BUILD_VERSION=$(BUILD_VERSION) \
	-t $(REPO)/$(BINARY_NAME) --push .

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

# A catch-all target to prevent errors when extra arguments are passed
%:
	@:
