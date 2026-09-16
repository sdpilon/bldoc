# compile-engine Specification

## Purpose

The `compile-engine` capability defines `make`'s real behavior: turning
a target's recorded dependencies (from the `manifest-store` capability)
into a disposable, mechanically-derived intermediate — raw-mode
concatenation or field-mode JSON, always fully recomputed.

## Requirements

### Requirement: Intermediate file location and naming
The compiled intermediate for a target SHALL be written under `.bldoc/`
in the current working directory. For a raw-mode target: with no
recorded extension, at `.bldoc/<target>` (raw bytes); with a recorded
extension, at `.bldoc/<target>.<ext>` (raw bytes). For a field-mode or
list-mode target: at `.bldoc/<target>.<ext>` when an extension is
recorded, otherwise at `.bldoc/<target>.json`. A recorded extension of
`yaml` or `yml` on a field-mode or list-mode target SHALL select native
YAML encoding for the intermediate's content; any other recorded
extension, or none, SHALL select JSON encoding — the file's name always
follows the recorded extension when one is set, independent of which
encoding was selected.

#### Scenario: Raw-mode output path
- **WHEN** `bldoc make README` compiles a raw-mode target with no recorded extension
- **THEN** the CLI writes the concatenated bytes to `.bldoc/README`

#### Scenario: Raw-mode output path with a declared extension
- **WHEN** `bldoc make PROJECTS` compiles a raw-mode target recorded with extension `yaml`
- **THEN** the CLI writes the concatenated bytes to `.bldoc/PROJECTS.yaml`

#### Scenario: Field-mode output path
- **WHEN** `bldoc make README` compiles a field-mode target with no recorded extension
- **THEN** the CLI writes a JSON object to `.bldoc/README.json`

#### Scenario: Field-mode output path with a YAML extension
- **WHEN** `bldoc make README` compiles a field-mode target recorded with extension `yaml`
- **THEN** the CLI writes a YAML-encoded object to `.bldoc/README.yaml`

#### Scenario: List-mode output path with no recorded extension
- **WHEN** `bldoc make PROJECTS` compiles a list-mode target with no recorded extension
- **THEN** the CLI writes a JSON array to `.bldoc/PROJECTS.json`

#### Scenario: List-mode output path with a YAML extension
- **WHEN** `bldoc make PROJECTS` compiles a list-mode target recorded with extension `yaml`
- **THEN** the CLI writes a YAML-encoded array to `.bldoc/PROJECTS.yaml`

### Requirement: `make` requires an existing target
`make <target>` SHALL reject a target name not recorded in the manifest.

#### Scenario: Unknown target rejected
- **WHEN** `bldoc make MISSING` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero

### Requirement: Raw-mode compilation concatenates dependencies
For a target whose recorded dependencies are all whole-file (unnamed),
`make` SHALL write the byte-for-byte concatenation of each dependency's
resolved content, in the order the dependencies were added — a
dependency's resolved content is its source file's entire bytes, or, if
the dependency carries an anchor, that anchor's resolved Markdown
section content (see the Markdown section resolution requirements
below).

#### Scenario: Raw-mode concatenation
- **WHEN** `README` has whole-file dependencies `a.txt` then `b.txt` recorded, and `make README` is run
- **THEN** `.bldoc/README` contains the bytes of `a.txt` immediately followed by the bytes of `b.txt`, identical to `cat a.txt b.txt`

#### Scenario: Raw-mode concatenation with an anchor-addressed dependency
- **WHEN** `README` has a whole-file (no `:field`) dependency recorded with anchor `purpose` in `spec.md`, and `make README` is run
- **THEN** `.bldoc/README` contains `spec.md`'s `purpose` heading's resolved section content

