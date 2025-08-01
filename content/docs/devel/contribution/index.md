---
title: "Contributions Guide"
linkTitle: "Contribution"
weight: 20
---


## Synopsis

Developers wishing to contribute, as well as those having no choice in the matter, will find
this documentation page full of useful information to guide you towards finalising your contribution.


## Git Commands

### Setup and Configure Git

```bash
# Set user name and email (for signoff).
git config --global --add user.name "User Name (dept)"
git config --global --add user.email "user.name@de.bosch.com"
```


### Working with Commits

#### Sign Off

```bash
# Add a signoff to a commit.
git commit -s -m "Commit message."

# Append a signoff to an existing commit.
git commit --amend -s --no-edit

# Add signoff to several commits.
git rebase --signoff HEAD~3
```


#### Cherry-Pick

Particularly useful when recovering from mistakes, or where you simply need to
get a commit onto your branch.

```bash
# Cherry pick between branches.
git checkout --track origin/devel
git checkout foo
git cherry-pick a253c6359aa85da5627caf2f746282ac0e53cea1
```


### Working with PRs

#### Rebase a PR and push to requesters branch.

```bash
# Checkout the PR locally.
git fetch origin pull/86/head:PR-BRANCH
git checkout PR-BRANCH

# Rebase the PR and resolve any conflicts.
git pull --rebase upstream main

# Push the PR back to the originatory fork/branch (force if necessary).
git push https://github.boschdevcloud.com/PR-USER/dse.modelc.git HEAD:devel
git push https://github.boschdevcloud.com/PR-USER/dse.modelc.git HEAD:devel --force
```

#### Squash a PR locally

```bash
# Checkout the PR locally.
git fetch origin pull/86/head:PR-BRANCH
git checkout PR-BRANCH

# Soft reset to the commit (SHA) _prior_ to the PR commits
$ git log --oneline
$ git reset --soft 5df673c
$ git commit -am "MDF API implementation."
$ git commit --amend --author="" --no-edit

# Rebase the PR and squash in any edits.
$ git add .
$ git commit --amend --no-adit
$ git pull --rebase upstream main
# or
$ git pull --rebase origin devel

# Set the author.
git commit --amend --author="John Doe (VM/EMT3) <John.Doe@de.bosch.com>" --no-edit

# Fix any issues, dd a signoff, merge to main/devel branch.
$ git commit --amend -s --no-edit
$ git checkout main
# or
$ git checkout devel
$ git pull . PR-MDF --rebase

# Push the PR back to the originatory fork/branch.
$ git push upstream HEAD:main
# and/or
$ git push origin HEAD:devel
```


#### Modify an existing PR on the requester fork/branch.

```bash
# Update local devel branch (push local changes prior).
$ git remote -v
origin  https://github.boschdevcloud.com/fsil/dse.modelc.git (fetch)
origin  https://github.boschdevcloud.com/fsil/dse.modelc.git (push)
$ git switch origin/devel
$ git pull --rebase    # Ensure you have no local commits at this point.

# Create a new branch and pull in the PR (no merge).
git switch -c USER-PR-BRANCH origin/devel
git pull https://github.boschdevcloud.com/USER/dse.modelc.git PR-BRANCH --rebase

# Edit and commit changes.
git commit -m"Fix PR issue."

 # Push back to the PR (and delete the local branch).
git push https://github.boschdevcloud.com/USER/dse.modelc.git HEAD:PR-BRANCH
git switch devel
git branch -D USER-PR-BRANCH devel
```


### Working with Remotes

#### Setup for OSS Contributions

When working directly with OSS repos it is helpful to maintain a clear set
of remotes (and avoid mistakes). The following example shows one approach:


```bash
# Add the origin remote as the development mirror (using clone).
$ git clone https://github.boschdevcloud.com/fsil/dse.modelc.git

# Add the upstream remote as the OSS repo.
git remote add upstream https://github.com/boschglobal/some.repo.git
git fetch upstream main
git fetch upstream main --tags

# Add a fsil.fork remote for your own PRs.
git remote add fsil.fork https://github.com/boschglobal/some.repo.git
git fetch fsil.fork
git switch -c PR-BRANCH fsil.fork/main

# Review the remote setup.
$ git remote -v
fsil.fork       https://github.boschdevcloud.com/USER/fsil.dse.modelc.git (fetch)
fsil.fork       https://github.boschdevcloud.com/USER/fsil.dse.modelc.git (push)
origin          https://github.boschdevcloud.com/fsil/dse.modelc.git (fetch)
origin          https://github.boschdevcloud.com/fsil/dse.modelc.git (push)
upstream        https://github.com/boschglobal/dse.modelc.git (fetch)
upstream        https://github.com/boschglobal/dse.modelc.git (push)
```


