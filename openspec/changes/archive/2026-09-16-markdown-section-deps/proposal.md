## Why

Field- and list-mode dependencies can only pull in a whole file or one
scalar value from a structured TOML/JSON/YAML document today. There is
no way to depend on just one section of a larger prose document (for
example, one `## Purpose` section of an OpenSpec `spec.md`) without
first hand-splitting that document into separate files. Markdown-section
dependencies close that gap by letting a dependency address a single
heading's content directly.

## What Changes

- New source-ref form `source#anchor` addressing one Markdown heading's
  own content by a GitHub-style slug, valid only for `.md`/`.markdown`
  sources — any other extension is rejected. `#anchor` and `:field-path`
  cannot be combined in the same source-ref.
- The anchor may optionally be a breadcrumb path of slugs
  (`parent-slug/child-slug`, root-first) to disambiguate a heading whose
  slug is not unique within the file; a plain single-segment anchor
  remains valid whenever it happens to be unique.
- The resolved value is the heading's own raw multi-line text, up to the
  next heading of any level — not tree-walked or scalar-constrained the
  way `:field-path` resolution is.
- New `--nested` flag on `add-dep`: expands capture to the heading plus
  all of its nested subsections, stopping at the next heading of
  equal-or-shallower level. Requires a `#anchor` source-ref; rejected on
  a raw-mode dependency (a whole-file dependency already exists for
  pulling in more content in raw mode).
- `add-dep` gains one narrow exception to its current no-filesystem-access
  behavior: when a source-ref carries `#anchor`, it reads the source file
  to check whether the given anchor matches more than one heading, and
  if so prints a non-blocking warning suggesting a breadcrumb path — the
  dependency is still recorded either way.
- `make` resolves an anchor-addressed dependency by parsing the source
  file's Markdown headings, locating the (optionally breadcrumb-scoped)
  match, and extracting its content. `make` SHALL still reject an
  ambiguous plain anchor at compile time even if `add-dep`'s warning was
  shown or missed — the warning does not weaken this guarantee.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `cli-shell`: `add-dep`'s source-ref grammar gains the `#anchor` form,
  and `add-dep` gains the `--nested` flag with its own validation rules.
- `manifest-store`: a dependency's recorded state gains its anchor (and
  whether it was recorded with `--nested`); `add-dep` validates an
  anchor's uniqueness against its source file and warns (without
  blocking) when ambiguous; `rm-dep` and duplicate-dependency matching
  account for the anchor as part of a source-ref's identity.
- `compile-engine`: `make` gains Markdown heading parsing and
  anchor-based section resolution, including its error conditions
  (heading not found, ambiguous anchor, non-Markdown source, `--nested`
  on a raw-mode dependency).

## Impact

- `internal/cli`: source-ref parsing (new `#anchor` grammar), `add-dep`'s
  flag surface (`--nested`) and its new file-reading validation path.
- `internal/manifest`: the `Dep` schema and its `bldoc.toml` round-trip
  gain anchor/nested fields; `rm-dep` and duplicate-dependency matching
  logic.
- `internal/compile`: a new Markdown heading parser and section-resolution
  path, exercised in parallel with the existing whole-file and
  `:field-path` resolution paths — neither of those, nor JSON/YAML/TOML
  encoding of the compiled intermediate, changes.
