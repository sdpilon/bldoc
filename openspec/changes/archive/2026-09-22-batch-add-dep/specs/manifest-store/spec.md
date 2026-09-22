## MODIFIED Requirements

### Requirement: `add-dep` warns on an ambiguous plain anchor
Given a source-ref with a single-segment `#anchor` (no breadcrumb) naming
a readable Markdown source file, `add-dep` SHALL read that file and
check whether the anchor's slug matches more than one heading. If it
does, `add-dep` SHALL print a non-blocking warning suggesting a
breadcrumb path, and SHALL still record the dependency. If the source
file cannot be read (missing, unreadable, or any other I/O error),
`add-dep` SHALL skip this check silently and still record the
dependency, consistent with `add-dep` otherwise never requiring a
source file to exist.

#### Scenario: Ambiguous anchor warned, not blocked
- **WHEN** the user runs `bldoc add-dep README@x spec.md#scenario-name` and `spec.md` has more than one heading whose slug is `scenario-name`
- **THEN** the CLI prints a warning suggesting a breadcrumb path, exits 0, and records the dependency with anchor `scenario-name`

#### Scenario: Unambiguous anchor produces no warning
- **WHEN** the user runs `bldoc add-dep README@x spec.md#purpose` and `spec.md` has exactly one heading whose slug is `purpose`
- **THEN** the CLI exits 0 with no warning, and records the dependency with anchor `purpose`

#### Scenario: Breadcrumb anchor is never checked for ambiguity
- **WHEN** the user runs `bldoc add-dep README@x spec.md#requirement-a/scenario-name`
- **THEN** the CLI does not attempt an ambiguity check (a breadcrumb is already a specific address) and records the dependency

#### Scenario: Ambiguity check skipped when the source file is unreadable
- **WHEN** the user runs `bldoc add-dep README@x missing.md#purpose` and `missing.md` does not exist
- **THEN** the CLI exits 0 with no warning and no error, and records the dependency with anchor `purpose`

### Requirement: `add-dep` records a dependency
Given a target-ref and source-ref that both parse and refer to an
existing target, `add-dep` SHALL append the dependency (source ref,
optional field path or anchor, optional `--nested` flag, optional format
template) to that target's recorded dependencies, after any dependencies
already recorded.

#### Scenario: Whole-file dependency recorded
- **WHEN** `bldoc add-dep README pyproject.toml` is run against an existing `README` target
- **THEN** the manifest records `pyproject.toml` as a whole-file dependency of `README`, appended after any dependencies already recorded

#### Scenario: Field-addressed dependency recorded with format
- **WHEN** `bldoc add-dep README@version --format "Python version must be %s to run this project." pyproject.toml@project.requires-python` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `version` of `README`, sourced from `project.requires-python` in `pyproject.toml`, with the given format template

#### Scenario: Anchor-addressed dependency recorded
- **WHEN** `bldoc add-dep README@summary spec.md#purpose` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `summary` of `README`, sourced from anchor `purpose` in `spec.md`, with `--nested` unset

#### Scenario: Anchor-addressed dependency recorded with `--nested`
- **WHEN** `bldoc add-dep README@section spec.md#requirements --nested` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `section` of `README`, sourced from anchor `requirements` in `spec.md`, with `--nested` set

### Requirement: Raw/field mode exclusivity enforced
For a target with an explicitly declared mode, `add-dep` SHALL reject
any dependency — including the first — whose field-addressing (bare vs.
`@field`) does not match that declared mode. For a target with no
declared mode, `add-dep` SHALL reject a target-ref whose field-addressing
does not match every dependency the target already has recorded, when
the target has at least one recorded dependency.

#### Scenario: Field-addressed dep rejected on a raw-mode target
- **WHEN** `bldoc add-dep README@version --format "..." pyproject.toml@project.version` is run and `README` already has a whole-file (unnamed) dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Whole-file dep rejected on a field-mode target
- **WHEN** `bldoc add-dep README other.toml` is run and `README` already has a field-addressed dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Field-addressed dep rejected as the first dependency of a declared raw-mode target
- **WHEN** `bldoc add-dep PROJECTS@name entry.yaml@name` is run and `PROJECTS` was created with `bldoc new PROJECTS --mode raw` and has no dependencies recorded yet
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Whole-file dep rejected as the first dependency of a declared field-mode target
- **WHEN** `bldoc add-dep PROJECTS entry.yaml` is run and `PROJECTS` was created with `bldoc new PROJECTS --mode field` and has no dependencies recorded yet
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `add-dep` requires record-scoped addressing on a list-mode target
For a target whose declared mode is `list`, `add-dep` SHALL reject a target-ref with no field component (a bare target name), and SHALL reject a field component containing more than one `.`. A field component with no `.` (e.g. `bldoc`) names a record-only dependency; a field component with exactly one `.` (e.g. `bldoc.description`) names a single field within that record.