### Requirement: A target with no dependencies compiles to an empty intermediate
`make` SHALL compile a target with no recorded dependencies and no
declared mode, or with a declared raw mode, to an empty raw-mode
intermediate. `make` SHALL compile a target with no recorded dependencies
and a declared field mode to an empty field-mode intermediate (an empty
object, encoded per the target's extension as described above). `make`
SHALL compile a target with no recorded dependencies and a declared list
mode to an empty list-mode intermediate (an empty array, encoded per the
target's extension as described above).

#### Scenario: Empty target
- **WHEN** `README` has no recorded dependencies and no declared mode, and `make README` is run
- **THEN** `.bldoc/README` is created empty (zero bytes)

#### Scenario: Empty declared field-mode target
- **WHEN** `PROJECTS` was created with `bldoc new PROJECTS --mode field`, has no recorded dependencies, and `make PROJECTS` is run
- **THEN** `.bldoc/PROJECTS.json` is created containing `{}`

#### Scenario: Empty declared list-mode target
- **WHEN** `PROJECTS` was created with `bldoc new PROJECTS --mode list`, has no recorded dependencies, and `make PROJECTS` is run
- **THEN** `.bldoc/PROJECTS.json` is created containing `[]`

### Requirement: Field-mode compilation resolves each dependency's value
For a target whose recorded dependencies are all field-addressed, `make`
SHALL resolve each dependency's raw value — the source file's entire
content if the dependency has no field-path or anchor, the value at the
field-path if it has one, or the anchor's resolved Markdown section
content if it has one (see the Markdown section resolution requirements
below) — apply the dependency's format template if present, and emit one
JSON object with one entry per field: `{"value": <rendered>, "raw": <raw
value>}`. `value` equals `raw` when no format template was recorded.

#### Scenario: Field resolved from a whole source file
- **WHEN** `README:summary` depends on `notes.txt` with no field-path and no format template, and `make README` is run
- **THEN** `.bldoc/README.json` contains `{"summary": {"value": "<notes.txt's content>", "raw": "<notes.txt's content>"}}`

#### Scenario: Field resolved from a structured field-path with a format template
- **WHEN** `README:version` depends on `pyproject.toml:project.requires-python` with format `"Python version must be %s to run this project."`, and that field-path resolves to `3.11`, and `make README` is run
- **THEN** `.bldoc/README.json` contains `{"version": {"value": "Python version must be 3.11 to run this project.", "raw": "3.11"}}`

#### Scenario: Field resolved from an anchor-addressed dependency
- **WHEN** `README:summary` depends on `spec.md#purpose` with no format template, and `make README` is run
- **THEN** `.bldoc/README.json` contains `{"summary": {"value": "<purpose heading's resolved section content>", "raw": "<same content>"}}`

### Requirement: Markdown section resolution supports `.md` and `.markdown` sources
`make` SHALL resolve an anchor-addressed dependency only when its source
file's extension is `.md` or `.markdown` (case-insensitive); any other
extension SHALL be rejected.

#### Scenario: Unsupported source format rejected for an anchor
- **WHEN** a dependency has an anchor into a source file whose extension is neither `.md` nor `.markdown`
- **THEN** `make` reports an error and exits non-zero

### Requirement: An anchor resolves by GitHub-style slug matching
`make` SHALL parse a Markdown source's ATX headings (`#` through `######`)
and compute each one's GitHub-style slug (lowercased, spaces replaced
with hyphens, punctuation stripped). A single-segment anchor SHALL match
any heading whose slug equals it, anywhere in the document. A breadcrumb
anchor (`parent-slug/child-slug`, root-first) SHALL match a heading only
when its chain of ancestor headings' slugs equals the breadcrumb's
segments in order.

#### Scenario: Plain anchor matches by slug
- **WHEN** a dependency's anchor is `purpose` and the source file has exactly one heading `## Purpose`
- **THEN** `make` resolves that heading's section content

#### Scenario: Breadcrumb anchor matches by ancestor chain
- **WHEN** a dependency's anchor is `requirement-x/scenario-y` and the source file has a heading `#### Scenario: Y` nested under `### Requirement: X`
- **THEN** `make` resolves that `#### Scenario: Y` heading's section content

### Requirement: An anchor with no matching heading is rejected
`make` SHALL reject an anchor (plain or breadcrumb) that matches no
heading in the source file.

#### Scenario: Heading not found
- **WHEN** a dependency's anchor matches no heading in its source file
- **THEN** `make` reports an error and exits non-zero

### Requirement: An ambiguous plain anchor is rejected at compile time
`make` SHALL reject a single-segment anchor that matches more than one
heading in the source file, regardless of whether `add-dep` already
warned about the ambiguity when the dependency was recorded — the
dependency must use a breadcrumb path to resolve.

#### Scenario: Ambiguous plain anchor rejected
- **WHEN** a dependency's plain anchor matches more than one heading in its source file
- **THEN** `make` reports an error and exits non-zero, independent of any warning `add-dep` printed when the dependency was recorded

#### Scenario: Breadcrumb anchor resolves where a plain anchor would be ambiguous
- **WHEN** a source file has two headings sharing the slug `scenario-name` under different parents, and a dependency's anchor is the breadcrumb `requirement-a/scenario-name`
- **THEN** `make` resolves the one heading matching that breadcrumb, with no ambiguity error

### Requirement: Anchor content capture defaults to the heading's own text
`make` SHALL resolve a non-`--nested` anchor's section content as the
Markdown source between the end of the matched heading's line and the
start of the next heading at any level, or the end of the file if there
is none — excluding the matched heading's own line.

#### Scenario: Default capture stops at the next heading regardless of level
- **WHEN** a dependency's anchor matches `## Requirements`, immediately followed in the source by `### Requirement: X`, and `--nested` was not recorded
- **THEN** the resolved content is empty (or only the heading's own intervening text, if any), not including `### Requirement: X` or anything after it

#### Scenario: Default capture reaches end of file
- **WHEN** a dependency's anchor matches the last heading in its source file
- **THEN** the resolved content is everything from after that heading's line to the end of the file

### Requirement: `--nested` expands capture to the heading plus its nested subsections
`make` SHALL resolve a `--nested` anchor's section content as the
Markdown source between the end of the matched heading's line and the
start of the next heading at a level equal to or shallower than the
matched heading, or the end of the file if there is none — including any
deeper nested headings' lines and content, and excluding the matched
heading's own line.

#### Scenario: Nested capture includes deeper subsections
- **WHEN** a dependency's anchor matches `## Requirements` with `--nested` recorded, and the source has `### Requirement: X` and `#### Scenario: Y` nested under it before the next `##`-or-shallower heading
- **THEN** the resolved content includes `### Requirement: X`, `#### Scenario: Y`, and their text, verbatim

#### Scenario: Nested capture stops at a sibling or shallower heading
- **WHEN** a dependency's anchor matches `## Requirements` with `--nested` recorded, and the source's next `##`-or-shallower heading is `## Impact`
- **THEN** the resolved content does not include `## Impact` or anything after it

### Requirement: Structured field-path resolution supports TOML, JSON, and YAML
`make` SHALL parse a dependency's source file as TOML, JSON, or YAML
based on its file extension (`.toml`, `.json`, `.yaml`/`.yml`
respectively), both when resolving a field-path on a field-mode or
list-mode record-field dependency, and when resolving the whole
top-level document of a list-mode record-only dependency. When resolving
via a field-path, `make` resolves the field-path as a dot-separated
sequence of keys into the parsed document.

#### Scenario: Unsupported source format rejected
- **WHEN** a dependency has a field-path into a source file whose extension is none of `.toml`, `.json`, `.yaml`, or `.yml`
- **THEN** `make` reports an error and exits non-zero

#### Scenario: Unsupported source format rejected for a record-only dependency
- **WHEN** a list-mode target's record-only dependency's source file extension is none of `.toml`, `.json`, `.yaml`, or `.yml`
- **THEN** `make` reports an error and exits non-zero

### Requirement: A malformed structured source file is rejected
`make` SHALL report an error and exit non-zero if a dependency's
field-path requires parsing its source file as TOML/JSON/YAML and that
file cannot be parsed, or if a list-mode record-only dependency's source
file cannot be parsed as TOML/JSON/YAML.

#### Scenario: Malformed source file rejected
- **WHEN** a dependency's field-path source file has the appropriate extension but is not valid TOML/JSON/YAML
- **THEN** `make` reports an error and exits non-zero

#### Scenario: Malformed record-only source file rejected
- **WHEN** a list-mode target's record-only dependency's source file has the appropriate extension but is not valid TOML/JSON/YAML
- **THEN** `make` reports an error and exits non-zero

### Requirement: An unresolvable field-path is rejected
`make` SHALL reject a field-path that does not resolve to a value in the
parsed source document (a missing key at any level).

#### Scenario: Missing key rejected
- **WHEN** a dependency's field-path names a key that does not exist in its source file's parsed content
- **THEN** `make` reports an error and exits non-zero

### Requirement: A field-path must resolve to a scalar value
`make` SHALL reject a field-path that resolves to a table/object or an
array/list, rather than a string, number, or boolean.

#### Scenario: Non-scalar field-path rejected
- **WHEN** a dependency's field-path resolves to a table or array in its source file
- **THEN** `make` reports an error and exits non-zero

### Requirement: Duplicate field names within a target are rejected
`make` SHALL reject a field-mode target whose recorded dependencies use
the same target-side field name more than once. `make` SHALL reject a
list-mode target that would produce the same field name twice within the
same record — whether from two record-field dependencies naming the same
field in that record, from a record-only dependency's merged key
colliding with a record-field dependency's field name in that record, or
from two record-only dependencies on the same record whose source
documents share a top-level key.

#### Scenario: Duplicate field name rejected
- **WHEN** a field-mode target has two recorded dependencies both naming the field `version`
- **THEN** `make` reports an error and exits non-zero

#### Scenario: Duplicate field name rejected within a list-mode record
- **WHEN** a list-mode target has two recorded dependencies both naming field `description` of record `bldoc`
- **THEN** `make` reports an error and exits non-zero

#### Scenario: Merged key colliding with an explicit field is rejected
- **WHEN** a list-mode target's record-only dependency for record `bldoc` merges a source document containing top-level key `description`, and the same target also has a record-field dependency naming field `description` of record `bldoc`
- **THEN** `make` reports an error and exits non-zero

#### Scenario: Same field name in different records is allowed
- **WHEN** a list-mode target has a record-field dependency naming field `description` of record `bldoc` and another naming field `description` of record `projx`
- **THEN** `make` succeeds and both records carry their own `description` field

### Requirement: `make` with no target compiles every target
Given no target argument, `make` SHALL compile every target recorded in
the manifest. Given no manifest file, or a manifest with no targets, it
SHALL do nothing and exit 0.

#### Scenario: All targets compiled
- **WHEN** the manifest has targets `README` and `CHANGELOG`, and `bldoc make` is run with no arguments
- **THEN** the CLI writes both `README`'s and `CHANGELOG`'s intermediates

#### Scenario: No targets is a no-op
- **WHEN** no manifest file exists, or the manifest has no targets, and `bldoc make` is run with no arguments
- **THEN** the CLI exits 0 and writes no intermediates

### Requirement: `make` always fully recompiles
`make` SHALL always recompute a target's intermediate from the current
state of its source files, overwriting any existing intermediate at
that target's path, regardless of whether the sources have changed
since the last `make`.

#### Scenario: Re-running make overwrites the intermediate
- **WHEN** `make README` is run twice in a row with no changes to `README`'s sources or dependencies between runs
- **THEN** both runs succeed and `.bldoc/README`'s (or `.bldoc/README.json`'s) content is rewritten identically both times

### Requirement: A record-only dependency's source must resolve to a flat scalar document
`make` SHALL reject a list-mode record-only dependency whose parsed
source document has a top-level value that is not a string, number, or
boolean (i.e. a nested table/object or an array/list at the top level).

#### Scenario: Non-scalar top-level value rejected
- **WHEN** a list-mode target's record-only dependency's source document has a top-level key whose value is a table or array
- **THEN** `make` reports an error and exits non-zero

#### Scenario: Flat scalar document accepted
- **WHEN** a list-mode target's record-only dependency's source document has only string, number, or boolean top-level values
- **THEN** `make` merges each of those keys into the dependency's record

### Requirement: List-mode compilation merges and groups dependencies into records
For a list-mode target, `make` SHALL resolve each recorded dependency
into its declared record: a record-only dependency merges every
top-level key of its parsed source document into that record, using each
key's value as-is; a record-field dependency resolves its value exactly
as a field-mode dependency would (the whole source file's content when
it has no field-path, or the value at its field-path when it has one),
applies its format template if present, and stores the rendered result
under its field name in that record. `make` SHALL emit one array
containing one object per distinct record, ordered by the position of
each record's first dependency in the target's recorded dependency
order. Within a record, `make` places no guarantee on the relative order
of its fields.

#### Scenario: Record-only dependency merges a whole document
- **WHEN** a list-mode target has a record-only dependency for record `bldoc` sourced from a document with top-level keys `path` and `remote`, and `make` is run
- **THEN** the compiled array contains an object for record `bldoc` with `path` and `remote` set to those keys' values

#### Scenario: Record-field dependency resolves like a field-mode field
- **WHEN** a list-mode target has a record-field dependency naming field `description` of record `bldoc`, sourced from a whole file with no field-path, and `make` is run
- **THEN** the compiled array's `bldoc` object has a `description` field equal to that file's content

#### Scenario: Record-field dependency applies its format template
- **WHEN** a list-mode target has a record-field dependency naming field `version` of record `bldoc`, sourced from a field-path resolving to `1.2`, with format template `"v%s"`, and `make` is run
- **THEN** the compiled array's `bldoc` object has a `version` field equal to `"v1.2"`

#### Scenario: Multiple records ordered by first appearance
- **WHEN** a list-mode target's first recorded dependency belongs to record `bldoc` and a later dependency is the first one belonging to record `projx`, and `make` is run
- **THEN** the compiled array's first entry is the `bldoc` record and its second entry is the `projx` record

### Requirement: List-mode field values are plain, not wrapped
Unlike field-mode's `{"value": ..., "raw": ...}` wrapper, `make` SHALL
emit each list-mode record field as its resolved value directly (the
format-rendered result if a format template was given, otherwise the raw
resolved value) — with no accompanying `raw` counterpart.

#### Scenario: List-mode field has no raw wrapper
- **WHEN** a list-mode target's record-field dependency has a format template applied, and `make` is run
- **THEN** the compiled array's corresponding field is the rendered string itself, not an object containing separate `value` and `raw` keys
