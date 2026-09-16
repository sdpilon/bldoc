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
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "mode: "+resolveMode(target))
			if target.Ext != "" {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "ext: "+target.Ext)
			}
			for _, dep := range target.Deps {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), formatDep(dep))
			}
			return nil
		},
	}
}

// resolveMode returns t's declared mode if it has one, otherwise the
// mode computed from its dependencies (matching the compile engine's own
// inference for a target with no declared mode).
func resolveMode(t manifest.Target) string {
	if t.Mode != "" {
		return t.Mode
	}
	if len(t.Deps) > 0 && t.Deps[0].Field != "" {
		return "field"
	}
	return "raw"
}

// formatDep renders a dependency as "source[:path|#anchor]" for a
// whole-file dependency, or "field <- source[:path|#anchor] [(nested)]
// [(format: \"...\")]" for a field-addressed one.
func formatDep(d manifest.Dep) string {
	source := d.Source
	switch {
	case d.Path != "":
		source = source + ":" + d.Path
	case d.Anchor != "":
		source = source + "#" + d.Anchor
	}
	if d.Field == "" {
		return source
	}
	result := fmt.Sprintf("%s <- %s", d.Field, source)
	if d.Nested {
		result += " (nested)"
	}
	if d.Format != "" {
		result += fmt.Sprintf(" (format: %q)", d.Format)
	}
	return result
}
