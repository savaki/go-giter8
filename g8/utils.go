package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/btnguyen2k/go-giter8/git"
	"github.com/btnguyen2k/go-giter8/template"
	"github.com/gobwas/glob"
	"github.com/savaki/properties"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
// - repo: format [<https://host:port/>]<username>/repo-name-ends-with.g8>, example https://github.com/btnguyen2k/go_echo-microservices-seed.g8
//     if <https://host:port/> is not provided, https://github.com/ is assumed
func exportGitRepo(gitBinary string, repo *url.URL) error {
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

// userAndRepoNames extract the <username/repo-name> part from repository url
func userAndRepoNames(url *url.URL) string {
	userAndRepoNames := url.Path
	if strings.HasPrefix(userAndRepoNames, "/") {
		userAndRepoNames = userAndRepoNames[1:]
	}
	return userAndRepoNames
}

// relativePathToTemp creates relative path to our temporary storage location
func relativePathToTemp(dirs ...string) string {
	subdir := strings.Join(dirs, "/")
	return fmt.Sprintf("%s/.go-giter8/%s", os.Getenv("HOME"), subdir)
}

// readFieldsFromFile reads fields/values from file ("..properties" format)
func readFieldsFromFile(path string) (map[string]string, error) {
	props, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return map[string]string{}, nil
	}

	// print out "description" if any
	if desc := props.GetString("description", ""); desc != "" {
		fmt.Println(desc)
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

// readFieldsFromGitRepo reads fields/values from "default.properties" resided in a git repo
func readFieldsFromGitRepo(repo *url.URL) (map[string]string, error) {
	var path string
	// assume giter8 format
	if repo.Scheme == "file" {
		path = repo.Path + "/src/main/g8/default.properties"
	} else {
		path = relativePathToTemp(userAndRepoNames(repo), "src/main/g8/default.properties")
	}
	return readFieldsFromFile(path)
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

func copyDir(srcDir, destDir string) error {
	prefixLen := len(srcDir)
	return filepath.Walk(srcDir, func(path string, f os.FileInfo, err error) error {
		relativePath := path[prefixLen:]
		dest := destDir + relativePath
		fmt.Println("\tcopying", path, "->", dest)
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
