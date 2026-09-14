package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"bldoc/internal/compile"
	"bldoc/internal/manifest"
)

// intermediateDir is where compiled intermediates are written, relative
// to the current working directory.
const intermediateDir = ".bldoc"

func newMakeCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "make [<target>]",
		Short:         "Compile a target, or every target if none is given",
		Args:          rangeArgs(0, 1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := manifest.Load(manifest.FileName)
			if err != nil {
				return reportErr(cmd, err)
			}

			targets := m.Targets
			if len(args) == 1 {
				t, err := manifest.ShowTarget(m, args[0])
				if err != nil {
					return reportErr(cmd, err)
				}
				targets = []manifest.Target{t}
			}

			for _, t := range targets {
				if err := compileTarget(t); err != nil {
					return reportErr(cmd, err)
				}
			}
			return nil
		},
	}
}

// compileTarget compiles t and writes its intermediate under
// intermediateDir: "<name>" (or "<name>.<ext>" if t.Ext is set) for a
// raw-mode result, "<name>.json" for a field-mode one.
func compileTarget(t manifest.Target) error {
	result, err := compile.Target(t)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(intermediateDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", intermediateDir, err)
	}

	if result.IsField {
		data, err := json.MarshalIndent(result.Fields, "", "  ")
		if err != nil {
			return fmt.Errorf("encoding intermediate for %q: %w", t.Name, err)
		}
		path := filepath.Join(intermediateDir, t.Name+".json")
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		return nil
	}

	name := t.Name
	if t.Ext != "" {
		name += "." + t.Ext
	}
	path := filepath.Join(intermediateDir, name)
	if err := os.WriteFile(path, result.Raw, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
