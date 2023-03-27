// ----------------------------------------------------------------------------
// the code here is about how we log stuff
// ----------------------------------------------------------------------------
package core

import (
	"fmt"
	"log"
)

type LogMode int

const (
	LogModeNONE  LogMode = 0
	LogModeDEBUG LogMode = 1
	LogModeTRACE LogMode = 2
)

func (ctx *context) debug(str string, params ...interface{}) {
	if ctx.logMode >= LogModeDEBUG {
		log.Printf(str, params...)
	}
}

func (ctx *context) trace(lineNum int, str string, params ...interface{}) {
	if ctx.logMode >= LogModeTRACE {
		lineNumStr := fmt.Sprintf("%4d", lineNum)
		log.Printf(">>> TRACE >>> ["+lineNumStr+"] "+str, params...)
	}
}
