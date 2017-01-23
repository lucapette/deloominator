SOURCE_FILES?=$$(go list ./... | grep -v '/deluminator/vendor/')
TEST_PATTERN?=.
TEST_OPTIONS?=

setup: ## Install all the build and lint dependencies
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kisielk/errcheck

test: ## Run all the tests
	go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

lint: ## Run all the linters
	go vet $(SOURCE_FILES)
	errcheck $(SOURCE_FILES)

ci: lint test ## Run all the tests and code checks

assets: ## Embed static assets
	go-bindata -o pkg/api/static.go -pkg api assets/...

build: assets ## Build a beta version of deluminator
	go build
	gofmt -w pkg/api/static.go

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: assets
