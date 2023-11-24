# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0


DOCKER_DIRS = docsy-builder plantuml


.PHONY: doc
default: doc
doc:


.PHONY: lint
lint:
	@for d in $(TOOLS); do ($(MAKE) -C $$d all ); done

.PHONY: docker
docker:
	for d in $(DOCKER_DIRS) ;\
	do \
		docker build -f tools/docker/$$d/Dockerfile \
				--tag $$d:latest ./tools/docker/$$d ;\
	done;

.PHONY: tools
tools: docker

.PHONY: clean
clean:
	@for d in $(TOOLS); do ($(MAKE) -C $$d clean ); done
