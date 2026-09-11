package cli

import (
	"strings"
	"testing"
)

func TestRm_MissingTarget(t *testing.T) {
	_, stderr, err := execute(t, "rm")
	if err == nil {
		t.Fatal("expected an error for missing target name")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestRm_ValidTarget(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")

	if _, stderr, err := execute(t, "rm", "README"); err != nil {
		t.Fatalf("rm: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected no targets remaining, got stdout=%q", stdout)
	}
}

func TestRm_UnknownTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "rm", "MISSING")
	if err == nil {
		t.Fatal("expected an error removing an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}
