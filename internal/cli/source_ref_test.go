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

func TestParseSourceRef_Anchor(t *testing.T) {
	ref, err := parseSourceRef("spec.md#purpose")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Source != "spec.md" || ref.Anchor != "purpose" || ref.Path != "" {
		t.Fatalf("got %+v", ref)
	}
}

func TestParseSourceRef_AnchorBreadcrumb(t *testing.T) {
	ref, err := parseSourceRef("spec.md#requirement-x/scenario-y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Source != "spec.md" || ref.Anchor != "requirement-x/scenario-y" {
		t.Fatalf("got %+v", ref)
	}
}

func TestParseSourceRef_AnchorAndPathCombined(t *testing.T) {
	_, err := parseSourceRef("spec.md:some.path#purpose")
	if err == nil {
		t.Fatal("expected an error combining ':field-path' and '#anchor'")
	}
}

func TestParseSourceRef_MultipleHash(t *testing.T) {
	_, err := parseSourceRef("spec.md#a#b")
	if err == nil {
		t.Fatal("expected an error for a source-ref with more than one '#'")
	}
}

func TestParseSourceRef_EmptyAnchor(t *testing.T) {
	_, err := parseSourceRef("spec.md#")
	if err == nil {
		t.Fatal("expected an error for an empty anchor")
	}
}

func TestParseSourceRef_EmptySourceBeforeHash(t *testing.T) {
	_, err := parseSourceRef("#purpose")
	if err == nil {
		t.Fatal("expected an error for an empty source before '#'")
	}
}
