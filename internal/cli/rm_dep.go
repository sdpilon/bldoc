package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRmDepCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "rm-dep <target-ref> <source-ref>",
		Short:         "Remove a dependency from a target",
		Args:          exactArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := parseTargetRef(args[0]); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "bldoc rm-dep: %v\n", err)
				return err
			}
			if _, err := parseSourceRef(args[1]); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "bldoc rm-dep: %v\n", err)
				return err
			}
			return notImplemented(cmd)
		},
	}
}
