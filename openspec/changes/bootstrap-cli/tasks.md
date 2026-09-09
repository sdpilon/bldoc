## 1. Project Scaffolding

- [ ] 1.1 Initialize the Go module (`go mod init bldoc`), create the `cmd/bldoc/main.go` entrypoint and `internal/` layout, and verify `go build ./...` succeeds
- [ ] 1.2 Add the Cobra dependency and verify it resolves cleanly in `go.mod`/`go.sum`
- [ ] 1.3 Set up an in-process Cobra command test harness (invoke `Execute()`, capture stdout/stderr/exit code) and verify it correctly reports a deliberately failing sample assertion before relying on it

## 2. Root Command & Global Behavior

- [ ] 2.1 Wire up the root command with `--version` and verify `bldoc --version` exits 0 and prints a version string
- [ ] 2.2 Verify `bldoc --help` and `bldoc <command> --help` exit 0 and print usage text for every subcommand
- [ ] 2.3 Verify an unrecognized command (e.g. `bldoc frobnicate`) exits non-zero with an unknown-command error

## 3. Addressing Grammar Parsing

- [ ] 3.1 Implement the shared target-ref parser (`target` or `target:field`, rejecting more than one `:`) with unit tests covering both forms and the malformed case
- [ ] 3.2 Implement the shared source-ref parser (`source` or `source:path`) with unit tests covering both forms

## 4. Command Implementation

- [ ] 4.1 Implement `new <target>` requiring exactly one positional argument; verify the missing-argument usage error and the valid-argument not-yet-implemented behavior
- [ ] 4.2 Implement `add-dep <target-ref> [--format <template>] <source-ref>` using the shared parsers; verify whole-file dependency parsing, field-addressed dependency parsing with `--format`, the malformed-target-ref rejection, and the `--format`-without-`:field` rejection
- [ ] 4.3 Implement `rm-dep <target-ref> <source-ref>` using the shared parsers with no `--format` flag; verify correct parsing and that a supplied `--format` is rejected as an unknown flag
- [ ] 4.4 Implement `make [<target>]` accepting zero or one positional argument; verify the explicit-target, no-target (build-all), and too-many-arguments scenarios
- [ ] 4.5 Implement `rm <target>` and `show <target>`, each requiring exactly one positional argument; verify the missing-argument usage error for each
- [ ] 4.6 Implement `list` accepting no positional arguments; verify the extra-argument rejection and the no-argument acceptance

## 5. Final Verification

- [ ] 5.1 Run the full test suite against every scenario listed in `specs/cli-shell/spec.md` (including the not-yet-implemented stderr message/exit-1 behavior across all seven commands) and verify `go test ./...` passes with no scenario left uncovered
