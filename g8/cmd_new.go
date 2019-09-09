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
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/btnguyen2k/go-giter8/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

var commandNew = cli.Command{
	Name:  "new",
	Usage: "create a new project",
	Flags: []cli.Flag{
		flagGit,
		flagVerbose,
	},
	Action: newAction,
}

func newAction(c *cli.Context) {
	opts := Opts(c)

	if opts.Repo == "" {
		fmt.Println("ERROR - no template repo specified")
	}

	// extract the repo
	err := exportRepo(opts.Git, opts.Repo)
	check(err)

	// prompt the user to override the default properties
	fields, err := readFields(opts.Repo)
	check(err)

	// render the contents
	err = newProject(opts.Repo, fields)
	check(err)
}

func newProject(repo string, fields map[string]string) error {
	target := template.Normalize(fields["name"])
	if target == "" {
		check(errors.New("no name parameter defined"))
	}

	codebase := Path(repo, "src/main/g8")
	prefix := len(codebase)
	return filepath.Walk(codebase, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		_, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		relative := path[prefix:] // path is absolute; let's strip off the prefix
		destBytes, err := template.Render([]byte(target+relative), fields)
		if err != nil {
			return err
		}
		dest := string(destBytes)

		// ensure the directory exists
		dirname := filepath.Dir(dest)
		if !exists(dirname) {
			fmt.Printf("creating directory, %s\n", dirname)
			os.MkdirAll(dirname, 0755)
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		output, err := template.Render(data, fields)
		if err != nil {
			return err
		}

		fmt.Printf("writing %s\n", dest)
		return ioutil.WriteFile(dest, output, f.Mode().Perm())
	})
}
