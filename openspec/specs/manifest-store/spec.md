# manifest-store Specification

## Purpose

The `manifest-store` capability defines the manifest file that records
every target and its dependencies — the durable, hand-authored recipe
that the CLI's `new`, `add-dep`, `rm-dep`, `rm`, `list`, and `show`
commands read and write.

## Requirements

### Requirement: Manifest file location
The manifest SHALL be a single file named `bldoc.toml` in the current
working directory.

#### Scenario: Manifest path
- **WHEN** any CLI command reads or writes manifest state
- **THEN** it does so against `bldoc.toml` in the current working directory

### Requirement: `new` creates the manifest if absent
`new` SHALL create `bldoc.toml` if it does not already exist, then
record a new target in it.

#### Scenario: First target creates the manifest
- **WHEN** `bldoc new README` is run in a directory with no `bldoc.toml`
- **THEN** the CLI creates `bldoc.toml` containing a `README` target with no dependencies

### Requirement: `new` rejects a duplicate target name
`new` SHALL reject a target name that already exists in the manifest.

#### Scenario: Duplicate target rejected
- **WHEN** `bldoc new README` is run and the manifest already has a `README` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `add-dep` requires an existing target
`add-dep` SHALL reject a target-ref naming a target not already recorded
in the manifest.

#### Scenario: Unknown target rejected
- **WHEN** `bldoc add-dep MISSING pyproject.toml` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `add-dep` records a dependency
Given a target-ref and source-ref that both parse and refer to an
existing target, `add-dep` SHALL append the dependency (source ref,
optional field path, optional format template) to that target's
recorded dependencies, after any dependencies already recorded.

#### Scenario: Whole-file dependency recorded
- **WHEN** `bldoc add-dep README pyproject.toml` is run against an existing `README` target
- **THEN** the manifest records `pyproject.toml` as a whole-file dependency of `README`, appended after any dependencies already recorded

#### Scenario: Field-addressed dependency recorded with format
- **WHEN** `bldoc add-dep README:version --format "Python version must be %s to run this project." pyproject.toml:project.requires-python` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `version` of `README`, sourced from `project.requires-python` in `pyproject.toml`, with the given format template

### Requirement: Raw/field mode exclusivity enforced
For a target with an explicitly declared mode, `add-dep` SHALL reject
any dependency — including the first — whose field-addressing (bare vs.
`:field`) does not match that declared mode. For a target with no
declared mode, `add-dep` SHALL reject a target-ref whose field-addressing
does not match every dependency the target already has recorded, when
the target has at least one recorded dependency.

#### Scenario: Field-addressed dep rejected on a raw-mode target
- **WHEN** `bldoc add-dep README:version --format "..." pyproject.toml:project.version` is run and `README` already has a whole-file (unnamed) dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Whole-file dep rejected on a field-mode target
- **WHEN** `bldoc add-dep README other.toml` is run and `README` already has a field-addressed dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Field-addressed dep rejected as the first dependency of a declared raw-mode target
- **WHEN** `bldoc add-dep PROJECTS:name entry.yaml:name` is run and `PROJECTS` was created with `bldoc new PROJECTS --mode raw` and has no dependencies recorded yet
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
- **WHEN** `bldoc add-dep PROJECTS:bldoc project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a record-only dependency on record `bldoc`, sourced from `project-headers/bldoc.yaml`

#### Scenario: Record-field dependency accepted
- **WHEN** `bldoc add-dep PROJECTS:bldoc.description entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a dependency on field `description` of record `bldoc`, sourced from `entry-description.md`

