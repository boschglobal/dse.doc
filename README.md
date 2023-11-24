<!--
Copyright 2023 Robert Bosch GmbH

SPDX-License-Identifier: Apache-2.0
-->

# Dynamic Simulation Environment - Documentation

## Introduction

Documentation project of the Dynamic Simulation Environment (DSE) Core Platform.


### Project Structure

```
L- content      Documentation content (used by Hugo to build the doc site).
L- doc          Additional documentation resources (e.g. yEd image sources).
L- licenses     Third Party Licenses.
L- static       Static content (e.g. images).
L- tools        Supporting tools.
```


## Usage

### Tools

Generated documentation is built using containerised tools. Those
tools can be built as follows:

```bash
$ git clone https://github.com/boschglobal/dse.doc.git
$ cd dse.doc
$ make tools
```

Alternatively, the latest Docker Images are available on ghcr.io and can be
used as follows:

```bash
$ export DOC_CDOCGEN_IMAGE=ghcr.io/boschglobal/dse-cdocgen:main
```


### Build

```bash
$ git clone https://github.com/boschglobal/dse.doc.git
$ cd dse.doc

# Pull is latest versions of linked documentation (i.e Hugo modules/git repos).
$ hugo mod get -u

# Build and serve documentation.
$ hugo server
$ hugo server --baseURL http://localhost:1313/dse.doc/
```

## Contribute

Please refer to the [CONTRIBUTING.md](./CONTRIBUTING.md) file.


## License

Dynamic Simulation Environment Documentation is open-sourced under the
Apache-2.0 license for code contributions and Creative Commons Attribution license
(CC-BY-4.0) for documentation contributions.
See the [LICENSE](LICENSE) and [NOTICE](./NOTICE) files for details.


### Third Party Licenses

[Third Party Licenses](licenses/)
