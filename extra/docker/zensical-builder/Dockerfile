# Copyright 2026 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0
#

ARG ZENSICAL_VERSION=latest

FROM zensical/zensical:${ZENSICAL_VERSION}

USER root
RUN apk add --no-cache bash git python3 py3-yaml

WORKDIR /docs
