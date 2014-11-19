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
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/go-giter8/git"
	"github.com/savaki/properties"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "giter8"
	app.Usage = "generate templates using github"
	app.Version = "0.1"
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

// ExportRepo(git, loyal3/service-template-finatra.g8) => nil
func exportRepo(gitpath, repo string) error {
	if exists(Path(repo)) {
		return nil
	}

	user := strings.Split(repo, "/")[0]
	client := git.New(gitpath, Path(user))
	client.Verbose = Verbose
	return client.Export(repo)
}

// path relative to our temporary storage location
func Path(dirs ...string) string {
	subdir := strings.Join(dirs, "/")
	return fmt.Sprintf("%s/.go-giter8/%s", os.Getenv("HOME"), subdir)
}

func readFields(repo string) (map[string]string, error) {
	// assume giter8 format
	path := Path(repo, "src/main/g8/default.properties")
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return map[string]string{}, nil
	}

	// ask the user for input on each of the fields
	fields := map[string]string{}
	for _, key := range p.Keys() {
		defaultValue := p.GetString(key, "")
		fmt.Printf("%s [%s]: ", key, defaultValue)

		var value string
		fmt.Scanln(&value)

		if strings.TrimSpace(value) == "" {
			fields[key] = defaultValue
		} else {
			fields[key] = value
		}
	}

	return fields, nil
}
