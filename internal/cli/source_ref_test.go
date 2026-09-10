package cli

import "testing"

func TestParseSourceRef_WholeFile(t *testing.T) {
	ref, err := parseSourceRef("pyproject.toml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Source != "pyproject.toml" || ref.Path != "" {
		t.Fatalf("got %+v", ref)
	}
}

func TestParseSourceRef_Path(t *testing.T) {
	ref, err := parseSourceRef("pyproject.toml:project.requires-python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Source != "pyproject.toml" || ref.Path != "project.requires-python" {
		t.Fatalf("got %+v", ref)
	}
}

func TestParseSourceRef_Malformed(t *testing.T) {
	_, err := parseSourceRef("a:b:c")
	if err == nil {
		t.Fatal("expected an error for a source-ref with more than one ':'")
	}
}
