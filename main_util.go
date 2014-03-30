package main

import (
	"errors"
	"os"
	"path/filepath"
)

func validateDirectory(srcDir string) (string, error) {
	var err error
	workingDir := filepath.Clean(srcDir)
	if workingDir, err = filepath.Abs(workingDir); err != nil {
		return "", err
	}
	if fileinfo, err := os.Stat(workingDir); err != nil {
		return "", err
	} else {
		if !fileinfo.IsDir() {
			return "", errors.New(workingDir + " is not a directory")
		}
	}
	return workingDir, nil
}
