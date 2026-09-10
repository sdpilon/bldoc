## MODIFIED Requirements

### Requirement: `new` requires a single target name
`new` SHALL require exactly one positional argument: the target name.

#### Scenario: Missing target name
- **WHEN** the user runs `bldoc new` with no arguments
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Valid target name
- **WHEN** the user runs `bldoc new README`
- **THEN** the CLI accepts the argument and records a new target named `README` in the manifest (see the `manifest-store` capability)

### Requirement: `rm-dep` parses target-ref and source-ref
`rm-dep` SHALL require a target-ref and a source-ref, using the same
grammar as `add-dep`, and SHALL NOT accept a `--format` flag.

#### Scenario: Removing a whole-file dependency
- **WHEN** the user runs `bldoc rm-dep README pyproject.toml`
- **THEN** the CLI parses the target-ref and source-ref and removes the matching dependency from the manifest (see the `manifest-store` capability)

#### Scenario: `--format` is not accepted
- **WHEN** the user runs `bldoc rm-dep` with a `--format` flag
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `list` takes no arguments
`list` SHALL reject any positional arguments.

#### Scenario: Extra argument rejected
- **WHEN** the user runs `bldoc list README`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: No arguments accepted
- **WHEN** the user runs `bldoc list`
- **THEN** the CLI accepts the invocation and lists the manifest's targets (see the `manifest-store` capability)

### Requirement: Valid invocations report not-yet-implemented
`make`, given syntactically valid arguments, SHALL print a message to
stderr indicating the command is not yet implemented and exit with
status 1. `new`, `add-dep`, `rm-dep`, `rm`, `list`, and `show` no longer
exhibit this behavior — each performs its real manifest operation
instead (see the `manifest-store` capability).

#### Scenario: Not-yet-implemented exit code
- **WHEN** `make` is run with valid arguments
- **THEN** the CLI exits with status 1 and prints a message to stderr naming the command as not yet implemented
