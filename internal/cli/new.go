package cli

import (
	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "new <target>",
		Short:         "Declare a new target",
		Args:          exactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.AddTarget(m, args[0]); err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.Save(manifest.FileName, m); err != nil {
				return reportErr(cmd, err)
			}
			return nil
		},
	}
}
