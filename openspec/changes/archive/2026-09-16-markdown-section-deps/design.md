## Context

See proposal.md - Why. Today's dependency resolution (`internal/compile`)
has exactly two paths: a whole-file byte read (`dep.Path == ""`), and a
structured TOML/JSON/YAML field-path walk (`dep.Path != ""`) that
requires its final value to be a scalar. Markdown section content is
neither — it's a multi-line block, and Markdown isn't a supported
structured format. This change adds a third, parallel resolution path
scoped to `.md`/`.markdown` sources, keyed off a new `Anchor` (and
`Nested` flag) on `manifest.Dep`, alongside the existing `Path`/`Field`/
`Format`.

## Goals / Non-Goals

**Goals:**
- Let a dependency address one Markdown heading's own content, or that
  heading plus its nested subsections, as a raw multi-line value.
- Keep the failure mode for an ambiguous anchor loud (a compile-time
  error), never a silent wrong-heading resolution.
- Fit into the existing raw/field/list resolution paths with minimal
  branching, reusing `resolveDepValue`'s existing whole-file/field-path
  fork rather than inventing a parallel dependency-resolution pipeline.

**Non-Goals:**
- No change to non-Markdown structured resolution (`:field-path` /
  `resolveFieldPath`) or to output encoding (JSON/YAML).
- No general "any text format with sections" abstraction — anchors are
  Markdown-only, addressed by ATX headings (`#`...`######`) specifically.
- No support for setext-style headings (`Title\n=====`) — ATX only, since
  that's what GitHub's own slug convention (which the matching rule is
  modeled on) is defined against.

## Decisions

### Grammar: `#anchor` as a third source-ref form, mutually exclusive with `:path`
`parseSourceRef` (`internal/cli/source_ref.go`) currently splits on `:`.
It gains a companion check for exactly one `#`: when present, the
segment before it is `Source`, the segment after is `Anchor`, and `Path`
must be empty (reject if the pre-`#` segment itself contains a `:`).
Rejected up front rather than letting `#` fall through to `Path` parsing,
because a Markdown section can't be scalar-resolved by `resolveFieldPath`
- keeping the two mechanisms syntactically exclusive avoids a dependency
that parses fine but can never resolve.

### Value resolution: parallel to whole-file read, not a variant of field-path
`resolveDepValue` (`internal/compile/compile.go`) gains a third branch:
`dep.Anchor != ""` resolves via a new `resolveAnchor(source, anchor,
nested)` in a new file (e.g. `internal/compile/markdown.go`), returning
a raw string exactly like the whole-file branch does - no `isScalar`
check, no tree-walk. This is why the change reuses `resolveDepValue`'s
existing fork rather than extending `resolveFieldPath`: the two return
different shapes (scalar leaf vs. raw block) for fundamentally different
source structures (key/value tree vs. document outline).

### Heading model: a flat list of (level, slug, breadcrumb, byte-range), not a tree
Parsing produces one flat, ordered list of headings, each carrying its
own nesting level (from the number of leading `#`), computed slug, and
its ancestor breadcrumb (built by tracking the current open heading at
each shallower level as the list is walked). Anchor resolution and both
capture modes are then pure operations over this flat list plus the
source's raw text: a plain anchor filters by slug equality; a breadcrumb
anchor filters by breadcrumb equality; content capture finds the next
list entry whose level is `<=` the match's level (nested) or the very
next entry regardless of level (default), and slices the raw text
between byte offsets. A real tree structure isn't needed - it would only
buy nothing here since neither matching nor capture ever needs to walk
upward from a heading (the breadcrumb is already precomputed per node,
and boundary-finding is a linear scan forward from the match).

### Ambiguity check timing: duplicated at `add-dep` (warn) and `make` (hard error), not shared code
`add-dep`'s ambiguity check and `make`'s are two separate call sites into
the same slug-matching logic (both need "does this plain anchor match >1
heading", so they can share the matching helper), but they react
differently on purpose: `add-dep` warns and proceeds (per proposal.md -
what changes), `make` always hard-errors on an ambiguous plain anchor,
regardless of what `add-dep` printed or whether `add-dep` even ran
against the current file content (the manifest can be hand-edited, or
the file can change between `add-dep` and `make` - see the
`manifest-store` capability's existing "malformed manifest" handling for
the established precedent that `make` never trusts unvalidated state).
`add-dep`'s check is best-effort convenience; it is not a substitute for
`make`'s validation.

### `--nested` gating happens at the CLI/manifest layer, not re-checked by `make`
`add-dep` rejects `--nested` without an anchor or on a raw-mode
(no-`:field`) dependency (cli-shell capability) before anything is ever
recorded. `make` does not re-validate this combination — it already
trusts other structurally-enforced invariants from a valid manifest
(e.g. raw/field mode exclusivity is never re-checked at compile time
either). A hand-edited manifest that violates this is already outside
the guarantees `make` makes for any other recorded-but-invalid
combination.

## Risks / Trade-offs

- **[Risk]** A plain anchor that's unambiguous today can become ambiguous
  later if a new colliding heading is added elsewhere in the file,
  breaking a previously-working dependency. → **Mitigation**: this is a
  loud `make`-time error naming the anchor and suggesting a breadcrumb,
  not a silent misresolution - the explicit trade-off accepted over
  GitHub-style auto-numbering, which would misresolve silently instead.
- **[Risk]** GitHub's own slug algorithm has edge cases (duplicate raw
  text handling, emoji, non-ASCII) that a from-scratch reimplementation
  could get subtly wrong, producing anchors that don't match what a user
  copies from a rendered GitHub page. → **Mitigation**: scope the initial
  implementation to the common case (ASCII letters/digits/hyphens/spaces,
  punctuation stripped) and treat exact GitHub parity as a follow-up if
  it turns out to matter in practice - this repo's own specs (the
  primary intended use case) don't exercise the exotic cases.
- **[Trade-off]** `add-dep` reading the filesystem for `#anchor`
  specifically breaks its otherwise-universal "pure manifest bookkeeping,
  no I/O" rule. Accepted because the alternative (deferring all feedback
  to `make`) makes the most likely mistake (a mistyped or duplicate
  anchor) invisible until compile time, and the check is explicitly
  best-effort (silently skipped if the file can't be read) rather than a
  hard dependency on file access at `add-dep` time.
