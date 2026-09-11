// Package manifest reads, writes, and mutates bldoc.toml — the
// committed record of every target and its dependencies.
package manifest

// Manifest is the root of bldoc.toml: an ordered list of targets.
type Manifest struct {
	Targets []Target `toml:"target"`
}

// Target is one entry in the manifest: a name and its ordered list of
// dependencies.
type Target struct {
	Name string `toml:"name"`
	Deps []Dep  `toml:"dep"`
}

// Dep is one dependency recorded on a target: a source ref (Source,
// optional Path into its structured content) and, for a field-mode
// dependency, the target-side Field name and optional Format template.
// Field is empty for a whole-file (raw-mode) dependency.
type Dep struct {
	Source string `toml:"source"`
	Path   string `toml:"path,omitempty"`
	Field  string `toml:"field,omitempty"`
	Format string `toml:"format,omitempty"`
}
