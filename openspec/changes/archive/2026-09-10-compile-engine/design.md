## Context

See `proposal.md` for motivation. `manifest-io` already provides
`internal/manifest` (`Load`, `Save`, `AddTarget`/`AddDep`/etc.) and the
`Manifest`/`Target`/`Dep` schema — `Dep.Source`, `Dep.Path` (field-path
into the source, empty for whole-file), `Dep.Field` (target-side field
name, empty for raw mode), `Dep.Format` (optional template). `make`'s
`RunE` in `internal/cli/make.go` currently just calls `notImplemented`.

## Goals / Non-Goals

**Goals:**
- Give `make` real behavior for both raw and field mode, per the
  `manifest-store` capability's existing mode-exclusivity guarantee (a
  target's deps are never mixed).
- Define exactly where the compiled intermediate lands on disk.

**Non-Goals:**
- Any drift/staleness detection, caching, or incremental compilation —
  `make` always fully recomputes (architecture vision: "idempotent, no
  drift detection").
- Fact-checking the intermediate against its sources after the fact
  (a possible future `bldoc verify`).
- Per-fact provenance metadata in the intermediate beyond `value`/`raw`
  (already decided against in the architecture vision).
- Partial-file/snippet addressing beyond a single scalar field-path
  (e.g. extracting one array element, or a whole nested table) — parked,
  same as `manifest-io`.
- Validating that a dependency's format template is a well-formed
  printf-style string with exactly one `%s` — a malformed template just
  produces whatever `fmt.Sprintf` yields.
- Managing the consuming project's `.gitignore` to exclude `.bldoc/` —
  left to the user, consistent with the intermediate being "disposable,
  uncommitted" rather than something bldoc enforces.

## Decisions

**Intermediate location: `.bldoc/<target>` (raw) or `.bldoc/<target>.json` (field), cwd-relative.**
Mirrors the manifest's own cwd-relative convention. A fixed, predictable
per-target path is simpler than a configurable output directory, and
nothing in the architecture vision calls for configurability here.

**Path-less field-mode dependency uses the whole source file as its raw value.**
`Dep.Path` and `Dep.Field` are independent: a field-mode dependency
(`Dep.Field != ""`) can have `Dep.Path == ""`, meaning "use this whole
file's content as the field's raw value" rather than a structured
field-path. This keeps one dependency-resolution rule for both cases:
raw = file content when `Path == ""`, or the field-path's resolved value
when `Path != ""`. No separate flag or special-case is needed to
distinguish them.

**Structured parsing is format-detected by file extension.** `.toml` →
TOML, `.json` → JSON, `.yaml`/`.yml` → YAML, anything else → error. Only
triggered when a dependency has a field-path (`Path != ""`); a
whole-file dependency never needs to parse its source, regardless of
extension. Each format is parsed into a generic `map[string]interface{}`
tree (via each library's own "decode into `interface{}`" mode) so one
field-path resolver walks all three uniformly.

**Field-path is a dot-separated sequence of map keys, scalars only.**
Consistent with the address grammar already established for `source:path`
in `bootstrap-cli` (dots inside a path segment, e.g.
`project.requires-python`, are literal key characters, not further
delimiters — splitting on `.` is unambiguous here since TOML/JSON/YAML
keys in this project's use cases don't themselves contain dots). A
resolved table/array is rejected rather than silently stringified,
since serializing an arbitrary nested structure into the "raw" string
field would blur the raw/field-mode intermediate schema's simplicity.

**Duplicate field names are a new compile-time check, not a `manifest-io`
fix.** `manifest-io`'s `AddDep` only rejects a duplicate `(source, path)`
pair, so two deps naming the same target-side field from different
sources are representable in the manifest today. Rather than reopening
`manifest-io` to add a stricter `AddDep` check (which would need its own
migration story for any manifest that already has this shape), `make`
itself validates field-name uniqueness across a target's dependencies
before compiling — a natural fit here, since JSON emission is the actual
place this ambiguity would otherwise be silently resolved by last-write-wins.

**A target with no dependencies compiles to an empty raw file.** A
target's mode is normally computed from its dependencies (per
`manifest-io`'s design); with zero dependencies there's nothing to
compute from. Treating the empty case as raw mode (an empty byte
string) avoids inventing a third mode or an error for a legitimately
empty target (e.g. one just created via `new` and not yet built out).

**New package: `internal/compile`.** Owns raw-mode concatenation,
field-mode resolution (source parsing, field-path walk, format
templating), duplicate-field validation, and intermediate
serialization — one function per target given its `manifest.Target`,
returning either raw bytes or a JSON-serializable map plus an error.
`internal/cli/make.go` becomes a thin wrapper: load the manifest,
resolve target(s), call into `internal/compile`, write the result under
`.bldoc/`.

**New dependency: a YAML library (e.g. `gopkg.in/yaml.v3`).** TOML
support already exists from `manifest-io`; JSON uses `encoding/json`
from the standard library. YAML has no standard-library support, so
this is the one new external dependency this change needs.

## Risks / Trade-offs

- **Extension-based format detection is fragile if a source file's
  extension doesn't match its actual content** (e.g. a `.json` file
  that's actually YAML) → Mitigation: this matches common tooling
  convention (most editors/tools pick a parser by extension too); a
  mismatch surfaces immediately as a "malformed source file" error
  rather than silently misparsing.
- **Rejecting non-scalar field-path results forecloses some real use
  cases** (e.g. wanting a whole nested table rendered as an inline
  snippet) → Mitigation: explicitly parked in the architecture vision;
  revisit only if a real change needs it, rather than speculatively
  building it now.
- **`.bldoc/` is created without any `.gitignore` guidance**, so a user
  could accidentally commit generated intermediates → Mitigation:
  accepted as a non-goal; a future change could have `bldoc new` (or a
  dedicated init step) offer to add a `.gitignore` entry if this proves
  to be a real recurring papercut.
