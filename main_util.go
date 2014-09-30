package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateDirectory(srcDir string) (string, error) {
	var err error
	workingDir := filepath.Clean(srcDir)
	if workingDir, err = filepath.Abs(workingDir); err != nil {
		return "", fmt.Errorf("Error resolving absolute file path %q: %v", srcDir, err)
	}
	if fileinfo, err := os.Stat(workingDir); err != nil {
		return "", fmt.Errorf("Error accessing %q: %v", workingDir, err)
	} else {
		if !fileinfo.IsDir() {
			return "", fmt.Errorf("%s is not a directory", workingDir)
		}
	}
	return workingDir, nil
}
