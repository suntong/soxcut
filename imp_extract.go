////////////////////////////////////////////////////////////////////////////
// Program: soxcut
// Purpose: sox wrapper tool
// Authors: Tong Sun (c) 2025-2025, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

// *** Sub-command: extract ***
// Exec implements the business logic of command `extract`
func (x *ExtractCommand) Exec(args []string) error {
	// err := ...
	// clis.WarnOn("extract::Exec", err)
	// or,
	// clis.AbortOn("extract::Exec", err)
	soxcut(args)
	return nil
}
