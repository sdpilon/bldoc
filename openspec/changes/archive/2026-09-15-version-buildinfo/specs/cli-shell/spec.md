## MODIFIED Requirements

### Requirement: Help and version behave normally
`--help`/`-h` on the root command or any subcommand SHALL print usage
information and exit 0. `--version` SHALL print the CLI's version and
exit 0. When the binary carries embedded VCS build info (as reported by
`runtime/debug.ReadBuildInfo()`), `--version`'s output SHALL include a
parenthetical suffix with the commit's short SHA, its commit timestamp,
and the word `modified` when the working tree had uncommitted changes at
build time; when no VCS build info is embedded, `--version` SHALL print
only the plain version string, with no parenthetical.

#### Scenario: Help exits successfully
- **WHEN** the user runs `bldoc --help` or `bldoc <command> --help`
- **THEN** the CLI prints usage information and exits 0

#### Scenario: Version exits successfully
- **WHEN** the user runs `bldoc --version`
- **THEN** the CLI prints a version string and exits 0

#### Scenario: Version includes VCS build info when available
- **WHEN** the user runs `bldoc --version` on a binary built from a git checkout with VCS stamping enabled
- **THEN** the printed version string includes a parenthetical with the commit's short SHA and its commit timestamp

#### Scenario: Version flags an uncommitted-change build
- **WHEN** the user runs `bldoc --version` on a binary built while the working tree had uncommitted changes
- **THEN** the printed version string's parenthetical includes `modified`

#### Scenario: Version omits the parenthetical when VCS info is unavailable
- **WHEN** the user runs `bldoc --version` on a binary built with no embedded VCS build info (e.g. `-buildvcs=false`, or no `.git` present at build time)
- **THEN** the CLI prints the plain version string with no parenthetical, and does not error
