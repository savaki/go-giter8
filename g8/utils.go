package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobwas/glob"
	"github.com/savaki/properties"

	"github.com/btnguyen2k/go-giter8/git"
	"github.com/btnguyen2k/go-giter8/template"
)

// exitIfError terminates application in case of error
func exitIfError(err error) {
	if err != nil {
		panic(err)
		// log.Fatalln(err)
	}
}

// pathExists checks if a file/directory pathExists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	switch {
	case err != nil:
		return false
	default:
		return fi.IsDir()
	}
}

// cleanDir removes directory and all its content
func cleanDir(tempDir string) error {
	if pathExists(tempDir) {
		return os.RemoveAll(tempDir)
	}
	return nil
}

// exportGitRepo exports a git repo to local directory
// - gitBinary: git's executable binary file
// - repo: format [<https://host:port/>]<username>/repo-name-ends-with.g8>[@branchOrTagName],
//     example https://github.com/btnguyen2k/go_echo-microservices-seed.g8@template-v0.4.r1
//     if <https://host:port/> is not provided, https://github.com/ is assumed
func exportGitRepo(opts *Options, repo *url.URL) error {
	if repo.Scheme == "file" {
		// @branch or @tag is not supported with local files
		return nil
	}

	gitBinary := opts.Git
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
	client.Verbose = opts.Verbose
	return client.Clone(repo.String(), tagOrBranchName(repo))
	// return client.Export(repo, tagOrBranchName(repo))
}

// userAndRepoNames extracts the <username/repo-name> part from repository url
func userAndRepoNames(url *url.URL) string {
	userAndRepoNames := url.Path
	if i := strings.LastIndex(userAndRepoNames, "@"); i >= 0 {
		userAndRepoNames = userAndRepoNames[0:i]
	}
	if strings.HasPrefix(userAndRepoNames, "/") {
		userAndRepoNames = userAndRepoNames[1:]
	}
	return userAndRepoNames
}

// tagOrBranchName extract the tag/branch name part from repository url.
// tagOrBranchName is suffixed to the repository url after @ character.
func tagOrBranchName(url *url.URL) string {
	tagOrBranchName := ""
	if i := strings.LastIndex(url.Path, "@"); i >= 0 {
		tagOrBranchName = url.Path[i+1:]
	}
	return tagOrBranchName
}

// relativePathToTemp creates relative path to our temporary storage location
func relativePathToTemp(dirs ...string) string {
	subdir := strings.Join(dirs, "/")
	return fmt.Sprintf("%s/.go-giter8/%s", os.Getenv("HOME"), subdir)
}

// readFieldsFromFile reads fields/values from file (".properties" format)
func readFieldsFromFile(opts *Options, path string) (map[string]string, error) {
	if opts.Verbose {
		fmt.Printf("Loading parameters from file %s...", path)
	}
	props, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		if opts.Verbose {
			fmt.Printf("error.\n")
		}
		return map[string]string{}, nil
	}
	if opts.Verbose {
		fmt.Printf("done.\n")
	}

	// print out "description" if any
	if desc := props.GetString("description", ""); desc != "" {
		fmt.Printf("\t%s\n", desc)
	}

	fmt.Printf("Customizing parameters: ")
	if opts.NoInputs {
		fmt.Printf("skipped.\n")
	} else {
		fmt.Printf("\n")
	}
	// ask the user for input on each of the fields
	fields := map[string]string{}
	scanner := bufio.NewScanner(os.Stdin)
	for _, key := range props.Keys() {
		if key == "" {
			continue
		}
		defaultValue := props.GetString(key, "")
		var value string
		if key != "verbatim" && key != "description" && !opts.NoInputs {
			// do not ask for input for 'system' fields
			fmt.Printf("\t%s [%s]: ", key, defaultValue)
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

// readFieldsFromG8Template reads fields/values from "default.properties" resided in a "giter8" template
func readFieldsFromG8Template(opts *Options, repo *url.URL) (map[string]string, error) {
	var path string
	// assume giter8 format
	if repo.Scheme == "file" {
		path = repo.Path + "/src/main/g8/default.properties"
	} else {
		path = relativePathToTemp(userAndRepoNames(repo), "src/main/g8/default.properties")
	}
	return readFieldsFromFile(opts, path)
}

func transformFilename(filename string, fields map[string]string) (string, error) {
	result, err := template.Render([]byte(filename), fields)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func mkdir(target string, mode os.FileMode) error {
	if !pathExists(target) {
		return os.MkdirAll(target, mode)
	}
	return nil
}

func copyDir(opts *Options, srcDir, destDir string) error {
	prefixLen := len(srcDir)
	return filepath.Walk(srcDir, func(path string, f os.FileInfo, err error) error {
		relativePath := path[prefixLen:]
		dest := destDir + relativePath
		if opts.Verbose {
			fmt.Printf("\tCopying %s -> %s\n", path, dest)
		}
		switch f.Mode() & os.ModeType {
		case os.ModeDir:
			return mkdir(dest, f.Mode())
		case os.ModeSymlink:
			if err := mkdir(filepath.Dir(dest), 0755); err != nil {
				return err
			}
			if link, err := os.Readlink(path); err != nil {
				return err
			} else {
				return os.Symlink(link, dest)
			}
		default:
			if err := mkdir(filepath.Dir(dest), 0755); err != nil {
				return err
			}
			if content, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				return ioutil.WriteFile(dest, content, f.Mode())
			}
		}
		return nil
	})
}

func isFileMatched(relativePath string, f os.FileInfo, matchList []string) bool {
	isWindows := strings.ToLower(runtime.GOOS) == "windows"
	if isWindows {
		relativePath = strings.TrimPrefix(relativePath, "\\")
	} else {
		relativePath = strings.TrimPrefix(relativePath, "/")
	}

	for _, pattern := range matchList {
		if isWindows {
			pattern = strings.ReplaceAll(pattern, "/", "\\\\")
		} else {
			pattern = strings.ReplaceAll(pattern, "\\", "/")
		}
		if matched, _ := filepath.Match(pattern, f.Name()); matched {
			return true
		}
		if g := glob.MustCompile(pattern); g != nil && g.Match(relativePath) {
			return true
		}
	}
	return false
}
