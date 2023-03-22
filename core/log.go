// ----------------------------------------------------------------------------
// the code here is about how we log stuff
// ----------------------------------------------------------------------------
package core

import (
	"log"
)

var (
	Verbose bool = false
)

func info(fmt string, params ...interface{}) {
	if Verbose {
		log.Printf(fmt, params...)
	}
}
