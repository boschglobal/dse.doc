# Copyright 2023 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

# Usage
# =====
# alias dse-plantuml='docker run --rm -it --volume $(pwd):/tmp/code ${DSE_PLANTUML_IMAGE:-ghcr.io/boschglobal/dse-plantuml:latest}'
# dse-plantuml -tpng index.md

FROM debian:bookworm-slim
LABEL maintainer="timothy.rule@de.bosch.com"

RUN set -eux; \
    apt-get -y update; \
    apt-get -y install --no-install-recommends \
        curl \
        default-jre \
        graphviz \
        ; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

RUN curl -fSL https://github.com/plantuml/plantuml/releases/download/v1.2023.11/plantuml-asl-1.2023.11.jar \
        -o /usr/local/bin/plantuml-asl.jar

WORKDIR /tmp/code
ENTRYPOINT ["java", "-Djava.awt.headless=true", "-jar", "/usr/local/bin/plantuml-asl.jar"]
