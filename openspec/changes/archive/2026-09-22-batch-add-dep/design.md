## Context

`add-dep` currently takes exactly two positional arguments — a
target-ref (`target` or `target:field`, or `target:record.field` on a
list-mode target) and a source-ref (`source`, `source:path`, or
`source#anchor`) — plus `--format` and `--nested` flags. `internal/cli`
parses those two arguments and calls `manifest.AddDep` once;
`internal/manifest.AddDep` validates and appends a single `Dep` to one
target, then the caller calls `manifest.Save` once. See proposal.md -
Why for the motivating case (a list-mode record with several fields)
and for why that batch form also forced a symbol rename: `:` inside a
ref has only ever meant "address a component within this ref," and the
batch form needs a distinct "map this field to that source" relationship
that would otherwise collide with it.

## Goals / Non-Goals

**Goals:**
- Let one `add-dep` invocation record several dependencies against one
  target (and, for a list-mode target, one record), without changing
  anything about how a single dependency is recorded today beyond the
  symbol rename below.
- Move ref-internal field addressing from `:` to `@` everywhere it
  appears (target-ref, source-ref field-path), freeing `:` to mean
  "field maps to source" in the batch form's pairs, with no ambiguity
  between the two roles.
- Reuse `manifest.AddDep`'s existing per-dependency validation as-is,
  called once per pair, rather than writing a parallel validation path.

**Non-Goals:**
- Adding a way to batch-add dependencies across more than one target in
  a single call.
- Per-pair `--format`/`--nested` overrides within the batch form (see
  proposal.md - What Changes).
- Any change to `rm-dep`'s argument count, `show`'s output shape, or the
  compile engine — this change is scoped to how dependencies get
  addressed and recorded, not how they're removed, displayed, or
  compiled. `rm-dep` inherits the `@` rename automatically because it
  reuses the same ref parsers.

## Decisions

### Rename the ref-internal field separator from `:` to `@`
`target:field`, `source:path`, and `target:record.field` all become
`target@field`, `source@path`, `target@record.field`. `#anchor` (and its
`/`-separated breadcrumb form) is untouched — it already meant "a
location within this source," the same role `@` now plays for
structured field-paths, so the two symbols already have compatible,
un-confusable meanings (`Path` and `Anchor` on `manifest.Dep` remain
mutually exclusive, unchanged by this rename).

Alternatives considered:
- **Do nothing; give the batch form's pairs a symbol other than `:`**
  (e.g. `field=source-ref`, using `=`). This would avoid a breaking
  change to the already-shipped `add-dep`/`rm-dep` syntax entirely.
  Rejected per the user's explicit preference: `:` reads as "key maps
  to value" (as in a key:value store), which is the relationship the
  batch pairs actually have, and forcing that intuition onto a
  different symbol just because `:` was taken for something else was
  the original source of the confusion this change is fixing.
- **`.` for field addressing** (`target.field`, `source.path`).
  Rejected: field-paths already use `.` internally to walk nested keys
  (`project.requires-python`), and file names routinely contain `.`
  too (`pyproject.toml`), so a bare dot can't unambiguously mark where
  "which file" ends and "which key inside it" begins.
- **`@` for field addressing** (chosen). No collision with anything
  already in the grammar: not used in file paths in practice, not used
  in field-path keys, not used in anchors/breadcrumbs. Reads naturally
  as "at" (`README@version` — the `version` slot of `README`;
  `pyproject.toml@project.version` — the `project.version` key of
  `pyproject.toml`).

### Dispatch on argument count, not a new flag or subcommand
`add-dep` becomes a variadic command (`cobra.MinimumNArgs(2)` in place
of `exactArgs(2)`). Exactly two arguments always means the existing
target-ref/source-ref form; three or more always means the new form.

Alternative considered: a `--batch` flag or a separate `add-deps`
command. Rejected because it forces the user to know up front they'll
want more than one dependency, and because argument-count dispatch has
no real ambiguity to resolve — the existing two-argument grammar and the
new pair grammar parse the same raw strings differently only when
exactly two arguments are given, which is precisely the case reserved
for the old form.

### The target argument must be bare in the variadic form
In the variadic form, `args[0]` is checked for `@` and rejected outright
if present, rather than being parsed as a target-ref and having its
`Field` (if any) silently ignored. A target-ref's `@field` has no
meaning in this form — every pair carries its own field — so a `@field`
on the target argument is almost certainly a mistake (e.g. muscle memory
from the two-argument form) and should fail loudly rather than silently
drop the field.

### Record-name argument only for list-mode targets
Whether `args[1]` is a record name or the first pair depends on the
target's declared mode, which is only known after loading the manifest.
Implementation order in `add_dep.go`'s `RunE`:
1. If `len(args) == 2`: existing code path, unchanged.
2. Else: validate `args[0]` has no `@`, load the manifest, look up the
   target (unknown-target error if missing — same error as today, just
   raised earlier), and branch on `target.Mode`:
   - `"list"`: `args[1]` is the record name (reject if it contains `@`,
     `:`, or `.`); `args[2:]` are pairs, parsed with the record name
     fixed.
   - anything else (`"raw"`, `"field"`, or `""` for an undeclared-mode
     legacy target): `args[1:]` are pairs, no record name.
