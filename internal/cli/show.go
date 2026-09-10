package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"bldoc/internal/manifest"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "show <target>",
		Short:         "Show a target's current dependencies",
		Args:          exactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			target, err := manifest.ShowTarget(m, args[0])
			if err != nil {
				return reportErr(cmd, err)
			}
			for _, dep := range target.Deps {
				fmt.Fprintln(cmd.OutOrStdout(), formatDep(dep))
			}
			return nil
		},
	}
}

// formatDep renders a dependency as "source[:path]" for a whole-file
// dependency, or "field <- source[:path] (format: \"...\")" for a
// field-addressed one.
func formatDep(d manifest.Dep) string {
	source := d.Source
	if d.Path != "" {
		source = source + ":" + d.Path
	}
	if d.Field == "" {
		return source
	}
	if d.Format != "" {
		return fmt.Sprintf("%s <- %s (format: %q)", d.Field, source, d.Format)
	}
	return fmt.Sprintf("%s <- %s", d.Field, source)
}
