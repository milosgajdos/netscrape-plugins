PACKAGES=$(shell go list -tags "${BUILDTAGS}" ./... | grep -v /vendor/)
PKG=${PACKAGES}
INTEGRATION_PACKAGE=${PKG}

TESTFLAGS_RACE=-race
GOFILES=$(shell find . -type f -name '*.go')
GO_TAGS=$(if $(BUILDTAGS),-tags "$(BUILDTAGS)",)

# Flags passed to `go test`
TESTFLAGS ?= -v $(TESTFLAGS_RACE)
TESTFLAGS_PARALLEL ?= 8

test: lint local-tests

start-dgraph:
	docker-compose -f docker-compose.yml up -d zero alpha ratel

clean:
	docker-compose -f docker-compose.yml down --volumes

.PHONY: lint
lint: ## Run lint checks (requires https://golangci-lint.run)
	@golangci-lint run -v --timeout 10m0s ./...

.PHONY: local-tests
local-tests: | start-dgraph
	@go test ${TESTFLAGS} -parallel ${TESTFLAGS_PARALLEL} ${INTEGRATION_PACKAGE}
