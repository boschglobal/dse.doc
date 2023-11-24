# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0


TOOLS := tools/plantuml


.PHONY: doc
default: doc
doc:


.PHONY: lint
lint:
	@for d in $(TOOLS); do ($(MAKE) -C $$d all ); done


.PHONY: tools
tools:
	@for d in $(TOOLS); do ($(MAKE) -C $$d all ); done


.PHONY: clean
clean:
	@for d in $(TOOLS); do ($(MAKE) -C $$d clean ); done
