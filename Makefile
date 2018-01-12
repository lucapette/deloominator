SOURCE_FILES?=$$(go list ./...)
TEST_PATTERN?=.
TEST_OPTIONS?=

setup: ## Install all the build and lint dependencies
	go get -u --insecure github.com/golang/dep/cmd/dep
	go get golang.org/x/tools/cmd/stringer
	go get github.com/gobuffalo/packr/...
	go get github.com/alecthomas/gometalinter
	gometalinter --install
	dep ensure

embed:
	packr

build-server: embed ## Build the API server
	go build cmd/deloominator.go

build-ui: ## Build the UI
	cd ui && npm run build

build: build-ui build-server ## Build a dev version of deloominator

test-server: embed ## Run API tests
	bin/run-test go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

test-ui: ## Run UI tests
	cd ui && npm run test

test: test-server test-ui ## Run all the tests

lint-server: embed ## Run golang linters
	gometalinter --vendor --disable-all \
	-e packr \
	--enable=vet \
	--enable=gofmt \
	--enable=errcheck \
	--enable=deadcode \
	--enable=staticcheck \
	--enable=gosimple \
	--enable=structcheck \
	--enable=maligned \
	--enable=unparam \
	./...

lint-ui: ## Run JS linters
	cd ui && npm run eslint
	cd ui && npm run prettier:check

lint: lint-server lint-ui ## Rull all linters

run-server: build-server ## Run the API server
	bin/run deloominator

run-ui: ## Run the UI application
	cd ui && npm run start

# For now, it does not make sense to build the UI on CI as there no tests that
# rely on that.
stub-ui:
	touch ui/dist/index.html ui/dist/App.js ui/dist/App.js.map

ci: stub-ui build-server lint test ## Run all the tests and code checks

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build

.PHONY: ui
