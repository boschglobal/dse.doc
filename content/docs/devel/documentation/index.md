---
title: "Documentation Systems"
linkTitle: "Documentation"
weight: 20
---

The documentation system is built by pulling content from individual Git Repos into
a Hugo/Docsy project. The Hugo build system then consolidates that content into
a single documentation system which is then published.

> Hint: Documentation in each Git Repo needs to follow the same layout and structure as
the Documentation Repo for the content to be merged correctly.



## Layout

### Repo Documentation

When writing documentation construct a layout as described in the following section.
The entire`doc/content` folder of your Git Repo will be mounted (effectively merged) into the `doc/content` folder of the documentation system, therefore pay special attention to the following naming conventions to avoid conflicts with content from other repos:

* `repo` - Typically the tail part of your repo name (i.e `fsil.runnable` becomes `runnable`).
* `topic` - A unique topic name. If you are describing the architecture of a model, then the topic can be the name of that model (which is also typically the `repo` name).
* `model` - If your repo represents a model, then use that name (which is also typically the `repo` name).
* `tool` - The name of a tool or script that you are documenting.

> Hint: No need to over think it, just use your repo name as the stem for the name ofyour page bundle folder (Hugo term).


```
doc/
└── content/
    └── apis/                       <-- *** mount point for apis content ***
        └── <repo>/...              <-- generated API documentation
    └── docs/                       <-- *** mount point for docs content ***
        └── arch/
            └── <topic>|<model>/
                └── index.md        <-- documentation, primary content for leaf bundle
                └── image.png
        └── devel/
            └── <repo>_<topic>/
                └── examples/       <-- example code, belongs to the bundle
                └── index.md        <-- documentation
                └── image.png       <-- exported images (e.g. from yEd)
        └── examples/<repo>/...     <-- example code, alternative location
        └── user/
            └── models/
                └── <model>/
                    └── index.md    <-- documentation
            └── tools/
                └── <tool>/
                    └── index.md    <-- documentation
    yed/                            <-- yEd image source files
    Makefile                        <-- generate implementation
Makefile                            <-- generate target
```

### Mount Configuration

Content from Git Repos is mounted into the documentation system using the `hugo.toml` file.
Additionally the `go.mod` needs to be updated to reference the correct version of the Git Repo.


#### File: hugo.toml

```toml
...
# DSE ModelC
# -----------
[[module.imports]]
  path = "github.com/boschglobal/dse.modelc"
  disable = false
  ignoreConfig = true
  [[module.imports.mounts]]
    source = "doc/content/apis"
    target = "content/apis"
  [[module.imports.mounts]]
    source = "doc/content/docs"
    target = "content/docs"
...
```


#### File: go.mod

For repos with version 2+ it is necessary to manually update the go.mod file.

> Note: Modifications to `go.mod` will eventually be scripted.

```go
module github.com/boschglobal/dse.doc

go 1.19

require (
        github.com/boschglobal/dse.clib v1.0.5 // indirect
        github.com/boschglobal/dse.modelc v2.0.0+incompatible
        github.com/boschglobal/dse.schemas v1.1.4 // indirect
        github.com/boschglobal/dse.standards v1.0.3 // indirect
        // ...
)
```


## Generation

Several Repos/Projects have generated documentation. This documentation is
updated with the following generalised process:

1. Update documentation in the source files.

2. Run the `make generate` Makefile target to update the generated content.

3. Commit the updated content and push the changes upstream.

4. When the documentation is ready, tag a "patch" release on the Repo, the next
    time the documentation system updates, it will fetch the updated content.


Each of the generated documentation formats/systems are explained in the
following sections.


### C based API Documentation

* Markdown format documentation, embedded in C comment blocks.
* CDocGen toolchain for generation of documentation (tools/cdocgen).
* PlantUML images generated from embedded diagrams (tools/plantuml).
* Examples read from source files, and also included in build for quality assurance of example code.

An example of API Doc Generation is available in this [Makefile](https://github.com/boschglobal/dse.modelc/blob/main/doc/Makefile).


### YAML based Schema Documentation

* OpenAPI schema definition with embedded documentation (including examples).
* Validation (swagger) and generation (widdershins).
* Additional templating (adding Hugo metadata) with sed.

An example of Schema Doc Generation is available in this [Makefile](https://github.com/boschglobal/dse.schemas/blob/main/Makefile).



## Hugo / Docsy

The documentation system is built using Hugo with the Docsy theme. Content is
written in Markdown format. General information is available at these links:

* [Hugo](https://gohugo.io/)
* [Docsy](https://www.docsy.dev/)
* Markdown :
  * [Basic Syntax](https://www.markdownguide.org/basic-syntax/)
  * [Extended Syntax](https://www.markdownguide.org/extended-syntax/)
  * [Markdown Cheatsheet with GitHub flavour](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)
* [Hugo Page Bundles](https://gohugo.io/content-management/page-bundles/)


Most of the content is pulled from other Git Repos and mounted into the
content directory. Therefore, the same documentation may be reused in
several documentation systems.


### Hugo Content Organisation

> Warning: Hugo content organisation is not always explained well. For instance, "leaf means it has no children", really does not explain what a _leaf_ is ... and since when do leaves have children? or not have them? This [post](https://discourse.gohugo.io/t/explanation-of-page-bundles/36686) is helpful if you suspect that the Hugo documentation seems to be missing the point.

In Hugo, content is organised as [Page Bundles](https://gohugo.io/content-management/page-bundles), which are, in turn collections of [Page Resources](https://gohugo.io/content-management/page-resources/).
Page bundles are simply a collection of related files (i.e. page resources) all placed in the same directory.
Page bundles may either be a Leaf Bundle (directory contains file `index.md`) or a Branch Bundle (directory contains file `_index.md`).
A simple Page can be difficult to work with, especially if your content includes images.

```
content/
└── docs
    ├── page.md           <-- page
    ├── leaf              <-- page bundle, specifically a leaf bundle
    │   ├── leaf.jpg
    │   └── index.md
    └── branch            <-- page bundle, specifically a branch bundle
        ├── page-1.md
        ├── page-2.md
        ├── branch.jpg
        └── _index.md
```

In terms of arranging content, a Branch Bundle may contain other Bundles (Branch or Leaf), where as a Leaf Bundle may not contain other Bundles. This arrangement is reflected in the structure of the site.

