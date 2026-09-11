package compile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestParseSource_TOML(t *testing.T) {
	path := writeTemp(t, "pyproject.toml", "[project]\nrequires-python = \"3.11\"\n")

	tree, err := parseSource(path)
	if err != nil {
		t.Fatalf("parseSource: %v", err)
	}
	v, err := resolveFieldPath(tree, "project.requires-python")
	if err != nil {
		t.Fatalf("resolveFieldPath: %v", err)
	}
	if v != "3.11" {
		t.Fatalf("expected 3.11, got %v", v)
	}
}

func TestParseSource_JSON(t *testing.T) {
	path := writeTemp(t, "package.json", `{"project": {"version": "1.2.3"}}`)

	tree, err := parseSource(path)
	if err != nil {
		t.Fatalf("parseSource: %v", err)
	}
	v, err := resolveFieldPath(tree, "project.version")
	if err != nil {
		t.Fatalf("resolveFieldPath: %v", err)
	}
	if v != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %v", v)
	}
}

func TestParseSource_YAML(t *testing.T) {
	path := writeTemp(t, "config.yaml", "project:\n  version: \"9.9.9\"\n")

	tree, err := parseSource(path)
	if err != nil {
		t.Fatalf("parseSource: %v", err)
	}
	v, err := resolveFieldPath(tree, "project.version")
	if err != nil {
		t.Fatalf("resolveFieldPath: %v", err)
	}
	if v != "9.9.9" {
		t.Fatalf("expected 9.9.9, got %v", v)
	}
}

func TestParseSource_UnsupportedExtension(t *testing.T) {
	path := writeTemp(t, "notes.txt", "hello")

	if _, err := parseSource(path); err == nil {
		t.Fatal("expected an error for an unsupported source extension")
	}
}

func TestParseSource_Malformed(t *testing.T) {
	path := writeTemp(t, "broken.toml", "this is not [ valid toml")

	if _, err := parseSource(path); err == nil {
		t.Fatal("expected an error parsing a malformed TOML file")
	}
}

func TestResolveFieldPath_MissingKey(t *testing.T) {
	tree := map[string]interface{}{"project": map[string]interface{}{"version": "1.0"}}

	if _, err := resolveFieldPath(tree, "project.missing"); err == nil {
		t.Fatal("expected an error for a missing key")
	}
}

func TestResolveFieldPath_NonScalarTable(t *testing.T) {
	tree := map[string]interface{}{"project": map[string]interface{}{"version": "1.0"}}

	if _, err := resolveFieldPath(tree, "project"); err == nil {
		t.Fatal("expected an error resolving a field-path to a table")
	}
}

func TestResolveFieldPath_NonScalarArray(t *testing.T) {
	tree := map[string]interface{}{"items": []interface{}{"a", "b"}}

	if _, err := resolveFieldPath(tree, "items"); err == nil {
		t.Fatal("expected an error resolving a field-path to an array")
	}
}
