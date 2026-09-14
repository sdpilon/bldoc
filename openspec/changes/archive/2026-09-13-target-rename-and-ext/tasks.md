## 1. Manifest schema & operations

- [x] 1.1 Add optional `Mode` and `Ext` fields to the `Target` struct in `internal/manifest`, and verify an existing manifest with neither key still loads unchanged (existing `io_test.go` cases pass as-is)
- [x] 1.2 Support recording `Mode` (validated as exactly `raw` or `field`) and `Ext` when a target is created, and verify via unit tests that a target created with either round-trips through `Save`/`Load` intact
- [x] 1.3 Extend `AddDep`'s raw/field mode-exclusivity check to branch on whether the target has a declared `Mode`: if set, reject any dependency (including the first) whose field-addressing doesn't match it; if unset, keep today's exact behavior (checked against existing recorded deps, first dependency unconstrained). Verify via unit tests covering: rejection on the first dependency of a declared-mode target, and the unchanged legacy behavior for a target with no declared mode
- [x] 1.4 Implement a `RenameTarget` operation — reject a missing old name, reject a duplicate new name, otherwise rename in place preserving dependency order, `Mode`, and `Ext` — and verify via unit tests for both rejection cases and the successful rename

## 2. CLI command surface

- [x] 2.1 Make `--mode` required on the `new` command, validated as exactly `raw` or `field`, and verify via a cli-package test covering: missing flag, invalid value, and a valid value being recorded
- [x] 2.2 Add `--ext` flag parsing to `new`, rejected immediately when combined with `--mode field`, and verify via a cli-package test covering both the rejection and the accepted `--mode raw --ext <ext>` case
- [x] 2.3 Add a new `rename` command (`internal/cli/rename.go`) requiring exactly two positional arguments, wired into `root.go`'s command list, and verify via a cli-package test covering missing-argument and extra-argument usage errors
- [x] 2.4 Wire `rename`'s success/error paths to `manifest.RenameTarget` using the existing `reportErr` convention, and verify via a cli-package test for the unknown-old-target error, the duplicate-new-name error, and the successful rename

## 3. Compile engine

- [x] 3.1 Update raw-mode output path resolution so a target with `Ext` set compiles to `.bldoc/<target>.<ext>` instead of `.bldoc/<target>`, and verify via a test that `bldoc make` on an ext-bearing raw-mode target writes to the extended path
- [x] 3.2 Update the empty-target case to branch on declared `Mode`: empty raw/undeclared compiles to an empty raw-mode intermediate (unchanged); empty declared-field compiles to `{}` at `.bldoc/<target>.json`. Verify via tests covering both branches
- [x] 3.3 Add a regression test confirming a non-empty field-mode target's output path (`.bldoc/<target>.json`) is unaffected by any `Ext` value

## 4. `show` output

- [x] 4.1 Update `show`'s output to include a target's mode (declared if set, otherwise computed from its dependencies) and its extension when set, and verify via cli-package tests covering: a declared-mode target, a legacy target with computed mode, and an ext-bearing target

## 5. Verification

- [x] 5.1 Run `go test ./...` and confirm all existing and new tests pass
- [x] 5.2 Manually exercise the new behavior end to end against `examples/demo/`: create targets with `bldoc new PROJECTS --mode field` and `bldoc new NOTES --mode raw --ext txt`, confirm a mismatched first `add-dep` is rejected on each, add matching dependencies, `bldoc make` both, and confirm `.bldoc/PROJECTS.json` and `.bldoc/NOTES.txt` contain the expected output; then `bldoc rename` one and confirm `bldoc show` on the new name reflects the same dependencies, mode, and extension
- [x] 5.3 Confirm `examples/demo/bldoc.toml` (written before this change, with no `mode`/`ext` keys) still loads and `bldoc make` still compiles it identically to before
