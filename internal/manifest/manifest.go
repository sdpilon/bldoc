// Package manifest reads, writes, and mutates bldoc.toml — the
// committed record of every target and its dependencies.
package manifest

// Manifest is the root of bldoc.toml: an ordered list of targets.
type Manifest struct {
	Targets []Target `toml:"target"`
}

// Target is one entry in the manifest: a name, its ordered list of
// dependencies, and optionally a declared Mode ("raw" or "field") and,
// for a raw-mode target, an output Ext. Mode is empty for a target
// created before this capability existed (or hand-authored without it),
// in which case its mode is inferred from its dependencies instead. Ext
// is only ever set together with Mode "raw".
type Target struct {
	Name string `toml:"name"`
	Mode string `toml:"mode,omitempty"`
	Ext  string `toml:"ext,omitempty"`
	Deps []Dep  `toml:"dep"`
}

// Dep is one dependency recorded on a target: a source ref (Source, and
// at most one of Path into its structured content or Anchor into one of
// its Markdown headings) and, for a field-mode dependency, the
// target-side Field name and optional Format template. Field is empty
// for a whole-file (raw-mode) dependency. Nested is only meaningful
// alongside Anchor: it expands section capture to include the matched
// heading's nested subsections.
type Dep struct {
	Source string `toml:"source"`
	Path   string `toml:"path,omitempty"`
	Anchor string `toml:"anchor,omitempty"`
	Nested bool   `toml:"nested,omitempty"`
	Field  string `toml:"field,omitempty"`
	Format string `toml:"format,omitempty"`
}
