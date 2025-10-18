// soxcut - sox cut tool

// Audio file manipulating with sox

package main

////////////////////////////////////////////////////////////////////////////
// Program: soxcut
// Purpose: sox cut tool
// Authors: Tong Sun (c) 2025-2025, All rights reserved
////////////////////////////////////////////////////////////////////////////

//go:generate sh soxcut_cliGen.sh
//go:generate emd gen -in README.beg.e.md -in README.e.md -in README.end.e.md -out README.md

import (
	"fmt"
	"os"

	"github.com/go-easygen/go-flags"
)

// for `go generate -x`
//go:generate sh soxcut_cliGen.sh

//////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	progname = "soxcut"
	version  = "0.1.0"
	date     = "2025-10-17"

	// Opts store all the configurable options
	Opts OptsT
)

var gfParser = flags.NewParser(&Opts, flags.Default)

////////////////////////////////////////////////////////////////////////////
// Function definitions

// ==========================================================================
// Function main
func main() {
	Opts.Version = showVersion
	Opts.Verbflg = func() {
		Opts.Verbose++
	}

	if _, err := gfParser.Parse(); err != nil {
		fmt.Println()
		gfParser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	fmt.Println()
	//DoSoxcut()
}

//==========================================================================
// support functions

func showVersion() {
	fmt.Fprintf(os.Stderr, "soxcut - sox cut tool, version %s\n", version)
	fmt.Fprintf(os.Stderr, "Built on %s\n", date)
	fmt.Fprintf(os.Stderr, "Copyright (C) 2025-2025, Tong Sun\n\n")
	fmt.Fprintf(os.Stderr, "Audio file manipulating with sox\n")
	os.Exit(0)
}
