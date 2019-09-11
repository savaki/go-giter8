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
	"github.com/btnguyen2k/go-giter8/git"
	"github.com/codegangsta/cli"
	"github.com/savaki/properties"
	"log"
	"os"
	"strings"
)

const (
	// Version of giter8
	Version = "0.2.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "giter8"
	app.Usage = "Generate templates from GitHub"
	app.Version = Version
	app.Commands = []cli.Command{
		commandNew,
	}
	app.Run(os.Args)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// helper to determine if path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// export a git repo to local directory
// - gitpath
// - repo: format <username>/repo-name-ends-with.g8>, example btnguyen2k/microservices-undertow-seed.g8
func exportRepo(gitpath, repo string) error {
	dirOut := Path(repo)
	if exists(Path(repo)) {
		// output directory exist, clean up for a fresh clone
		err := os.RemoveAll(dirOut)
		if err != nil {
			return err
		}
	}
	tokens := strings.Split(repo, "/")
	if len(tokens) != 2 || !strings.HasSuffix(tokens[1], ".g8") {
		return errors.New("repository must be in format <username>/repo-name-ends-with.g8>")
	}

	user := tokens[0]
	client := git.New(gitpath, Path(user))
	client.Verbose = Verbose
	return client.Export(repo)
}

// path relative to our temporary storage location
func Path(dirs ...string) string {
	subdir := strings.Join(dirs, "/")
	return fmt.Sprintf("%s/.go-giter8/%s", os.Getenv("HOME"), subdir)
}

// read fields and values from "default.properties"
func readFields(repo string) (map[string]string, error) {
	// assume giter8 format
	path := Path(repo, "src/main/g8/default.properties")
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return map[string]string{}, nil
	}

	// print out template description if exists
	desc := p.GetString("description", "")
	if desc != "" {
		fmt.Println(desc)
	}

	// ask the user for input on each of the fields
	fields := map[string]string{}
	for _, key := range p.Keys() {
		if key == "" {
			continue
		}
		defaultValue := p.GetString(key, "")
		var value string
		if key != "verbatim" && key != "description" {
			// do not ask for input for "system" fields
			fmt.Printf("%s [%s]: ", key, defaultValue)
			fmt.Scanln(&value)
		}
		if strings.TrimSpace(value) == "" {
			fields[key] = defaultValue
		} else {
			fields[key] = value
		}
	}

	return fields, nil
}
