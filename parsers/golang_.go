// ----------------------------------------------------------------------------
// the code here is about what a Go code parser should be able to do
// ----------------------------------------------------------------------------
package parsers

import (
	"strings"

	core "github.com/ninjawule/otodok/core"
)

func init() {
	core.RegisterAdapter(".go", &golangParser{})
}

type golangParser struct{}

// RemoveTrailingComment should remove the inline trailing comment, if any, e.g.:
// - "some instruction // some comment" => "some instruction "
// - "	// some comment "               => "	"
// - "// some comment "                 => ""
func (thisParser *golangParser) RemoveTrailingComment(line string) string {
	if index := strings.Index(line, "//"); index > -1 {
		return line[:index]
	}

	if index := strings.Index(line, "/*"); index > -1 && strings.Contains(line, "*/") {
		return line[:index]
	}

	return line
}

// IsBlockCommentStart : are we detecting a starting comment block, e.g. (in Go): "/* bla bla"
func (thisParser *golangParser) IsBlockCommentStart(line string) bool {
	return strings.HasPrefix(line, "/*")
}

// IsBlockCommentStart : are we detecting a starting comment block, e.g. (in Go): "bla bla */"
func (thisParser *golangParser) IsBlockCommentEnd(line string) bool {
	return strings.HasSuffix(line, "*/")
}

// IsFunctionDeclStart : are we detecting a starting function declaration,
// e.g. (in Go): "func doThis(..." or "func (r *receiver) doThis(..."
func (thisParser *golangParser) IsFunctionDeclStart(line string) bool {
	return strings.HasPrefix(line, "func ")
}

// IsNewScopeOpening : are we witnessing the opening of a new scope ? e.g. (in Go): " ... {"
func (thisParser *golangParser) IsNewScopeOpening(line string) bool {
	return strings.HasSuffix(line, "{")
}

// IsCurrentScopeClosing : are we witnessing the closing of the current scope ? e.g. (in Go): " ... }"
func (thisParser *golangParser) IsCurrentScopeClosing(line string) bool {
	return strings.HasPrefix(line, "}")
}

// ParseFunctionDetails must return the name of the function declared on the given line,
// and if the function is a method, the "class" it belongs to
//
//nolint:gocritic // TODO fix potential -1 index here
func (thisParser *golangParser) ParseFunctionDetails(line string) (funcName core.FuncName, ownerType core.TypeName) {
	// we either have: 1) "func funcName(..." -> "funcName(..."
	// or            : 2) "func (r *ownerType) funcName(..." -> "(r *ownerType) funcName(..."
	funcLine := line[len("func "):]

	// now we either have: a) "funcName(..." -> "funcName(..."
	// or                : b) "(r *ownerType) funcName(..." -> "(r *ownerType) funcName(..."
	if isMethod := funcLine[0] == '('; isMethod {
		// from: b) "(r *ownerType) funcName(..."
		ownerType = core.TypeName(funcLine[strings.Index(funcLine, " ")+1 : strings.Index(funcLine, ")")])
		// "(r *ownerType) funcName(..." -> "*ownerType) funcName(..."
		funcLine = funcLine[strings.Index(funcLine, string(ownerType)):]
		// from: "*ownerType) funcName(...""
		funcName = core.FuncName(funcLine[strings.Index(funcLine, " ")+1 : strings.Index(funcLine, "(")])
	} else {
		// from: a) "funcName(..."
		funcName = core.FuncName(funcLine[:strings.Index(funcLine, "(")])
	}

	return
}
