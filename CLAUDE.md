# bldoc

A CLI dev tool that tracks source-of-truth dependencies (whole files, or
specific fields within them) and compiles them into disposable,
mechanically-derived intermediates — conceptually similar to `make`,
scoped only to the "store → intermediate" tiers. See `openspec/` for the
full architecture and this project's OpenSpec changes.

## Git conventions

- Remote: GitHub, public.
- Work on each OpenSpec change happens on its own branch or worktree,
  merged to `main` via a required pull request — never a direct push to
  `main`.
- Every PR follows `.github/PULL_REQUEST_TEMPLATE.md` and names the
  OpenSpec change it implements.

## Commit messages

Conventional Commits (`type(scope): summary`), single line, no body
unless truly necessary, no AI-tool footer or session link. Scopes for
this project:

- `cli` — command surface / argument parsing
- `manifest` — manifest file format, read/write
- `compile` — the `make` compile engine (raw mode / field mode / list mode)
- `build` — module/dependency/tooling setup
- `openspec` — OpenSpec change artifacts themselves (proposal/design/specs/tasks)
- `docs` — README and other non-OpenSpec repo docs
