package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newNewCmd() *cobra.Command {
	var mode string
	var ext string
	cmd := &cobra.Command{
		Use:           "new <target>",
		Short:         "Declare a new target",
		Args:          exactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if mode != "raw" && mode != "field" && mode != "list" {
				return reportErr(cmd, fmt.Errorf("--mode is required and must be \"raw\", \"field\", or \"list\" (got %q)", mode))
			}

			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.AddTarget(m, args[0], mode, ext); err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.Save(manifest.FileName, m); err != nil {
				return reportErr(cmd, err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&mode, "mode", "", "target mode: raw, field, or list (required)")
	cmd.Flags().StringVar(&ext, "ext", "", "output file extension")
	return cmd
}
