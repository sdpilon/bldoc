## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: `--ext` requires `--mode raw`
**Reason**: `--ext` is no longer raw-mode-exclusive — `field`- and `list`-mode targets can now declare it too, to select both the compiled intermediate's filename suffix and its output serialization (see the `compile-engine` capability).
**Migration**: No change for existing raw-mode usage. A previously-rejected invocation like `bldoc new PROJECTS --mode field --ext yaml` is now accepted; see the new "`--ext` is accepted with any mode" requirement below.

## ADDED Requirements

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
