## Why

`bldoc` has no code yet. A prior exploration session worked out the tool's
full intended architecture (a committed manifest tracking per-target
dependencies, a `make` step that compiles them into a disposable
intermediate, strict raw-mode/field-mode output separation, a minimal
value+raw intermediate schema) but building all of that as one change
would bundle CLI-surface decisions (language, argument conventions, error
handling) together with manifest-format and compile-engine decisions that
are largely independent of each other. Staging the work lets each layer's
open questions get resolved on their own rather than all at once.

## What Changes

- Introduce the `bldoc` CLI entry point/executable.
- Add command stubs for `new`, `add-dep`, `rm-dep`, and `make` (final
  command names/shapes are this change's own design question) that parse
  their expected arguments but do not yet read or write a manifest, track
  dependencies, or compile anything.
- Establish baseline CLI conventions: help output, handling of
  unknown commands/flags, exit codes, and a consistent response for a
  recognized-but-not-yet-implemented command.
- Explicitly defer to a later change: the manifest file format and
  location, actual dependency tracking, target compilation (raw mode and
  field mode), and the intermediate output schema.

## Capabilities

### New Capabilities

- `cli-shell`: the `bldoc` command-line surface — command names, argument
  parsing, help/error/exit-code conventions, and stub ("not yet
  implemented") behavior for `new`, `add-dep`, `rm-dep`, and `make` —
  prior to any real manifest or compilation behavior existing behind it.

### Modified Capabilities

(none — this is a new project with no existing capabilities)

## Impact

Adds a new CLI executable to what is currently an empty repository.
Language/runtime/argument-parsing library are not yet chosen; that
choice, along with the exact command/flag shapes, belongs to this
change's own `design.md`, not to this proposal.
