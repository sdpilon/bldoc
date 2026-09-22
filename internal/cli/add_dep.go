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
		Use:           "add-dep <target-ref> <source-ref> | add-dep <target> [<record-name>] <pair> <pair>...",
		Short:         "Declare one or more dependencies on a target",
		Args:          cobra.MinimumNArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 {
				return runAddDepSingle(cmd, args, format, nested)
			}
			return runAddDepBatch(cmd, args, format, nested)
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "template used to render the field's value")
	cmd.Flags().BoolVar(&nested, "nested", false, "include the anchor's nested subsections")
	return cmd
}

// runAddDepSingle is the original two-argument add-dep form: a single
// target-ref and source-ref, unchanged in shape from before the batch
// form existed (its refs now use '@' rather than ':', per the ref
// package's rename).
func runAddDepSingle(cmd *cobra.Command, args []string, format string, nested bool) error {
	targetRef, err := parseTargetRef(args[0])
	if err != nil {
		return reportErr(cmd, err)
	}
	sourceRef, err := parseSourceRef(args[1])
	if err != nil {
		return reportErr(cmd, err)
	}
	if format != "" && targetRef.Field == "" {
		return reportErr(cmd, fmt.Errorf("--format requires a target-ref with a @field (got %q)", args[0]))
	}
	if nested && sourceRef.Anchor == "" {
		return reportErr(cmd, fmt.Errorf("--nested requires a source-ref with a '#anchor' (got %q)", args[1]))
	}
	if nested && targetRef.Field == "" {
		return reportErr(cmd, fmt.Errorf("--nested requires a target-ref with a @field (got %q)", args[0]))
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
}

// runAddDepBatch is the variadic add-dep form: a bare target, an
// explicit record name when the target is list-mode, and two or more
// pairs. It validates every pair against one loaded manifest and only
// calls manifest.Save once every pair has succeeded, so a failure on
// any pair leaves bldoc.toml exactly as it was.
func runAddDepBatch(cmd *cobra.Command, args []string, format string, nested bool) error {
	target := args[0]
	if err := validateBareTarget(target); err != nil {
		return reportErr(cmd, err)
	}

	m, err := manifest.Load(manifest.FileName)
	if err != nil {
		return reportErr(cmd, err)
	}
	t, err := manifest.ShowTarget(m, target)
	if err != nil {
		return reportErr(cmd, err)
	}

	rest := args[1:]
	recordName := ""
	if t.Mode == "list" {
		recordName, rest = rest[0], rest[1:]
		if err := validateRecordName(recordName); err != nil {
			return reportErr(cmd, err)
		}
	}
	if len(rest) < 2 {
		if recordName != "" {
			return reportErr(cmd, fmt.Errorf("add-dep batch form requires at least two pairs after the target and record name %q (got %d)", recordName, len(rest)))
		}
		return reportErr(cmd, fmt.Errorf("add-dep batch form requires at least two pairs after the target (got %d)", len(rest)))
	}

	pairs := make([]pair, 0, len(rest))
	for _, a := range rest {
		p, err := parsePair(a)
		if err != nil {
			return reportErr(cmd, err)
		}
		pairs = append(pairs, p)
	}

	for _, p := range pairs {
		if format != "" && p.Field == "" {
			return reportErr(cmd, fmt.Errorf("--format requires a field-addressed pair (got bare pair %q)", formatPairSource(p)))
		}
		if nested && p.Source.Anchor == "" {
			return reportErr(cmd, fmt.Errorf("--nested requires every pair's source-ref to have a '#anchor' (got %q)", formatPairSource(p)))
		}
		if nested && p.Field == "" {
			return reportErr(cmd, fmt.Errorf("--nested requires a field-addressed pair (got bare pair %q)", formatPairSource(p)))
		}
	}

	for _, p := range pairs {
		warnAmbiguousAnchor(cmd, p.Source)

		field := p.Field
		if t.Mode == "list" {
			if field == "" {
				field = recordName
			} else {
				field = recordName + "." + field
			}
		}
		dep := manifest.Dep{
			Source: p.Source.Source,
			Path:   p.Source.Path,
			Anchor: p.Source.Anchor,
			Nested: nested,
			Field:  field,
			Format: format,
		}
		if err := manifest.AddDep(m, target, dep); err != nil {
			return reportErr(cmd, err)
		}
	}

	if err := manifest.Save(manifest.FileName, m); err != nil {
		return reportErr(cmd, err)
	}
	return nil
}

func formatPairSource(p pair) string {
	s := p.Source.Source
	switch {
	case p.Source.Path != "":
		s += "@" + p.Source.Path
	case p.Source.Anchor != "":
		s += "#" + p.Source.Anchor
	}
	if p.Field == "" {
		return s
	}
	return p.Field + ":" + s
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