3. Reject if fewer than two pairs remain after removing the target and
   any record-name argument.

A target with no declared `Mode` (a pre-existing manifest predating the
mode-declaration requirement) is treated the same as `raw`/`field` here:
no record-name argument. Each pair is still parsed (field-vs-bare) by
the rule below, and `AddDep`'s existing undeclared-mode inference (the
first dependency recorded sets the target's effective mode) governs
exactly as it does for two `add-dep` calls in sequence — no new logic
needed.

### Pair parsing: split on the first `:`, unconditionally
Because a source-ref can never itself contain `:` (field-path addressing
now uses `@`), a pair's shape is decided by simple presence of `:`, with
no need to know the target's mode first:

- A pair containing `:` is split on its *first* occurrence into `field`
  (everything before) and a remainder, which is parsed as a source-ref
  with the existing `parseSourceRef` (itself may still contain one
  `@path` or one `#anchor`, exactly as today).
- A pair containing no `:` is parsed directly as a source-ref — a
  whole-file/anchor dependency, or, within a list-mode record, a
  record-only dependency.

Whether a given pair's shape is *allowed* for the target's mode is not
decided here at all — it's decided by the same mode-exclusivity and
list-mode addressing checks `manifest.AddDep` already runs for a single
dependency, called once per pair. A bare pair against a field-mode
target, for instance, is rejected the same way a bare `add-dep <target>
<source>` call is rejected against a field-mode target today: by
`AddDep`, not by the CLI's pair parser. This removes an entire class of
CLI-side "what does this pair mean for this mode" branching that a
mode-aware parser would otherwise need — the parser only ever answers
"does this pair have a field," and validity is `AddDep`'s job, as it
already is.

This mirrors today's target-ref parsing (`ref.go`), which already
splits `target@field` the same way (previously `target:field`), and
reuses `parseSourceRef` unchanged — only a new call site loops over
`args[n:]` instead of parsing one source-ref.

### Atomicity via "mutate in memory, save once" — no working-copy needed
`manifest.AddDep` is not changed. The CLI loads the manifest once,
calls `manifest.AddDep` in a loop (one call per pair, against the same
`*manifest.Manifest` value), and calls `manifest.Save` exactly once,
only after every pair has succeeded. If any call fails, the function
returns the error immediately without calling `Save` — the on-disk
`bldoc.toml` was never touched, so it's automatically left exactly as
it was, with no need to snapshot or roll back an in-memory copy. This is
exactly how today's single-pair `add-dep` already behaves; the batch
form only changes the number of `AddDep` calls made before the one
`Save`.

Alternative considered: validate all pairs against a cloned manifest
first, then apply and save. Rejected as unnecessary complexity — the
loop-then-save-once approach gives the same atomicity guarantee for
free, since nothing is persisted until the end regardless.

### `--format`/`--nested` are invocation-level flags
Both flags are read once per `add-dep` invocation (as today) and passed
unchanged into every `AddDep` call the loop makes. No new flag syntax is
introduced for per-pair overrides.

## Risks / Trade-offs

- **Breaking change to every existing `add-dep`/`rm-dep` invocation
  using `:field`/`:path`.** Any script, alias, or muscle-memory
  invocation using the old `target:field`/`source:path` syntax breaks
  immediately once this ships — there is no transition period or
  dual-syntax support. → Accepted deliberately: this is a personal,
  pre-1.0 dev tool with no external users, and the whole point of this
  change is to fix a notation the user already found confusing while
  it's cheap to fix (before more surface area is built on top of it).
  No dual-syntax fallback is planned; `parseTargetRef`/`parseSourceRef`
  simply use `@` going forward.
- **A user forgets the record name on a list-mode batch call and gets a
  confusing error.** E.g. `bldoc add-dep PROJECTS description:a.md
  path:b.yaml` on a list-mode target parses `description` as the record
  name (accepted — it's a bare word with no `@`/`:`/`.`) and
  `path:b.yaml` as the only remaining pair, then fails the "at least two
  pairs" check. → Mitigation: the usage error for "fewer than two
  pairs" should name what was parsed (record name and pair count) so
  the missing-record-name mistake is diagnosable from the error text,
  not just a bare "need 2 pairs" message. This is a task-level detail,
  not a spec requirement.
- **Partial-success expectations.** A user might expect a batch call to
  record whichever pairs were valid and just report the rest as errors
  (a "best effort" batch). This change deliberately does the opposite
  (all-or-nothing) to keep `add-dep`'s existing "manifest unchanged on
  error" guarantee intact rather than introducing a new, weaker
  guarantee for the batch form specifically. → Documented explicitly in
  proposal.md and the cli-shell/manifest-store delta specs; no
  mitigation needed beyond that, since this is the intended behavior.
