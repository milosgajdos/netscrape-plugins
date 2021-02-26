# syntax = docker/dockerfile:1-experimental

FROM --platform=${BUILDPLATFORM} golang:1.15.8-alpine AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
            GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ./...

FROM base AS integration-tests
RUN --mount=target=. \
            --mount=type=cache,target=/root/.cache/go-build \
            go test -v ./...

FROM base AS unit-tests
RUN --mount=target=. \
            --mount=type=cache,target=/root/.cache/go-build \
            go test -test.short -v ./...

FROM golangci/golangci-lint:v1.37-alpine AS lint-base

FROM base AS lint
RUN --mount=target=. \
            --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
            --mount=type=cache,target=/root/.cache/go-build \
            --mount=type=cache,target=/root/.cache/golangci-lint \
            golangci-lint run --timeout 10m0s ./...