#### Pushing to OSS Remotes

```bash
# Pull in upstream changes (avoid merge).
$ git branch
* devel
  main
$ pull upstream main --rebase

# Push changes from 'devel' to 'main'.
git push upstream devel:main

# Alternative, push from the current HEAD location to 'main'.
git push upstream HEAD:main
```


## GitHub

### Dependabot Auto Merge

Setup your GitHub Repo for Dependabot:

1. Settings -> General -> Pull Requests -> Allow auto-merge [select]
2. Rules -> Ruleset -> [New branch ruleset] ->
    1. Ruleset Name : default (or main)
    2. Bypass list : Repository admin & Maintain roles
    3. Target branches : default (or main)
    4. Rules : Require status checks to pass (these are workflow names)
        - Lint Code Base
        - build (linux-amd64)
        - test_go
    5. [Create]
    6. Enforcement status -> [Active] -> [Save changes]
3. Issues -> Labels -> New Label -> add "automerge", "dependabot"


Add the following files to configure Dependabot and add an Auto-Merge workflow.

#### .github/dependabot.yml
```yaml
---
version: 2
updates:
  - package-ecosystem: "gomod"
    directories:
      - "**/*"
    schedule:
      interval: "daily"
    groups:
      gomod-automerge:
        update-types:
          - "patch"
          - "minor"
    labels:
      - "automerge"
      - "dependabot"
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "daily"
    groups:
      actions-merge:
        update-types:
          - "patch"
          - "minor"
          - "major"
    labels:
      - "dependabot"
```

#### .github/workflows/dependabot.yaml

> Hint:  Make sure to match `github.repository` to your repo name.

```yaml
---
name: Dependabot auto-merge
on:  # yamllint disable-line rule:truthy
  pull_request:

permissions:
  contents: write
  pull-requests: write

jobs:
  dependabot:
    runs-on: ubuntu-latest
    if: github.event.pull_request.user.login == 'dependabot[bot]' && github.repository == 'boschglobal/REPO_NAME'
    steps:
      - name: Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@v2
        with:
          github-token: "${{ secrets.GITHUB_TOKEN }}"
      - name: Enable auto-merge for Dependabot PRs
        if: contains(github.event.pull_request.labels.*.name, 'automerge') && (steps.metadata.outputs.update-type == 'version-update:semver-patch' || steps.metadata.outputs.update-type == 'version-update:semver-minor')
        run: gh pr merge --squash --auto "$PR_URL"
        env:
          PR_URL: ${{github.event.pull_request.html_url}}
          GH_TOKEN: ${{secrets.GITHUB_TOKEN}}
```

#### .github/workflows/super_linter.yaml

> Hint: Make sure to adjust your super-linter checkout to match this snippet:

```yaml
---
name: Super Linter

...

jobs:
  build:
    name: Lint Code Base
    runs-on: [ubuntu-latest]
    steps:
      - uses: actions/checkout@vX
        with:
          fetch-depth: 0
...
```


## Maintain a Stable Fork

