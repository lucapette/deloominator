SOURCE_FILES?=$$(go list ./... | grep -v '/deloominator/vendor/')
TEST_PATTERN?=.
TEST_OPTIONS?=

setup: ## Install all the build and lint dependencies
	go get -u golang.org/x/tools/cmd/stringer
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

embed:
	go-bindata -prefix "pkg/db/migrations" -o pkg/db/migrations.go -pkg db pkg/db/migrations
	go-bindata -o pkg/api/static.go -pkg api ui/dist/index.html ui/dist/App.js ui/dist/App.js.map
	gofmt -s -w pkg/api/static.go
	gofmt -s -w pkg/db/migrations.go

build-api: embed ## Build the API server
	go build cmd/deloominator.go

build-ui: ## Build the UI
	cd ui && yarn build

build: build-ui build-api ## Build a dev version of deloominator

test: embed ## Run all the tests
	bin/run-test go test $(TEST_OPTIONS) -cover $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

lint: embed ## Run all the linters
	gometalinter --vendor --disable-all \
	--enable=vet \
	--enable=gofmt \
	--enable=errcheck \
	./...

run-api: build-api ## Run the API server
	bin/run deloominator

run-ui: ## Run the UI application
	cd ui && yarn start

# For now, it doesn't make sense to build the UI on travis as there
# no tests that rely on that.
stub-ui:
	touch ui/dist/index.html ui/dist/App.js ui/dist/App.js.map

ci: stub-ui build-api lint test ## Run all the tests and code checks

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build

.PHONY: ui
