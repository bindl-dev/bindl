# Contributing

Thank you for your interest in Bindl. We welcome and encourage public contributions.

Prior to submitting code changes, we ask that you communicate the proposal and reach an agreement (typically in GitHub Issues) to help reviewers and increase the likelihood of accepting the changeset.

## Development environment

Ironically, as a project for dependency management, we have some dependencies which are not contained by `bindl.yaml`. While they are not _strictly necessary_, they are a convenient compromise. We are certainly open to alternatives. Until then, the project makes use of:

- **GNU Make** version 4.3
- [**direnv**](https://direnv.net)

The good news is that other dependencies are listed in this project's own `bindl.yaml`. To set them up, start with bootstrapping Bindl.

### Bootstrapping Bindl

If you haven't already have Bindl in your `$PATH`, use one of the following methods:

- `make bin/bindl` to use latest release
- `make bin/bindl-min` to compile from currently checked-out state

## Ensuring continuous integration

The project uses GitHub Actions to run continuous integration tasks. The CI pipeline invokes various Makefile rules to reduce development environment vs CI variability.

```bash
# To run tests:
make test/all
# or look inside Makefile on which specific tests you'd like to run.

# To run linters:
make lint
# and to run auto-fixer:
make lint/fix

# To make development builds:
make bin/bindl-dev
# This specifically invokes goreleaser to emulate builds on release
```
