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
	"fmt"
	"log"
	"os"
	"os/exec"
)

func (g *Git) Clone(url string) error {
	if g.Verbose {
		log.Printf("git clone %s\n", url)
	}
	if _, err := os.Stat(g.Target); os.IsNotExist(err) {
		os.MkdirAll(g.Target, 0755)
	}

	cmd := exec.Command(g.Git, "clone", url)
	cmd.Dir = g.Target
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (g *Git) Export(repo string) error {
	err := g.Clone(https(repo))
	if err != nil {
		log.Printf("repo not available over https; attempting to clone repo via ssh")
		err = g.Clone(ssh(repo))
	}

	return err
}

func https(repo string) string {
	return fmt.Sprintf("https://github.com/%s.git", repo)
}

func ssh(repo string) string {
	return fmt.Sprintf("git@github.com:%s.git", repo)
}
