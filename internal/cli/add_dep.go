package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"bldoc/internal/compile"
	"bldoc/internal/manifest"
)

func newAddDepCmd() *cobra.Command {
	var format string
	var nested bool
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
			if nested && sourceRef.Anchor == "" {
				return reportErr(cmd, fmt.Errorf("--nested requires a source-ref with a '#anchor' (got %q)", args[1]))
			}
			if nested && targetRef.Field == "" {
				return reportErr(cmd, fmt.Errorf("--nested requires a target-ref with a :field (got %q)", args[0]))
			}

			warnAmbiguousAnchor(cmd, sourceRef)

			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}
			dep := manifest.Dep{
				Source: sourceRef.Source,
				Path:   sourceRef.Path,
				Anchor: sourceRef.Anchor,
				Nested: nested,
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
	cmd.Flags().BoolVar(&nested, "nested", false, "include the anchor's nested subsections")
	return cmd
}

// warnAmbiguousAnchor is a best-effort check: for a plain (non-
// breadcrumb) anchor, it reads sourceRef's source file and prints a
// non-blocking warning to cmd's error output if the anchor matches more
// than one heading. A breadcrumb anchor is never checked (it already
// names a specific heading), and any error reading or parsing the
// source file is ignored — add-dep never requires a source file to
// exist, and this check is advisory only.
func warnAmbiguousAnchor(cmd *cobra.Command, sourceRef SourceRef) {
	if sourceRef.Anchor == "" || strings.Contains(sourceRef.Anchor, "/") {
		return
	}
	suggestions, err := compile.AmbiguousAnchorSuggestions(sourceRef.Source, sourceRef.Anchor)
	if err != nil || len(suggestions) == 0 {
		return
	}
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: anchor %q matches %d headings; disambiguate with a breadcrumb path, e.g. %q\n", sourceRef.Anchor, len(suggestions), suggestions[0])
}
