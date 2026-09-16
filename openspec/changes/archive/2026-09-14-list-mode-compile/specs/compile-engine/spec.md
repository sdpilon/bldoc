## MODIFIED Requirements

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

### Requirement: Structured field-path resolution supports TOML, JSON, and YAML
`make` SHALL parse a dependency's source file as TOML, JSON, or YAML
based on its file extension (`.toml`, `.json`, `.yaml`/`.yml`
respectively), both when resolving a field-path on a field-mode or
list-mode record-field dependency, and when resolving the whole
top-level document of a list-mode record-only dependency.

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

## ADDED Requirements

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
