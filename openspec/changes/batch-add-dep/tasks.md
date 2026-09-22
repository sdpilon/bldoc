## 1. Rename the ref field-addressing symbol from `:` to `@`

- [ ] 1.1 In `internal/cli/ref.go` (`parseTargetRef`) and `internal/cli/source_ref.go` (`parseSourceRef`), change the split character from `:` to `@` (source-ref's `#anchor` handling and `:field-path`-vs-`#anchor` conflict check both move to `@`), and update their error messages to name `@`; verify by updating `ref_test.go`/`source_ref_test.go` to use `@` and confirming they pass
- [ ] 1.2 Update every existing test in `internal/cli` (`add_dep_test.go`, `rm_dep_test.go`, and any other test using `target:field`/`source:path` syntax) to use `@`, and verify `go test ./...` passes
- [ ] 1.3 Update `examples/` and any README/CLI-help text using `:field`/`:path` syntax to `@field`/`@path`, per CLAUDE.md's pre-PR doc-consistency sweep

## 2. Pair parsing

- [ ] 2.1 Add a `pair` type and a `parsePair(s string) (pair, error)` helper in `internal/cli` that splits `s` on its first `:` (if present) into a field name and a source-ref parsed with `parseSourceRef`, or parses `s` directly as a bare source-ref if no `:` is present — with no mode-awareness, per design.md's "split on the first `:`, unconditionally" decision; verify with unit tests covering a bare source-ref, `field:source`, `field:source@path`, `field:source#anchor`
- [ ] 2.2 Add a bare-name validator (rejects `@`, `:`, and `.` for the record-name argument; rejects `@` for the target argument) and cover it with unit tests for valid and invalid inputs

## 3. `add-dep` variadic dispatch

- [ ] 3.1 Change `add-dep`'s `cobra.Command.Args` from `exactArgs(2)` to a minimum of 2, and add a branch in `RunE`: `len(args) == 2` keeps today's code path unchanged (now using `@`); `len(args) >= 3` enters the new batch path — verify the existing two-argument test suite (`add_dep_test.go`, updated for `@` in task 1.2) still passes
- [ ] 3.2 In the batch path, reject a first argument containing `@` with a usage error before loading the manifest, and verify with a test asserting a non-zero exit and unchanged manifest
- [ ] 3.3 In the batch path, load the manifest, look up the target (usage error on unknown target, same message shape as today), and branch on `target.Mode == "list"` to decide whether `args[1]` is a record name; verify with a test that an unknown target in the batch form errors the same way as the two-argument form
- [ ] 3.4 Reject an invocation left with fewer than two pairs after removing the target and any record-name argument, including what was parsed (record name, pair count) in the error message; verify with a test for the fewer-than-two-pairs case on both a list-mode and a non-list-mode target

## 4. Recording pairs

- [ ] 4.1 Build each pair's `manifest.Dep` (reusing the same construction the two-argument path already does: `Source`, `Path`, `Anchor` from the pair's source-ref; `Nested` and `Format` from the invocation-level flags; `Field` from the pair, combined with the record name via `manifest.SplitRecordField`'s inverse for list-mode) and call `manifest.AddDep` once per pair against one loaded `*manifest.Manifest`, in argument order — letting `AddDep`'s existing mode-exclusivity checks reject a pair whose shape doesn't fit the target's mode
- [ ] 4.2 Call `manifest.Save` exactly once, only after every `AddDep` call in the loop has succeeded; on any `AddDep` error, return immediately without saving, and verify with a test that a later-pair failure (duplicate, mode-exclusivity, or malformed) leaves the manifest file byte-for-byte unchanged from before the call, even when earlier pairs in the same invocation were individually valid
- [ ] 4.3 Run `warnAmbiguousAnchor` for each anchor-addressed pair, same as the two-argument path does today, and verify a batch call with two anchor pairs produces a warning for each ambiguous one

## 5. List-mode record batching

- [ ] 5.1 Verify (with a test) the exact scenario from proposal.md: `add-dep PROJECTS bldoc description:entry-description.md path:project-headers/bldoc.yaml` against a `list`-mode target records both fields under record `bldoc`, matching what `show` would report as two separate `add-dep PROJECTS@bldoc.<field>` calls would have produced
- [ ] 5.2 Verify a record-only pair (no `:`) combined with a field pair in the same batch call records correctly, and that `--format` combined with a record-only pair in a batch call is rejected

## 6. Flag behavior

- [ ] 6.1 Verify `--format` and `--nested` given on a batch call apply to every pair recorded by that call (one test per flag, each with two pairs)

## 7. Documentation

- [ ] 7.1 Update any `bldoc add-dep`/`rm-dep` usage text or examples in the project's README (or other top-level docs, per CLAUDE.md's pre-PR doc-consistency sweep) that show `:field`/`:path` syntax or describe `add-dep` as taking exactly one dependency per call
