package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMake_TooManyTargets(t *testing.T) {
	_, stderr, err := execute(t, "make", "README", "ROADMAP")
	if err == nil {
		t.Fatal("expected a usage error for more than one target")
	}
	if strings.Contains(stderr, "not yet implemented") {
		t.Fatalf("expected a usage error, not the not-yet-implemented stub, got stderr=%q", stderr)
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestMake_UnknownTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "make", "MISSING")
	if err == nil {
		t.Fatal("expected an error for an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestMake_ExplicitTarget_RawMode(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	if err := os.WriteFile("a.txt", []byte("hello "), 0o644); err != nil {
		t.Fatalf("WriteFile a.txt: %v", err)
	}
	if err := os.WriteFile("b.txt", []byte("world"), 0o644); err != nil {
		t.Fatalf("WriteFile b.txt: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README", "a.txt"); err != nil {
		t.Fatalf("add-dep a.txt: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README", "b.txt"); err != nil {
		t.Fatalf("add-dep b.txt: %v", err)
	}

	if _, stderr, err := execute(t, "make", "README"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	data, err := os.ReadFile(filepath.Join(".bldoc", "README"))
	if err != nil {
		t.Fatalf("reading .bldoc/README: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", data)
	}
}

func TestMake_ExplicitTarget_FieldMode(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")
	if err := os.WriteFile("pyproject.toml", []byte("[project]\nrequires-python = \"3.11\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile pyproject.toml: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README@version", "--format", "Python version must be %s to run this project.", "pyproject.toml@project.requires-python"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "make", "README"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	data, err := os.ReadFile(filepath.Join(".bldoc", "README.json"))
	if err != nil {
		t.Fatalf("reading .bldoc/README.json: %v", err)
	}
	var got map[string]struct {
		Value string `json:"value"`
		Raw   string `json:"raw"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshaling intermediate: %v", err)
	}
	f, ok := got["version"]
	if !ok {
		t.Fatalf("expected a version field, got %+v", got)
	}
	if f.Raw != "3.11" {
		t.Fatalf("expected raw 3.11, got %q", f.Raw)
	}
	want := "Python version must be 3.11 to run this project."
	if f.Value != want {
		t.Fatalf("expected value %q, got %q", want, f.Value)
	}
}

func TestMake_RawMode_ExtSuffixesOutputPath(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, _, err := execute(t, "new", "NOTES", "--mode", "raw", "--ext", "yaml"); err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := os.WriteFile("a.yaml", []byte("- a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile a.yaml: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "NOTES", "a.yaml"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "make", "NOTES"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	data, err := os.ReadFile(filepath.Join(".bldoc", "NOTES.yaml"))
	if err != nil {
		t.Fatalf("reading .bldoc/NOTES.yaml: %v", err)
	}
	if string(data) != "- a\n" {
		t.Fatalf("expected %q, got %q", "- a\n", data)
	}
	if _, err := os.Stat(filepath.Join(".bldoc", "NOTES")); err == nil {
		t.Fatal("expected no extensionless .bldoc/NOTES file to be written")
	}
}

func TestMake_EmptyDeclaredFieldMode_CompilesToEmptyObject(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, _, err := execute(t, "new", "PROJECTS", "--mode", "field"); err != nil {
		t.Fatalf("new: %v", err)
	}

	if _, stderr, err := execute(t, "make", "PROJECTS"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	data, err := os.ReadFile(filepath.Join(".bldoc", "PROJECTS.json"))
	if err != nil {
		t.Fatalf("reading .bldoc/PROJECTS.json: %v", err)
	}
	if strings.TrimSpace(string(data)) != "{}" {
		t.Fatalf("expected an empty JSON object, got %q", data)
	}
}

func TestMake_FieldMode_OutputPathUnaffectedByExt(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")
	if err := os.WriteFile("pyproject.toml", []byte("[project]\nrequires-python = \"3.11\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile pyproject.toml: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README@version", "pyproject.toml@project.requires-python"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "make", "README"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	if _, err := os.Stat(filepath.Join(".bldoc", "README.json")); err != nil {
		t.Fatalf("expected .bldoc/README.json to exist: %v", err)
	}
}

func TestMake_FieldMode_ExtYAML_WritesValidYAML(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, _, err := execute(t, "new", "README", "--mode", "field", "--ext", "yaml"); err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := os.WriteFile("notes.txt", []byte("some notes"), 0o644); err != nil {
		t.Fatalf("WriteFile notes.txt: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README@summary", "notes.txt"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "make", "README"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	data, err := os.ReadFile(filepath.Join(".bldoc", "README.yaml"))
	if err != nil {
		t.Fatalf("reading .bldoc/README.yaml: %v", err)
	}
	var out map[string]map[string]string
	if err := yaml.Unmarshal(data, &out); err != nil {
		t.Fatalf("compiled output is not valid YAML: %v\noutput:\n%s", err, data)
	}
	if out["summary"]["value"] != "some notes" {
		t.Fatalf("expected summary.value == \"some notes\", got %+v", out)
	}
}

func TestMake_ListMode_ExtYAML_WritesValidYAML(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, _, err := execute(t, "new", "PROJECTS", "--mode", "list", "--ext", "yaml"); err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := os.WriteFile("bldoc.yaml", []byte("path: ~/bldoc\n"), 0o644); err != nil {
		t.Fatalf("WriteFile bldoc.yaml: %v", err)
	}
	if err := os.WriteFile("description.md", []byte("a CLI tool"), 0o644); err != nil {
		t.Fatalf("WriteFile description.md: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "PROJECTS@bldoc", "bldoc.yaml"); err != nil {
		t.Fatalf("add-dep record-only: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "PROJECTS@bldoc.description", "description.md"); err != nil {
		t.Fatalf("add-dep record.field: %v", err)
	}

	if _, stderr, err := execute(t, "make", "PROJECTS"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	data, err := os.ReadFile(filepath.Join(".bldoc", "PROJECTS.yaml"))
	if err != nil {
		t.Fatalf("reading .bldoc/PROJECTS.yaml: %v", err)
	}
	var out []map[string]interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		t.Fatalf("compiled output is not valid YAML: %v\noutput:\n%s", err, data)
	}
	if len(out) != 1 || out[0]["path"] != "~/bldoc" || out[0]["description"] != "a CLI tool" {
		t.Fatalf("expected one merged record, got %+v", out)
	}
}

func TestMake_NoTarget_CompilesAll(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	newTarget(t, "CHANGELOG", "raw")
	if err := os.WriteFile("a.txt", []byte("a"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README", "a.txt"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "make"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	if _, err := os.Stat(filepath.Join(".bldoc", "README")); err != nil {
		t.Fatalf("expected .bldoc/README to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(".bldoc", "CHANGELOG")); err != nil {
		t.Fatalf("expected .bldoc/CHANGELOG to exist: %v", err)
	}
}

func TestMake_NoTarget_NoTargetsIsNoOp(t *testing.T) {
	t.Chdir(t.TempDir())

	if _, stderr, err := execute(t, "make"); err != nil {
		t.Fatalf("make: err=%v stderr=%q", err, stderr)
	}

	if _, err := os.Stat(".bldoc"); err == nil {
		t.Fatal("expected no .bldoc directory to be created for an empty manifest")
	}
}
