# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0


SRC           := $(shell find . -type f -name '*.go' -print)
BIN_DIR       := $(CURDIR)/bin
CLI_TOOLS     ?= $(BIN_DIR)/cdocgen
GOFLAGS       :=
DOCKER_DIRS   := build/cdocgen
DOCKER_PREFIX := dse


.PHONY: all
all: build test


.PHONY: build
build: cmd

cmd: $(CLI_TOOLS)

$(CLI_TOOLS): $(SRC)
	go build $(GOFLAGS) -o $@ ./cmd/$(subst $(BIN_DIR)/,,$@)


.PHONY: test
test:
	go test ./...
	go test ./... -tags e2e


.PHONY: lint
lint:
	go fmt -x ./...
	PATH=~/go/bin:$$PATH; goimports -w -v .
	docker run -t --rm \
		-v $$(pwd):/app \
		-v ~/.cache/golangci-lint/v1.54.2:/root/.cache \
		-w /app \
		golangci/golangci-lint:latest \
		golangci-lint run -v ./...
