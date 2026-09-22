## MODIFIED Requirements

### Requirement: `add-dep` parses target-ref, source-ref, and optional format
`add-dep` SHALL require a target-ref, a source-ref, and accept an
optional `--format` flag. A target-ref SHALL be either a bare target
name or `target@field`. A source-ref SHALL be a bare file path,
`path@field-path`, or `path#anchor` — the anchor either a single
GitHub-style slug or a breadcrumb path of slugs (`parent-slug/child-slug`,
root-first). A source-ref SHALL NOT combine `@field-path` and `#anchor`.

#### Scenario: Whole-file dependency
- **WHEN** the user runs `bldoc add-dep README pyproject.toml`
- **THEN** the CLI parses `README` as the target with no field, and `pyproject.toml` as the source with no field-path

#### Scenario: Field-addressed dependency with format
- **WHEN** the user runs `bldoc add-dep README@version --format "Python version must be %s to run this project." pyproject.toml@project.requires-python`
- **THEN** the CLI parses `README` as the target with field `version`, `pyproject.toml` as the source with field-path `project.requires-python`, and captures the format string

#### Scenario: Malformed target-ref is rejected
- **WHEN** the user runs `add-dep` with a target-ref containing more than one `@`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Anchor-addressed dependency
- **WHEN** the user runs `bldoc add-dep README@summary spec.md#purpose`
- **THEN** the CLI parses `README` as the target with field `summary`, and `spec.md` as the source with anchor `purpose`

#### Scenario: Anchor-addressed dependency with a breadcrumb path
- **WHEN** the user runs `bldoc add-dep README@scenario spec.md#requirement-x/scenario-y`
- **THEN** the CLI parses `spec.md` as the source with anchor breadcrumb `requirement-x/scenario-y`

#### Scenario: Combining a field-path and an anchor is rejected
- **WHEN** the user runs `bldoc add-dep README@summary spec.md@some.path#purpose`
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `add-dep` parses an optional `--nested` flag
`add-dep` SHALL accept an optional `--nested` flag. `add-dep` SHALL
reject `--nested` when the source-ref has no `#anchor`, and SHALL reject
`--nested` when the target-ref has no `@field` (i.e. on a whole-file/
raw-mode dependency).

#### Scenario: Nested flag accepted with an anchor-addressed field dependency
- **WHEN** the user runs `bldoc add-dep README@section spec.md#requirements --nested`
- **THEN** the CLI parses `README` as the target with field `section`, `spec.md#requirements` as the anchor-addressed source, and records that `--nested` was given

#### Scenario: Nested flag rejected without an anchor
- **WHEN** the user runs `bldoc add-dep README@section spec.md --nested` with no `#anchor` on the source-ref
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Nested flag rejected on a whole-file dependency
- **WHEN** the user runs `bldoc add-dep README spec.md#requirements --nested` with no `@field` on the target-ref
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `--format` requires a field-addressed target-ref
`add-dep` SHALL reject `--format` when the target-ref has no `@field`.

#### Scenario: Format on a whole-file dependency is rejected
- **WHEN** the user runs `bldoc add-dep README --format "..." pyproject.toml` with no `@field` on the target-ref
- **THEN** the CLI reports a usage error and exits non-zero

## ADDED Requirements

### Requirement: `add-dep` accepts a variadic multi-pair form
`add-dep` SHALL accept, in addition to its existing two-positional-
argument form, a variadic form invoked with three or more positional
arguments: `add-dep <target> [<record-name>] <pair> <pair> [<pair>...]`.
In this form the first argument SHALL be a bare target name (no `@`) —
the CLI SHALL reject a first argument containing `@` in a three-or-more-
argument invocation. When the named target's declared mode is `list`,
the second argument SHALL be a bare record name (no `@`, `:`, or `.`),
and every remaining argument is a pair; for any other declared mode (or
an undeclared mode), every argument after the target is a pair with no
record-name argument. The form SHALL require at least two pairs; the
CLI SHALL reject an invocation that, after removing the target and any
record-name argument, leaves fewer than two pairs.

#### Scenario: Two raw-mode pairs recorded in one call
- **WHEN** the user runs `bldoc add-dep README pyproject.toml go.mod` against a `raw`-mode `README` target
- **THEN** the CLI parses `README` as a bare target and records two whole-file dependencies, on `pyproject.toml` and `go.mod`, in that order

#### Scenario: Two field-mode pairs recorded in one call
- **WHEN** the user runs `bldoc add-dep README version:pyproject.toml@project.version summary:spec.md#purpose` against a `field`-mode `README` target
- **THEN** the CLI parses `README` as a bare target and records a dependency on field `version` sourced from `pyproject.toml@project.version`, and a dependency on field `summary` sourced from `spec.md#purpose`

