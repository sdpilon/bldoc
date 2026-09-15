## MODIFIED Requirements

### Requirement: `new` requires an explicit mode
`new` SHALL require a `--mode` flag whose value is exactly `raw`, `field`, or `list`, and SHALL reject the invocation when it is omitted or set to any other value, leaving the manifest unchanged.

#### Scenario: Mode omitted rejected
- **WHEN** `bldoc new README` is run with no `--mode` flag
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Invalid mode value rejected
- **WHEN** `bldoc new README --mode nonsense` is run
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Valid mode recorded
- **WHEN** `bldoc new PROJECTS --mode field` is run
- **THEN** the manifest records target `PROJECTS` with mode `field` and no dependencies

#### Scenario: Valid list mode recorded
- **WHEN** `bldoc new PROJECTS --mode list` is run
- **THEN** the manifest records target `PROJECTS` with mode `list` and no dependencies

### Requirement: `new` accepts an optional output extension
`new` SHALL accept an optional `--ext <ext>` flag regardless of the declared mode, and record it on the created target.

#### Scenario: Extension recorded
- **WHEN** `bldoc new PROJECTS --mode raw --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `raw` and extension `yaml`

#### Scenario: Extension recorded with field mode
- **WHEN** `bldoc new PROJECTS --mode field --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `field` and extension `yaml`

#### Scenario: Extension recorded with list mode
- **WHEN** `bldoc new PROJECTS --mode list --ext yaml` is run
- **THEN** the manifest records target `PROJECTS` with mode `list` and extension `yaml`

#### Scenario: No extension by default
- **WHEN** `bldoc new README --mode raw` is run with no `--ext`
- **THEN** the manifest records target `README` with mode `raw` and no extension set

## ADDED Requirements

### Requirement: `add-dep` requires record-scoped addressing on a list-mode target
For a target whose declared mode is `list`, `add-dep` SHALL reject a target-ref with no field component (a bare target name), and SHALL reject a field component containing more than one `.`. A field component with no `.` (e.g. `bldoc`) names a record-only dependency; a field component with exactly one `.` (e.g. `bldoc.description`) names a single field within that record.

#### Scenario: Bare target-ref rejected on a list-mode target
- **WHEN** `bldoc add-dep PROJECTS entry.yaml` is run and `PROJECTS` was created with `bldoc new PROJECTS --mode list`
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Record-only dependency accepted
- **WHEN** `bldoc add-dep PROJECTS:bldoc project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a record-only dependency on record `bldoc`, sourced from `project-headers/bldoc.yaml`

#### Scenario: Record-field dependency accepted
- **WHEN** `bldoc add-dep PROJECTS:bldoc.description entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records a dependency on field `description` of record `bldoc`, sourced from `entry-description.md`

#### Scenario: Over-qualified field component rejected
- **WHEN** `bldoc add-dep PROJECTS:bldoc.description.extra entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

### Requirement: `--format` requires a record-field-scoped target-ref on a list-mode target
For a target whose declared mode is `list`, `add-dep` SHALL reject `--format` when the target-ref's field component has no `.` (i.e. names a record-only dependency rather than a single record field).

#### Scenario: Format rejected on a record-only dependency
- **WHEN** `bldoc add-dep PROJECTS:bldoc --format "..." project-headers/bldoc.yaml` is run against a `list`-mode `PROJECTS` target
- **THEN** the CLI reports an error and exits non-zero, and the manifest is unchanged

#### Scenario: Format accepted on a record-field dependency
- **WHEN** `bldoc add-dep PROJECTS:bldoc.description --format "%s" entry-description.md` is run against a `list`-mode `PROJECTS` target
- **THEN** the manifest records the dependency on field `description` of record `bldoc` with the given format template
