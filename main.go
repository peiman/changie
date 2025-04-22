// main.go
//
// This is the entry point for the changie application.
// It implements a clean separation between the main function and the application logic,
// which improves testability and maintains a clear separation of concerns.

// Package main is the entry point for the changie application.
//
// The main package contains minimal logic and delegates command execution
// to the cmd package. This approach enhances testability by ensuring that
// most of the application logic is contained in packages that can be
// imported and tested independently.
package main

import (
	"fmt"
	"os"

	"github.com/peiman/changie/cmd"
)

// run executes the application and returns an exit code.
//
// This function delegates to cmd.Execute() for command execution
// and converts any errors into appropriate exit codes:
// - 0: Success
// - 1: Error
//
// Using a separate run function from main allows for easier testing
// and proper error handling without directly calling os.Exit in tests.
func run() int {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

// main is intentionally not covered by tests because it's the program's entry point.
// All logic is tested via the run() function and other commands. The main function's
// sole purpose is to call run() and exit accordingly. Attempting to cover main directly
// would require integration tests or running the built binary separately.
func main() {
	os.Exit(run())
}
