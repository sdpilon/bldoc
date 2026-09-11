## 1. Source files

- [x] 1.1 Create `examples/demo/intro.txt` and `examples/demo/license.txt` (short, arbitrary text) and verify their concatenation is what `NOTES` should compile to
- [x] 1.2 Create `examples/demo/pyproject.toml` with a `project.requires-python` key, `examples/demo/package.json` with a `name` key, and `examples/demo/config.yaml` with a nested `owner.team` key, and verify each file parses as valid TOML/JSON/YAML respectively

## 2. Manifest construction

- [x] 2.1 From within `examples/demo/`, build the manifest via the CLI: `bldoc new NOTES`, `bldoc add-dep NOTES intro.txt`, `bldoc add-dep NOTES license.txt` — verify `bldoc show NOTES` lists both dependencies in order
- [x] 2.2 Continue via the CLI: `bldoc new SUMMARY`, `bldoc add-dep SUMMARY:version --format "Python version must be %s to run this project." pyproject.toml:project.requires-python`, `bldoc add-dep SUMMARY:name package.json:name`, `bldoc add-dep SUMMARY:team config.yaml:owner.team` — verify `bldoc show SUMMARY` lists all three field-addressed dependencies, one with a format template
- [x] 2.3 Record the exact sequence of commands from 2.1-2.2 for reuse in the README (task 4.1)

## 3. Compiled-output hygiene

- [x] 3.1 Add `examples/demo/.gitignore` ignoring `.bldoc/` and verify `git status` at the repo root reports nothing untracked under `examples/demo/.bldoc/` after running `make` (task 3.2)

## 4. Verification and documentation

- [x] 4.1 Write `examples/demo/README.md`: state how to build/run `bldoc` against the demo (e.g. `go run ../../cmd/bldoc make` from `examples/demo/`), and list the exact `new`/`add-dep` commands recorded in task 2.3 so the manifest can be reproduced by hand
- [x] 4.2 Run `bldoc make` (or `go run ../../cmd/bldoc make`) from `examples/demo/` and verify it exits zero, producing `.bldoc/NOTES` (raw concatenation of `intro.txt` + `license.txt`) and `.bldoc/SUMMARY.json` (three fields, `version`'s `value` rendered through its format template, `name` and `team` passed through raw)
