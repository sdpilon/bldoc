## 1. Implementation

- [x] 1.1 In `internal/cli/root.go`, read `runtime/debug.ReadBuildInfo()` when building the root command and, when VCS settings (`vcs.revision`, `vcs.time`, `vcs.modified`) are present, append a parenthetical to the reported version string: short SHA (first 12 chars of `vcs.revision`), the `vcs.time` value, and `modified` when `vcs.modified == "true"`; when VCS settings are absent, leave the version string unchanged; verify via unit test covering all three cases (full info present, `modified` true, no VCS info at all)
- [x] 1.2 Verify `bldoc --version`'s existing help/version tests still pass unchanged, and add a test asserting the plain (no-VCS-info) fallback matches today's exact output format

## 2. Verification

- [x] 2.1 Run `go build ./...`, `go vet ./...`, `gofmt -l .`, and `go test ./...`, confirming all clean
- [x] 2.2 Manually verify against a real build: `go install ./cmd/bldoc` from this branch, then `bldoc --version`, confirming the printed commit SHA/timestamp match `go version -m $(which bldoc)`'s `vcs.revision`/`vcs.time`
