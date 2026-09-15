# cli-shell Specification

## Purpose

The `cli-shell` capability defines the `bldoc` command-line interface's
command surface, argument grammar, and baseline behavior, independent of
any manifest or compilation logic, which later capabilities introduce.

## Requirements

### Requirement: Recognized command surface
The CLI SHALL recognize exactly eight commands: `new`, `add-dep`,
`rm-dep`, `make`, `rm`, `rename`, `list`, and `show`. Any other first
argument SHALL be rejected as an unknown command.

#### Scenario: Unknown command is rejected
- **WHEN** the user runs `bldoc frobnicate`
- **THEN** the CLI prints an unknown-command error and exits non-zero

#### Scenario: Recognized command is accepted
- **WHEN** the user runs any of the eight recognized commands with valid arguments
- **THEN** the CLI does not report an unknown-command error

### Requirement: `new` requires a single target name
`new` SHALL require exactly one positional argument: the target name.

#### Scenario: Missing target name
- **WHEN** the user runs `bldoc new` with no arguments
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Valid target name
- **WHEN** the user runs `bldoc new README`
- **THEN** the CLI accepts the argument and records a new target named `README` in the manifest (see the `manifest-store` capability)

### Requirement: `new` requires a `--mode` flag
`new` SHALL require a `--mode` flag whose value is `raw`, `field`, or `list`, in addition to its positional target name.

#### Scenario: Missing mode flag
- **WHEN** the user runs `bldoc new README` with no `--mode` flag
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Invalid mode value
- **WHEN** the user runs `bldoc new README --mode nonsense`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Valid mode flag
- **WHEN** the user runs `bldoc new PROJECTS --mode field`
- **THEN** the CLI parses `PROJECTS` as the target name and `field` as the mode (see the `manifest-store` capability)

#### Scenario: Valid list mode flag
- **WHEN** the user runs `bldoc new PROJECTS --mode list`
- **THEN** the CLI parses `PROJECTS` as the target name and `list` as the mode (see the `manifest-store` capability)

### Requirement: `--ext` is accepted with any mode
`new` SHALL accept an optional `--ext <ext>` flag regardless of the declared `--mode`.

#### Scenario: Extension accepted with raw mode
- **WHEN** the user runs `bldoc new PROJECTS --mode raw --ext yaml`
- **THEN** the CLI parses `PROJECTS` as the target name, `raw` as the mode, and `yaml` as the extension, and records all three (see the `manifest-store` capability)

#### Scenario: Extension accepted with field mode
- **WHEN** the user runs `bldoc new PROJECTS --mode field --ext yaml`
- **THEN** the CLI parses `PROJECTS` as the target name, `field` as the mode, and `yaml` as the extension, and records all three (see the `manifest-store` capability)

#### Scenario: Extension accepted with list mode
- **WHEN** the user runs `bldoc new PROJECTS --mode list --ext yaml`
- **THEN** the CLI parses `PROJECTS` as the target name, `list` as the mode, and `yaml` as the extension, and records all three (see the `manifest-store` capability)

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
- **THEN** the CLI compiles `README`'s recorded dependencies into its intermediate (see the `compile-engine` capability)

#### Scenario: Make with no target
- **WHEN** the user runs `bldoc make` with no arguments
- **THEN** the CLI compiles every target recorded in the manifest into its intermediate (see the `compile-engine` capability)

#### Scenario: Make rejects more than one target
- **WHEN** the user runs `bldoc make README ROADMAP`
- **THEN** the CLI reports a usage error and exits non-zero

### Requirement: `rename` requires exactly two positional arguments
`rename` SHALL require exactly two positional arguments: the old target name, then the new target name.

#### Scenario: Missing arguments
- **WHEN** the user runs `bldoc rename` or `bldoc rename README` with fewer than two arguments
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Extra argument rejected
- **WHEN** the user runs `bldoc rename README NOTES EXTRA`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Valid rename invocation
- **WHEN** the user runs `bldoc rename README NOTES`
- **THEN** the CLI accepts the arguments and renames the target (see the `manifest-store` capability)

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
