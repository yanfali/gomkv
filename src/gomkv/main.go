package main

import (
	"path/filepath"
	"fmt"
	"os"
)

var workingDir string

func init() {
	var err error
	workingDir, err = filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	workingDir = filepath.Clean(workingDir)
	fmt.Println(workingDir)
	mkv, err := filepath.Glob(workingDir + "/*.mkv")
	m4v, err := filepath.Glob(workingDir + "/*.m4v")
	files := append(mkv, m4v...)
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No mkv/m4v files found in path. Exiting.\n")
		os.Exit(1)
	}
}

func main() {
}
