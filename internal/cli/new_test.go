package cli

import (
	"strings"
	"testing"
)

func TestNew_MissingTarget(t *testing.T) {
	_, stderr, err := execute(t, "new")
	if err == nil {
		t.Fatal("expected an error for missing target name")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestNew_ValidTarget(t *testing.T) {
	t.Chdir(t.TempDir())

	if _, stderr, err := execute(t, "new", "README"); err != nil {
		t.Fatalf("expected new to succeed, got err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if strings.TrimSpace(stdout) != "README" {
		t.Fatalf("expected README to be recorded, got stdout=%q", stdout)
	}
}

func TestNew_DuplicateTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, _, err := execute(t, "new", "README"); err != nil {
		t.Fatalf("first new: %v", err)
	}

	_, stderr, err := execute(t, "new", "README")
	if err == nil {
		t.Fatal("expected an error creating a duplicate target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}