Maintain a Stable Fork with the following references:
* fork (branch): Sync with upstream repo. This branch represents the fork.
* stable (tag): Reference for the stable forked repo.
* PR/* (branch): PR branches, each based from 'stable'.
* main (branch): The operational branch, based on 'stable', and including
  commits from the PR branches.

### Update the stable reference

```bash
$ git fetch --tags upstream
$ git checkout v3.43.3
$ git tag -f stable
```

### Rebase PR Branches

```bash
$ git branch --list PR/*
$ git checkout PR/requires_check_nil
$ git rebase stable
# Repeat for each PR ... correct errors as they occur.
```

### Rebuild the "main" branch

```bash
$ git checkout stable
$ git switch -C main
$ git pull . PR/requires_check_nil
# Repeat for each PR ...

# Update module URLs.
find . -name '*.go' | xargs sed -i 's,github.com/go-task/task,github.com/boschglobal/task,g'
find . -name 'go.mod' | xargs sed -i 's,github.com/go-task/task,github.com/boschglobal/task,g'
$ git add .
$ git commit -s -m"Update module URLs."

# Reset main.
$ git switch -C main
$ git tag -f latest
```

### Update the Fork/Origin

```bash
$ git push origin stable -f
$ git push origin main -f
$ git push origin latest -f

# Push updated PRs ...
$ git push origin PR/requires_check_nil -f
```


## Contributions

### Checklist

> Hint: This list covers all items that are checked for every contribution. If you write code
properly, use the automation, and document appropriately ... then the list will already be satisfied.

> Notice: Items in bold text are mandatory in all circumstances.

- Individual/Small Contributions:
  - [ ] Sign-off from approved contributor. Sign-off implies the following:
    - [ ] __No copy/pasted code__!
    - [ ] Approval to contribute (from your organisation/company if required).
    - [ ] Acceptance of existing licence(s) used by the repository.
    - [ ] No unauthorised IP included with the commit.
    - [ ] No personal data or images.
    - [ ] No confidential information.
  - [ ] Concise commits (squash smaller changes).
  - [ ] __Linear history (no merges)__.
  - [ ] Copyright statement included with new files.
  - [ ] All items from Code Quality Checks.
- New Dependencies included in Contributions:
  - [ ] Suitability of dependency license for integration with _this_ repository.
  - [ ] Documentation of dependency license and copyright notices.
  - [ ] Evaluation of dependency risk (especially to project stability).
  - [ ] Archiving of dependency source code (mitigate continuity risk).
  - [ ] OSS Scan completed.
  - [ ] All items from Code Quality Checks.
- New Repositories:
  - [ ] Review of architectural concept (as it relates to other repositories).
  - [ ] License selection for code and documentation.
  - [ ] Approval to create new repository (from your organisation/company if required).
  - [ ] Copyright statement included with new files.
  - [ ] Documentation (including README file).
    - [ ] User documentation.
    - [ ] Architectural documentation.
    - [ ] Developer documentation.
    - [ ] Integration with documentation publishing systems (where appropriate).
  - [ ] __Clean history!__ Prior to publishing the repository, the history should be cleaned (i.e. squash to single commit) to remove any intermediate data from the history.
  - [ ] OSS Scan completed.
  - [ ] Contributor agreement (permission to publish).
  - [ ] All items from Code Quality Checks.
- Code Quality Checks:
  - [ ] No compile warnings or errors.
  - [ ] No linter warnings or errors (no circumvention).
  - [ ] __Memory leak checks__ (often enabled in unit test frameworks).
  - [ ] All tests pass.
  - [ ] Consistent variable naming (with existing code).
  - [ ] Consistent commenting style (with existing code).
  - [ ] Consistent code formatting (with existing code).
  - [ ] Documentation updated, especially generated documentation.
  - [ ] Software Design:
    - [ ] Compatibility - if possible changes _should be_ backwards compatible.
    - [ ] Platform - avoid precompile flags (i.e. `#ifdef/#endif`) for managing platform specific code.
    - [ ] Logic - avoid excessive nesting (i.e. 3+ levels) and complex conditional statements.
    - [ ] Macros - avoid use.
  - [ ] Operational Reviews:
    - [ ] Logging:
      - [ ] Normal operation produces a minimum of logging, sufficient to describe the configuration of the system in a support ticket, and especially sufficient to identify configuration errors.
      - [ ] Information logging allows normal operation of system (i.e. no debug logging).
      - [ ] __No logging output during unit tests__.
    - [ ] Fault conditions are reported, and where possible, resolution hints are given.
    - [ ] FIXME items are either resolved or opened as Issues. __No FIXME items should remain in the code__.
    - [ ] TODO items are opened as Issues if appropriate.



### Supporting Automation

#### Build System

Most code quality checks are supported via automation in repositories. These
same checks will be used by CI systems. Use the automation before submitting PRs.

```bash
# Clean the build environment.
$ make clean
# If you added new dependencies, an additional clean is required.
$ make cleanall

# Build the software.
$ make
# Resolve any compiler warnings.

# Test the software.
$ make test
# Resolve any compiler warnings from the test code.

# Run static checks.
$ make super-linter

# Generate additional content, especially if documentation or examples were
# modified as they may/will be used as supplemental documentation content.
$ make generate
# Remember to commit changes for generated content.
```


#### Code Formatting

Code (C language) can be formatted to the prevailing style with the support of automation.

```bash
# Setup the code formatting command (you can put this in your .bash_aliases file).
$ export DSE_CLANG_FORMAT_IMAGE=ghcr.io/boschglobal/dse-clang-format:main
$ alias dse-clang-format='docker run --rm -it --volume $(pwd):/tmp/code ${DSE_CLANG_FORMAT_IMAGE}'

# Format a source file.
$ dse-clang-format dse/modelc/gateway.h
```


### Platform Specific Code

As much as possible please avoid platform specific code, and in particular the use of `#ifdef/#endif` as
these make software particularly difficult to maintain and understand.

When using platform specific code becomes necessary, please use of of the following design patterns:


#### Use the existing <platform.h>

The DSE C Library includes a header file for managing smaller platform specific issues.
https://github.com/boschglobal/dse.clib/blob/main/dse/platform.h


#### Use the Linker

Use the linker to switch between code files containing the platform specific implementations, and if
necessary describe the interface of those implementations with a header file. Not only does this
technique remove the need for _any_ `#ifdef/#endif` in your code, it will also likely make
the underlying design more robust, easy to test ... and faster to compile.

The DSE Model C Library uses this technique in its Message Queue implementation:

```cmake
add_library(adapter OBJECT
...
    $<$<BOOL:${UNIX}>:transport/mq_posix.c>
...
```

and associated source code:
https://github.com/boschglobal/dse.modelc/blob/main/dse/modelc/adapter/transport/mq_posix.c



### OSS Compliance Items


#### Establishing a Mirror

A mirroring strategy for dependences can be largely automated using GitHub actions. The following
represents such an implementation:


<details>
<summary>mirror.yaml (workflow)</summary>
name: OSS Mirrors
on:
  workflow_dispatch:
  schedule:
    - cron: 0 10 * * 0
jobs:
  build:
    runs-on: [fsil-bpc]
    strategy:
      matrix:
        # List each pair of source/target repos to be mirrored.
        include:
          - source: "https://github.com/redis/hiredis"
            target: "fsil-oss-mirrors/hiredis.git"
            branch: "master"

          - source: "https://github.com/libevent/libevent"
            target: "fsil-oss-mirrors/libevent.git"
            branch: "master"

          - source: "https://github.com/yaml/libyaml"
            target: "fsil-oss-mirrors/libyaml.git"
            branch: "master"

          - source: "https://github.com/msgpack/msgpack-c"
            target: "fsil-oss-mirrors/msgpack-c.git"
            branch: "c_master"

          - source: "https://gitlab.com/cmocka/cmocka"
            target: "fsil-oss-mirrors/cmocka.git"
            branch: "master"

          - source: "https://github.com/dlfcn-win32/dlfcn-win32"
            target: "fsil-oss-mirrors/dlfcn-win32.git"
            branch: "master"

    steps:
      - name: git-sync branch
        uses: wei/git-sync@v3
        with:
          source_repo: ${{ matrix.source }}
          source_branch: ${{ matrix.branch }}
          destination_branch: ${{ matrix.branch }}
          destination_repo:  "https://${{ secrets.GHE_USER }}:${{ secrets.GHE_TOKEN }}@github.boschdevcloud.com/${{ matrix.target }}"
      - name: git-sync tags
        uses: wei/git-sync@v3
        with:
          source_repo:  ${{ matrix.source }}
          source_branch: "refs/tags/*"
          destination_branch: "refs/tags/*"
          destination_repo: "https://${{ secrets.GHE_USER }}:${{ secrets.GHE_TOKEN }}@github.boschdevcloud.com/${{ matrix.target }}"

</details>



#### Configuration for SCANS

Projects will include an automation to assist in the generation of OSS Scan Packages. These OSS Scan Packages
include the source from all dependences and may be used to run OSS Compliance Scans.

> Important: If you are adding a new dependency to a repository, please ensure that it is also added to the
automation which generates OSS Scan Packages.

```bash
# Generate the OSS Scan Package.
$ make oss

# Evaluate the OSS Scan Package content (e.g check for your new dependency).
$ ls dse/__oss__/
dlfcnwin32/  dse.clib/  dse_ncodec/  event/  hiredis/  msgpackc/  yaml/
```


#### Documentation of License (and notices)

Dependency's of a repository need to have their license and associated notices
documented. The easiest way to achieve that is to simply copy each artifact into
a `licenses` directory.

> Important: All transitive dependencies need to be documented.

```
licenses/
└── barrust/
    └── LICENSE                 <-- Code directly included in repo.
└── dlfcnwin32/
    └── COPYING
└── dse.clib/
    └── NOTICE
    └── LICENSE
└── dse_ncodec/
    └── LICENSE
    └── NOTICE
└── event/
    └── LICENSE
└── flatbuffers/                <-- Transitive dependency.
    └── LICENSE.txt
└── flatcc/                     <-- Transitive dependency.
    └── LICENSE
    └── NOTICE
└── hiredis/
    └── COPYING
└── msgpackc/
    └── LICENSE_1_0.txt
    └── NOTICE
└── yaml/
    └── License
```
