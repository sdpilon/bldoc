## Context

See proposal.md for motivation. This design revises part of manifest-io's
existing ADR ("mode is computed, not stored") rather than starting clean,
so the reasoning for departing from it matters:

That ADR's concern was a stored mode flag drifting out of sync with the
deps through some bug — computing it instead makes the invariant
self-enforcing. But the actual enforcement requirement it produced has a
gap: `add-dep`'s exclusivity check only applies "when the target has at
least one recorded dependency" — the first `add-dep` on a fresh target is
completely unconstrained, silently deciding the target's mode with no
acknowledgment. That's the specific failure mode this design closes.

The key realization: an explicit, required, `new`-time mode is not a
second source of truth shadowing the deps (the thing the original ADR
was avoiding) — it's the thing the deps get checked *against*, which is a
strictly stronger version of what `add-dep` already effectively does
today (lock in a mode from the first accepted dependency onward). It
closes the gap without reintroducing drift risk, provided declared mode
is always authoritative over inference once it exists.

Two other existing constraints still apply: target names are unvalidated
strings with no other per-target metadata beyond dependencies, and this
change must not force a migration on every existing `bldoc.toml` (many of
which have no notion of a stored mode at all).

## Goals / Non-Goals

**Goals:**
- Make mode an explicit, required choice at `new` time for every newly
  created target, closing the first-`add-dep`-is-unconstrained gap.
- Preserve existing manifests and their behavior unchanged — a target
  with no stored mode keeps today's inferred-from-deps behavior forever,
  no forced backfill.
- Let a raw-mode target's compiled output carry a chosen extension.
- Add `rename` as a minimal, manifest-only operation.

**Non-Goals:**
- Validating that a raw-mode target's concatenated bytes actually conform
  to its declared `ext` (e.g. parse as YAML). `ext` is a filename hint,
  not a compile mode.
- Migrating existing manifests to add an explicit mode to their targets.
- Changing or clearing a target's `mode` or `ext` after creation via a
  dedicated command.
- Cleaning up `.bldoc/` output on `rename` or on any other mutation.

## Decisions

**`--mode <raw|field>` is required on `new`, with no default.** An
optional flag with a default would silently re-admit exactly the
"decided by accident" failure this change exists to close — the whole
point is that mode is never decided implicitly again. Alternative
considered: optional, defaulting to `raw`. Rejected for that reason.

**A target's declared `Mode` is authoritative from the first `add-dep`
onward; a target with no declared `Mode` keeps today's inferred
behavior unchanged.** `AddDep`'s exclusivity check branches on whether
`Mode` is set: if set, every dependency (including the first) must match
it; if unset, the existing "matches every dependency already recorded,
once at least one exists" check applies exactly as it does today.
Alternative considered: infer and backfill an explicit mode onto every
existing target the first time its manifest is loaded. Rejected — it's a
write-on-read surprise for a safety improvement that doesn't fix any bug
in existing manifests, and existing manifests have no failure mode this
change needs to retroactively repair.

**`--ext` is validated against `--mode` immediately at `new` time, not
deferred to `add-dep`.** Since mode is now known at creation, `new`
rejects `--ext` together with `--mode field` right there — strictly
better than discovering the conflict later on whatever `add-dep` first
happens to be field-addressed. This also means the earlier idea of a
separate `add-dep`-time rejection for "field-addressed dep on an
ext-bearing target" is no longer needed: a declared-raw-mode target with
`ext` set already rejects any field-addressed dependency via the
mode-exclusivity check above, with no `ext`-specific rule required.

**An empty declared-field-mode target compiles to `{}`, not an empty raw
file.** Today's unconditional "empty target compiles to empty raw-mode
intermediate" rule only exists because computed-mode has no way to know a
still-empty target was meant to be field-mode. With mode declared upfront,
that ambiguity is gone, so the empty case can honor it: empty raw/
undeclared → empty raw bytes (unchanged); empty declared-field → `{}` at
`.bldoc/<target>.json`.

**`ext` is a plain optional TOML field on the target, `mode` likewise.**
Both omitted when unset — purely additive to the schema, an existing
`bldoc.toml` with neither key parses unchanged.

**`rename` never touches `.bldoc/`.** `rm` already never deletes a
removed target's compiled output — the manifest layer has no existing
responsibility for `.bldoc/` cleanup. `rename` leaving a stale
`.bldoc/<old>[.ext]` behind is consistent with that existing boundary.

**`rename` is a pure identity-field mutation, preserving `Mode` and
`Ext` untouched.** No other manifest structure references a target by
name. Validation order matches `new`/`rm`'s existing pattern: reject if
`<old>` doesn't exist, reject if `<new>` already does, otherwise rename
in place.

**`show` displays a target's mode (declared or computed) and extension
when set.** `show`'s whole purpose is surfacing a target's recorded
configuration; a mystery mode would be exactly the kind of ambiguity this
change is meant to remove.

## Risks / Trade-offs

- **`new`'s CLI surface change is a real breaking change for anyone with
  existing scripts or muscle memory** calling `bldoc new <target>` with
  no flag → Mitigation: the usage error names the missing flag plainly
  (matches every other usage-error message's existing style); this is a
  one-time adjustment, not an ongoing cost, and manifest *files* remain
  fully compatible.
- **Two different enforcement paths now exist for mode-exclusivity**
  (declared-mode targets vs. legacy inferred-mode targets) → more
  branching in `AddDep` than a single unconditional rule → Mitigation:
  the branch is a single `if Mode != ""` check; both paths are already
  independently spec'd with their own scenarios, so the fork is testable
  and explicit rather than implicit.
- **`ext` is naming-only, unvalidated** → a target could declare
  `--ext yaml` while concatenating files that aren't valid YAML →
  Mitigation: accepted, matches raw-mode's existing no-validation
  philosophy end to end; a future validating mode can be layered on
  separately without breaking this one.
- **Stale `.bldoc/` output after `rename`** (or after `rm`, already true
  today) → Mitigation: accepted as pre-existing behavior, not made worse
  by this change.
- **`mode`/`ext` immutable after creation** → getting either wrong means
  `rm` + `new` + re-adding every dependency → Mitigation: accepted;
  identical to the pre-existing cost of a wrongly-typed target name.

## Migration Plan

No manifest migration needed — `mode` and `ext` are both optional,
additive fields; a `bldoc.toml` written before this change has neither
and keeps behaving exactly as it does today. The only adjustment is at
the CLI-usage level: any script or workflow that calls `bldoc new
<target>` with no `--mode` will need updating to pass one.
