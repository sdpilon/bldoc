## MODIFIED Requirements

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
- **WHEN** `bldoc add-dep README:version --format "Python version must be %s to run this project." pyproject.toml:project.requires-python` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `version` of `README`, sourced from `project.requires-python` in `pyproject.toml`, with the given format template

#### Scenario: Anchor-addressed dependency recorded
- **WHEN** `bldoc add-dep README:summary spec.md#purpose` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `summary` of `README`, sourced from anchor `purpose` in `spec.md`, with `--nested` unset

#### Scenario: Anchor-addressed dependency recorded with `--nested`
- **WHEN** `bldoc add-dep README:section spec.md#requirements --nested` is run against an existing `README` target
- **THEN** the manifest records a dependency on field `section` of `README`, sourced from anchor `requirements` in `spec.md`, with `--nested` set

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
- **WHEN** `bldoc add-dep README:section spec.md#requirements --nested` is run and `README` already has field `section` recorded with source-ref `spec.md#requirements` (without `--nested`)
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `show` reports a target's recorded dependencies
Given a target name recorded in the manifest, `show` SHALL print that
target's recorded dependencies (source ref, field path or anchor if any,
whether `--nested` was recorded, format template if any) in the order
they were added.

#### Scenario: Dependencies shown
- **WHEN** `bldoc show README` is run and `README` has two dependencies recorded, added in a specific order
- **THEN** the CLI prints both dependencies in the order they were added, including each one's field path and format template if present

#### Scenario: Anchor dependency shown with nested flag
- **WHEN** `bldoc show README` is run and `README` has a dependency recorded with anchor `requirements` in `spec.md` and `--nested` set
- **THEN** the CLI's output for that dependency includes the anchor `requirements` and indicates `--nested` was recorded

## ADDED Requirements

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
- **WHEN** the user runs `bldoc add-dep README:x spec.md#scenario-name` and `spec.md` has more than one heading whose slug is `scenario-name`
- **THEN** the CLI prints a warning suggesting a breadcrumb path, exits 0, and records the dependency with anchor `scenario-name`

#### Scenario: Unambiguous anchor produces no warning
- **WHEN** the user runs `bldoc add-dep README:x spec.md#purpose` and `spec.md` has exactly one heading whose slug is `purpose`
- **THEN** the CLI exits 0 with no warning, and records the dependency with anchor `purpose`

#### Scenario: Breadcrumb anchor is never checked for ambiguity
- **WHEN** the user runs `bldoc add-dep README:x spec.md#requirement-a/scenario-name`
- **THEN** the CLI does not attempt an ambiguity check (a breadcrumb is already a specific address) and records the dependency

#### Scenario: Ambiguity check skipped when the source file is unreadable
- **WHEN** the user runs `bldoc add-dep README:x missing.md#purpose` and `missing.md` does not exist
- **THEN** the CLI exits 0 with no warning and no error, and records the dependency with anchor `purpose`
