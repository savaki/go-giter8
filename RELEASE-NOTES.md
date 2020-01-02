# go-giter8 release notes

## 2020-01-01 - v0.4.0

- Support scaffolding & `scaffold` command.


## 2019-12-30 - v0.3.2

- Enhance "verbatim" list: support files under specific directories.


## 2019-11-25 - vc

- Fix bug: input with spaces is not treated as a whole string.


## 2019-09-19 - v0.3.0

- Generate template output from any git repository.
- Remove temp directory after template output is generated successfully.
- Generate template output from local directory (protocol `file://`).
- Other fixes and improvements.


## 2019-09-11 - v0.2.0.1

- Migrate to use package `github.com/urfave/cli`.


## 2019-09-10 - v0.2.0

- Forked from [savaki/go-giter8](https://github.com/savaki/go-giter8).
- Fixed bug: `unrecognized import path "code.google.com/p/go-uuid/uuid"`.
- Clearly document that currently `go-giter8` supports only templates from GitHub.
- Removed non-identifier transform functions to be compatible with package `text/template`.
