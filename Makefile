PACKAGES=$(shell go list -tags "${BUILDTAGS}" ./... | grep -v /vendor/)
PKG=${PACKAGES}
INTEGRATION_PACKAGE=${PKG}

TESTFLAGS_RACE=-race
GOFILES=$(shell find . -type f -name '*.go')
GO_TAGS=$(if $(BUILDTAGS),-tags "$(BUILDTAGS)",)

# Flags passed to `go test`
TESTFLAGS ?= -v $(TESTFLAGS_RACE)
TESTFLAGS_PARALLEL ?= 8

test: lint unit-tests

.PHONY: lint
lint: export DOCKER_BUILDKIT=1
lint:
	@docker build . --target lint

.PHONY: unit-tests
unit-tests: export DOCKER_BUILDKIT=1
unit-tests:
	@docker build . --target unit-tests

.PHONY: integration-tests
integration-tests: export DOCKER_BUILDKIT=1
integration-tests: | start-dgraph
	@docker build . --target integration-tests

.PHONY: local-tests
local-tests:
	@go test ${TESTFLAGS} -parallel ${TESTFLAGS_PARALLEL} ${INTEGRATION_PACKAGE}
