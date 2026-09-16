package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bldoc/internal/manifest"
)

func TestTarget_EmptyDepsCompilesToEmptyRaw(t *testing.T) {
	result, err := Target(manifest.Target{Name: "README"})
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if result.IsField {
		t.Fatal("expected raw-mode result for an empty target")
	}
	if len(result.Raw) != 0 {
		t.Fatalf("expected empty raw bytes, got %q", result.Raw)
	}
}

func TestTarget_EmptyDeclaredRawModeCompilesToEmptyRaw(t *testing.T) {
	result, err := Target(manifest.Target{Name: "NOTES", Mode: "raw"})
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if result.IsField {
		t.Fatal("expected raw-mode result for an empty declared-raw target")
	}
	if len(result.Raw) != 0 {
		t.Fatalf("expected empty raw bytes, got %q", result.Raw)
	}
}

func TestTarget_EmptyDeclaredFieldModeCompilesToEmptyObject(t *testing.T) {
	result, err := Target(manifest.Target{Name: "PROJECTS", Mode: "field"})
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if !result.IsField {
		t.Fatal("expected field-mode result for an empty declared-field target")
	}
	if result.Fields == nil {
		t.Fatal("expected a non-nil empty Fields map (marshals to {}, not null)")
	}
	if len(result.Fields) != 0 {
		t.Fatalf("expected no fields, got %+v", result.Fields)
	}
}

func TestTarget_RawModeConcatenatesInOrder(t *testing.T) {
	dir := t.TempDir()
	aPath := filepath.Join(dir, "a.txt")
	bPath := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(aPath, []byte("hello "), 0o644); err != nil {
		t.Fatalf("WriteFile a: %v", err)
	}
	if err := os.WriteFile(bPath, []byte("world"), 0o644); err != nil {
		t.Fatalf("WriteFile b: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{Source: aPath}, {Source: bPath}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if result.IsField {
		t.Fatal("expected raw-mode result")
	}
	if string(result.Raw) != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", result.Raw)
	}
}

func TestTarget_FieldModeWholeFile(t *testing.T) {
	dir := t.TempDir()
	notesPath := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(notesPath, []byte("some notes"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{Source: notesPath, Field: "summary"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if !result.IsField {
		t.Fatal("expected field-mode result")
	}
	f, ok := result.Fields["summary"]
	if !ok {
		t.Fatalf("expected a summary field, got %+v", result.Fields)
	}
	if f.Value != "some notes" || f.Raw != "some notes" {
		t.Fatalf("expected value==raw==%q, got %+v", "some notes", f)
	}
}

func TestTarget_RawModeAnchor(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Title\n\n## Purpose\n\nPurpose text.\n\n## Requirements\n\nreqs\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{Source: specPath, Anchor: "purpose"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if result.IsField {
		t.Fatal("expected raw-mode result")
	}
	if !strings.Contains(string(result.Raw), "Purpose text.") {
		t.Fatalf("got %q", result.Raw)
	}
	if strings.Contains(string(result.Raw), "Requirements") {
		t.Fatalf("expected only the purpose section, got %q", result.Raw)
	}
}

func TestTarget_FieldModeAnchor(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Title\n\n## Purpose\n\nPurpose text.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{Source: specPath, Anchor: "purpose", Field: "summary"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	f, ok := result.Fields["summary"]
	if !ok {
		t.Fatalf("expected a summary field, got %+v", result.Fields)
	}
	if !strings.Contains(f.Value, "Purpose text.") || f.Value != f.Raw {
		t.Fatalf("got %+v", f)
	}
}

func TestTarget_FieldModeAnchorNested(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	content := "# Title\n\n## Requirements\n\n### Requirement: X\nbody\n\n## Impact\n\nimpact\n"
	if err := os.WriteFile(specPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{Source: specPath, Anchor: "requirements", Nested: true, Field: "section"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	f := result.Fields["section"]
	if !strings.Contains(f.Value, "Requirement: X") {
		t.Fatalf("expected nested capture to include the subsection, got %+v", f)
	}
	if strings.Contains(f.Value, "impact") {
		t.Fatalf("expected nested capture to stop before the sibling section, got %+v", f)
	}
}

func TestTarget_AnchorUnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	notesPath := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(notesPath, []byte("# Purpose\nbody\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{Source: notesPath, Anchor: "purpose", Field: "summary"}},
	}
	if _, err := Target(target); err == nil {
		t.Fatal("expected an error resolving an anchor against a non-Markdown source")
	}
}

func TestTarget_FieldModeStructuredWithFormat(t *testing.T) {
	dir := t.TempDir()
	pyprojectPath := filepath.Join(dir, "pyproject.toml")
	if err := os.WriteFile(pyprojectPath, []byte("[project]\nrequires-python = \"3.11\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{{
			Source: pyprojectPath,
			Path:   "project.requires-python",
			Field:  "version",
			Format: "Python version must be %s to run this project.",
		}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	f, ok := result.Fields["version"]
	if !ok {
		t.Fatalf("expected a version field, got %+v", result.Fields)
	}
	if f.Raw != "3.11" {
		t.Fatalf("expected raw 3.11, got %q", f.Raw)
	}
	want := "Python version must be 3.11 to run this project."
	if f.Value != want {
		t.Fatalf("expected value %q, got %q", want, f.Value)
	}
}

func TestTarget_DuplicateFieldNameRejected(t *testing.T) {
	dir := t.TempDir()
	aPath := filepath.Join(dir, "a.toml")
	bPath := filepath.Join(dir, "b.toml")
	if err := os.WriteFile(aPath, []byte("x = \"1\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile a: %v", err)
	}
	if err := os.WriteFile(bPath, []byte("x = \"2\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile b: %v", err)
	}

	target := manifest.Target{
		Name: "README",
		Deps: []manifest.Dep{
			{Source: aPath, Path: "x", Field: "version"},
			{Source: bPath, Path: "x", Field: "version"},
		},
	}
	if _, err := Target(target); err == nil {
		t.Fatal("expected an error for duplicate field names")
	}
}
