# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0
#
#
# Usage (via Makefile targets – do not invoke docker directly):
#   make build    – build this image as dse-zensical-local:<version>
#   make fetch    – populate content/ from sibling repositories
#   make site     – build the static site into site/
#   make serve    – live-preview at http://localhost:8000
#   make validate – strict build (fails on broken links)
#
# Official image: https://hub.docker.com/r/zensical/zensical

FROM zensical/zensical:latest

USER root
RUN apk add --no-cache git python3 py3-yaml

WORKDIR /docs
