package cli

import (
	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "rename <old> <new>",
		Short:         "Rename a target",
		Args:          exactArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.RenameTarget(m, args[0], args[1]); err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.Save(manifest.FileName, m); err != nil {
				return reportErr(cmd, err)
			}
			return nil
		},
	}
}
