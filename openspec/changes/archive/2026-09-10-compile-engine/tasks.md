## 1. Source Parsing Foundations

- [x] 1.1 Add a YAML dependency (e.g. `gopkg.in/yaml.v3`) and verify it resolves cleanly in `go.mod`/`go.sum`
- [x] 1.2 Create the `internal/compile` package with a source parser that decodes TOML/JSON/YAML by file extension into a generic map, and a field-path resolver (dot-separated key walk, scalar-only results), with unit tests covering all three formats, a missing key, a non-scalar result, a malformed file, and an unsupported extension

## 2. Compilation Core

- [x] 2.1 Implement raw-mode compilation — byte-for-byte concatenation of a target's whole-file dependencies in added order — with unit tests including the zero-dependency (empty output) case
- [x] 2.2 Implement field-mode compilation — resolve each dependency's raw value (whole source file when path-less, else the field-path's resolved value), apply the format template when present, assemble `{field: {"value": ..., "raw": ...}}` — with unit tests covering a whole-file field dependency, a structured field-path dependency, and the no-format-template `value == raw` case
- [x] 2.3 Implement duplicate-field-name validation across a target's recorded dependencies, with a unit test covering the rejection case

## 3. CLI Wiring

- [x] 3.1 Wire `make <target>` to compile the target and write its intermediate under `.bldoc/` (`.bldoc/<target>` for raw mode, `.bldoc/<target>.json` for field mode); verify the unknown-target rejection and both successful-compilation scenarios from `specs/compile-engine/spec.md`
- [x] 3.2 Wire `make` with no target argument to compile every target recorded in the manifest; verify the all-targets-compiled and no-targets-no-op scenarios
- [x] 3.3 Update `make_test.go`'s existing not-yet-implemented assertions to reflect real compiled behavior, consistent with the modified `cli-shell` requirement (the not-yet-implemented requirement itself is removed)

## 4. Final Verification

- [x] 4.1 Run the full test suite against every scenario in `specs/compile-engine/spec.md` and the modified `cli-shell` scenarios, and verify `go test ./...`, `go vet ./...`, and `gofmt -l .` are all clean with no scenario left uncovered
