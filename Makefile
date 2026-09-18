# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0


###############
## Docker Images.
ZENSICAL_VERSION ?= 0.0.62
BUILDER_IMAGE     ?= zensical-builder:latest

###############
## Build parameters.
ZENSICAL_PORT    ?= 8000
ZENSICAL_CONFIG  ?= zensical.toml
BUILT_CONFIG     ?= build/zensical.generated.toml
REPOS_FILE        ?= repos.txt
BUILD_DIR         ?= build
REPOS_DIR         ?= $(BUILD_DIR)/repos
STAGE_DIR         ?= $(BUILD_DIR)/staged
SOURCE_MAP        ?= $(BUILD_DIR)/source-map.tsv
SITE_DIR          ?= $(BUILD_DIR)/site
DOCKER_DIRS       ?= cdocgen plantuml zensical-builder

###############
## Scripts
GEN_NAV_SCRIPT    ?= scripts/gen_nav.py
METADATA_SCRIPT   ?= scripts/metadata-update.py
FIX_LINKS_SCRIPT  ?= scripts/fix-source-links.py
CLONE_SCRIPT      ?= scripts/pull-doc.sh
COPY_SCRIPT       ?= scripts/copy-doc.sh


.DEFAULT := help

.PHONY: help print-zensical-version pull docker build run clean cleanall

help:
	@echo "Available targets:"
	@echo "  build         Build the documentation server."
	@echo "  run           Run the documentation server."
	@echo "  pull          Pull doc content from repos."
	@echo "  docker        Build Docker images."
	@echo "  clean         Clean the build dir."
	@echo "  cleanall      Remove all cached build/repo content."
	@echo "Local development commands:"
	@echo "  make build"
	@echo "  make run"

print-zensical-version:
	@printf '%s\n' "$(ZENSICAL_VERSION)"

pull:
	REPOS_DIR="$(REPOS_DIR)" bash $(CLONE_SCRIPT) $(REPOS_FILE)

docker:
	@for d in $(DOCKER_DIRS); do \
		docker build -f extra/docker/$$d/Dockerfile \
			--build-arg ZENSICAL_VERSION=$(ZENSICAL_VERSION) \
			--tag $$d:latest ./extra/docker/$$d; \
	done

build:
	@if [ ! -d "$(REPOS_DIR)" ]; then \
		$(MAKE) pull; \
	fi

	docker run --rm \
		--entrypoint sh \
		--user $$(id -u):$$(id -g) \
		-v "$(CURDIR):/docs" \
		-w /docs \
		$(BUILDER_IMAGE) \
		-c '\
			mkdir -p $(REPOS_DIR) $(STAGE_DIR) $(SITE_DIR) && \
			python3 /docs/$(METADATA_SCRIPT) && \
			bash /docs/$(COPY_SCRIPT) $(REPOS_DIR) $(STAGE_DIR) && \
			if [ -f content/index.md ]; then \
				cp -a content/index.md $(STAGE_DIR)/index.md; \
				printf "%s\\t%s\\t%s\\n" \
					"index.md" \
					"dse.doc" \
					"content/index.md" >> $(SOURCE_MAP); \
			fi'

	docker run --rm \
		--entrypoint python3 \
		--user $$(id -u):$$(id -g) \
		-v "$(CURDIR):/docs" \
		-w /docs \
		$(BUILDER_IMAGE) \
		/docs/$(GEN_NAV_SCRIPT) \
		--input $(ZENSICAL_CONFIG) \
		--output $(BUILT_CONFIG) \
		--site-dir $(STAGE_DIR) \
		--site-output-dir $(SITE_DIR)

	docker run --rm \
		--user $$(id -u):$$(id -g) \
		-v "$(CURDIR):/docs" \
		$(BUILDER_IMAGE) \
		build -f $(BUILT_CONFIG)

	docker run --rm \
		--entrypoint python3 \
		--user $$(id -u):$$(id -g) \
		-v "$(CURDIR):/docs" \
		-w /docs \
		$(BUILDER_IMAGE) \
		/docs/$(FIX_LINKS_SCRIPT) \
		$(SITE_DIR) \
		$(SOURCE_MAP)

run:
	docker run --rm -it \
		--entrypoint sh \
		-v "$(CURDIR)/$(SITE_DIR):/tmp/root/dse.doc:ro" \
		-p $(ZENSICAL_PORT):$(ZENSICAL_PORT) \
		$(BUILDER_IMAGE) \
		-c 'printf "%s\n" \
			"<!doctype html>" \
			"<html><head><meta http-equiv=\"refresh\" content=\"0; url=/dse.doc/\"></head></html>" \
			> /tmp/root/index.html && \
		exec python3 -m http.server \
			$(ZENSICAL_PORT) \
			--bind 0.0.0.0 \
			--directory /tmp/root'

clean:
	rm -rf $(BUILD_DIR)/

cleanall: clean
	rm -rf .cache/ .zensical/

cleandocker:
	@for d in $(DOCKER_DIRS); do \
		docker image rm -f $$d:latest 2>/dev/null || true; \
	done
	docker image prune -f
	docker container prune -f
	docker volume prune -f