package cli

import "github.com/spf13/cobra"

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "show <target>",
		Short:         "Show a target's current dependencies",
		Args:          exactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return notImplemented(cmd)
		},
	}
}
