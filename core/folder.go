// ----------------------------------------------------------------------------
// the code here is about how to apply otodok to a folder and subfolders
// ----------------------------------------------------------------------------
package core

import (
	"fmt"
	"os"
	"path"
)

const (
	goExt = ".go"
)

func HandleDirectory(relativePath string, excludedPaths []string) error {
	// getting the current working dir
	wd, errWd := os.Getwd()
	if errWd != nil {
		return fmt.Errorf("could not determine the current wd: %w", errWd)
	}

	// computing the targeted absolute path
	absolutePath := path.Join(wd, relativePath)

	// mapping the excluded paths
	excludedPathsMap := map[string]bool{}
	for _, excludedPath := range excludedPaths {
		excludedPathsMap[excludedPath] = true
	}

	return doHandleDirectory(absolutePath, excludedPathsMap)
}

func doHandleDirectory(absolutePath string, excludedPaths map[string]bool) error {
	info("looking for go source code in '%s'", absolutePath)

	// reading the given path
	entries, errDir := os.ReadDir(absolutePath)
	if errDir != nil {
		return fmt.Errorf("error while listing files at path '%s'. Cause: %w", absolutePath, errDir)
	}

	// let's list the files
	for _, entry := range entries {
		if entry.IsDir() && !excludedPaths[entry.Name()] {
			if errHandle := doHandleDirectory(path.Join(absolutePath, entry.Name()), excludedPaths); errHandle != nil {
				return errHandle
			}
		} else {
			if path.Ext(entry.Name()) == goExt {
				info("looking %s", entry.Name())
			}
		}
	}

	return nil
}
