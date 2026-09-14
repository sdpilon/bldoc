## MODIFIED Requirements

### Requirement: Intermediate file location and naming
The compiled intermediate for a target SHALL be written under `.bldoc/`
in the current working directory: a raw-mode target with no recorded
extension at `.bldoc/<target>` (raw bytes); a raw-mode target with a
recorded extension at `.bldoc/<target>.<ext>` (raw bytes); and a
field-mode target's intermediate at `.bldoc/<target>.json` (a single
JSON object), regardless of any recorded extension.

#### Scenario: Raw-mode output path
- **WHEN** `bldoc make README` compiles a raw-mode target with no recorded extension
- **THEN** the CLI writes the concatenated bytes to `.bldoc/README`

#### Scenario: Raw-mode output path with a declared extension
- **WHEN** `bldoc make PROJECTS` compiles a raw-mode target recorded with extension `yaml`
- **THEN** the CLI writes the concatenated bytes to `.bldoc/PROJECTS.yaml`

#### Scenario: Field-mode output path
- **WHEN** `bldoc make README` compiles a field-mode target
- **THEN** the CLI writes the JSON object to `.bldoc/README.json`

### Requirement: A target with no dependencies compiles to an empty intermediate
`make` SHALL compile a target with no recorded dependencies and no
declared mode, or with a declared raw mode, to an empty raw-mode
intermediate. `make` SHALL compile a target with no recorded dependencies
and a declared field mode to an empty field-mode intermediate (a JSON
object with no entries).

#### Scenario: Empty target
- **WHEN** `README` has no recorded dependencies and no declared mode, and `make README` is run
- **THEN** `.bldoc/README` is created empty (zero bytes)

#### Scenario: Empty declared field-mode target
- **WHEN** `PROJECTS` was created with `bldoc new PROJECTS --mode field`, has no recorded dependencies, and `make PROJECTS` is run
- **THEN** `.bldoc/PROJECTS.json` is created containing `{}`