#### Scenario: Over-qualified field component rejected
- **WHEN** `bldoc add-dep PROJECTS:bldoc.description.extra entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `--format` requires a record-field-scoped target-ref on a list-mode target
For a target whose declared mode is `list`, `add-dep` SHALL reject `--format` when the target-ref's field component has no `.` (i.e. names a record-only dependency rather than a single record field).

#### Scenario: Format rejected on a record-only dependency
- **WHEN** `bldoc add-dep PROJECTS:bldoc --format "..." project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Format accepted on a record-field dependency
- **WHEN** `bldoc add-dep PROJECTS:bldoc.description --format "%s" entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records the dependency on field `description` of record `bldoc` with the given format template

### Requirement: `add-dep` rejects a duplicate dependency
`add-dep` SHALL reject a target-ref/source-ref pair already recorded on
the target.

#### Scenario: Duplicate dependency rejected
- **WHEN** `bldoc add-dep README pyproject.toml` is run and `README` already has `pyproject.toml` recorded as a whole-file dependency
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `rm-dep` requires an existing target and recorded dependency
`rm-dep` SHALL reject a target-ref naming a target not in the manifest,
and SHALL reject a source-ref not recorded as a dependency of that
target.

#### Scenario: Unknown target rejected
- **WHEN** `bldoc rm-dep MISSING pyproject.toml` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero

#### Scenario: Unrecorded dependency rejected
- **WHEN** `bldoc rm-dep README other.toml` is run and `README` has no `other.toml` dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `rm-dep` removes a recorded dependency
Given a target-ref and source-ref that match a recorded dependency,
`rm-dep` SHALL remove that dependency from the manifest, preserving the
order of the target's remaining dependencies.

#### Scenario: Dependency removed
- **WHEN** `bldoc rm-dep README pyproject.toml` is run and `README` has `pyproject.toml` recorded as a whole-file dependency
- **THEN** the manifest no longer records `pyproject.toml` as a dependency of `README`, and `README`'s other recorded dependencies are unchanged and keep their relative order

### Requirement: `rm` requires an existing target
`rm` SHALL reject a target name not recorded in the manifest.

#### Scenario: Unknown target rejected
- **WHEN** `bldoc rm MISSING` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero

### Requirement: `rm` deletes a target and its dependencies
Given a target name recorded in the manifest, `rm` SHALL remove that
target and all of its recorded dependencies from the manifest.

#### Scenario: Target removed
- **WHEN** `bldoc rm README` is run and a `README` target exists
- **THEN** the manifest no longer contains a `README` target or any of its dependencies

### Requirement: `list` reports every target
`list` SHALL print every target name recorded in the manifest, one per
line, in the order the targets were declared (i.e. the order `new` was
called for each). Given a manifest with no targets, or no manifest file
at all, it SHALL print no target names.

#### Scenario: Targets listed
- **WHEN** `bldoc list` is run against a manifest where `README` was declared before `CHANGELOG`
- **THEN** the CLI prints `README` followed by `CHANGELOG`

#### Scenario: No targets
- **WHEN** `bldoc list` is run and no manifest file exists, or the manifest has no targets
- **THEN** the CLI exits 0 and prints no target names

### Requirement: `show` requires an existing target
`show` SHALL reject a target name not recorded in the manifest.

#### Scenario: Unknown target rejected
- **WHEN** `bldoc show MISSING` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero

### Requirement: `show` reports a target's recorded dependencies
Given a target name recorded in the manifest, `show` SHALL print that
target's recorded dependencies (source ref, field path if any, format
template if any) in the order they were added.

#### Scenario: Dependencies shown
- **WHEN** `bldoc show README` is run and `README` has two dependencies recorded, added in a specific order
- **THEN** the CLI prints both dependencies in the order they were added, including each one's field path and format template if present

### Requirement: Malformed manifest file is rejected
If `bldoc.toml` exists but cannot be parsed, any command that needs to
read it SHALL report an error and exit non-zero rather than proceeding
with partial or default state.

#### Scenario: Unparseable manifest
- **WHEN** any manifest-reading command is run and `bldoc.toml` exists but is not valid TOML
- **THEN** the CLI reports an error and exits non-zero, and does not modify the file

### Requirement: `new` requires an explicit mode
`new` SHALL require a `--mode` flag whose value is exactly `raw`, `field`, or `list`, and SHALL reject the invocation when it is omitted or set to any other value, leaving the manifest unchanged.

#### Scenario: Mode omitted rejected
- **WHEN** `bldoc new README` is run with no `--mode` flag
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Invalid mode value rejected
- **WHEN** `bldoc new README --mode nonsense` is run
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Valid mode recorded
- **WHEN** `bldoc new PROJECTS --mode field` is run
- **THEN** the manifest records target `PROJECTS` with mode `field` and no dependencies

#### Scenario: Valid list mode recorded
- **WHEN** `bldoc new PROJECTS --mode list` is run
- **THEN** the manifest records target `PROJECTS` with mode `list` and no dependencies

### Requirement: `new` accepts an optional output extension
`new` SHALL accept an optional `--ext <ext>` flag regardless of the declared mode, and record it on the created target.

#### Scenario: Extension recorded
- **WHEN** `bldoc new PROJECTS --mode raw --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `raw` and extension `yaml`

