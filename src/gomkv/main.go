package main

import (
	"path/filepath"
	"fmt"
	"os"
	"exec"
	"handbrake"
)

var workingDir string
var files []string

func init() {
	var err error
	workingDir, err = filepath.Abs(".")
	if err != nil {
		panic(err)
	}
}

func main() {
	workingDir = filepath.Clean(workingDir)
	fmt.Println(workingDir)
	mkv, err := filepath.Glob(workingDir + "/*.mkv")
	if err != nil {
		panic(err)
	}
	m4v, err := filepath.Glob(workingDir + "/*.m4v")
	if err != nil {
		panic(err)
	}
	files = append(mkv, m4v...)
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No mkv/m4v files found in path. Exiting.\n")
		os.Exit(1)
	}
	std, err := exec.Command("HandBrakeCLI", "-t0", "-i", files[0])
	meta := handbrake.ParseOutput(std.Err)
	fmt.Print(meta)
}
