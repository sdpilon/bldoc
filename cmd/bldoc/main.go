// Command bldoc is the CLI entry point.
package main

import (
	"os"

	"bldoc/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
