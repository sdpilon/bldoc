package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "list",
		Short:         "List all targets",
		Args:          exactArgs(0),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			for _, name := range manifest.ListTargets(m) {
				fmt.Fprintln(cmd.OutOrStdout(), name)
			}
			return nil
		},
	}
}
