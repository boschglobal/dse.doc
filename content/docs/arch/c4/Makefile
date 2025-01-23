# Copyright 2025 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

WORKSPACE 			?= dse_workspace

LITE_NAME			= c4lite
LITE_CONTAINER  	= structurizr/lite
CLI_NAME			= c4cli
CLI_CONTAINER  	   	= structurizr/cli
PLANTUML_CONTAINER 	?= dse-plantuml:latest

EXPORT_FORMATS = plantuml plantuml/c4plantuml mermaid dot


.PHONY: lite
lite: stop
	docker run -d --rm --name $(LITE_NAME) \
		-p 8080:8080 \
		-e STRUCTURIZR_WORKSPACE_FILENAME=$(WORKSPACE) \
		-v $$(pwd)/structurizr:/usr/local/structurizr \
		$(LITE_CONTAINER)


.PHONY: stop
stop:
	@-docker kill $(LITE_NAME)


.PHONY: generate
generate: generate_clean cli_export
	docker run -it --rm \
		--volume $$(pwd):/tmp \
		--workdir /tmp \
		$(PLANTUML_CONTAINER) \
			-o /tmp/images structurizr/export/plantuml
	cd images; for file in *.png; do mv "$$file" "$${file#structurizr-}"; done;


generate_clean:
	rm -rf structurizr/export
	rm -rf images
	rm -rf hsperfdata_root

cli_export: $(EXPORT_FORMATS)

$(EXPORT_FORMATS):
	docker run -i --rm --name $(CLI_NAME) \
		-v $$(pwd)/structurizr:/usr/local/structurizr \
		$(CLI_CONTAINER) \
		export \
			-w $(WORKSPACE).dsl \
			-format $@ \
			-output export/$(notdir $@)


.PHONY: clean
clean: stop generate_clean
