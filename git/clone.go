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
