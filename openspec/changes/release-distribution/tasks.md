## 1. Version injection

- [ ] 1.1 Change `internal/cli/root.go`'s `version` from a `const` to a package-level `var version = "0.0.0-dev"`, and verify `go build ./... && go vet ./...` still pass
- [ ] 1.2 Verify a plain local build still reports the default: `go build -o /tmp/bldoc-local ./cmd/bldoc && /tmp/bldoc-local --version` prints `bldoc version 0.0.0-dev`

## 2. GoReleaser configuration

- [ ] 2.1 Add `.goreleaser.yaml` at the repo root defining a build matrix of `{linux, darwin, windows} x {amd64, arm64}` for `./cmd/bldoc`, with `ldflags` templated as `-s -w -X bldoc/internal/cli.version={{.Version}}`
- [ ] 2.2 Configure archive naming/format in `.goreleaser.yaml`: `.tar.gz` for linux/darwin, `.zip` for windows
- [ ] 2.3 Enable GoReleaser's checksum file generation covering every archive
- [ ] 2.4 Verify locally with `goreleaser release --snapshot --clean` (no tag push, no publish) and confirm `dist/` contains six archives (one per platform/arch) plus a checksums file
- [ ] 2.5 Verify a built snapshot binary reports the snapshot version when run with `--version` (e.g. extract the linux/amd64 archive and run `./bldoc --version`)

## 3. Release workflow

- [ ] 3.1 Add `.github/workflows/release.yml` triggered on `push: tags: ['v*']` only, with `permissions: contents: write` scoped to the release job
- [ ] 3.2 In the workflow, checkout with full tag history (`fetch-depth: 0`), set up Go via `actions/setup-go@v5` with `go-version-file: go.mod`, then run GoReleaser's release action against the pushed tag
- [ ] 3.3 Verify the workflow does not trigger on a plain push to `main` or on a pull request (inspect the `on:` trigger; no live tag push needed for this check)
- [ ] 3.4 Push a real test tag (e.g. `v0.0.1-test.1`) to a scratch branch/fork or dry-run the workflow, and verify a GitHub Release is created with six platform archives, a checksums file, and each binary reporting `--version` matching the tag (minus leading `v`)
- [ ] 3.5 Delete the test tag/release created in 3.4 once verified, so it doesn't linger as a real release

## 4. Documentation

- [ ] 4.1 Add `RELEASING.md` at the repo root documenting: how to cut a release (create and push a `vX.Y.Z` tag), what the workflow produces, and how to install `bldoc` from a downloaded release archive (download, extract, place on `PATH`)

## 5. Final verification

- [ ] 5.1 Run the full existing test/lint suite (`go build ./...`, `go test ./...`, `go vet ./...`, `gofmt -l .`, `golangci-lint run`) and confirm nothing broke from the `version` var change
- [ ] 5.2 Re-read `specs/release-distribution/spec.md` scenario by scenario and confirm each is satisfied by the snapshot build (2.4-2.5) and the tag-triggered run (3.4)
