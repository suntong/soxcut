////////////////////////////////////////////////////////////////////////////
// Program: soxcut
// Purpose: sox wrapper tool
// Authors: Tong Sun (c) 2025-2025, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"os"

	"github.com/go-easygen/go-flags/clis"
)

// *** Sub-command: splice ***

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The SpliceCommand type defines all the configurable options from cli.
type SpliceCommand struct {
	FileList string `short:"l" long:"list" env:"SOXCUT_FILELIST" description:"the list file containing sources to splice (mandatory)" required:"true"`
}

var spliceCommand SpliceCommand

////////////////////////////////////////////////////////////////////////////
// Function definitions

func init() {
	gfParser.AddCommand("splice",
		"splice sources for smooth transition",
		`Example:
  soxcut splice -l <listFile> [-o <outputFile>] [sox_effects...]
  soxcut splice -l audio-files.lst -E 200 -L 100 -- compand 0.3,1 6:-70,-60,-20,-10,-5,-5 0 -90 0.1 rate 96k pad 0.5 15

`,
		&spliceCommand)
}

func (x *SpliceCommand) Execute(args []string) error {
	fmt.Fprintf(os.Stderr, "splice sources for smooth transition\n")
	// fmt.Fprintf(os.Stderr, "Copyright (C) 2025-2025, Tong Sun\n\n")
	clis.Setup("soxcut::splice", Opts.Verbose)
	clis.Verbose(1, "Doing Splice, with %+v, %+v", Opts, args)
	// fmt.Println(x.FileList)
	return x.Exec(args)
}

// // Exec implements the business logic of command `splice`
// func (x *SpliceCommand) Exec(args []string) error {
// 	// err := ...
// 	// clis.WarnOn("splice::Exec", err)
// 	// or,
// 	// clis.AbortOn("splice::Exec", err)
// 	return nil
// }
