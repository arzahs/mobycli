BINARY=mobycli

.DEFAULT_GOAL := help

dep:
	dep ensure -v -vendor-only

build: ## Build the binary
	CGO_ENABLED=0 go build $(LDFLAGS) -v -a -installsuffix cgo -o $(BINARY) .

run:
	go run .

test: ## Run the unit tests
	go test -race -v $(shell go list ./... | grep -v /vendor/)

lint: ## Lint all files
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

vet: ## Run the vet tool
	go vet $(shell go list ./... | grep -v /vendor/)

clean: ## Clean up build artifacts
	go clean

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'