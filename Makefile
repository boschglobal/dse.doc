# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0


CONTAINERS = cdocgen docsy-builder plantuml


.PHONY: doc
default: doc
doc:

.PHONY: docker
docker:
	for d in $(CONTAINERS) ;\
	do \
		docker build -f extra/docker/$$d/Dockerfile \
				--tag $$d:test ./extra/docker/$$d ;\
	done;
