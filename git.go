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

func newGit(gopath string, packageURL string) *git {
	packageRoot := packageGitDir(packageURL)
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

func packageGitDir(name string) string {
	tokens := strings.Split(name, "/")
	if len(tokens) > 3 {
		return strings.Join(tokens[:3], "/")
	}
	return name
}