#### Scenario: List-mode record recorded with two field pairs
- **WHEN** the user runs `bldoc add-dep PROJECTS bldoc description:entry-description.md path:project-headers/bldoc.yaml` against a `list`-mode `PROJECTS` target
- **THEN** the CLI parses `PROJECTS` as the target, `bldoc` as the record name, and records two dependencies scoped to record `bldoc`: field `description` sourced from `entry-description.md`, and field `path` sourced from `project-headers/bldoc.yaml`

#### Scenario: List-mode record-only pair alongside a field pair
- **WHEN** the user runs `bldoc add-dep PROJECTS bldoc defaults.yaml description:entry-description.md` against a `list`-mode `PROJECTS` target
- **THEN** the CLI records a record-only dependency on record `bldoc` sourced from `defaults.yaml`, and a dependency on field `description` of record `bldoc` sourced from `entry-description.md`

#### Scenario: `@` in the target argument is rejected
- **WHEN** the user runs `bldoc add-dep README@version pyproject.toml go.mod` (three arguments, first argument containing `@`)
- **THEN** the CLI reports a usage error and exits non-zero, and the manifest is unchanged

#### Scenario: Fewer than two pairs after target and record name is rejected
- **WHEN** the user runs `bldoc add-dep PROJECTS bldoc description:entry-description.md` (three arguments) against a `list`-mode `PROJECTS` target, leaving only one pair after the record name
- **THEN** the CLI reports a usage error and exits non-zero, and the manifest is unchanged

#### Scenario: Two-argument invocation still uses the existing single-dependency form
- **WHEN** the user runs `bldoc add-dep README pyproject.toml` (two arguments)
- **THEN** the CLI parses it using the existing target-ref/source-ref grammar, not the variadic form, exactly as before this capability existed

### Requirement: Pair grammar in the variadic `add-dep` form
A pair SHALL be either a bare source-ref (`source`, `source@path`, or
`source#anchor`, parsed with the existing source-ref grammar) — naming a
whole-file dependency, an anchor-addressed dependency, or, within a
list-mode record, a record-only dependency — or `<field>:<source-ref>`,
where `<field>` is a single field name (no `.`, `@`, or `:`) and the
remainder is parsed as a source-ref with the existing grammar — naming a
field-mode dependency or, within a list-mode record, a single record
field. A pair is split on its first `:`, if any; because a source-ref
never itself contains `:`, this split is unambiguous regardless of the
target's mode. A pair's field-vs-bare shape is validated against the
target's mode by the same mode-exclusivity and list-mode addressing
rules that already govern a single `add-dep` call (see the
`manifest-store` capability), not by the CLI's argument parsing.

#### Scenario: Field name split from a source-ref carrying its own path
- **WHEN** the user runs `bldoc add-dep README version:pyproject.toml@project.version other:go.mod`
- **THEN** the CLI parses the first pair's field as `version` and source-ref as `pyproject.toml@project.version` (source `pyproject.toml`, path `project.version`), and the second pair's field as `other` and source-ref as `go.mod`

#### Scenario: Field name split from a source-ref carrying an anchor
- **WHEN** the user runs `bldoc add-dep README summary:spec.md#purpose section:spec.md#requirements`
- **THEN** the CLI parses the first pair's field as `summary` and source-ref as `spec.md#purpose` (source `spec.md`, anchor `purpose`)

#### Scenario: A bare pair on a field-mode target is rejected by mode exclusivity
- **WHEN** the user runs `bldoc add-dep README version:pyproject.toml@project.version go.mod` against a `field`-mode `README` target, where the second pair has no `field:` prefix
- **THEN** the CLI reports a mode-exclusivity error and exits non-zero, and the manifest is unchanged

### Requirement: `--format` and `--nested` apply to every pair in a batch call
Given the variadic `add-dep` form, `--format` and `--nested`, if
supplied, SHALL apply uniformly to every pair recorded by that
invocation. There is no per-pair override in this form.

#### Scenario: `--nested` applied to every anchor pair in a batch call
- **WHEN** the user runs `bldoc add-dep README section:spec.md#requirements sub:spec.md#configuration --nested`
- **THEN** the CLI records both dependencies with `--nested` set

#### Scenario: `--format` applied to every field pair in a batch call
- **WHEN** the user runs `bldoc add-dep README --format "%s (pinned)" version:pyproject.toml@project.version tool:pyproject.toml@project.name`
- **THEN** the CLI records both dependencies with the same format template
