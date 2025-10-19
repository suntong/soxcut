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

// *** Sub-command: extract ***

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The ExtractCommand type defines all the configurable options from cli.
type ExtractCommand struct {
	FileI string `short:"i" long:"input" env:"SOXCUT_FILEI" description:"the source to cut from (mandatory)" required:"true"`
	FileS string `short:"s" long:"segments" env:"SOXCUT_FILES" description:"the segments definition file (mandatory)" required:"true"`
}

var extractCommand ExtractCommand

////////////////////////////////////////////////////////////////////////////
// Function definitions

func init() {
	gfParser.AddCommand("extract",
		"extract segments from source and splice them for smooth transition",
		`Example:
  soxcut extract -i <inputFile> -s <segmentsFile> -o <outputFile> [sox_effects...]
  soxcut extract -i input1.wav -s timings.txt -o output.mp3 -f="-C 128"
  soxcut extract -i audio.flac -s timings.txt -o final.opus -f="-C 16" -v -- gain -n highpass 80 pad 0 5

`,
		&extractCommand)
}

func (x *ExtractCommand) Execute(args []string) error {
	fmt.Fprintf(os.Stderr, "extract segments from source and splice them for smooth transition\n")
	// fmt.Fprintf(os.Stderr, "Copyright (C) 2025-2025, Tong Sun\n\n")
	clis.Setup("soxcut::extract", Opts.Verbose)
	clis.Verbose(1, "Doing Extract, with %+v, %+v", Opts, args)
	// fmt.Println(x.FileI, x.FileS)
	return x.Exec(args)
}

// // Exec implements the business logic of command `extract`
// func (x *ExtractCommand) Exec(args []string) error {
// 	// err := ...
// 	// clis.WarnOn("extract::Exec", err)
// 	// or,
// 	// clis.AbortOn("extract::Exec", err)
// 	return nil
// }
