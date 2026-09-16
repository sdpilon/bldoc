## MODIFIED Requirements

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

## ADDED Requirements

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
