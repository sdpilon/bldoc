package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.0-dev"

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "bldoc",
		Version: version,
	}
	root.AddCommand(newNewCmd())
	root.AddCommand(newAddDepCmd())
	root.AddCommand(newRmDepCmd())
	root.AddCommand(newMakeCmd())
	root.AddCommand(newRmCmd())
	root.AddCommand(newShowCmd())
	root.AddCommand(newListCmd())
	return root
}

// notImplemented prints the standard stub message for a command whose
// arguments were valid but whose real behavior doesn't exist yet.
func notImplemented(cmd *cobra.Command) error {
	fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: not yet implemented\n", cmd.Name())
	return errors.New("not yet implemented")
}

// reportErr prints err to the command's stderr in the standard "bldoc
// <command>: <message>" form and returns it, so a RunE can `return
// reportErr(cmd, err)`.
func reportErr(cmd *cobra.Command, err error) error {
	fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: %v\n", cmd.Name(), err)
	return err
}

// exactArgs returns a cobra.PositionalArgs validator requiring exactly n
// positional arguments, printing its own usage-error message on failure
// (SilenceErrors/SilenceUsage on the root means nothing else will).
func exactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: expected exactly %d argument(s), got %d\n", cmd.Name(), n, len(args))
			return fmt.Errorf("wrong number of arguments")
		}
		return nil
	}
}

// rangeArgs returns a cobra.PositionalArgs validator requiring between
// min and max (inclusive) positional arguments.
func rangeArgs(min, max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < min || len(args) > max {
			fmt.Fprintf(cmd.ErrOrStderr(), "bldoc %s: expected between %d and %d argument(s), got %d\n", cmd.Name(), min, max, len(args))
			return fmt.Errorf("wrong number of arguments")
		}
		return nil
	}
}
