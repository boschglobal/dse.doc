# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

# See https://github.com/klakegg/docker-hugo for details of base image.

FROM klakegg/hugo:ext-debian-ci
LABEL maintainer="timothy.rule@de.bosch.com"

SHELL ["/bin/bash", "-c"]

RUN set -eux;\
    git config --global --add safe.directory /src ;\
    npm install -D autoprefixer;\
    npm install -D postcss-cli;\
    npm install -D postcss;\
    apt-get clean;\
    rm -rf /var/lib/apt/lists/*
