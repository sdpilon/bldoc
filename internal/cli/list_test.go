package cli

import (
	"strings"
	"testing"
)

func TestList_ExtraArgumentRejected(t *testing.T) {
	_, stderr, err := execute(t, "list", "README")
	if err == nil {
		t.Fatal("expected a usage error for an extra positional argument")
	}
	if strings.Contains(stderr, "not yet implemented") {
		t.Fatalf("expected a usage error, not the not-yet-implemented stub, got stderr=%q", stderr)
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestList_NoArguments(t *testing.T) {
	_, stderr, err := execute(t, "list")
	if err == nil {
		t.Fatal("expected an error since the command is not yet implemented")
	}
	if !strings.Contains(stderr, "bldoc list: not yet implemented") {
		t.Fatalf("expected not-yet-implemented message, got stderr=%q", stderr)
	}
}
