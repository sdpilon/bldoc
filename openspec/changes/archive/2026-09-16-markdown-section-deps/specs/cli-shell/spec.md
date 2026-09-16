## MODIFIED Requirements

### Requirement: `add-dep` parses target-ref, source-ref, and optional format
`add-dep` SHALL require a target-ref, a source-ref, and accept an
optional `--format` flag. A target-ref SHALL be either a bare target name
or `target:field`. A source-ref SHALL be a bare file path, `path:field-path`,
or `path#anchor` — the anchor either a single GitHub-style slug or a
breadcrumb path of slugs (`parent-slug/child-slug`, root-first). A
source-ref SHALL NOT combine `:field-path` and `#anchor`.

#### Scenario: Whole-file dependency
- **WHEN** the user runs `bldoc add-dep README pyproject.toml`
- **THEN** the CLI parses `README` as the target with no field, and `pyproject.toml` as the source with no field-path

#### Scenario: Field-addressed dependency with format
- **WHEN** the user runs `bldoc add-dep README:version --format "Python version must be %s to run this project." pyproject.toml:project.requires-python`
- **THEN** the CLI parses `README` as the target with field `version`, `pyproject.toml` as the source with field-path `project.requires-python`, and captures the format string

#### Scenario: Malformed target-ref is rejected
- **WHEN** the user runs `add-dep` with a target-ref containing more than one `:`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Anchor-addressed dependency
- **WHEN** the user runs `bldoc add-dep README:summary spec.md#purpose`
- **THEN** the CLI parses `README` as the target with field `summary`, and `spec.md` as the source with anchor `purpose`

#### Scenario: Anchor-addressed dependency with a breadcrumb path
- **WHEN** the user runs `bldoc add-dep README:scenario spec.md#requirement-x/scenario-y`
- **THEN** the CLI parses `spec.md` as the source with anchor breadcrumb `requirement-x/scenario-y`

#### Scenario: Combining a field-path and an anchor is rejected
- **WHEN** the user runs `bldoc add-dep README:summary spec.md:some.path#purpose`
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `rm-dep` parses target-ref and source-ref
`rm-dep` SHALL require a target-ref and a source-ref, using the same
grammar as `add-dep`, and SHALL NOT accept a `--format` or `--nested`
flag.

#### Scenario: Removing a whole-file dependency
- **WHEN** the user runs `bldoc rm-dep README pyproject.toml`
- **THEN** the CLI parses the target-ref and source-ref and removes the matching dependency from the manifest (see the `manifest-store` capability)

#### Scenario: `--format` is not accepted
- **WHEN** the user runs `bldoc rm-dep` with a `--format` flag
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: `--nested` is not accepted
- **WHEN** the user runs `bldoc rm-dep` with a `--nested` flag
- **THEN** the CLI reports a usage error and exits non-zero

## ADDED Requirements

### Requirement: `add-dep` parses an optional `--nested` flag
`add-dep` SHALL accept an optional `--nested` flag. `add-dep` SHALL
reject `--nested` when the source-ref has no `#anchor`, and SHALL reject
`--nested` when the target-ref has no `:field` (i.e. on a whole-file/
raw-mode dependency).

#### Scenario: Nested flag accepted with an anchor-addressed field dependency
- **WHEN** the user runs `bldoc add-dep README:section spec.md#requirements --nested`
- **THEN** the CLI parses `README` as the target with field `section`, `spec.md#requirements` as the anchor-addressed source, and records that `--nested` was given

#### Scenario: Nested flag rejected without an anchor
- **WHEN** the user runs `bldoc add-dep README:section spec.md --nested` with no `#anchor` on the source-ref
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Nested flag rejected on a whole-file dependency
- **WHEN** the user runs `bldoc add-dep README spec.md#requirements --nested` with no `:field` on the target-ref
- **THEN** the CLI reports a usage error and exits non-zero
