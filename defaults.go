package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var GoPath = os.Getenv("GOPATH")

const Package = "github.com/gokyle/gowik"

var ErrNoValidSourceDir = fmt.Errorf("no valid source directories")

// ReqFiles stores the files needed to start the wiki.
var ReqFiles = []string{
	"templates/body.html",
	"templates/edit.html",
	"templates/footer.html",
	"templates/head.html",
	"templates/navbar.html",
	"pages/index.md",
	"assets/css/bootstrap.css",
}

func initDefaultPaths() {
	paths := make([]string, 0)
	GoPaths := strings.Split(GoPath, ":")

	for _, p := range GoPaths {
		paths = append(paths, filepath.Join(p, "src", Package))
	}

	var err error
	if missing := CheckPaths(Wiki.WikiDir); len(missing) > 0 {
		err = CopyMissing(missing, paths)
	}
	if err != nil {
		panic("missing files: " + err.Error())
	}
}

func CheckPaths(base string) []string {
	missing := make([]string, 0)
	for _, f := range ReqFiles {
		fPath := filepath.Join(base, f)
		if !FileExists(fPath) {
			missing = append(missing, f)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("[!] %s is missing:\n", base)
		for _, f := range missing {
			fmt.Printf("\t[*] %s\n", f)
		}
	}
	return missing
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func CopyFile(source, dest string) (err error) {
	fmt.Printf("[+] copying %s -> %s\n", source, dest)
	if !FileExists(filepath.Dir(dest)) {
		err = os.MkdirAll(filepath.Dir(dest), 0755)
		if err != nil {
			return
		}
	}
	sFile, err := os.Open(source)
	if err != nil {
		return
	}
	defer sFile.Close()
	dFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dFile.Close()
	_, err = io.Copy(dFile, sFile)
	return
}

func CopyMissing(missing, paths []string) (err error) {
	var srcPath string
	for _, p := range paths {
		if srcMissing := CheckPaths(p); len(srcMissing) == 0 {
			srcPath = p
			break
		}
	}
	if srcPath == "" {
		err = ErrNoValidSourceDir
		return
	}
	for _, missingFile := range missing {
		err = CopyFile(filepath.Join(srcPath, missingFile),
			filepath.Join(Wiki.WikiDir, missingFile))
		if err != nil {
			return
		}
	}
	return
}
