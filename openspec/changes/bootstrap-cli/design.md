## Context

See proposal.md - Why. This change stands up `bldoc`'s CLI command surface
only; no manifest file, dependency tracking, or compilation exists yet.

## Goals / Non-Goals

**Goals:**
- A working `bldoc` binary implementing the full command grammar for all
  seven verbs, including the `target[:field]` / `source[:path]` addressing
  syntax, as pure argument parsing.
- Consistent, predictable behavior for recognized-but-unimplemented
  commands, plus normal `--help`/`--version` output.

**Non-Goals:**
- Reading, writing, or otherwise touching a manifest file.
- Any dependency tracking or compilation (raw-mode or field-mode).
- Choosing the manifest's file format or location.

## Decisions

**Go + Cobra.** Cobra is the standard for multi-verb Go CLIs (`kubectl`,
`gh`, `hugo`), and gives per-command flag sets plus help/usage/error
conventions for free — a saving that compounds as more verbs are added in
later changes. Alternative considered: stdlib `flag` with manual
dispatch — zero dependencies, more consistent with this project's general
minimalism, but would mean hand-rolling and re-verifying help/usage/error
formatting by hand as the verb set grows. Cobra was judged worth its one
dependency.

**Command surface**: `new`, `add-dep`, `rm-dep`, `make`, `rm`, `list`,
`show`.
- `new <target>`
- `add-dep <target-ref> [--format <template>] <source-ref>`
- `rm-dep <target-ref> <source-ref>`
- `make [<target>]` — omitted target means "every target in the manifest"
- `rm <target>`
- `list`
- `show <target>`

`<target-ref>` is `target` or `target:field`; `<source-ref>` is `source`
or `source:path`. An unnamed target-ref implies raw (whole-file) mode; a
`:field` target-ref implies field mode — this grammar mirrors the
raw-mode/field-mode split from the tool's overall design.

**`--format` requires a field-addressed target-ref.** Format templates
only make sense for named fields (field mode); raw mode concatenates
bytes verbatim. This is a static argument-combination check — it needs no
manifest state — so it belongs in this change rather than a later one.

**Stub behavior**: any recognized command given syntactically valid
arguments prints `bldoc <verb>: not yet implemented` to stderr and exits
`1`. This distinguishes "understood but not built yet" from a usage error
(malformed arguments) or success.

**Layout**: `cmd/bldoc/main.go` as the entrypoint, command definitions
under `internal/`. Standard idiomatic Go CLI layout.

**Module path**: `module bldoc` — not a GitHub path, since the repo stays
local-only for now and the module isn't meant to be imported by other Go
code. Trivially renamable later via `go mod edit -module` if that changes.

**Testing**: invoke each Cobra command's `Execute()` in-process and assert
exit code/output, rather than shelling out to a built binary.

## Risks / Trade-offs

- Taking on Cobra now, before any real functionality exists → mitigated
  by it being one stable, widely-used dependency rather than a deep tree,
  with the payoff growing as more verbs land in later changes.
- Locking in the addressing grammar before the manifest/compile engine
  exists risks a mismatch if a later change's design surfaces a reason to
  change it → mitigated by grammar parsing living in one small, easily
  revisited place, fully separate from manifest logic.

## Open Questions

- Exact wording of the "not yet implemented" message and of help/version
  text is left to implementation — doesn't affect the spec, which only
  requires that such a message exists.
