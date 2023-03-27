// ----------------------------------------------------------------------------
// how we scopes, i.e. a file, a function or a method, as a docbit container
// ----------------------------------------------------------------------------
//
//nolint:cyclop // no need
package core

import (
	"fmt"
)

// the name for a type
type TypeName string

// the name for a function
type FuncName string

// // the unique ID for a function
// type scopeUID string

type scope struct {
	dir               string    // the directory where we find this scope
	file              string    // the file where we find this scope
	fnName            FuncName  // the name of the function
	fnReceiver        TypeName  // the name of the receiver, if this is a method
	currentBit        *docbit   // the currently built documentation bit
	allDocbits        []*docbit // all the documentation bits contained in this scope
	inBlockComment    bool      // true if, at the line currently scanned, we're inside a block comment
	inFuncDeclaration bool      // true if, at the line currently scanned, we're inside a function declaration block
	openNestedScopes  int       // the number of open nested scopes we have at the moment of scanning
}

func (thisScope *scope) String() string {
	if thisScope.fnReceiver != "" {
		return fmt.Sprintf("%s/%s/%s#%s", thisScope.dir, thisScope.file, thisScope.fnReceiver, thisScope.fnName)
	}

	return fmt.Sprintf("%s/%s/%s", thisScope.dir, thisScope.file, thisScope.fnName)
}

// func buildscopeUID(dir string, file string, fnName FuncName, receiver TypeName) scopeUID {
// 	return scopeUID(fmt.Sprintf("%s-%s-%s", dir, fnName, receiver))
// }

// func (thisScope *scope) getUID() scopeUID {
// 	if thisScope.uid == "" {
// 		thisScope.uid = buildscopeUID(thisScope.dir, thisScope.fnName, thisScope.fnReceiver)
// 	}

// 	return thisScope.uid
// }

// func (thisScope *scope) isRoot() bool {
// 	return thisScope.fnName == ""
// }

// func (thisScope *scope) isFunction() bool {
// 	return thisScope.fnName != "" && !thisScope.isMehod()
// }

// func (thisScope *scope) isMehod() bool {
// 	return thisScope.fnReceiver != ""
// }

// func (thisScope *scope) closeCurrentBit() {
// 	if thisScope.currentBit != nil {
// 		// adding the bit to the pile we have
// 		thisScope.allDocbits = append(thisScope.allDocbits, thisScope.currentBit)

// 		// this bit's building is done
// 		thisScope.currentBit = nil
// 	}
// }

// ----------------------------------------------------------------------------
// scope parsing methods
// ----------------------------------------------------------------------------

func (thisScope *scope) addDocbit(content string, lineNum int) error {
	if thisScope.currentBit == nil {
		// init of the current docbit + adding to the pile of all the docbits contained in this scope
		thisScope.currentBit = &docbit{}
		thisScope.allDocbits = append(thisScope.allDocbits, thisScope.currentBit)
	}

	thisScope.currentBit.addLine(content, lineNum)

	return nil
}

func (thisScope *scope) closeDocbit() {
	thisScope.currentBit = nil
}

func (thisScope *scope) setInBlockComment(val bool) error {
	thisScope.inBlockComment = val

	return nil
}

func (thisScope *scope) setInFuncDeclaration(val bool) {
	thisScope.inFuncDeclaration = val
}
