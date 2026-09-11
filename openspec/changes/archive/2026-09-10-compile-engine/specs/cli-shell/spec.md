## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: Valid invocations report not-yet-implemented
**Reason**: This requirement scoped the not-yet-implemented stub
behavior down to `make` as each other command gained real behavior in
`manifest-io`. This change gives `make` real behavior too, so no
recognized command exhibits this behavior any longer.
**Migration**: None — there is no replacement requirement. Every
recognized command's valid-argument behavior is now specified by its own
requirement (see `cli-shell`'s per-command requirements and the
`manifest-store`/`compile-engine` capabilities).
