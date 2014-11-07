go-giter8
=========

go-giter8 is a command line tool based on n8han/giter8 to generate files and directories from templates published on any git repository.  It's implemented in Go and can produce output for any purpose.

The motivation came after noticing the complexity of modifying the original giter8 codebase.  The additional features provided by go-giter8 are:

* go-giter8 shells out to the real git binary for git access rather than simulating it via libraries
* multiple filters are now supported for directories

# Installation

You can install go-giter8 using the standard go deployment tool.  This will install g8 as ```$GOPATH/bin/g8```.

```
go get github.com/savaki/go-giter8/g8
```

***Future*** - need to set up brew install

# Upgrading 

At any point you can upgrade g8 using the following:

```
go get -u github.com/savaki/go-giter8/g8
```

# Usage

Template repositories must reside in git and be named with the suffix ```.g8```.  The syntax of go-giter8 is slightly different from the original giter8.

## New Project

To create a new project from template, for example, [loyal3/service-template-finatra.g8](https://github.com/loyal3/service-template-finatra.g8):

```
$ g8 new loyal3/service-template-finatra
```

The ```.g8``` suffix is assumed.  You can also specify a full repository name:

```
$ g8 git@github.com:loyal3/service-template-finatra.g8.git
```

or

```
$ g8 https://github.com/loyal3/service-template-finatra.g8.git
```

g8 uses your git binary underneath the hood so any settings you've applied to git will also be picked up by g8.

# Formatting Template Fields

go-giter8 has built-in support for formatting template fields. Formatting options can be added when referencing fields. For example, the name field can be formatted in upper camel case with:

```
$name;format="Camel"$
```

The formatting options are:

    upper    | uppercase       : all uppercase letters
    lower    | lowercase       : all lowercase letters
    cap      | capitalize      : uppercase first letter
    decap    | decapitalize    : lowercase first letter
    start    | start-case      : uppercase the first letter of each word
    word     | word-only       : remove all non-word letters (only a-zA-Z0-9_)
    Camel    | upper-camel     : upper camel case (start-case, word-only)
    camel    | lower-camel     : lower camel case (start-case, word-only, decapitalize)
    hyphen   | hyphenate       : replace spaces with hyphens
    norm     | normalize       : all lowercase with hyphens (lowercase, hyphenate)
    snake    | snake-case      : replace spaces and dots with underscores
    packaged | package-dir     : replace dots with slashes (net.databinder -> net/databinder)
    random   | generate-random : appends random characters to the given string