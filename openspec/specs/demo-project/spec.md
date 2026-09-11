# demo-project Specification

## Purpose

The `demo-project` capability is a checked-in, ready-to-run example under `examples/demo/` that lets anyone try `bldoc`'s real compile behavior immediately, without pointing the tool at a real project.

## Requirements

### Requirement: Demo project structure
`examples/demo/` SHALL contain a `bldoc.toml` manifest, a `README.md`, a `.gitignore`, and every source file referenced by that manifest's dependencies.

#### Scenario: Demo directory present
- **WHEN** `examples/demo/` is inspected
- **THEN** it contains `bldoc.toml`, `README.md`, `.gitignore`, and each source file referenced as a dependency in `bldoc.toml`

### Requirement: Demo manifest covers raw-mode and field-mode
The checked-in `examples/demo/bldoc.toml` SHALL declare at least one raw-mode target with two or more whole-file dependencies, and at least one field-mode target with dependencies sourced from a TOML file, a JSON file, and a YAML file, with at least one of those dependencies using a format template.

#### Scenario: Raw-mode target present
- **WHEN** `examples/demo/bldoc.toml` is inspected
- **THEN** it declares a target whose recorded dependencies are all whole-file (unnamed), with two or more dependencies

#### Scenario: Field-mode target present
- **WHEN** `examples/demo/bldoc.toml` is inspected
- **THEN** it declares a target with field-addressed dependencies resolving from a `.toml`, a `.json`, and a `.yaml`/`.yml` source file, and at least one of those dependencies records a format template

### Requirement: Demo compiles successfully
Running `bldoc make` from within `examples/demo/` SHALL exit zero and write an intermediate under `.bldoc/` for every target declared in `examples/demo/bldoc.toml`.

#### Scenario: make succeeds against the demo
- **WHEN** `bldoc make` is run with `examples/demo/` as the working directory
- **THEN** the command exits zero and `.bldoc/` contains an intermediate for each target declared in `examples/demo/bldoc.toml`

### Requirement: Demo output is disposable
`examples/demo/.gitignore` SHALL exclude `.bldoc/`, the demo's own compiled output directory, from version control.

#### Scenario: Compiled output ignored
- **WHEN** `bldoc make` is run against the demo and `.bldoc/` is created
- **THEN** `git status` run from the repository root does not report any file under `examples/demo/.bldoc/` as untracked

### Requirement: Demo documents its own construction
`examples/demo/README.md` SHALL document the sequence of `bldoc new` and `bldoc add-dep` commands that produce the checked-in `bldoc.toml`, and SHALL document how to run `bldoc make` against the demo.

#### Scenario: README documents commands and usage
- **WHEN** `examples/demo/README.md` is read
- **THEN** it lists the `bldoc new`/`bldoc add-dep` commands that construct the checked-in manifest's targets and dependencies, and it states how to invoke `bldoc make` against the demo
