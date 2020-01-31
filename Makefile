.PHONY: help
.DEFAULT_GOAL := help

build: ## Build awsProfileSwitcher
	go build -v .

lint: ## Lint
	docker run --rm -v $(pwd):/data cytopia/golint .

test: ## Run unit tests
	go test -v .

clean: ## Clean DIR
	rm awsProfileSwitcher

help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
