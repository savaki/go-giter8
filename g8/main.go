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
	"bufio"
	"errors"
	"fmt"
	"github.com/btnguyen2k/go-giter8/git"
	"github.com/savaki/properties"
	"github.com/urfave/cli"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	// Version of go-giter8
	Version = "0.3.2"
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

// helper to determine if relativePathToTemp exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// export a git repo to local directory
// - gitBinary: relativePathToTemp to git's executable binary file
// - repo: format [<https://host:port/>]<username>/repo-name-ends-with.g8>, example https://github.com/btnguyen2k/go_echo-microservices-seed.g8
//     if <https://host:port/> is not provided, https://github.com/ is assumed
func exportRepo(gitBinary string, repo *url.URL) error {
	if repo.Scheme == "file" {
		return nil
	}

	userAndRepoNames := userAndRepoNames(repo)
	dirOut := relativePathToTemp(userAndRepoNames)
	err := cleanDir(dirOut)
	if err != nil {
		return err
	}

	tokens := strings.Split(userAndRepoNames, "/")
	if len(tokens) != 2 || !strings.HasSuffix(tokens[1], ".g8") {
		return errors.New("repository must be in format [<https://host:port/>]<username>/repo-name-ends-with.g8>")
	}

	user := tokens[0]
	client := git.New(gitBinary, relativePathToTemp(user))
	client.Verbose = Verbose
	return client.Export(repo)
}

// cleanDir removes directory and all its content
func cleanDir(tempDir string) error {
	if exists(tempDir) {
		return os.RemoveAll(tempDir)
	}
	return nil
}

// userAndRepoNames extract the <username/repo-name> part from repository url
func userAndRepoNames(url *url.URL) string {
	userAndRepoNames := url.Path
	if strings.HasPrefix(userAndRepoNames, "/") {
		userAndRepoNames = userAndRepoNames[1:]
	}
	return userAndRepoNames
}

// relativePathToTemp create relative path to our temporary storage location
func relativePathToTemp(dirs ...string) string {
	subdir := strings.Join(dirs, "/")
	return fmt.Sprintf("%s/.go-giter8/%s", os.Getenv("HOME"), subdir)
}

// read fields and values from "default.properties"
func readFields(repo *url.URL) (map[string]string, error) {
	var path string
	// assume giter8 format
	if repo.Scheme == "file" {
		path = repo.Path + "/src/main/g8/default.properties"
	} else {
		path = relativePathToTemp(userAndRepoNames(repo), "src/main/g8/default.properties")
	}
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
	scanner := bufio.NewScanner(os.Stdin)
	for _, key := range p.Keys() {
		if key == "" {
			continue
		}
		defaultValue := p.GetString(key, "")
		var value string
		if key != "verbatim" && key != "description" {
			// do not ask for input for 'system' fields
			fmt.Printf("%s [%s]: ", key, defaultValue)
			if scanner.Scan() {
				value = scanner.Text()
			}
		}
		if strings.TrimSpace(value) == "" {
			fields[key] = defaultValue
		} else {
			fields[key] = value
		}
	}

	return fields, nil
}
