.PHONY: help
.DEFAULT_GOAL := help

ifndef GOOS
	GOOS := linux
endif
ifndef GOARCH
	GOARCH := amd64
endif

ENV=export CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)

build: ## Build awsProfileSwitcher
	$(ENV) &&  go build -v .

lint: ## Lint
	docker run --rm -v $(pwd):/data cytopia/golint .

test: ## Run unit tests
	go test -v .

clean: ## Clean DIR
	rm awsProfileSwitcher

help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
