// ----------------------------------------------------------------------------
// how we model documentation bits
// ----------------------------------------------------------------------------
package core

// a documentation bit is 1 ou N consecutive functional comments, that starts with a specific key
type docbit struct {
	content []*docline
	// involvedIn map[scopeUID]*scope
}

func (db *docbit) addLine(content string, lineNum int) {
	db.content = append(db.content, &docline{
		content: content,
		lineNum: lineNum,
	})
}

type docline struct {
	content string
	lineNum int
}
