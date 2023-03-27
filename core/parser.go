// ----------------------------------------------------------------------------
// the code here is about what a generic code parser should be able to do
// ----------------------------------------------------------------------------
package core

var (
	allParsers = map[string]IParser{}
)

func RegisterAdapter(ext string, parser IParser) {
	// only doing this once
	if allParsers[ext] == nil {
		allParsers[ext] = parser
	}
}

type IParser interface {
	// RemoveTrailingComment should remove the inline trailing comment, if any, e.g. (in Go):
	// - "some instruction // some comment" => "some instruction "
	// - "	// some comment "               => "	"
	// - "// some comment "                 => ""
	RemoveTrailingComment(line string) string

	// IsBlockCommentStart : are we detecting a starting comment block, e.g. (in Go): "/* bla bla"
	IsBlockCommentStart(line string) bool

	// IsBlockCommentEnd : are we detecting a finishing comment block, e.g. (in Go): "bla bla */"
	IsBlockCommentEnd(line string) bool

	// IsFunctionDeclStart : are we detecting a starting function declaration,
	// e.g. (in Go): "func doThis(...)" or "func (r *receiver) doThis(...)"
	IsFunctionDeclStart(line string) bool

	// IsNewScopeOpening : are we witnessing the opening of a new scope ? e.g. (in Go): " ... {"
	IsNewScopeOpening(line string) bool

	// IsCurrentScopeClosing : are we witnessing the closing of the current scope ? e.g. (in Go): " ... }"
	IsCurrentScopeClosing(line string) bool

	// ParseFunctionDetails must return the name of the function declared on the given line,
	// and if the function is a method, the "class" it belongs to
	ParseFunctionDetails(line string) (funcName FuncName, ownerType TypeName)
}
