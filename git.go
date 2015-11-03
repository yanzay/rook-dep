package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type git struct {
	workTree string
	gitDir   string
}

func newGit(gopath string, packageRoot string) *git {
	return &git{
		workTree: fmt.Sprintf("--work-tree=%s/src/%s", gopath, packageRoot),
		gitDir:   fmt.Sprintf("--git-dir=%s/src/%s/.git", gopath, packageRoot),
	}
}

func (g *git) checkout(target string) error {
	cmd := exec.Command("git", g.workTree, g.gitDir, "checkout", target)
	return cmd.Run()
}

func (g *git) currentCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	hash, err := cmd.Output()
	return strings.TrimSpace(string(hash)), err
}

func (g *git) params() string {
	return fmt.Sprintf("%s %s", g.workTree, g.gitDir)
}
