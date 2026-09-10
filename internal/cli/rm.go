package cli

import "github.com/spf13/cobra"

func newRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "rm <target>",
		Short:         "Delete a target",
		Args:          exactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return notImplemented(cmd)
		},
	}
}
