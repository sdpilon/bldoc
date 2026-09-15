## 1. Manifest layer: `list` mode and `record[.field]` addressing

- [x] 1.1 Extend the accepted `--mode` values from `raw`/`field` to include `list` wherever mode is validated, and verify via unit test that `bldoc new PROJECTS --mode list` records mode `list`, while an unrecognized mode is still rejected
- [x] 1.2 Add a `splitRecordField(field string) (record, field string, ok bool)` helper in `internal/manifest` that splits on `.`: zero dots means record-only (`field == ""`), one dot means `record.field`, and empty input or more than one dot means `ok == false`; verify via unit tests covering all four cases
- [x] 1.3 Extend `AddDep`'s mode branch (`internal/manifest/ops.go:65`) with a `case "list"` that rejects any dependency whose `Field` fails `splitRecordField`; verify via unit tests matching the `manifest-store` delta spec's addressing scenarios (bare ref rejected, record-only accepted, record.field accepted, over-qualified rejected)
- [x] 1.4 In that same `case "list"` branch, reject a dependency carrying `--format` when its split field has no `.` (record-only); verify via unit tests matching the `manifest-store` delta spec's format-restriction scenarios (rejected on record-only, accepted on record.field)

## 2. CLI layer: flag acceptance

- [x] 2.1 Add `list` to `new`'s accepted `--mode` values (`internal/cli/new.go`); verify via test that `bldoc new PROJECTS --mode list` succeeds and other invalid values are still rejected
- [x] 2.2 Remove `new`'s `--ext`-rejected-with-`--mode field` branch so `--ext` is accepted unconditionally regardless of mode; verify via unit tests that `--mode field --ext yaml` and `--mode list --ext yaml` are both accepted and recorded, per the `cli-shell` delta spec

## 3. Compile layer: resolving list-mode dependencies

- [x] 3.1 Extract the scalar-vs-table/array check already inside `resolveFieldPath`'s switch (`internal/compile/source.go:63`) into a standalone `isScalar(v interface{}) bool`, and reuse it in place; verify existing field-path scalar-rejection tests still pass unchanged
- [x] 3.2 Add a whole-document validator that applies `isScalar` to every top-level value of a parsed document, rejecting the first non-scalar found; verify via unit test with a document containing a nested table/array top-level value (rejected) and an all-scalar document (accepted)
- [x] 3.3 Implement `compileList(target manifest.Target) (Result, error)` in `internal/compile`: for each dependency, split its target-side field via the manifest helper from 1.2; a record-only dependency parses its source with the existing `parseSource` and merges its validated top-level keys into that record; a record.field dependency resolves exactly like today's field-mode dependency (whole-file content or field-path, optional `--format`) into that one field; verify via unit tests covering whole-document merge, single-field resolution with a format template applied, and malformed/unsupported-format source rejection for a record-only dependency, matching the `compile-engine` delta spec's scenarios
- [x] 3.4 Within `compileList`, track record order by first appearance and detect per-record duplicate field names (record-only merge colliding with itself, with another record-only merge, or with an explicit record.field dependency), erroring on the first collision; verify via unit tests matching the `compile-engine` delta spec's duplicate-field scenarios, including that the same field name in two different records is allowed
- [x] 3.5 Wire `compile.Target` to dispatch to `compileList` when the target's declared `Mode` is `list` (alongside today's raw/field dispatch); verify via unit test that a list-mode target with zero dependencies compiles to an empty list-mode result

## 4. Compile layer: serialization

- [x] 4.1 Add `IsList bool`, `RecordOrder []string`, and `Records map[string]map[string]interface{}` to `compile.Result`; verify via unit test that `compileList`'s return value populates all three correctly for a multi-record target
- [x] 4.2 Implement `compile.Encode(result Result, ext string) ([]byte, error)`: pass `Raw` through unchanged for a raw-mode result; for field- or list-mode results, encode as JSON by default or as YAML (via `yaml.v3`) when `ext` is `yaml`/`yml` case-insensitively — building the ordered `[]map[string]interface{}` from `RecordOrder`/`Records` for list mode; verify via unit tests covering JSON and YAML output for both field- and list-mode results, parsing the YAML output back with a real parser to confirm it's valid
- [x] 4.3 Update `internal/cli/make.go`'s `compileTarget` to call `compile.Encode` instead of its inline `json.MarshalIndent`, and to resolve the intermediate's filename suffix from `Ext` for field/list-mode targets (not just raw-mode ones), per the `compile-engine` delta spec's file-naming rules; verify via an integration test that `bldoc make` on a field-mode target created with `--ext yaml` writes `.bldoc/<target>.yaml` containing valid YAML

## 5. Verification and doc sweep

- [x] 5.1 Run `go test ./...` and confirm no regressions in existing raw-mode or field-mode behavior for targets that don't declare `--mode list`
- [x] 5.2 Add or confirm CLI-level tests in `internal/cli` covering the `cli-shell` delta spec's new scenarios (`--mode list` accepted, `--ext` accepted with any mode); verify via `go test ./internal/cli/...`
- [x] 5.3 Sweep `examples/demo/README.md` and `CLAUDE.md` for any description of `--mode`/`--ext` that's now stale given list mode's addition, and update it in this same change rather than as a follow-up; verify by re-reading the affected sections against the shipped behavior
- [x] 5.4 Run the same checks CI runs (`.github/workflows/ci.yml`'s `checks` job: build, vet, fmt check, test) and confirm they pass locally before opening the PR
