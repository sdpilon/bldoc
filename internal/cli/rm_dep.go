package cli

import (
	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newRmDepCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "rm-dep <target-ref> <source-ref>",
		Short:         "Remove a dependency from a target",
		Args:          exactArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetRef, err := parseTargetRef(args[0])
			if err != nil {
				return reportErr(cmd, err)
			}
			sourceRef, err := parseSourceRef(args[1])
			if err != nil {
				return reportErr(cmd, err)
			}

			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.RemoveDep(m, targetRef.Target, sourceRef.Source, sourceRef.Path, sourceRef.Anchor); err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.Save(manifest.FileName, m); err != nil {
				return reportErr(cmd, err)
			}
			return nil
		},
	}
}
