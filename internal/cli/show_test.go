package cli

import (
	"strings"
	"testing"
)

func TestShow_MissingTarget(t *testing.T) {
	_, stderr, err := execute(t, "show")
	if err == nil {
		t.Fatal("expected an error for missing target name")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestShow_ValidTarget(t *testing.T) {
	_, stderr, err := execute(t, "show", "README")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc show: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}
