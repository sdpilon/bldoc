# bldoc demo

A checked-in, ready-to-run example project. Use it to try `bldoc` end to
end without pointing the tool at a real project.

## Run it

From the repository root:

```
go build -o bldoc ./cmd/bldoc
cd examples/demo
../../bldoc make
```

(or, without a separate build step: `go run ../../cmd/bldoc make` from
`examples/demo/`)

This compiles both targets declared in `bldoc.toml`:

- `.bldoc/NOTES` — the raw-mode concatenation of `intro.txt` and `license.txt`
- `.bldoc/SUMMARY.json` — field-mode values pulled from `pyproject.toml`,
  `package.json`, and `config.yaml`

`.bldoc/` is gitignored — it's disposable and always fully recomputed, so
feel free to delete it and re-run `make` at any time.

## How the manifest was built

The checked-in `bldoc.toml` was produced by running these commands from
within `examples/demo/`:

```
bldoc new NOTES
bldoc add-dep NOTES intro.txt
bldoc add-dep NOTES license.txt

bldoc new SUMMARY
bldoc add-dep SUMMARY:version --format "Python version must be %s to run this project." pyproject.toml:project.requires-python
bldoc add-dep SUMMARY:name package.json:name
bldoc add-dep SUMMARY:team config.yaml:owner.team
```

`NOTES` demonstrates raw mode (whole-file dependencies, concatenated
byte-for-byte). `SUMMARY` demonstrates field mode: one field resolved
from TOML with a format template applied, and two fields resolved from
JSON and YAML with no format (passed through raw).

Delete `bldoc.toml` and re-run the commands above to reproduce it by hand.
