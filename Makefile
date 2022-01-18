BUILD_OPTS?="-ldflags=-s -w"
.DEFAULT_GOAL := help

build: ## Build binary file
	go build ${BUILD_OPTS} -o ./wfreq .

test: build ## Test with Kafka's metamorphosis
	./wfreq -i ./tests/metamorphosis.txt

help: ## Show help menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

