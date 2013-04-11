package main

import (
	"exec"
	"fmt"
	"handbrake"
	"os"
	"path/filepath"
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
	for _, file := range files {
		std, err := exec.Command("HandBrakeCLI", "-t0", "-i", file)
		if err != nil {
			panic(err)
		}
		meta := handbrake.ParseOutput(std.Err)
		fmt.Println(meta)
		result, err := handbrake.FormatCLIOutput(meta)
		fmt.Println(result)
	}
}
