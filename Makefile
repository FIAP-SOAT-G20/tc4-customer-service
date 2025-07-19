.DEFAULT_GOAL := help

export PATH := $(shell go env GOPATH)/bin:$(PATH)


# Variables
LAMBDA_DIR     = .
BINARY_NAME    = bin/bootstrap
ZIP_NAME       = dist/function.zip
BUILD_OS       = linux
BUILD_ARCH     = amd64
MAIN_FILE=main.go

VERSION=$(shell git describe --tags --always --dirty)
TEST_PATH=./internal/...
TEST_COVERAGE_FILE_NAME=coverage.out
MONGO_URI = mongodb://admin:admin@localhost:27017/fastfood_10soat_g22_tc4?authSource=admin
LAMBDA_INPUT_FILE=test/data/api_gateway_proxy_request_event_payload_customer_not_found.json

# Go commands
AWSLAMBDARPCCMD ?= awslambdarpc
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) $MAIN_FILE
GOTEST=ENVIRONMENT=test $(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOTIDY=$(GOCMD) mod tidy
SHCMD=sh

# Looks at comments using ## on targets and uses them to produce a help output.
.PHONY: help
help: ALIGN=22
help: ## 📜 Print this message
	@echo "Usage: make <command>"
	@awk -F '::? .*## ' -- "/^[^':]+::? .*## /"' { printf "  make '$$(tput bold)'%-$(ALIGN)s'$$(tput sgr0)' - %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo

.PHONY: fmt
fmt: ## 🗂️  Format the code
	@echo  "🟢 Formatting the code..."
	$(GOCMD) fmt ./...
	@echo

.PHONY: build
build: fmt ## 🔨 Build the application
	@echo  "🟢 Building the application..."
	#$(GOBUILD) -v -gcflags='all=-N -l' -o bin/$(APP_NAME) $(MAIN_FILE)
	GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) $(GOBUILD) -ldflags="-s -w" -o $(LAMBDA_DIR)/$(BINARY_NAME) $(LAMBDA_DIR)/main.go
	@echo


.PHONY: package
package: build ## 📦 Package the binary into a .zip file for Lambda deployment
	@echo "📦 Packaging Lambda binary into zip..."
	mkdir -p ./dist
	zip -j $(LAMBDA_DIR)/$(ZIP_NAME) $(LAMBDA_DIR)/$(BINARY_NAME)
	@echo

.PHONY: start-lambda
start-lambda:  build  ## ▶  Start the lambda application locally to prepare to receive requests
	@echo "🟢 Starting lambda ..."
	# @_LAMBDA_SERVER_PORT=3300 $(GOCMD) run $(LAMBDA_DIR)/main.go
	@$(GOCMD) run $(LAMBDA_DIR)/main.go
	@echo

.PHONY: trigger-lambda
trigger-lambda: ## ⚡  Trigger lambda with the input file stored in variable $LAMBDA_INPUT_FILE
	@echo "🟢 Triggering lambda with event: $(LAMBDA_INPUT_FILE)"
	@PATH="$(shell go env GOPATH)/bin:$$PATH" \
		'$(AWSLAMBDARPCCMD)' -a localhost:3300 -e $(LAMBDA_INPUT_FILE)
	@echo

mock: ## Generate mocks
	@echo  "🟢 Generating mocks..."
	@go install go.uber.org/mock/mockgen@latest
	@mkdir -p internal/core/port/mocks
	@rm -rf internal/core/port/mocks/*
	@for file in internal/core/port/*.go; do \
		mockgen -source=$$file -destination=internal/core/port/mocks/`basename $$file _port.go`_mock.go -package=mocks; \
	done


.PHONY: test
test: lint ## 🧪 Run tests
	@echo  "🟢 Running tests..."
	@$(GOFMT) ./...
	@$(GOVET) ./...
	@$(GOTIDY)
	$(GOTEST) $(TEST_PATH) -race -v
	@echo

.PHONY: coverage
coverage: ## 🧪 Run tests with coverage
	@echo  "🟢 Running tests with coverage..."
# remove files that are not meant to be tested
	$(GOTEST) $(TEST_PATH) -coverprofile=$(TEST_COVERAGE_FILE_NAME).tmp
	@cat $(TEST_COVERAGE_FILE_NAME).tmp | grep -v "_mock.go" | grep -v "_request.go" | grep -v "_response.go" \
	| grep -v "_gateway.go" | grep -v "_datasource.go" | grep -v "_presenter.go" | grep -v "middleware" \
	| grep -v "config" | grep -v "route" | grep -v "util" | grep -v "database" \
	| grep -v "server" | grep -v "logger" | grep -v "httpclient" > $(TEST_COVERAGE_FILE_NAME)
	@rm $(TEST_COVERAGE_FILE_NAME).tmp
	$(GOCMD) tool cover -html=$(TEST_COVERAGE_FILE_NAME)
	@echo

.PHONY: clean
clean: ## 🧹 Clean up binaries and coverage files
	@echo "🔴 Cleaning up..."
	$(GOCLEAN)
	rm -f $(APP_NAME)
	rm -f $(TEST_COVERAGE_FILE_NAME)
	rm -f $(LAMBDA_DIR)/$(BINARY_NAME) $(LAMBDA_DIR)/$(ZIP_NAME)
	@echo


.PHONY: lint
lint: ## 🔍 Run linter
	@echo "🟢 Running linter..."
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.7 run --out-format colored-line-number
	@echo

.PHONY: install
install: ## 📦 Install dependencies
	@echo "🟢 Installing dependencies..."
	go mod download
	@go install github.com/blmayer/awslambdarpc@latest
	@echo

.PHONY: compose-up
compose-up: ## ▶  Start local database with docker compose
	@echo "🟢 Starting development environment..."
	docker compose pull
	docker-compose up -d --wait --build
	@echo

.PHONY: compose-down
compose-down: ## ■  Stops local database with docker compose
	@echo "🔴 Stopping development environment..."
	docker-compose down
	@echo

.PHONY: compose-clean
compose-clean: ## 🧹 Clean the application with docker compose, removing volumes and images
	@echo "🔴 Cleaning the application..."
	docker compose down --volumes --rmi all
	@echo

.PHONY: scan
scan: ## 🔍 Run security scan
	@echo  "🟠 Running security scan..."
	@go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 -show verbose ./...
	@go run github.com/aquasecurity/trivy/cmd/trivy@latest image --severity HIGH,CRITICAL $(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_APP):latest
	@echo
