# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

# Usage
# =====
# alias dse-cdocgen='docker run --rm -it --volume $(pwd):/tmp/code ${DSE_CDOCGEN_IMAGE:-dse-cdocgen:latest}'
# dse-cdocgen --input test/testdata/header.h --output header.md --title Header --linktitle Header

FROM golang:bookworm AS builder
ADD ./go /src
RUN cd /src && make build


FROM debian:bookworm-slim
LABEL maintainer="timothy.rule@de.bosch.com"

RUN set -eux; \
    apt-get -y update; \
    apt-get -y install --no-install-recommends \
        clang \
        ; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/bin /usr/local/bin

WORKDIR /tmp/code
ENTRYPOINT ["/usr/local/bin/cdocgen"]
