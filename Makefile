BINPATH ?= build

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

SEED_COUNT         ?= 1000
SEED_PDF_COUNT     ?= 100
SEED_DATA_COUNT    ?= 100
SEED_PATH         ?= /test/path
SEED_COLLECTION_ID ?= test-collection

.PHONY: all
all: audit test build

.PHONY: audit
audit: ## Runs checks for security vulnerabilities on dependencies (including transient ones)
	dis-vulncheck

.PHONY: build
build: ## Builds binary of application code and stores in bin directory as dis-legacy-cache-purger
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dis-legacy-cache-purger

.PHONY: convey
convey: ## Runs unit test suite and outputs results on http://127.0.0.1:8080/
	goconvey ./...

.PHONY: debug
debug: ## Used to run code locally in debug mode
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dis-legacy-cache-purger
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dis-legacy-cache-purger

.PHONY: delimiter-%
delimiter-%:
	@echo '===================${GREEN} $* ${RESET}==================='

.PHONY: fmt
fmt: ## Run Go formatting on code
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test: ## Runs unit tests including checks for race conditions and returns coverage
	go test -race -cover ./...

.PHONY: test-component
test-component:
	exit

.PHONY: seed
seed: ## Populate MongoDB with CacheTime documents
	mongosh --eval "var count=$(SEED_COUNT); var pdfCount=$(SEED_PDF_COUNT); var dataCount=$(SEED_DATA_COUNT); var path='$(SEED_PATH)'; var collectionID='$(SEED_COLLECTION_ID)';" scripts/seed.js

.PHONY: help
help: ## Show help page for list of make targets
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
