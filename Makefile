# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0


TOOLS := tools/cdocgen tools/plantuml
CONTAINERS = docsy-builder


.PHONY: doc
default: doc
doc:


.PHONY: lint
lint:
	@for d in $(TOOLS); do ($(MAKE) -C $$d all ); done

.PHONY: docker
docker:
	for d in $(CONTAINERS) ;\
	do \
		docker build -f tools/$$d/Dockerfile \
				--tag $$d:latest ./tools/$$d ;\
	done;

.PHONY: tools
tools:
	@for d in $(TOOLS); do ($(MAKE) -C $$d all ); done

.PHONY: clean
clean:
	@for d in $(TOOLS); do ($(MAKE) -C $$d clean ); done
