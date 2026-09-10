package cli

import (
	"strings"
	"testing"
)

func TestMake_ExplicitTarget(t *testing.T) {
	_, stderr, err := execute(t, "make", "README")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc make: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}

func TestMake_NoTarget(t *testing.T) {
	_, stderr, err := execute(t, "make")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc make: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}

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
