// ----------------------------------------------------------------------------
// we use an object here to avoid global variables
// ----------------------------------------------------------------------------
package core

import (
	"fmt"
	"path"
)

// ----------------------------------------------------------------------------
// the parsing context
// ----------------------------------------------------------------------------

type context struct {
	// configuration data
	lang      IParser
	excluded  map[string]bool
	godir     string
	logMode   LogMode
	docbitKey string
	baseURL   string

	// processing data
	currentRootScope *scope
	currentFuncScope *scope
	includedScopes   []*scope
	// includedScopes   map[scopeUID]*scope
}

func (ctx *context) withCurrentDir(dirPath string) *context {
	// changing the dir means changing the current root scope
	ctx.currentRootScope = &scope{
		dir: dirPath,
	}

	return ctx
}

func (ctx *context) setCurrentFile(filePath string) {
	// changing the dir means changing the current root scope, using the previous one for the directory
	ctx.currentRootScope = &scope{
		dir:  ctx.currentRootScope.dir,
		file: filePath,
	}

	// since we're changing the file, we have have to adapt the parser to it
	ctx.lang = allParsers[path.Ext(filePath)]

	// let's add this root scope to all the known scopes
	// ctx.includedScopes[ctx.currentRootScope.getUID()] = ctx.currentRootScope
	ctx.includedScopes = append(ctx.includedScopes, ctx.currentRootScope)
}

func (ctx *context) getCurrentRootOrFunctionScope() *scope {
	if ctx.currentFuncScope != nil {
		return ctx.currentFuncScope
	}

	return ctx.currentRootScope
}

func (ctx *context) openNewFunctionScope(funcName FuncName, typeName TypeName) error {
	if ctx.currentFuncScope != nil {
		return fmt.Errorf("cannot open a new function scope when there's already one open (%s)", ctx.currentFuncScope)
	}

	ctx.currentFuncScope = &scope{
		dir:        ctx.currentRootScope.dir,
		file:       ctx.currentRootScope.file,
		fnName:     funcName,
		fnReceiver: typeName,
	}

	// let's add this root scope to all the known scopes
	ctx.includedScopes = append(ctx.includedScopes, ctx.currentFuncScope)

	return nil
}

func (ctx *context) closeCurrentFunctionScope() {
	ctx.currentFuncScope = nil
}
