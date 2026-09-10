package cli

import (
	"strings"
	"testing"
)

func TestAddDep_WholeFile(t *testing.T) {
	_, stderr, err := execute(t, "add-dep", "README", "pyproject.toml")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc add-dep: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}

func TestAddDep_FieldWithFormat(t *testing.T) {
	_, stderr, err := execute(t, "add-dep", "README:version", "--format", "Python version must be %s.", "pyproject.toml:project.requires-python")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc add-dep: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}

func TestAddDep_MalformedTargetRef(t *testing.T) {
	_, stderr, err := execute(t, "add-dep", "README:version:extra", "pyproject.toml")
	if err == nil {
		t.Fatal("expected a usage error for a malformed target-ref")
	}
	if strings.Contains(stderr, "not yet implemented") {
		t.Fatalf("expected a usage error, not the not-yet-implemented stub, got stderr=%q", stderr)
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestAddDep_FormatWithoutField(t *testing.T) {
	_, stderr, err := execute(t, "add-dep", "README", "--format", "x", "pyproject.toml")
	if err == nil {
		t.Fatal("expected a usage error when --format is used without a :field target-ref")
	}
	if strings.Contains(stderr, "not yet implemented") {
		t.Fatalf("expected a usage error, not the not-yet-implemented stub, got stderr=%q", stderr)
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}
