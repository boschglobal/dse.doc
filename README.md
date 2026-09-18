<!--
Copyright 2026 Robert Bosch GmbH

SPDX-License-Identifier: Apache-2.0
-->

# Dynamic Simulation Environment - Documentation

## Introduction

Documentation project of the Dynamic Simulation Environment (DSE) Core Platform.
Built with [Zensical](https://zensical.org) (latest stable) using a fully
containerised Docker + Makefile workflow — no local toolchain installation required.


### Project Structure

```text
dse.doc
├── content
│   ├── docs                    <-- Documentation content aggregated from DSE repositories.
│   ├── stylesheets             <-- Bosch brand-colour overrides.
│   └── index.md                <-- Documentation site landing page.
├── doc                         <-- Additional documentation resources (e.g. yEd sources).
├── extra
│   └── docker                  <-- Dockerfiles for documentation build tools.
├── overrides                   <-- Zensical theme customisation directory.
├── scripts                     <-- Project maintenance and build scripts.
│   ├── copy-doc.sh             <-- Stage content into the generated docs tree.
│   ├── gen_nav.py              <-- Generate the navigation configuration.
│   ├── metadata-update.py      <-- Update generated documentation metadata.
│   └── pull-doc.sh             <-- Clone documentation trees from source repositories.
├── .github
│   └── workflows               <-- CI workflows for building and publishing documentation.
├── Dockerfile                  <-- Local development image.
├── Makefile                    <-- Build, run, and clean targets.
├── repos.txt                   <-- Source repository manifest.
└── zensical.toml               <-- Zensical site configuration.
```


## Usage

```bash
# Clone the repository.
git clone https://github.com/boschglobal/dse.doc.git
cd dse.doc

# Build the documentation images.
make docker

# Pull documentation and build the site.
make pull
make build

# Serve the site at http://localhost:8000.
make run

# Remove generated build artifacts.
make clean

# Remove generated artifacts and cached content.
make cleanall
```


## Configuration

### sources.yaml

Defines which sub-directories from each sibling repository are copied into
`content/` before the Zensical build runs.  Edit this file to add, remove, or
remap content sources.  `repo: "."` entries are skipped (content lives in
`dse.doc` itself).

> **Note:** Zensical does not yet have native multi-repo ("Subprojects")
> support.  The `sources.yaml` + `scripts/fetch_sources.py` approach is the recommended
> workaround until that feature ships.  Track progress at
> https://zensical.org/about/roadmap/#subprojects


## Contribute

Please refer to the [CONTRIBUTING.md](./CONTRIBUTING.md) file.


## License

Dynamic Simulation Environment Documentation is open-sourced under the
Apache-2.0 license for code contributions and Creative Commons Attribution license
(CC-BY-4.0) for documentation contributions.
See the [LICENSE](LICENSE) and [NOTICE](./NOTICE) files for details.


### Third Party Licenses

[Third Party Licenses](licenses/)
