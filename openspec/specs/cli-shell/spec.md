# cli-shell Specification

## Purpose

The `cli-shell` capability defines the `bldoc` command-line interface's
command surface, argument grammar, and baseline behavior, independent of
any manifest or compilation logic, which later capabilities introduce.

## Requirements

### Requirement: Recognized command surface
The CLI SHALL recognize exactly seven commands: `new`, `add-dep`,
`rm-dep`, `make`, `rm`, `list`, and `show`. Any other first argument SHALL
be rejected as an unknown command.

#### Scenario: Unknown command is rejected
- **WHEN** the user runs `bldoc frobnicate`
- **THEN** the CLI prints an unknown-command error and exits non-zero

#### Scenario: Recognized command is accepted
- **WHEN** the user runs any of the seven recognized commands with valid arguments
- **THEN** the CLI does not report an unknown-command error

### Requirement: `new` requires a single target name
`new` SHALL require exactly one positional argument: the target name.

#### Scenario: Missing target name
- **WHEN** the user runs `bldoc new` with no arguments
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Valid target name
- **WHEN** the user runs `bldoc new README`
- **THEN** the CLI accepts the argument and records a new target named `README` in the manifest (see the `manifest-store` capability)

### Requirement: `add-dep` parses target-ref, source-ref, and optional format
`add-dep` SHALL require a target-ref, a source-ref, and accept an
optional `--format` flag. A target-ref SHALL be either a bare target name
or `target:field`. A source-ref SHALL be either a bare file path or
`path:field-path`.

#### Scenario: Whole-file dependency
- **WHEN** the user runs `bldoc add-dep README pyproject.toml`
- **THEN** the CLI parses `README` as the target with no field, and `pyproject.toml` as the source with no field-path

#### Scenario: Field-addressed dependency with format
- **WHEN** the user runs `bldoc add-dep README:version --format "Python version must be %s to run this project." pyproject.toml:project.requires-python`
- **THEN** the CLI parses `README` as the target with field `version`, `pyproject.toml` as the source with field-path `project.requires-python`, and captures the format string

#### Scenario: Malformed target-ref is rejected
- **WHEN** the user runs `add-dep` with a target-ref containing more than one `:`
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `--format` requires a field-addressed target-ref
`add-dep` SHALL reject `--format` when the target-ref has no `:field`.

#### Scenario: Format on a whole-file dependency is rejected
- **WHEN** the user runs `bldoc add-dep README --format "..." pyproject.toml` with no `:field` on the target-ref
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `rm-dep` parses target-ref and source-ref
`rm-dep` SHALL require a target-ref and a source-ref, using the same
grammar as `add-dep`, and SHALL NOT accept a `--format` flag.

#### Scenario: Removing a whole-file dependency
- **WHEN** the user runs `bldoc rm-dep README pyproject.toml`
- **THEN** the CLI parses the target-ref and source-ref and removes the matching dependency from the manifest (see the `manifest-store` capability)

#### Scenario: `--format` is not accepted
- **WHEN** the user runs `bldoc rm-dep` with a `--format` flag
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `make` accepts an optional target
`make` SHALL accept zero or one positional argument. Given a target name,
it refers to compiling that target; given no argument, it refers to
compiling every target in the manifest.

#### Scenario: Make with an explicit target
- **WHEN** the user runs `bldoc make README`
- **THEN** the CLI accepts the argument and reports the command as not yet implemented

#### Scenario: Make with no target
- **WHEN** the user runs `bldoc make` with no arguments
- **THEN** the CLI accepts the invocation as a build-all request and reports the command as not yet implemented

#### Scenario: Make rejects more than one target
- **WHEN** the user runs `bldoc make README ROADMAP`
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `rm` and `show` require a single target
`rm` and `show` SHALL each require exactly one positional argument: the
target name.

#### Scenario: Missing target
- **WHEN** the user runs `bldoc rm` or `bldoc show` with no arguments
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

### Requirement: Help and version behave normally
`--help`/`-h` on the root command or any subcommand SHALL print usage
information and exit 0. `--version` SHALL print the CLI's version and
exit 0.

#### Scenario: Help exits successfully
- **WHEN** the user runs `bldoc --help` or `bldoc <command> --help`
- **THEN** the CLI prints usage information and exits 0

#### Scenario: Version exits successfully
- **WHEN** the user runs `bldoc --version`
- **THEN** the CLI prints a version string and exits 0
