// ----------------------------------------------------------------------------
// the code here is about how to apply otodok to a file
// ----------------------------------------------------------------------------
package core

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (ctx *context) parseFile(filePath string, lineHandlingFn lineHandlingFunc) error {
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO remove
	if filePath != "av_request__ws_.go" && filePath != "purchase_order__ws_.go" && filePath != "av_request__ws_options.go" {
		return nil
	}

	// we're handling a new file
	ctx.setCurrentFile(filePath)

	// but maybe we shouldn't
	if ctx.lang == nil {
		ctx.debug("not parsing file '%s' because of its unhandled extension", filePath)

		return nil
	}

	ctx.debug("--------------------------------------------------------------------------------")
	ctx.debug("reading file: %s", filePath)

	// go for the reading
	fileBytes, errRead := os.ReadFile(path.Join(ctx.godir, ctx.currentRootScope.dir, filePath))
	if errRead != nil {
		return fmt.Errorf("error while readling file (%s). Cause: %w", filePath, errRead)
	}

	// getting all the line from the file, unparsed
	fileLines := strings.Split(string(fileBytes), "\n")

	// handling each line
	for i, rawLine := range fileLines {
		// nothing good comes from keeping the leading & trailing spaces - hence the TrimSpace
		if errLine := lineHandlingFn(ctx, strings.TrimSpace(rawLine), i+1); errLine != nil {
			return fmt.Errorf("error while handling line n°%d of file (%s). Cause: %w", i+1, filePath, errLine)
		}
	}

	return nil
}
