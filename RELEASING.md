# Releasing bldoc

Releases are built and published automatically by
[GoReleaser](https://goreleaser.com/), driven by a GitHub Actions
workflow (`.github/workflows/release.yml`) that runs whenever a tag
matching `v*` is pushed.

## Cutting a release

1. Make sure `main` is at the commit you want to release.
2. Create and push a tag named `vX.Y.Z` (semantic versioning, e.g.
   `v0.1.0`):

   ```sh
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. The `Release` workflow picks up the tag push and runs GoReleaser,
   which:
   - builds `bldoc` for linux, macOS, and windows, on amd64 and arm64
     (six binaries total)
   - packages each binary in a `.tar.gz` archive (`.zip` on windows)
   - generates a `checksums.txt` covering every archive
   - publishes all of it as assets on a new GitHub Release for that tag

No other push or pull request triggers a release — only a `v*` tag
push does.

## Installing a released binary

1. Go to the repository's [Releases](../../releases) page and find the
   release for the version you want.
2. Download the archive for your platform, e.g.
   `bldoc_0.1.0_linux_amd64.tar.gz` or `bldoc_0.1.0_windows_amd64.zip`.
3. Extract it:

   ```sh
   tar -xzf bldoc_0.1.0_linux_amd64.tar.gz
   ```

4. Move the extracted `bldoc` binary (or `bldoc.exe` on windows)
   somewhere on your `PATH`, e.g.:

   ```sh
   mv bldoc /usr/local/bin/
   ```

5. Confirm it works and reports the expected version:

   ```sh
   bldoc --version
   ```

(Optional) verify the download against `checksums.txt` before
extracting it:

```sh
sha256sum -c --ignore-missing checksums.txt
```
