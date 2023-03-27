// ----------------------------------------------------------------------------
// the code here is about how to apply otodok to a folder and subfolders
// ----------------------------------------------------------------------------
package core

import (
	"fmt"
	"os"
	"path"
)

func ParseDirectory(godir string, excludedPaths []string, docbitKey string, logMode LogMode, baseURL string) error {
	// using a context to keep track of all the work done along the way
	ctx := &context{
		excluded:       map[string]bool{},
		godir:          godir,
		logMode:        logMode,
		docbitKey:      docbitKey,
		baseURL:        baseURL,
		includedScopes: []*scope{},
	}

	// mapping the excluded paths
	for _, excludedPath := range excludedPaths {
		ctx.excluded[excludedPath] = true
	}

	// first, we go for all the doc comments and the function declarations
	if errParse := ctx.withCurrentDir("").parseDirectory("", findFunctionDeclarationsAndSpecsInLine); errParse != nil {
		return errParse
	}

	ctx.debug("================================================================================")

	usefulScopes := []*scope{}

	for _, scope := range ctx.includedScopes {
		if len(scope.allDocbits) > 0 {
			usefulScopes = append(usefulScopes, scope)

			ctx.debug("")
			ctx.debug("Scope: %s", scope)

			for _, db := range scope.allDocbits {
				for _, line := range db.content {
					ctx.debug("%s [%d]", line.content, line.lineNum)
				}

				ctx.debug(" + ")
			}
		}
	}

	ctx.debug("================================================================================")

	// // then, we go for all the function uses
	// if errParse := ctx.withCurrentDir("").parseDirectory("", buildFunctionUtilizationsTree); errParse != nil {
	// 	return errParse
	// }

	return nil
}

func (ctx *context) parseDirectory(relativePath string, lineHandlingFn lineHandlingFunc) error {
	currentDir := path.Join(ctx.currentRootScope.dir, relativePath)

	ctx.debug("looking for source code in '%s'", currentDir)

	// reading the given path
	entries, errDir := os.ReadDir(path.Join(ctx.godir, ctx.currentRootScope.dir, relativePath))
	if errDir != nil {
		return fmt.Errorf("error while listing files at path '%s'. Cause: %w", relativePath, errDir)
	}

	// let's list the files
	for _, entry := range entries {
		if entry.IsDir() {
			if !ctx.excluded[entry.Name()] {
				if errHandle := ctx.withCurrentDir(currentDir).parseDirectory(entry.Name(), lineHandlingFn); errHandle != nil {
					return errHandle
				}
			}
		} else {
			if errHandle := ctx.withCurrentDir(currentDir).parseFile(entry.Name(), lineHandlingFn); errHandle != nil {
				return errHandle
			}
		}
	}

	return nil
}
