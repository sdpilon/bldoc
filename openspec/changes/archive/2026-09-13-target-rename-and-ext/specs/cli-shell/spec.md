## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: `new` requires a `--mode` flag
`new` SHALL require a `--mode` flag whose value is `raw` or `field`, in addition to its positional target name.

#### Scenario: Missing mode flag
- **WHEN** the user runs `bldoc new README` with no `--mode` flag
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Invalid mode value
- **WHEN** the user runs `bldoc new README --mode nonsense`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Valid mode flag
- **WHEN** the user runs `bldoc new PROJECTS --mode field`
- **THEN** the CLI parses `PROJECTS` as the target name and `field` as the mode (see the `manifest-store` capability)

### Requirement: `--ext` requires `--mode raw`
`new` SHALL accept an optional `--ext <ext>` flag, and SHALL reject it when given together with `--mode field`.

#### Scenario: Extension rejected with field mode
- **WHEN** the user runs `bldoc new PROJECTS --mode field --ext yaml`
- **THEN** the CLI reports a usage error and exits non-zero

#### Scenario: Extension accepted with raw mode
- **WHEN** the user runs `bldoc new PROJECTS --mode raw --ext yaml`
- **THEN** the CLI parses `PROJECTS` as the target name, `raw` as the mode, and `yaml` as the extension, and records all three (see the `manifest-store` capability)

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
