package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"

	"github.com/btnguyen2k/go-giter8/template"
)

var commandScaffold = cli.Command{
	Name:        "scaffold",
	ShortName:   "sf",
	Usage:       "Generate files from a scaffolding",
	Description: "Generate files from a giter8 scaffold",
	// Flags: []cli.Flag{
	// 	flagGit,
	// 	flagVerbose,
	// },
	Action: scaffoldAction,
}

// handle command "scaffold"
func scaffoldAction(c *cli.Context) {
	exitIfError(generateScaffold(Opts(c)))
}

func generateScaffold(opts *Options) error {
	scaffoldName := strings.TrimSpace(opts.ScaffoldName)
	if scaffoldName == "" {
		return errors.New("ERROR - no scaffold name specified")
	}
	srcDir := ".g8/" + scaffoldName
	if !isDir(srcDir) {
		return errors.New("ERROR - [" + srcDir + "] not readable, or not a directory")
	}
	// must stand at project's root directory
	destDir := "."

	// prompt the user to override the default properties
	fields, err := readFieldsFromFile(opts, srcDir+"/default.properties")
	exitIfError(err)

	delete(fields, "description") // remove system field "description"
	var verbatim []string
	if val, ok := fields["verbatim"]; ok && val != "" {
		verbatim = regexp.MustCompile("[,;:\\s]+").Split(fields["verbatim"], -1)
	}
	delete(fields, "verbatim") // remove system field "verbatim"

	fmt.Println("Generating scaffold " + scaffoldName + "...")
	prefixLen := len(srcDir)
	return filepath.Walk(srcDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || f.Name() == "default.properties" {
			return nil
		}

		relativePath := path[prefixLen:]
		destFileName, err := transformFilename(destDir+relativePath, fields)
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
