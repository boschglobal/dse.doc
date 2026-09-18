---
title: "Documentation Systems"
linkTitle: "Documentation"
weight: 20
---

The documentation system is built by aggregating content from individual Git repositories into a Zensical documentation site. During the build, content from the source repositories is staged into the generated content tree and then rendered by Zensical into the published documentation.

> Hint: Documentation in each Git repository should follow the same layout and structure as the documentation repository so that the staged content merges cleanly.

## Layout

### Repo Documentation

When writing documentation, use the layout described below. The content under the repository's doc/content directory is staged into the generated documentation tree, so keep folder names unique and consistent to avoid collisions with content from other repositories.

* `repo` - Typically the tail part of the repository name (for example, `dse.modelc` becomes `modelc`).
* `topic` - A unique topic name. If you are describing the architecture of a model, the topic can be the model name.
* `model` - If the repository represents a model, use that name.
* `tool` - The name of a tool or script that you are documenting.

> Hint: A simple convention is to use the repository name as the stem for the bundle or page folder.

```text
doc/
└── content/
    └── apis/                       <-- mount point for generated API documentation
        └── <repo>/...
    └── docs/                       <-- mount point for narrative documentation
        └── arch/
            └── <topic>|<model>/
                └── index.md
                └── image.png
        └── devel/
            └── <repo>_<topic>/
                └── examples/
                └── index.md
                └── image.png
        └── examples/<repo>/...
        └── user/
            └── models/
                └── <model>/
                    └── index.md
            └── tools/
                └── <tool>/
                    └── index.md
    yed/                            <-- yEd source files
    Makefile                        <-- generation helpers
Makefile                            <-- top-level build target
```

### Content Staging Configuration

Content from Git repositories is staged into the generated documentation tree by the fetch step defined in the repository's sources manifest and helper scripts. Zensical then reads from the staged content tree configured in the site settings.

#### File: sources.yaml

The content aggregation is driven by the repository manifest and the fetch script. This is where source repositories and target paths are mapped.

#### File: zensical.toml

Zensical uses the site configuration in zensical.toml. Snippets are resolved from the staged content tree by setting the snippets base path to the generated staging directory.

```toml
[project.markdown_extensions.pymdownx.snippets]
base_path = ["build/staged"]
check_paths = true
```

## Generation

Generated documentation is updated through the repository build workflow:

1. Update the source documentation in the relevant repository.
2. Run the fetch and staging step to copy the content into the generated tree.
3. Run the documentation build to render the site.
4. Commit the updated content and push it upstream.

The most common build targets are:

* `make fetch` - stage content from sibling repositories into the build tree.
* `make site` - build the Zensical documentation site.
* `make validate` - run the strict validation build.

Each generated documentation format or workflow is described in the following sections.

### C based API Documentation

* Markdown format documentation embedded in C comment blocks.
* CDocGen for generation of documentation.
* PlantUML images generated from embedded diagrams.
* Example source files are also included in the build as a quality check.

An example of API doc generation is available in this [Makefile](https://github.com/boschglobal/dse.modelc/blob/main/doc/Makefile).

### YAML based Schema Documentation

* OpenAPI schema definitions with embedded documentation.
* Validation and generation with Swagger and Widdershins.
* Additional templating for metadata and documentation integration.

An example of schema doc generation is available in this [Makefile](https://github.com/boschglobal/dse.schemas/blob/main/Makefile).

## Zensical

The documentation site is built with Zensical. Content is written in Markdown and can include code examples directly from staged source files by using Zensical snippet syntax.

Example snippet usage:

```md
--8<-- "apis/modelc/examples/model_interface.c"
```

This form is preferred over old Hugo-style include blocks. The snippet path is resolved relative to the staged content tree configured in zensical.toml, so the path must match the location visible under the configured base path.

Most content is pulled from other Git repositories and staged into the build tree before the Zensical render step. This makes it possible to reuse the same documentation in multiple documentation systems while keeping the source content in the originating repository.

