// ----------------------------------------------------------------------------
// the code here is about how to parse a line found in source code
// ----------------------------------------------------------------------------
package core

import (
	"fmt"
	"strings"
)

// the signature for the functions that handle the scanning of code lines
type lineHandlingFunc func(ctx *context, line string, num int) error

// ----------------------------------------------------------------------------
// the next function should help finding the functional documention bits
// and gather them into scopes (root scopes, or functions / methods)
// ----------------------------------------------------------------------------
//
//nolint:cyclop,revive // no need
func findFunctionDeclarationsAndSpecsInLine(ctx *context, line string, num int) error {
	// stopper, used for debugging purposes
	if strings.HasSuffix(line, "// OTODOK_STOP") {
		return fmt.Errorf("stopping at line %d", num)
	}

	// parsing: block comments
	if ctx.lang.IsBlockCommentStart(line) { // "/* blabla"
		ctx.trace(num, "Starting block comment")

		return ctx.getCurrentRootOrFunctionScope().setInBlockComment(true)
	}

	if ctx.lang.IsBlockCommentEnd(line) { // "blabla */"
		ctx.trace(num, "Ending block comment")

		return ctx.getCurrentRootOrFunctionScope().setInBlockComment(false)
	}

	if ctx.getCurrentRootOrFunctionScope().inBlockComment { // "/* ... [here] ... */"
		ctx.trace(num, "Still in block comment")

		return nil
	}

	// parsing: documentation comment (docbit)
	if strings.HasPrefix(line, ctx.docbitKey) { // § Some specifications right here
		content := line[len(ctx.docbitKey):]
		ctx.trace(num, "Found doc content: %s", content)

		return ctx.getCurrentRootOrFunctionScope().addDocbit(content, num)
	}

	// at this point, we're not dealing with a docbit here, so we close the current one, if any
	ctx.getCurrentRootOrFunctionScope().closeDocbit()

	// let's also get rid of the trailing comment, if any - then trim the spaces
	line = strings.TrimSpace(ctx.lang.RemoveTrailingComment(line)) // "... { // bla bla" => "... {"

	// parsing: function declarations (simple ones at least, not anonymous ones, for now)
	if ctx.lang.IsFunctionDeclStart(line) { // "func ..."
		fnName, typeName := ctx.lang.ParseFunctionDetails(line)
		ctx.trace(num, "Starting function %s (owner: %s)", fnName, typeName)

		if errOpen := ctx.openNewFunctionScope(fnName, typeName); errOpen != nil {
			return errOpen
		}

		ctx.getCurrentRootOrFunctionScope().setInFuncDeclaration(true) // "func ... \n"
	}

	// we loop here until the function scope is officially open and we can start scanning the function's content
	if ctx.getCurrentRootOrFunctionScope().inFuncDeclaration {
		if ctx.lang.IsNewScopeOpening(line) { // "\n ... {"
			ctx.trace(num, "Function scope officially opened (%s)", ctx.getCurrentRootOrFunctionScope())
			ctx.getCurrentRootOrFunctionScope().setInFuncDeclaration(false)

			return nil
		}

		ctx.trace(num, "Still in function declaration")

		return nil
	}

	// if we're closing a scope and opening a new one in the same line, let's handle it
	if ctx.lang.IsCurrentScopeClosing(line) && ctx.lang.IsNewScopeOpening(line) {
		ctx.trace(num, "Closing } then { opening a nested scope in '%s' (total open: %d)", ctx.getCurrentRootOrFunctionScope(),
			ctx.getCurrentRootOrFunctionScope().openNestedScopes)

		return nil
	}

	// if a new scope opens here, we know it's not a function or method in the root
	if ctx.lang.IsNewScopeOpening(line) {
		ctx.getCurrentRootOrFunctionScope().openNestedScopes++
		ctx.trace(num, "{ Opening a nested scope in '%s' (total open: %d)", ctx.getCurrentRootOrFunctionScope(),
			ctx.getCurrentRootOrFunctionScope().openNestedScopes)

		return nil
	}

	// if a scope closes here, then:
	if ctx.lang.IsCurrentScopeClosing(line) {
		if ctx.getCurrentRootOrFunctionScope().openNestedScopes >= 1 {
			ctx.getCurrentRootOrFunctionScope().openNestedScopes--
			ctx.trace(num, "Closing } a nested scope in '%s' (total open: %d)", ctx.getCurrentRootOrFunctionScope(),
				ctx.getCurrentRootOrFunctionScope().openNestedScopes)
		} else {
			ctx.trace(num, "Closing a function's scope (%s)", ctx.getCurrentRootOrFunctionScope())
			ctx.closeCurrentFunctionScope()
		}

		return nil
	}

	// nothing else - for now

	return nil
}

// ----------------------------------------------------------------------------
// the next function should help chain the scopes to each other, relying on
// the job done by the first function
// ----------------------------------------------------------------------------
// func buildFunctionUtilizationsTree(ctx *context, line string, num int) error {
// 	return nil
// }
