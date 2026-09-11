## 1. CI Workflow

- [x] 1.1 Add `.github/workflows/ci.yml` triggered on pull requests and pushes to `main`, running `go build ./...`, `go test ./...`, `go vet ./...`, and a `gofmt -l` check that fails if any file is unformatted, and verify it passes on a local dry run of each command
- [x] 1.2 Add a `golangci-lint` job to the same workflow (or a separate one) using the official `golangci-lint-action`, with a `.golangci.yml` config setting a reasonable default linter set for this project
- [x] 1.3 Push the workflow on a branch and open a PR, and verify all jobs actually run and pass in GitHub Actions on that PR

## 2. Documentation

- [x] 2.1 Update the project's `CLAUDE.md` if it documents anything about verification/testing that this workflow now automates or supersedes, and verify no other repo doc references a manual-only check that CI now covers
