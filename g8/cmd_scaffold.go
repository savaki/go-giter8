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
	Flags: []cli.Flag{
		flagNoInputs,
		flagVerbose,
	},
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

	// load parameters
	fields, err := readFieldsFromFile(opts, srcDir+"/default.properties")
	exitIfError(err)

	delete(fields, "description") // remove system field "description"
	var verbatim []string
	if val, ok := fields["verbatim"]; ok && val != "" {
		verbatim = regexp.MustCompile("[,;:\\s]+").Split(fields["verbatim"], -1)
	}
	delete(fields, "verbatim") // remove system field "verbatim"

	// generate scaffold
	fmt.Printf("Generating scaffold %s...", scaffoldName)
	if opts.Verbose {
		fmt.Println()
	}
	prefixLen := len(srcDir)
	err = filepath.Walk(srcDir, func(path string, f os.FileInfo, err error) error {
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
