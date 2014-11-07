package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/go-giter8/st"
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
	err := ExportRepo(opts.Git, opts.Repo)
	check(err)

	// prompt the user to override the default properties
	fields, err := ReadFields(opts.Repo)
	check(err)

	// render the contents
	err = createProject(opts.Repo, fields)
	check(err)
}

func createProject(repo string, fields map[string]string) error {
	target := st.Normalize(fields["name"])
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
		destBytes, err := st.Render([]byte(target+relative), fields)
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

		output, err := st.Render(data, fields)
		if err != nil {
			return err
		}

		fmt.Printf("writing %s\n", dest)
		return ioutil.WriteFile(dest, output, f.Mode().Perm())
	})
}
