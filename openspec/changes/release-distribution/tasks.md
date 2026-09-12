## 1. Version injection

- [x] 1.1 Change `internal/cli/root.go`'s `version` from a `const` to a package-level `var version = "0.0.0-dev"`, and verify `go build ./... && go vet ./...` still pass
- [x] 1.2 Verify a plain local build still reports the default: `go build -o /tmp/bldoc-local ./cmd/bldoc && /tmp/bldoc-local --version` prints `bldoc version 0.0.0-dev`

## 2. GoReleaser configuration

- [x] 2.1 Add `.goreleaser.yaml` at the repo root defining a build matrix of `{linux, darwin, windows} x {amd64, arm64}` for `./cmd/bldoc`, with `ldflags` templated as `-s -w -X bldoc/internal/cli.version={{.Version}}`
- [x] 2.2 Configure archive naming/format in `.goreleaser.yaml`: `.tar.gz` for linux/darwin, `.zip` for windows
- [x] 2.3 Enable GoReleaser's checksum file generation covering every archive
- [x] 2.4 Verify locally with `goreleaser release --snapshot --clean` (no tag push, no publish) and confirm `dist/` contains six archives (one per platform/arch) plus a checksums file
- [x] 2.5 Verify a built snapshot binary reports the snapshot version when run with `--version` (e.g. extract the linux/amd64 archive and run `./bldoc --version`)

## 3. Release workflow

- [x] 3.1 Add `.github/workflows/release.yml` triggered on `push: tags: ['v*']` only, with `permissions: contents: write` scoped to the release job
- [x] 3.2 In the workflow, checkout with full tag history (`fetch-depth: 0`), set up Go via `actions/setup-go@v5` with `go-version-file: go.mod`, then run GoReleaser's release action against the pushed tag
- [x] 3.3 Verify the workflow does not trigger on a plain push to `main` or on a pull request (inspect the `on:` trigger; no live tag push needed for this check)
- [x] 3.4 Push a real test tag (e.g. `v0.0.1-test.1`) to a scratch branch/fork or dry-run the workflow, and verify a GitHub Release is created with six platform archives, a checksums file, and each binary reporting `--version` matching the tag (minus leading `v`) — repo owner pushed `v0.0.1-test.1`; workflow run [34706964017](https://github.com/sdpilon/bldoc/actions/runs/34706964017) succeeded; the resulting [Release](https://github.com/sdpilon/bldoc/releases/tag/v0.0.1-test.1) has all 6 archives + `checksums.txt`, every checksum matches its asset's digest, and the downloaded linux/amd64 binary reports `bldoc version 0.0.1-test.1` exactly
- [ ] 3.5 Delete the test tag/release created in 3.4 once verified, so it doesn't linger as a real release — **needs manual cleanup**: this session's git push access is scoped to the feature branch and rejects a tag-delete push (HTTP 403 on `git push origin --delete v0.0.1-test.1`), and no GitHub MCP tool deletes a tag or release either; the repo owner needs to delete the `v0.0.1-test.1` tag and its Release from the GitHub UI

## 4. Documentation

- [x] 4.1 Add `RELEASING.md` at the repo root documenting: how to cut a release (create and push a `vX.Y.Z` tag), what the workflow produces, and how to install `bldoc` from a downloaded release archive (download, extract, place on `PATH`)

## 5. Final verification

- [x] 5.1 Run the full existing test/lint suite (`go build ./...`, `go test ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run`) and confirm nothing broke from the `version` var change — build/test/vet/fmt all clean; local `golangci-lint` binary can't run at all (pre-existing Go-version mismatch with `go.mod`'s `go 1.27.1`, unrelated to this change) — CI's `golangci-lint-action@v9` installs a current version and is unaffected
- [x] 5.2 Re-read `specs/release-distribution/spec.md` scenario by scenario and confirm each is satisfied by the snapshot build (2.4-2.5) and the tag-triggered run (3.4) — all 8 scenarios confirmed: tag-triggered release and non-tag-push exclusion (workflow `on:` block + the real run), all 6 platform/arch combinations, archive contents, checksums covering every archive, released version matching the tag exactly, local-build version unaffected, and Release assets attached (all verified against the real `v0.0.1-test.1` run)
