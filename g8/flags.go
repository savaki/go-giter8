package main

import (
	"github.com/codegangsta/cli"
)

const (
	fieldGit     = "git"
	fieldVerbose = "verbose"
)

var (
	flagGit     = cli.StringFlag{fieldGit, "/usr/bin/git", "path to the git binary", "GIT"}
	flagVerbose = cli.BoolFlag{fieldVerbose, "additional debugging", "VERBOSE"}
)

var Verbose bool

type Options struct {
	Verbose bool
	Git     string
	Repo    string
}

func Opts(c *cli.Context) Options {
	Verbose = c.Bool(fieldVerbose)

	return Options{
		Verbose: Verbose,
		Git:     c.String(fieldGit),
		Repo:    c.Args().First(),
	}
}