#### Scenario: Bare target-ref rejected on a list-mode target
- **WHEN** `bldoc add-dep PROJECTS entry.yaml` is run and `PROJECTS` was created with `bldoc new PROJECTS --mode list`
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Record-only dependency accepted
- **WHEN** `bldoc add-dep PROJECTS@bldoc project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a record-only dependency on record `bldoc`, sourced from `project-headers/bldoc.yaml`

#### Scenario: Record-field dependency accepted
- **WHEN** `bldoc add-dep PROJECTS@bldoc.description entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a dependency on field `description` of record `bldoc`, sourced from `entry-description.md`

#### Scenario: Over-qualified field component rejected
- **WHEN** `bldoc add-dep PROJECTS@bldoc.description.extra entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `--format` requires a record-field-scoped target-ref on a list-mode target
For a target whose declared mode is `list`, `add-dep` SHALL reject `--format` when the target-ref's field component has no `.` (i.e. names a record-only dependency rather than a single record field).

#### Scenario: Format rejected on a record-only dependency
- **WHEN** `bldoc add-dep PROJECTS@bldoc --format "..." project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Format accepted on a record-field dependency
- **WHEN** `bldoc add-dep PROJECTS@bldoc.description --format "%s" entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records the dependency on field `description` of record `bldoc` with the given format template

### Requirement: `add-dep` rejects a duplicate dependency
`add-dep` SHALL reject a target-ref/source-ref pair already recorded on
the target. The anchor (including any breadcrumb) is part of the
source-ref for this comparison; `--nested` and `--format` are not — a
dependency already recorded with the same target-ref and source-ref is a
duplicate regardless of whether `--nested` or `--format` differ from the
new invocation.

#### Scenario: Duplicate dependency rejected
- **WHEN** `bldoc add-dep README pyproject.toml` is run and `README` already has `pyproject.toml` recorded as a whole-file dependency
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Duplicate anchor dependency rejected regardless of `--nested`
- **WHEN** `bldoc add-dep README@section spec.md#requirements --nested` is run and `README` already has field `section` recorded with source-ref `spec.md#requirements` (without `--nested`)
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

## ADDED Requirements

### Requirement: Batch `add-dep` validates every pair before recording any
Given a variadic `add-dep` invocation recording two or more pairs
against a single target (and, for a list-mode target, a single record),
the manifest store SHALL validate every pair — parsing, target
existence, mode exclusivity, list-mode record/field addressing, and
duplicate-dependency rejection, exactly as it validates a single
`add-dep` call — before writing any of them to `bldoc.toml`. If any
pair fails validation, the CLI SHALL report an error and exit non-zero,
and the manifest SHALL be left exactly as it was before the invocation,
with none of the invocation's pairs recorded, regardless of whether
earlier pairs in the same call would have validated successfully on
their own.

#### Scenario: All pairs recorded when every pair is valid
- **WHEN** `bldoc add-dep README pyproject.toml go.mod` is run against an existing `raw`-mode `README` target with neither dependency already recorded
- **THEN** the manifest records both `pyproject.toml` and `go.mod` as whole-file dependencies of `README`, in that order

#### Scenario: No pairs recorded when a later pair is invalid
- **WHEN** `bldoc add-dep README pyproject.toml go.mod` is run against an existing `raw`-mode `README` target that already has `go.mod` recorded as a dependency
- **THEN** the CLI reports a duplicate-dependency error and exits non-zero, and the manifest does not record `pyproject.toml` either, even though it alone would have been valid

#### Scenario: No pairs recorded when the target's mode rejects one pair
- **WHEN** a batch `add-dep` call against a `field`-mode target includes one field-addressed pair and one bare (whole-file) pair
- **THEN** the CLI reports a mode-exclusivity error and exits non-zero, and the manifest records neither pair

#### Scenario: No pairs recorded when a duplicate pair appears within the same call
- **WHEN** a batch `add-dep` call lists the same source-ref twice for the same target-ref
- **THEN** the CLI reports a duplicate-dependency error and exits non-zero, and the manifest records neither occurrence

### Requirement: Batch `add-dep` records a list-mode record across its pairs
Given a variadic `add-dep` invocation against a `list`-mode target with
an explicit record name, the manifest store SHALL record every pair
scoped to that one record: a bare-source-ref pair as a record-only
dependency on that record, and a `field:source-ref` pair as a
dependency on that field of that record — using the same record/field
addressing and `--format` restriction (rejecting `--format` on a
record-only pair) as the existing single-dependency `list`-mode
requirements.

#### Scenario: Record and field dependencies recorded together
- **WHEN** `bldoc add-dep PROJECTS bldoc description:entry-description.md path:project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a dependency on field `description` of record `bldoc` sourced from `entry-description.md`, and a dependency on field `path` of record `bldoc` sourced from `project-headers/bldoc.yaml`

#### Scenario: `--format` rejected on a record-only pair within a batch call
- **WHEN** a batch `add-dep` call against a `list`-mode target includes a bare (record-only) pair together with `--format`
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged
