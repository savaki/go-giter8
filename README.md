# go-giter8

`btnguyen2k/go-giter8` is a fork of [`savaki/go-giter8`](https://github.com/savaki/go-giter8),
which is a command line tool to generate files and directories from templates published on git repository.
It's implemented in Go and can produce output for any purpose.

Features & TODO:
- [x] Generate files and directories from templates published on [GitHub](https://github.com) repository.
- [ ] Generate template output from any git repository.
- [ ] Generate template output from local directory (protocol `file://`).
- [ ] Support scaffolding.

Latest version: [v0.2.0](RELEASE-NOTES.md).

## Installation

You can install `go-giter8` using the standard go deployment tool. This will install g8 as ```$GOPATH/bin/g8```.

```
go get github.com/btnguyen2k/go-giter8/g8
```

## Upgrading 

At any point you can upgrade `g8` using the following:

```
go get -u github.com/btnguyen2k/go-giter8/g8
```

## Usage

Template repositories must reside in [GitHub](https://github.com) and be named with the suffix ```.g8```.  The syntax of `go-giter8` is slightly different from the original [giter8](https://github.com/n8han/giter8).

### New Project

To create a new project from template, for example, [btnguyen2k/microservices-undertow-seed.g8](https://github.com/btnguyen2k/microservices-undertow-seed.g8):

```
$ g8 new btnguyen2k/microservices-undertow-seed.g8
```

`g8` uses your git binary underneath the hood so any settings you've applied to git will also be picked up by `g8`.

# Formatting Template Fields

`go-giter8` has built-in support for formatting template fields. Formatting options can be added when referencing fields. For example, the name field can be formatted in upper camel case with:

```
$name;format="Camel"$
```

If `name` field has value `myName`, the above formatting will transform to `MyName`.

The formatting options are:

    FuncName | Alternative Name
    ---------|--------------------------------------------------------------------------------
    upper    | uppercase       : all uppercase letters
    lower    | lowercase       : all lowercase letters
    cap      | capitalize      : uppercase first letter
    decap    | decapitalize    : lowercase first letter
    start    |                 : uppercase the first letter of each word
    word     |                 : remove all non-word letters (only a-zA-Z0-9_)
    Camel    |                 : upper camel case (start-case, word-only)
    camel    |                 : lower camel case (start-case, word-only, decapitalize)
    hyphen   | hyphenate       : replace spaces with hyphens
    norm     | normalize       : all lowercase with hyphens (lowercase, hyphenate)
    snake    |                 : replace spaces and dots with underscores
    packaged |                 : replace dots with slashes (net.databinder -> net/databinder)
    random   |                 : appends random characters to the given string

Fields are defined in `src/main/g8/default.properties` file:

```
description = This template generates a microservices project using Undertow framework.
verbatim    = .DS_Store *.java .gitlab-ci.yml api.conf api_samples.conf *.xml release.sh

name               = microservices-undertow-seed
shortname          = mus
desc               = Microservices project template using Undertow
organization       = com.github.btnguyen2k
app_author         = Thanh Nguyen <btnguyen2k@gmail.com>
app_version        = 0.1.0
scala_version      = 2.13.0
timezone           = Asia/Ho_Chi_Minh
```

Special fields:
- `description`: description of the template. It will be excluded from substitution list.
- `verbatim`:  list of file patterns (separated by space, comma, semi-colon or colon) such as `*.gif,*.png *.ico`. Files matching verbatim pattern are excluded from string template processing. `verbatim` field will be excluded from substitution list.
- `name`: it is used as the name of a project being created. `go-giter8` creates a project directory based off that name (normalized) that will contain the template output.

## Giter8 template

For information on giter8 templates, please see [http://www.foundweekends.org/giter8/](http://www.foundweekends.org/giter8/).
