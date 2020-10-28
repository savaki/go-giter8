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

package git

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func (g *Git) Clone(url, tagOrBranchName string) error {
	if i := strings.LastIndex(url, "@"); i >= 0 {
		url = url[0:i]
	}
	args := []string{"clone", url}
	if tagOrBranchName != "" {
		args = append(append(args, "--branch"), tagOrBranchName)
	}
	if g.Verbose {
		log.Printf("git %s\n", strings.Join(args, " "))
	}
	if _, err := os.Stat(g.TargetDir); os.IsNotExist(err) {
		os.MkdirAll(g.TargetDir, 0755)
	}

	cmd := exec.Command(g.GitBinary, args...)
	cmd.Dir = g.TargetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// func (g *Git) Export(repo *url.URL) error {
// 	return g.Clone(repo.String() + ".git")
// }
