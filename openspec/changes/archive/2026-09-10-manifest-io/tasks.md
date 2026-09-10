## 1. Manifest Package Foundations

- [x] 1.1 Add a TOML dependency and verify it resolves cleanly in `go.mod`/`go.sum`
- [x] 1.2 Create the `internal/manifest` package with schema structs for a target and its dependencies (source, optional field path, optional target-side field, optional format template) and verify `go build ./...` succeeds
- [x] 1.3 Implement `Load`/`Save` against `bldoc.toml` in the current working directory, with unit tests covering a round-trip (save then load) and a malformed-file parse error

## 2. Manifest Mutation Operations

- [x] 2.1 Implement `AddTarget` (creates `bldoc.toml` if absent, rejects a duplicate target name) with unit tests for both the create-manifest and duplicate-rejection cases
- [x] 2.2 Implement `AddDep` — existing-target check, raw/field mode exclusivity enforcement, duplicate-dependency rejection ((source, path) identity), append preserving order, storing the optional format template — with unit tests covering each rejection case and both the whole-file and field-addressed success cases
- [x] 2.3 Implement `RemoveDep` — existing-target check, unrecorded-dependency rejection, order-preserving removal of the remaining dependencies — with unit tests covering both rejection cases and the successful-removal case
- [x] 2.4 Implement `RemoveTarget` — existing-target check, deletes the target and all its recorded dependencies — with unit tests covering the rejection and successful-removal cases
- [x] 2.5 Implement a `ListTargets` accessor returning target names in declaration order, empty for a missing manifest or a manifest with no targets, with unit tests for both
- [x] 2.6 Implement a `ShowTarget` accessor — existing-target check, returns the target's dependencies in the order they were added — with unit tests covering the rejection and successful cases

## 3. CLI Command Wiring

- [x] 3.1 Wire `new` to `AddTarget`; verify the duplicate-target and valid-target scenarios from `specs/manifest-store/spec.md`
- [x] 3.2 Wire `add-dep` to `AddDep`; verify the unknown-target, mode-exclusivity (both directions), duplicate-dependency, and both success scenarios
- [x] 3.3 Wire `rm-dep` to `RemoveDep`; verify the unknown-target, unrecorded-dependency, and successful-removal scenarios
- [x] 3.4 Wire `rm` to `RemoveTarget`; verify the unknown-target and successful-removal scenarios
- [x] 3.5 Wire `list` to `ListTargets`; verify the targets-listed and no-targets scenarios
- [x] 3.6 Wire `show` to `ShowTarget`; verify the unknown-target and dependencies-shown scenarios
- [x] 3.7 Update the existing not-yet-implemented tests for `new`, `add-dep`, `rm-dep`, `rm`, and `list` so only `make`'s test still exercises the not-yet-implemented path, per the modified `cli-shell` requirement

## 4. Final Verification

- [x] 4.1 Run the full test suite against every scenario in `specs/manifest-store/spec.md` and the modified scenarios in `specs/cli-shell/spec.md`, and verify `go test ./...` passes with no scenario left uncovered
