// soxcut - sox cut tool
//
// Audio file manipulating with sox

package main

////////////////////////////////////////////////////////////////////////////
// Program: soxcut
// Purpose: sox cut tool
// Authors: Tong Sun (c) 2025-2025, All rights reserved
////////////////////////////////////////////////////////////////////////////

import (
//  	"fmt"
//  	"os"

// "github.com/go-easygen/go-flags"
)

// Template for main starts here

//  // for `go generate -x`
//  //go:generate sh soxcut_cliGen.sh

//////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

//  var (
//          progname  = "soxcut"
//          version   = "0.1.0"
//          date = "2025-10-18"

//  	// Opts store all the configurable options
//  	Opts OptsT
//  )
//
//  var gfParser = flags.NewParser(&Opts, flags.Default)

////////////////////////////////////////////////////////////////////////////
// Function definitions

//==========================================================================
// Function main
//  func main() {
//  	Opts.Version = showVersion
//  	Opts.Verbflg = func() {
//  		Opts.Verbose++
//  	}
//
//  	if _, err := gfParser.Parse(); err != nil {
//  		fmt.Println()
//  		gfParser.WriteHelp(os.Stdout)
//  		os.Exit(1)
//  	}
//  	fmt.Println()
//  	//DoSoxcut()
//  }
//
//  //==========================================================================
//  // support functions
//
//  func showVersion() {
//   	fmt.Fprintf(os.Stderr, "soxcut - sox cut tool, version %s\n", version)
//  	fmt.Fprintf(os.Stderr, "Built on %s\n", date)
//   	fmt.Fprintf(os.Stderr, "Copyright (C) 2025-2025, Tong Sun\n\n")
//  	fmt.Fprintf(os.Stderr, "Audio file manipulating with sox\n")
//  	os.Exit(0)
//  }
// Template for main ends here

// DoSoxcut implements the business logic of command `soxcut`
//  func DoSoxcut() error {
//  	return nil
//  }

// Template for type define starts here

// The OptsT type defines all the configurable options from cli.
type OptsT struct {
	Force   bool   `short:"f" long:"force" description:"Force"`
	Verbflg func() `short:"v" long:"verbose" description:"Verbose mode (Multiple -v options increase the verbosity)"`
	Verbose int
	Version func() `short:"V" long:"version" description:"Show program version and exit"`
}

// Template for type define ends here

// Template for "splice" CLI handling starts here
////////////////////////////////////////////////////////////////////////////
// Program: soxcut
// Purpose: sox cut tool
// Authors: Tong Sun (c) 2025-2025, All rights reserved
////////////////////////////////////////////////////////////////////////////

//  package main

//  import (
//  	"fmt"
//  	"os"
//
//  	"github.com/go-easygen/go-flags/clis"
//  )

// *** Sub-command: splice ***

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The SpliceCommand type defines all the configurable options from cli.
//  type SpliceCommand struct {
//  	DurExcess	int	`short:"E" long:"excess" env:"SOXCUT_DUREXCESS" description:"excess duration of the cross-fade overlap in ms" default:"500"`
//  	DurLeeway	int	`short:"L" long:"leeway" env:"SOXCUT_DURLEEWAY" description:"leeway duration for finding best splice point in ms" default:"200"`
//  	FileI	string	`short:"i" long:"input" env:"SOXCUT_FILEI" description:"the source to cut from (mandatory)" required:"true"`
//  	FileS	string	`short:"s" long:"segments" env:"SOXCUT_FILES" description:"the segments definition file (mandatory)" required:"true"`
//  	FileO	string	`short:"o" long:"output" env:"SOXCUT_FILEO" description:"the final output file" default:"output.mp3"`
//  }

//
//  var spliceCommand SpliceCommand
//
//  ////////////////////////////////////////////////////////////////////////////
//  // Function definitions
//
//  func init() {
//  	gfParser.AddCommand("splice",
//  		"sox splice for smooth transition",
//  		"",
//  		&spliceCommand)
//  }
//
//  func (x *SpliceCommand) Execute(args []string) error {
//   	fmt.Fprintf(os.Stderr, "sox splice for smooth transition\n")
//   	// fmt.Fprintf(os.Stderr, "Copyright (C) 2025-2025, Tong Sun\n\n")
//   	clis.Setup("soxcut::splice", Opts.Verbose)
//   	clis.Verbose(1, "Doing Splice, with %+v, %+v", Opts, args)
//   	// fmt.Println(x.DurExcess, x.DurLeeway, x.FileI, x.FileS, x.FileO)
//  	return x.Exec(args)
//  }
//
// // Exec implements the business logic of command `splice`
// func (x *SpliceCommand) Exec(args []string) error {
// 	// err := ...
// 	// clis.WarnOn("splice::Exec", err)
// 	// or,
// 	// clis.AbortOn("splice::Exec", err)
// 	return nil
// }
// Template for "splice" CLI handling ends here
