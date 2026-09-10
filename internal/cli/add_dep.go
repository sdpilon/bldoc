package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newAddDepCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:           "add-dep <target-ref> <source-ref>",
		Short:         "Declare a dependency on a target",
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
			if format != "" && targetRef.Field == "" {
				return reportErr(cmd, fmt.Errorf("--format requires a target-ref with a :field (got %q)", args[0]))
			}

			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			dep := manifest.Dep{
				Source: sourceRef.Source,
				Path:   sourceRef.Path,
				Field:  targetRef.Field,
				Format: format,
			}
			if err := manifest.AddDep(m, targetRef.Target, dep); err != nil {
				return reportErr(cmd, err)
			}
			if err := manifest.Save(manifest.FileName, m); err != nil {
				return reportErr(cmd, err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "template used to render the field's value")
	return cmd
}
