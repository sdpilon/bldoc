---
name: land
description: Use when a bldoc OpenSpec change (or other completed thread of work) is ready to land on main — the user says "land this", "land changes", "ship it", or asks to open/merge the PR for finished work.
---

Land a completed thread of work on `main`: branch, commit, integrate `origin/main`, run local checks, archive the OpenSpec change, open a PR, wait for CI, and squash-merge.

**Input**: Optionally specify a change name. If omitted, infer from conversation context; if ambiguous, ask.

**Steps**

1. **Branch**

   If HEAD is already on a non-`main` topic branch (e.g. from a worktree set up for this change), use it as-is. Otherwise, `git fetch origin` and create a new branch off `origin/main` named for the change (e.g. `<change-name>`).

2. **Stage only this thread's changes**

   Run `git status` and stage only the files this thread actually touched. Never `git add -A` / `git add .` — the working tree may hold unrelated in-progress files (e.g. local dogfooding scratch files) that must be left untouched.

3. **Commit**

   Use Conventional Commits (`type(scope): summary`), single line, no body unless truly necessary. Pick the scope from CLAUDE.md's list (`cli`, `manifest`, `compile`, `build`, `openspec`, `docs`) based on what changed; ask if genuinely ambiguous.

4. **Integrate `origin/main`**

   `git fetch origin && git merge origin/main`. If this produces a conflict, **stop without resolving it** — report the conflicting files and ask the user how to proceed. Do not guess at a resolution.

5. **Archive the OpenSpec change**

   Run the `openspec-archive-change` skill for this change (syncs delta specs into `openspec/specs/` and moves the change folder to `openspec/changes/archive/`) _before_ opening the PR — not as a follow-up. A deliberately-skipped `design.md` will surface as an "incomplete artifact" warning; confirm past it when the skip was a conscious call from the propose step, don't treat it as a real gap.

6. **Run local checks**

   These must all pass before opening a PR — they mirror CI exactly:

   ```
   go build ./...
   go test ./...
   go vet ./...
   gofmt -l .
   ```

   `gofmt -l .` must print nothing; any listed file means it's not CI-clean. If anything fails, stop and fix it (or report the failure) before continuing — do not open a PR on a red local check.

7. **Push and open the PR**

   Push the topic branch to `origin`. Confirm `gh` is authenticated (`gh auth status`) before the first `gh` call; if your environment requires a particular auth wrapper (e.g. a secrets-manager-backed `gh` invocation), apply it consistently to every `gh` call in this workflow.

   Write the PR body to a temp file following `.github/PULL_REQUEST_TEMPLATE.md` (`## Change`, `## Summary`, `## Test plan` sections; name the OpenSpec change under `## Change`, list the commands from step 6 under `## Test plan`), then create the PR with `--body-file <path>`. Conventional-commit-style title; no AI-tool footer or session link (per this project's CLAUDE.md — the title/body get the plain commit-message treatment, not the generic PR-footer default).

   Do not separately ask whether to monitor this PR's CI — waiting for checks through to merge is this skill's whole purpose and was already agreed to by invoking it. Just proceed to step 8.

8. **Wait for checks and merge readiness**

   Poll `gh pr checks <pr> --watch` (or repeat `gh pr checks`) until both required jobs — `Build, test, vet, fmt` and `golangci-lint` — report success. Also check `gh pr view <pr> --json mergeStateStatus,reviewDecision` for any required-review or merge-queue blocker.

   If any check fails, a required review is outstanding, or merge state is otherwise blocked: **stop**, report the exact blocker (job name, failure link, or review requirement) back to the user, and do not attempt to fix, override, or bypass it.

9. **Squash-merge**

   `gh pr merge <pr> --squash --delete-branch`. Then `git fetch origin` and confirm the squash commit is present on `origin/main`. Switch the local checkout back to `main`, pull, and delete the local topic branch (and remove its worktree, if one was used, once no longer needed).

10. **Report**

    Report back: the commit that landed on `origin/main`, the PR number/URL, and CI status — with real links, not placeholders. If landing stopped at any step, report the precise blocker instead.

**Guardrails**

- Never resolve a merge conflict against `origin/main` unilaterally — stop and ask.
- Never stage unrelated untracked files.
- Never open a PR on a failing local check.
- Never bypass, skip, or force through a failing/pending CI check or a required review.
- Confirm `gh` is authenticated before the first call, and apply your environment's auth wrapper (if any) consistently to every `gh` call.
- No AI-tool footer/session link in commit messages or PR descriptions (this project's CLAUDE.md).