#### Scenario: Extension recorded with field mode
- **WHEN** `bldoc new PROJECTS --mode field --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `field` and extension `yaml`

#### Scenario: Extension recorded with list mode
- **WHEN** `bldoc new PROJECTS --mode list --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `list` and extension `yaml`

#### Scenario: No extension by default
- **WHEN** `bldoc new README --mode raw` is run with no `--ext`
- **THEN** the manifest records target `README` with mode `raw` and no extension set

### Requirement: `rename` requires an existing old target and rejects a duplicate new name
`rename` SHALL reject an old-name not recorded in the manifest, and SHALL reject a new-name that already names an existing target.

#### Scenario: Unknown old target rejected
- **WHEN** `bldoc rename MISSING NEW` is run and no `MISSING` target exists
- **THEN** the CLI reports an error and exits non-zero

#### Scenario: Duplicate new name rejected
- **WHEN** `bldoc rename README NOTES` is run and a `NOTES` target already exists
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `rename` renames a target in place
Given a valid old name and an available new name, `rename` SHALL rename the target, preserving its recorded dependencies, their order, its declared mode (if any), and its extension (if any) unchanged.

#### Scenario: Target renamed
- **WHEN** `bldoc rename README NOTES` is run, `README` exists, and no `NOTES` target exists
- **THEN** the manifest no longer contains a `README` target, contains a `NOTES` target in its place, and `NOTES`'s dependencies (with their order), mode, and extension are identical to `README`'s before the rename

### Requirement: `show` reports a target's extension
Given a target name recorded in the manifest, `show` SHALL report the target's extension if one is set.

#### Scenario: Extension shown
- **WHEN** `bldoc show PROJECTS` is run and `PROJECTS` has extension `yaml` recorded
- **THEN** the CLI's output includes `PROJECTS`'s extension `yaml`

#### Scenario: No extension to show
- **WHEN** `bldoc show README` is run and `README` has no extension recorded
- **THEN** the CLI's output does not report an extension for `README`

### Requirement: `show` reports a target's mode
Given a target name recorded in the manifest, `show` SHALL report the target's mode: its declared mode if one was recorded at creation, otherwise the mode computed from its recorded dependencies.

#### Scenario: Declared mode shown
- **WHEN** `bldoc show PROJECTS` is run and `PROJECTS` was created with `--mode field`
- **THEN** the CLI's output reports `PROJECTS`'s mode as `field`

#### Scenario: Computed mode shown for a target with no declared mode
- **WHEN** `bldoc show README` is run and `README` has no declared mode but has a whole-file dependency recorded
- **THEN** the CLI's output reports `README`'s mode as `raw`

### Requirement: A target with no declared mode continues to use computed mode
A target with no `mode` recorded — including one written to `bldoc.toml` before this capability existed — SHALL continue to have its mode inferred from its recorded dependencies' field-addressing, governed by `add-dep`'s existing exclusivity check for targets with no declared mode.

#### Scenario: Pre-existing target unaffected
- **WHEN** a `bldoc.toml` contains a target with no `mode` key and one whole-file dependency already recorded, and a second whole-file `add-dep` is run against it
- **THEN** the CLI accepts it exactly as it would have before this capability existed
