package cli

import "github.com/spf13/cobra"

func newNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "new <target>",
		Short:         "Declare a new target",
		Args:          exactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return notImplemented(cmd)
		},
	}
}
