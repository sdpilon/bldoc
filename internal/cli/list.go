package cli

import "github.com/spf13/cobra"

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "list",
		Short:         "List all targets",
		Args:          exactArgs(0),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return notImplemented(cmd)
		},
	}
}
