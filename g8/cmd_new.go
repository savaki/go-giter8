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
	"github.com/btnguyen2k/go-giter8/template"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

var commandNew = cli.Command{
	Name:        "new",
	ShortName:   "n",
	Usage:       "Create a new project",
	Description: "Create a new project from giter8 template located on GitHub, repository must be in format <username>/<repo-name-ends-with.g8>",
	Flags: []cli.Flag{
		flagGit,
		flagVerbose,
	},
	Action: newAction,
}

// handle command "new"
func newAction(c *cli.Context) {
	opts := Opts(c)

	if opts.Repo == "" {
		exitIfError(errors.New("ERROR - no template repo specified"))
	}

	// extract the repo
	repo, err := url.Parse(opts.Repo)
	exitIfError(err)
	if repo.Scheme == "" {
		repo.Scheme = "https"
	}
	if repo.Host == "" && repo.Scheme != "file" {
		repo.Host = "github.com" // template are fetched from github by default
	}
	err = exportGitRepo(opts.Git, repo)
	exitIfError(err)

	// prompt the user to override the default properties
	fields, err := readFieldsFromGitRepo(repo)
	exitIfError(err)

	// render the contents
	err = newProject(repo, fields)
	exitIfError(err)

	if repo.Scheme != "file" {
		err = cleanDir(relativePathToTemp(userAndRepoNames(repo)))
		exitIfError(err)
	}
}

// create new project from template
func newProject(repo *url.URL, fields map[string]string) error {
	delete(fields, "description") // remove system field "description"

	target := template.Normalize(fields["name"])
	if target == "" {
		exitIfError(errors.New("no [name] parameter defined"))
	}

	var verbatim []string
	if val, ok := fields["verbatim"]; ok && val != "" {
		verbatim = regexp.MustCompile("[,;:\\s]+").Split(fields["verbatim"], -1)
	}
	delete(fields, "verbatim") // remove system field "verbatim"

	var codeBase, scaffoldsBase string
	if repo.Scheme == "file" {
		codeBase = repo.Path + "/src/main/g8"
		scaffoldsBase = repo.Path + "/src/main/scaffolds"
	} else {
		codeBase = relativePathToTemp(userAndRepoNames(repo), "src/main/g8")
		scaffoldsBase = relativePathToTemp(userAndRepoNames(repo), "src/main/scaffolds")
	}

	// copy scaffolds
	if pathExists(scaffoldsBase) {
		fmt.Println("Copying scaffolds...")
		if err := copyDir(scaffoldsBase, target+"/.g8"); err != nil {
			return err
		}
	}

	// create project structure from template
	fmt.Println("Generating project...")
	prefixLen := len(codeBase)
	return filepath.Walk(codeBase, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || f.Name() == "default.properties" {
			return nil
		}

		relativePath := path[prefixLen:] // relativePathToTemp is absolute; let's strip off the prefix
		// transform filename
		destFileName, err := transformFilename(target+relativePath, fields)
		if err != nil {
			return err
		}
		// ensure the directory exists
		if err = mkdir(filepath.Dir(destFileName), 0755); err != nil {
			return err
		}

		fmt.Println("\tgenerating", destFileName)

		// load file content
		inContent, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		outContent := inContent

		// transform content if not in "verbatim" list
		if !isFileMatched(relativePath, f, verbatim) {
			outContent, err = template.Render(inContent, fields)
			if err != nil {
				return err
			}
		}
		return ioutil.WriteFile(destFileName, outContent, f.Mode())
	})
}
