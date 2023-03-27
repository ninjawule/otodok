package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	core "github.com/ninjawule/otodok/core"
	_ "github.com/ninjawule/otodok/parsers"
)

func main() {
	// reading the arguments
	var goDirPath, excluded, docbitKey, baseURL string
	var verbose, verboser bool
	flag.StringVar(&goDirPath, "dir", "", "required: the folder containing the go source code to extract the documentation from")
	flag.StringVar(&docbitKey, "key", "// §", "the key used at the start of the documentation lines")
	flag.BoolVar(&verbose, "v", false, "activates the verbose mode")
	flag.BoolVar(&verboser, "vv", false, "activates the verbose mode")
	flag.StringVar(&excluded, "exclude", "vendor", "the folders to exclude, separated by a comma, e.g. dir1,dir2,dir3")
	flag.StringVar(&baseURL, "url", "", "required: the base URL used to help find each file in its repo "+
		"(e.g. https://github.com/ninjawule/otodok/blob/main/)")
	flag.Parse()

	// controlling the args
	if goDirPath == "" || baseURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// logging mode
	logMode := core.LogModeNONE
	if verbose {
		logMode = core.LogModeDEBUG
	} else if verboser {
		logMode = core.LogModeTRACE
	}

	// calling the main function
	if err := core.ParseDirectory(goDirPath, strings.Split(excluded, ","), docbitKey, logMode, baseURL); err != nil {
		panic(fmt.Errorf("error while handling directory '%s': %w", goDirPath, err))
	}
}
