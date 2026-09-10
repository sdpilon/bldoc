package cli

import "github.com/spf13/cobra"

func newMakeCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "make [<target>]",
		Short:         "Compile a target, or every target if none is given",
		Args:          rangeArgs(0, 1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return notImplemented(cmd)
		},
	}
}
