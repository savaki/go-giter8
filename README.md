# go-giter8

`btnguyen2k/go-giter8` is a fork of [`savaki/go-giter8`](https://github.com/savaki/go-giter8),
which is a command line tool to generate files and directories from templates published on git repository.
It's implemented in Go and can produce output for any purpose.

Features & TODO:
- [x] Generate template output from any git repository.
- [x] Generate template output from local directory (protocol `file://`).
- [ ] Support scaffolding.

Latest version: [v0.3.0](RELEASE-NOTES.md).

## Installation

You can install `go-giter8` using the standard go deployment tool. This will install g8 as ```$GOPATH/bin/g8```:

```
go get github.com/btnguyen2k/go-giter8/g8
```

or you can specified a specific version:

```
go get github.com/btnguyen2k/go-giter8/g8@v0.3.0
```

## Upgrading 

At any point you can upgrade `g8` using the following:

```
go get -u github.com/btnguyen2k/go-giter8/g8
```

or you can specified a specific version:

```
go get -u github.com/btnguyen2k/go-giter8/g8@v0.3.0
```

## Usage

Template repositories must reside in a git repository and be named with the suffix ```.g8```.  The syntax of `go-giter8` is slightly different from the original [giter8](https://github.com/n8han/giter8).

### New Project

To create a new project from template, for example, [btnguyen2k/go_echo-microservices-seed.g8](https://github.com/btnguyen2k/go_echo-microservices-seed.g8):

```
$ g8 new btnguyen2k/go_echo-microservices-seed.g8
```

By default, template is cloned from [Github](https://github.com). Use full repo url to create project from template resided in other git server:

```
$ g8 new https://gitlab.com/btnguyen2k/go_echo-microservices-seed.g8
```

`g8` uses your git binary underneath the hood so any settings you've applied to git will also be picked up by `g8`.

`go-giter8` can also generate template output from local directory (useful when testing template before publishing):

```
$ g8 new file:///home/btnguyen2k/workspace/go_echo-microservices-seed.g8
```

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
description = This template generates a microservices project in Go using Echo framework.
verbatim    = .DS_Store .gitlab-ci.yml release.sh

name         = go_echo-microservices-seed
shortname    = gems
desc         = Microservices project template for Go using Echo
organization = com.github.btnguyen2k
app_author   = Thanh Nguyen <btnguyen2k@gmail.com>
app_version  = 0.1.0
timezone     = Asia/Ho_Chi_Minh
```

Special fields:
- `description`: description of the template. It will be excluded from substitution list.
- `verbatim`:  list of file patterns (separated by space, comma, semi-colon or colon) such as `*.gif,*.png *.ico`. Files matching verbatim pattern are excluded from string template processing. `verbatim` field will be excluded from substitution list.
- `name`: it is used as the name of a project being created. `go-giter8` creates a project directory based off that name (normalized) that will contain the template output.

## Giter8 template

For information on giter8 templates, please see [http://www.foundweekends.org/giter8/](http://www.foundweekends.org/giter8/).
