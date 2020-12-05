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
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"

	"github.com/btnguyen2k/go-giter8/template"
)

var commandNew = cli.Command{
	Name:        "new",
	ShortName:   "n",
	Usage:       "Create a new project",
	Description: "Create a new project from giter8 template located on GitHub, repository must be in format <username>/<repo-name-ends-with.g8>",
	Flags: []cli.Flag{
		flagGit,
		flagNoInputs,
		flagVerbose,
	},
	Action: newAction,
}

// handle command "new"
func newAction(c *cli.Context) {
	opts := Opts(c)
	if strings.TrimSpace(opts.Repo) == "" {
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
	repo.Path = strings.TrimSuffix(repo.Path, "/")
	exitIfError(exportGitRepo(opts, repo))

	// load parameters
	fields, err := readFieldsFromG8Template(opts, repo)
	exitIfError(err)

	// render the contents
	exitIfError(newProject(opts, repo, fields))

	if repo.Scheme != "file" {
		exitIfError(cleanDir(relativePathToTemp(userAndRepoNames(repo))))
	}
}

// create new project from template
func newProject(opts *Options, repo *url.URL, fields map[string]string) error {
	delete(fields, "description") // remove system field "description"

	target := template.Normalize(fields["name"])
	if target == "" {
		return errors.New("no [name] parameter defined")
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
		fmt.Printf("Copying scaffolds...")
		if opts.Verbose {
			fmt.Println()
		}
		if err := copyDir(opts, scaffoldsBase, target+"/.g8"); err != nil {
			return err
		}
		if !opts.Verbose {
			fmt.Printf("done.\n")
		}
	}

	// create project structure from template
	fmt.Printf("Generating project...")
	if opts.Verbose {
		fmt.Println()
	}
	prefixLen := len(codeBase)
	err := filepath.Walk(codeBase, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || f.Name() == "default.properties" {
			return nil
		}

		relativePath := path[prefixLen:] // relativePathToTemp is absolute; let's strip off the prefix
		var destFileName = target + relativePath
		// transform filename if not in "verbatim" list
		if !isFileMatched(relativePath, f, verbatim) {
			destFileName, err = transformFilename(destFileName, fields)
			if err != nil {
				return err
			}
		}

		// ensure the directory exists
		if err = mkdir(filepath.Dir(destFileName), 0755); err != nil {
			return err
		}

		if opts.Verbose {
			fmt.Printf("\tgenerating %s...", destFileName)
		}
		// load file content
		inContent, err := ioutil.ReadFile(path)
		if err != nil {
			if opts.Verbose {
				fmt.Printf("error.\n")
			}
			return err
		}
		outContent := inContent

		// transform content if not in "verbatim" list
		if !isFileMatched(relativePath, f, verbatim) {
			outContent, err = template.Render(inContent, fields)
			if err != nil {
				if opts.Verbose {
					fmt.Printf("error.\n")
				}
				return err
			}
		}
		err = ioutil.WriteFile(destFileName, outContent, f.Mode())
		if opts.Verbose {
			if err != nil {
				fmt.Printf("error.\n")
			} else {
				fmt.Printf("done.\n")
			}
		}
		return err
	})
	if !opts.Verbose {
		if err != nil {
			fmt.Printf("error.\n")
		} else {
			fmt.Printf("done.\n")
		}
	}
	return err
}
