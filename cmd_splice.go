////////////////////////////////////////////////////////////////////////////
// Program: soxcut
// Purpose: sox cut tool
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
	DurExcess int    `short:"E" long:"excess" env:"SOXCUT_DUREXCESS" description:"excess duration of the cross-fade overlap in ms" default:"500"`
	DurLeeway int    `short:"L" long:"leeway" env:"SOXCUT_DURLEEWAY" description:"leeway duration for finding best splice point in ms" default:"200"`
	FileI     string `short:"i" long:"input" env:"SOXCUT_FILEI" description:"the source to cut from (mandatory)" required:"true"`
	FileS     string `short:"s" long:"segments" env:"SOXCUT_FILES" description:"the segments definition file (mandatory)" required:"true"`
	FileO     string `short:"o" long:"output" env:"SOXCUT_FILEO" description:"the final output file" default:"output.mp3"`
}

var spliceCommand SpliceCommand

////////////////////////////////////////////////////////////////////////////
// Function definitions

func init() {
	gfParser.AddCommand("splice",
		"sox splice for smooth transition",
		"",
		&spliceCommand)
}

func (x *SpliceCommand) Execute(args []string) error {
	fmt.Fprintf(os.Stderr, "sox splice for smooth transition\n")
	// fmt.Fprintf(os.Stderr, "Copyright (C) 2025-2025, Tong Sun\n\n")
	clis.Setup("soxcut::splice", Opts.Verbose)
	clis.Verbose(1, "Doing Splice, with %+v, %+v", Opts, args)
	// fmt.Println(x.DurExcess, x.DurLeeway, x.FileI, x.FileS, x.FileO)
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
