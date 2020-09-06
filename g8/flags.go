// The MIT License (MIT)
//
// Copyright (c) 2014 Matt Ho
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"github.com/urfave/cli"
)

const (
	fieldGit      = "git"
	fieldNoInputs = "no-inputs"
	fieldVerbose  = "verbose"
)

var (
	flagGit     = cli.StringFlag{Name: fieldGit, Value: "/usr/bin/git", Usage: "relativePathToTemp to the git binary", EnvVar: "GIT"}
	flagQuiet   = cli.BoolFlag{Name: fieldNoInputs, Usage: "accept all default values, do not ask for input"}
	flagVerbose = cli.BoolFlag{Name: fieldVerbose, Usage: "additional debugging"}
)

// var Verbose bool

type Options struct {
	NoInputs     bool
	Verbose      bool
	Git          string
	Repo         string
	ScaffoldName string
}

func Opts(c *cli.Context) *Options {
	// Verbose = c.Bool(fieldVerbose)
	return &Options{
		NoInputs:     c.Bool(fieldNoInputs),
		Verbose:      c.Bool(fieldVerbose),
		Git:          c.String(fieldGit),
		Repo:         c.Args().First(),
		ScaffoldName: c.Args().First(),
	}
}
