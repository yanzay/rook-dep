package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	Depfile     = "Rookfile"
	DepfileLock = "Rookfile.lock"
)

type Package struct {
	name    string
	version string
	url     string
	hash    string
}

func (p *Package) String() string {
	return fmt.Sprintf("%s: %s", p.name, p.version)
}

func main() {
	var err error
	if len(os.Args) > 2 && os.Args[1] == "search" {
		err = search(os.Args[2])
	} else {
		err = install()
	}
	printResult(err)
}

func install() error {
	var err error
	fileName, wrap := restoreCheck()
	packages, err := parseFile(fileName)
	if err != nil {
		return err
	}
	for i, pack := range packages {
		var err error
		packages[i].hash, err = getPackage(pack, wrap(pack.version))
		if err != nil {
			return err
		}
	}
	saveLock(packages)
	return nil
}

func restoreCheck() (string, func(string) string) {
	if _, err := os.Stat(DepfileLock); os.IsExist(err) {
		return DepfileLock, noWrap
	}
	return Depfile, wrapTag
}

func noWrap(str string) string {
	return str
}

func wrapTag(str string) string {
	return fmt.Sprintf("tags/%s", str)
}

func saveLock(packages []Package) error {
	var content string
	for _, pack := range packages {
		content += fmt.Sprintf("%s: %s\n", pack.name, pack.hash)
	}
	return ioutil.WriteFile("Rookfile.lock", []byte(content), os.ModePerm)
}

func printResult(err error) {
	fmt.Println("---")
	if err != nil {
		fmt.Printf("Rook failed: %s\n", err)
	} else {
		fmt.Printf("Rook success!\n")
	}
}

func parseFile(name string) ([]Package, error) {
	var packages []Package
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ":")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Unable to parse: %s", line)
		}
		name := strings.TrimSpace(tokens[0])
		version := strings.TrimSpace(tokens[1])
		url := transformName(name)
		packages = append(packages, Package{name: name, version: version, url: url})
	}
	return packages, nil
}

func transformName(name string) string {
	if strings.HasPrefix(name, "golang/") {
		name = strings.TrimPrefix(name, "golang/")
		return fmt.Sprintf("golang.org/x/%s", name)
	}
	return fmt.Sprintf("github.com/%s", name)
}

func getPackage(pack Package, target string) (string, error) {
	fmt.Println(pack.String())
	var err error
	goget := exec.Command("go", "get", "-d", "-u", pack.url)
	err = goget.Run()
	if err != nil {
		return "", fmt.Errorf("Error getting package %s", pack.name)
	}

	gitClient := newGit(os.Getenv("GOPATH"), pack.url)
	if pack.version != "latest" {
		err = gitClient.checkout(target)
		if err != nil {
			return "", fmt.Errorf("Error checking out to %s version tag %s: %s", pack.name, pack.version, err)
		}
	}
	hash, err := gitClient.currentCommitHash()
	if err != nil {
		return "", fmt.Errorf("Error getting last commit hash: %s", err)
	}

	return hash, nil
}
