// Package cli implements bldoc's command-line interface: the root
// command, its subcommands, and their shared argument-parsing helpers.
package cli

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is the CLI's reported version. It is overwritten at release
// build time via -ldflags "-X bldoc/internal/cli.version=...", so it
// must stay a var rather than a const.
var version = "0.0.0-dev"

// NewRootCmd builds the bldoc root command with every subcommand
// attached.
func NewRootCmd() *cobra.Command {
	info, _ := debug.ReadBuildInfo()
	root := &cobra.Command{
		Use:     "bldoc",
		Version: buildVersion(version, info),
	}
	root.AddCommand(newNewCmd())
	root.AddCommand(newAddDepCmd())
	root.AddCommand(newRmDepCmd())
	root.AddCommand(newMakeCmd())
	root.AddCommand(newRmCmd())
	root.AddCommand(newRenameCmd())
	root.AddCommand(newShowCmd())
	root.AddCommand(newListCmd())
	return root
}

// buildVersion enriches version with embedded VCS build info from info
// (commit short SHA, commit time, and a "modified" flag) when info
// carries a "vcs.revision" setting — the same information `go version
// -m <binary>` reports. It returns version unchanged when info is nil
// or carries no "vcs.revision" setting (e.g. built with
// -buildvcs=false, or from a source tree with no VCS present).
func buildVersion(version string, info *debug.BuildInfo) string {
	if info == nil {
		return version
	}

	var revision, vcsTime string
	modified := false
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			vcsTime = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if revision == "" {
		return version
	}

	short := revision
	if len(short) > 12 {
		short = short[:12]
	}
	detail := "commit " + short
	if vcsTime != "" {
		detail += ", " + vcsTime
	}
	if modified {
		detail += ", modified"
	}
	return fmt.Sprintf("%s (%s)", version, detail)
}

// reportErr prints err to the command's stderr in the standard "bldoc
// <command>: <message>" form and returns it, so a RunE can `return
// reportErr(cmd, err)`.
func reportErr(cmd *cobra.Command, err error) error {
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: %v\n", cmd.Name(), err)
	return err
}

// exactArgs returns a cobra.PositionalArgs validator requiring exactly n
// positional arguments, printing its own usage-error message on failure
// (SilenceErrors/SilenceUsage on the root means nothing else will).
func exactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: expected exactly %d argument(s), got %d\n", cmd.Name(), n, len(args))
			return fmt.Errorf("wrong number of arguments")
		}
		return nil
	}
}

// rangeArgs returns a cobra.PositionalArgs validator requiring between
// minArgs and maxArgs (inclusive) positional arguments.
func rangeArgs(minArgs, maxArgs int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < minArgs || len(args) > maxArgs {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: expected between %d and %d argument(s), got %d\n", cmd.Name(), minArgs, maxArgs, len(args))
			return fmt.Errorf("wrong number of arguments")
		}
		return nil
	}
}
