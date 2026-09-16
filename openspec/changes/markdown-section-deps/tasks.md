## 1. Source-ref and flag grammar

- [ ] 1.1 Extend `parseSourceRef` (`internal/cli/source_ref.go`) to recognize `source#anchor` (single-segment or `parent/child` breadcrumb), rejecting a source-ref that combines `:field-path` and `#anchor`; verify with new cases in `internal/cli/source_ref_test.go`
- [ ] 1.2 Add the `--nested` flag to `add-dep` (`internal/cli/add_dep.go`), rejecting it when the source-ref has no `#anchor` and when the target-ref has no `:field`; verify with new cases in `internal/cli/add_dep_test.go`
- [ ] 1.3 Add `--nested` rejection to `rm-dep` (mirroring its existing `--format` rejection) in `internal/cli/rm_dep.go`; verify with a new case in `internal/cli/rm_dep_test.go`

## 2. Manifest schema and `add-dep` recording

- [ ] 2.1 Add `Anchor` and `Nested` fields to `manifest.Dep` and its TOML round-trip; verify with a round-trip test in `internal/manifest/io_test.go`
- [ ] 2.2 Include `Anchor` (not `Nested`/`Format`) in `add-dep`'s duplicate-dependency identity check; verify with new cases in `internal/manifest/ops_test.go` covering a same-anchor duplicate with differing `--nested`
- [ ] 2.3 Update `show`'s dependency listing to print a dependency's anchor and whether `--nested` was recorded; verify with new cases in `internal/cli/show_test.go`

## 3. Markdown heading engine

- [ ] 3.1 Implement an ATX heading parser (level, raw title, byte offset of line end) in a new `internal/compile/markdown.go`; verify with unit tests covering `#` through `######` and non-heading `#` occurrences (e.g. inside code fences)
- [ ] 3.2 Implement GitHub-style slug computation (lowercase, spaces to hyphens, punctuation stripped); verify with unit tests including this repo's own duplicate-heading cases (`openspec/specs/cli-shell/spec.md`'s and `openspec/specs/manifest-store/spec.md`'s duplicate `Scenario:` headings)
- [ ] 3.3 Compute each heading's breadcrumb (root-first ancestor slug chain) while parsing; verify with unit tests against nested `##`/`###`/`####` fixtures
- [ ] 3.4 Implement anchor matching: plain-anchor lookup by slug (returning not-found, unique match, or ambiguous), and breadcrumb lookup by ancestor-chain equality; verify with unit tests covering all three plain-anchor outcomes and breadcrumb disambiguation of a genuinely duplicate slug
- [ ] 3.5 Implement content-range resolution for both capture modes — next heading at any level (default) and next heading at an equal-or-shallower level (`--nested`), including end-of-file as the boundary when there is no next heading; verify with unit tests covering both modes, including a nested capture that spans multiple sibling subsections

## 4. Compile-engine integration

- [ ] 4.1 Add a third branch to `resolveDepValue` (`internal/compile/compile.go`) dispatching `dep.Anchor != ""` to the new Markdown resolver, rejecting sources whose extension isn't `.md`/`.markdown`; verify with new cases in `internal/compile/compile_test.go`
- [ ] 4.2 Verify anchor resolution end-to-end for raw mode (concatenation), field mode (`{"value", "raw"}` wrapping and `--format`), and list mode (record-field resolution) with new cases in `internal/compile/compile_test.go` and `internal/compile/compile_list_test.go`

## 5. Ambiguity warning at `add-dep`

- [ ] 5.1 Implement `add-dep`'s best-effort ambiguity check: for a plain (non-breadcrumb) anchor, read and parse the source file, print a non-blocking warning naming a breadcrumb suggestion if the slug matches more than one heading, and silently skip the check if the file can't be read; verify with new cases in `internal/cli/add_dep_test.go` covering the warned, clean, breadcrumb-skips-check, and unreadable-file cases

## 6. Whole-suite verification

- [ ] 6.1 Run `go test ./...` and confirm all packages pass
- [ ] 6.2 Dogfood manually against this repo's own `openspec/specs/*/spec.md` (which has real duplicate headings) with `bldoc add-dep`, `bldoc make`, and `bldoc show`, confirming the ambiguity warning, breadcrumb disambiguation, and both capture modes behave as specced
- [ ] 6.3 Sweep `examples/demo/README.md` and any other docs describing `add-dep`'s source-ref grammar for staleness against the new `#anchor`/`--nested` forms, updating them in the same PR if any are found
