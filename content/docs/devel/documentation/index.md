---
title: "Documentation Systems"
linkTitle: "Documentation"
weight: 20
---

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


### Content Organisation

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


## Doc Generation

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


### YAML based Schema Documentation

* OpenAPI schema definition with embedded documentation (including examples).
* Validation (swagger) and generation (widdershins).
* Additional templating (adding Hugo metadata) with sed.
