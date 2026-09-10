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
`add-dep` SHALL reject a target-ref whose field-addressing (bare vs.
`:field`) does not match every dependency the target already has
recorded, when the target has at least one recorded dependency.

#### Scenario: Field-addressed dep rejected on a raw-mode target
- **WHEN** `bldoc add-dep README:version --format "..." pyproject.toml:project.version` is run and `README` already has a whole-file (unnamed) dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Whole-file dep rejected on a field-mode target
- **WHEN** `bldoc add-dep README other.toml` is run and `README` already has a field-addressed dependency recorded
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

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
