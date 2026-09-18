# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

# Usage
# =====
# alias dse-cdocgen='docker run --rm -it --volume $(pwd):/tmp/code ${DSE_CDOCGEN_IMAGE:-ghcr.io/boschglobal/dse-cdocgen:latest}'
# dse-cdocgen --input test/testdata/header.h --output header.md --title Header --linktitle Header

FROM golang:bookworm AS builder
ADD ./go /src
RUN cd /src && make build

FROM zensical/zensical:latest

USER root
RUN apk add --no-cache git python3 py3-yaml

COPY --from=builder /src/bin/cdocgen /usr/local/bin/cdocgen

WORKDIR /docs

ENTRYPOINT ["/usr/local/bin/cdocgen"]