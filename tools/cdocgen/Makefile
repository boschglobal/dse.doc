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
all: build test docker


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


.PHONY: docker
docker:
	for d in $(DOCKER_DIRS) ;\
	do \
		docker build -f $$d/Dockerfile \
			--tag $(DOCKER_PREFIX)-$$(basename $$d):latest . ;\
	done;


.PHONY: clean
clean:
	@rm -rf '$(BIN_DIR)'
	for d in $(DOCKER_DIRS) ;\
	do \
		docker images -q $(DOCKER_PREFIX)-$$(basename $$d) | xargs -r docker rmi -f ;\
		docker images -q */*/$(DOCKER_PREFIX)-$$(basename $$d) | xargs -r docker rmi -f ;\
	done;
	docker images -qf dangling=true | xargs -r docker rmi

