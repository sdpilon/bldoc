## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: `new` requires an explicit mode
`new` SHALL require a `--mode` flag whose value is exactly `raw` or `field`, and SHALL reject the invocation when it is omitted or set to any other value, leaving the manifest unchanged.

#### Scenario: Mode omitted rejected
- **WHEN** `bldoc new README` is run with no `--mode` flag
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Invalid mode value rejected
- **WHEN** `bldoc new README --mode nonsense` is run
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Valid mode recorded
- **WHEN** `bldoc new PROJECTS --mode field` is run
- **THEN** the manifest records target `PROJECTS` with mode `field` and no dependencies

### Requirement: `new` accepts an optional output extension
Given `--mode raw`, `new` SHALL accept an optional `--ext <ext>` flag and record it on the created target.

#### Scenario: Extension recorded
- **WHEN** `bldoc new PROJECTS --mode raw --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `raw` and extension `yaml`

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
