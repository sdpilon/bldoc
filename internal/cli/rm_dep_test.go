package cli

import (
	"strings"
	"testing"
)

func TestRmDep_Valid(t *testing.T) {
	_, stderr, err := execute(t, "rm-dep", "README", "pyproject.toml")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc rm-dep: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}

func TestRmDep_FormatFlagRejected(t *testing.T) {
	_, stderr, err := execute(t, "rm-dep", "README", "--format", "x", "pyproject.toml")
	if err == nil {
		t.Fatal("expected an error: rm-dep does not accept --format")
	}
	if strings.Contains(stderr, "not yet implemented") {
		t.Fatalf("expected a usage/flag error, not the not-yet-implemented stub, got stderr=%q", stderr)
	}
}
