SOURCE_FILES?=$$(go list ./... | grep -v '/deluminator/vendor/')

setup: ## Install all the build and test dependencies
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kisielk/errcheck

test: ## Run all the tests
	go test $(SOURCE_FILES) -timeout=30s

ci: ## Run all the tests and code checks
	errcheck $(SOURCE_FILES)

assets: ## Embed static assets
	go-bindata -o api/static.go -pkg api assets/...

build: assets ## Build a beta version of deluminator
	go build
	gofmt -w api/static.go

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: assets
