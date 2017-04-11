SOURCE_FILES?=$$(go list ./... | grep -v '/deloominator/vendor/')
TEST_PATTERN?=.
TEST_OPTIONS?=

setup: ## Install all the build and lint dependencies
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kisielk/errcheck

test: build-api ## Run all the tests
	go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

lint: ## Run all the linters
	go vet $(SOURCE_FILES)
	errcheck $(SOURCE_FILES)

ci: build-ui build-api lint test ## Run all the tests and code checks

build-api:
	go-bindata -o pkg/api/static.go -pkg api ui/dist/index.html ui/dist/App.js ui/dist/App.js.map
	go build cmd/deloominator.go

build-ui: ## Build the UI
	cd ui && yarn build

build: build-ui build-api ## Build a dev version of deloominator
	gofmt -w pkg/api/static.go

run-api: build-api ## Run the API server
	bin/run deloominator

run-ui: ## Run the UI application
	cd ui && yarn start

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build

.PHONY: ui
