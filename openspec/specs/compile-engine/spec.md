# compile-engine Specification

## Purpose

The `compile-engine` capability defines `make`'s real behavior: turning
a target's recorded dependencies (from the `manifest-store` capability)
into a disposable, mechanically-derived intermediate — raw-mode
concatenation or field-mode JSON, always fully recomputed.

## Requirements

### Requirement: Intermediate file location and naming
The compiled intermediate for a target SHALL be written under `.bldoc/`
in the current working directory: a raw-mode target's intermediate at
`.bldoc/<target>` (raw bytes), and a field-mode target's intermediate at
`.bldoc/<target>.json` (a single JSON object).

#### Scenario: Raw-mode output path
- **WHEN** `bldoc make README` compiles a raw-mode target
- **THEN** the CLI writes the concatenated bytes to `.bldoc/README`

#### Scenario: Field-mode output path
- **WHEN** `bldoc make README` compiles a field-mode target
- **THEN** the CLI writes the JSON object to `.bldoc/README.json`

### Requirement: `make` requires an existing target
`make <target>` SHALL reject a target name not recorded in the manifest.

#### Scenario: Unknown target rejected
- **WHEN** `bldoc make MISSING` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero

### Requirement: Raw-mode compilation concatenates dependencies
For a target whose recorded dependencies are all whole-file (unnamed),
`make` SHALL write the byte-for-byte concatenation of each dependency's
source file, in the order the dependencies were added.

#### Scenario: Raw-mode concatenation
- **WHEN** `README` has whole-file dependencies `a.txt` then `b.txt` recorded, and `make README` is run
- **THEN** `.bldoc/README` contains the bytes of `a.txt` immediately followed by the bytes of `b.txt`, identical to `cat a.txt b.txt`

### Requirement: A target with no dependencies compiles to an empty intermediate
`make` SHALL compile a target with no recorded dependencies to an empty
raw-mode intermediate.

#### Scenario: Empty target
- **WHEN** `README` has no recorded dependencies and `make README` is run
- **THEN** `.bldoc/README` is created empty (zero bytes)

### Requirement: Field-mode compilation resolves each dependency's value
For a target whose recorded dependencies are all field-addressed, `make`
SHALL resolve each dependency's raw value — the source file's entire
content if the dependency has no field-path, or the value at the
field-path if it has one — apply the dependency's format template if
present, and emit one JSON object with one entry per field:
`{"value": <rendered>, "raw": <raw value>}`. `value` equals `raw` when
no format template was recorded.

#### Scenario: Field resolved from a whole source file
- **WHEN** `README:summary` depends on `notes.txt` with no field-path and no format template, and `make README` is run
- **THEN** `.bldoc/README.json` contains `{"summary": {"value": "<notes.txt's content>", "raw": "<notes.txt's content>"}}`

#### Scenario: Field resolved from a structured field-path with a format template
- **WHEN** `README:version` depends on `pyproject.toml:project.requires-python` with format `"Python version must be %s to run this project."`, and that field-path resolves to `3.11`, and `make README` is run
- **THEN** `.bldoc/README.json` contains `{"version": {"value": "Python version must be 3.11 to run this project.", "raw": "3.11"}}`

### Requirement: Structured field-path resolution supports TOML, JSON, and YAML
`make` SHALL parse a dependency's source file as TOML, JSON, or YAML
based on its file extension (`.toml`, `.json`, `.yaml`/`.yml`
respectively) when the dependency has a field-path, and resolve the
field-path as a dot-separated sequence of keys into the parsed document.

#### Scenario: Unsupported source format rejected
- **WHEN** a dependency has a field-path into a source file whose extension is none of `.toml`, `.json`, `.yaml`, or `.yml`
- **THEN** `make` reports an error and exits non-zero

### Requirement: A malformed structured source file is rejected
`make` SHALL report an error and exit non-zero if a dependency's
field-path requires parsing its source file as TOML/JSON/YAML and that
file cannot be parsed.

#### Scenario: Malformed source file rejected
- **WHEN** a dependency's field-path source file has the appropriate extension but is not valid TOML/JSON/YAML
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
`make` SHALL reject a target whose recorded dependencies use the same
target-side field name more than once.

#### Scenario: Duplicate field name rejected
- **WHEN** a target has two recorded dependencies both naming the field `version`
- **THEN** `make` reports an error and exits non-zero

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
