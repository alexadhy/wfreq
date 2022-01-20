BUILD_OPTS?="-ldflags=-s -w"
.DEFAULT_GOAL := help

build: format ## Build binary file
	go build ${BUILD_OPTS} -o ./wfreq-svc ./cmd/service
	go build ${BUILD_OPTS} -o ./wfreq-cli ./cmd/client

format: ## Format code
	goimports -w -local github.com/alexadhy/wfreq ./cmd
	goimports -w -local github.com/alexadhy/wfreq ./internal

help: ## Show help menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

