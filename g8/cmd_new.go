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
func transformFilename(filename string, fields map[string]string) (string, error) {
    result, err := template.Render([]byte(filename), fields)
    if err != nil {
        return "", err
    }
    return string(result), nil
}

// create new project from template
func newProject(repo string, fields map[string]string) error {
    delete(fields, "description") // remove system field "description"

    target := template.Normalize(fields["name"])
    if target == "" {
        check(errors.New("no [name] parameter defined"))
    }

    var verbatim []string
    if val, ok := fields["verbatim"]; ok && val != "" {
        verbatim = regexp.MustCompile("[,;:\\s]+").Split(fields["verbatim"], -1)
    }
    delete(fields, "verbatim") // remove system field "verbatim"

    codebase := Path(repo, "src/main/g8")
    prefixLen := len(codebase)
    return filepath.Walk(codebase, func(path string, f os.FileInfo, err error) error {
        if f.IsDir() {
            return nil
        }

        relative := path[prefixLen:] // path is absolute; let's strip off the prefix
        // transform filename
        destFileName, err := transformFilename(target+relative, fields)
        if err != nil {
            return err
        }

        // ensure the directory exists
        dirname := filepath.Dir(destFileName)
        if !exists(dirname) {
            fmt.Printf("creating directory, %s\n", dirname)
            os.MkdirAll(dirname, 0755)
        }

        fmt.Printf("generating %s\n", destFileName)

        // load file content
        inContent, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }
        outContent := inContent
        // transform content if not in verbatim list
        if !isVerbatim(f, verbatim) {
            outContent, err = template.Render(inContent, fields)
            if err != nil {
                return err
            }
        }
        return ioutil.WriteFile(destFileName, outContent, f.Mode().Perm())
    })
}

func isVerbatim(f os.FileInfo, verbatim []string) bool {
    for _, v := range verbatim {
        matched, _ := filepath.Match(v, f.Name())
        if matched {
            return true
        }
    }
    return false
}
