package main

import (
	"flag"
	"os"
	"strings"

	"github.com/ninjawule/otodok/core"
)

func main() {
	// reading the arguments
	var goDirPath, excluded string
	flag.StringVar(&goDirPath, "f", "", "required: the folder containing the go source code to extract the documentation from")
	flag.BoolVar(&core.Verbose, "v", false, "activates the verbose mode")
	flag.StringVar(&excluded, "x", "vendor", "the folders to exclude, separated by a comma, e.g. dir1,dir2,dir3")
	flag.Parse()

	// controlling the args
	if goDirPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// calling the main function
	core.HandleDirectory(goDirPath, strings.Split(excluded, ","))
}
