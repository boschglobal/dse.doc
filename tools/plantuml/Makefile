# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

DOCKER_PREFIX := dse
DOCKER_IMAGE := $(shell basename $(shell pwd))

.PHONY: all
all: docker


.PHONY: docker
docker:
	docker build -f ./Dockerfile \
		--tag $(DOCKER_PREFIX)-$(DOCKER_IMAGE):latest .


.PHONY: clean
clean:
	docker images -q $(DOCKER_PREFIX)-$(DOCKER_IMAGE) | xargs -r docker rmi -f
	docker images -q */*/$(DOCKER_PREFIX)-$(DOCKER_IMAGE) | xargs -r docker rmi -f
	docker images -qf dangling=true | xargs -r docker rmi

