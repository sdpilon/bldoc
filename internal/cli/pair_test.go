package cli

import "testing"

func TestParsePair_Bare(t *testing.T) {
	p, err := parsePair("pyproject.toml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Field != "" || p.Source.Source != "pyproject.toml" {
		t.Fatalf("got %+v", p)
	}
}

func TestParsePair_FieldAndSource(t *testing.T) {
	p, err := parsePair("version:pyproject.toml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Field != "version" || p.Source.Source != "pyproject.toml" {
		t.Fatalf("got %+v", p)
	}
}

func TestParsePair_FieldAndSourceWithPath(t *testing.T) {
	p, err := parsePair("version:pyproject.toml@project.version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Field != "version" || p.Source.Source != "pyproject.toml" || p.Source.Path != "project.version" {
		t.Fatalf("got %+v", p)
	}
}

func TestParsePair_FieldAndSourceWithAnchor(t *testing.T) {
	p, err := parsePair("summary:spec.md#purpose")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Field != "summary" || p.Source.Source != "spec.md" || p.Source.Anchor != "purpose" {
		t.Fatalf("got %+v", p)
	}
}

func TestParsePair_EmptyFieldRejected(t *testing.T) {
	_, err := parsePair(":pyproject.toml")
	if err == nil {
		t.Fatal("expected an error for an empty field before ':'")
	}
}

func TestParsePair_MalformedSourceRefPropagates(t *testing.T) {
	_, err := parsePair("field:a@b@c")
	if err == nil {
		t.Fatal("expected the underlying source-ref parse error to propagate")
	}
}

func TestValidateBareTarget(t *testing.T) {
	if err := validateBareTarget("README"); err != nil {
		t.Fatalf("unexpected error for a bare target: %v", err)
	}
	if err := validateBareTarget("README@version"); err == nil {
		t.Fatal("expected an error for a target carrying '@field'")
	}
}

func TestValidateRecordName(t *testing.T) {
	if err := validateRecordName("bldoc"); err != nil {
		t.Fatalf("unexpected error for a bare record name: %v", err)
	}
	for _, bad := range []string{"bldoc@x", "bldoc:x", "bldoc.x"} {
		if err := validateRecordName(bad); err == nil {
			t.Fatalf("expected an error for record name %q", bad)
		}
	}
}
