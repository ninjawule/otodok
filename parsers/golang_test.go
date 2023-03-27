package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTrailingComment(t *testing.T) {
	p := &golangParser{}
	assert.Equal(t, "some instruction ", p.RemoveTrailingComment("some instruction // some comment"))
	assert.Equal(t, "   ", p.RemoveTrailingComment("   // some comment "))
	assert.Equal(t, "", p.RemoveTrailingComment("// some comment "))
	assert.Equal(t, "", p.RemoveTrailingComment("/* bla bla */"))
	assert.Equal(t, "yo ", p.RemoveTrailingComment("yo /* bla bla */"))
}

func TestIsBlockCommentStart(t *testing.T) {
	p := &golangParser{}
	assert.Equal(t, true, p.IsBlockCommentStart("/*"))
	assert.Equal(t, true, p.IsBlockCommentStart("/* bla bla bla"))
	assert.Equal(t, false, p.IsBlockCommentStart("// bla bla bla"))
}

func TestIsBlockCommentEnd(t *testing.T) {
	p := &golangParser{}
	assert.Equal(t, false, p.IsBlockCommentEnd("/* bla bla bla"))
	assert.Equal(t, true, p.IsBlockCommentEnd("*/"))
	assert.Equal(t, true, p.IsBlockCommentEnd("bla bla bla */"))
}

func TestParseFunctionDetails(t *testing.T) {
	p := &golangParser{}
	funcName, ownerType := p.ParseFunctionDetails("func init()")
	assert.Equal(t, "init", funcName)
	assert.Equal(t, "", ownerType)
	funcName, ownerType = p.ParseFunctionDetails("func (r *Contract) Validate()")
	assert.Equal(t, "Validate", funcName)
	assert.Equal(t, "*Contract", ownerType)
}
