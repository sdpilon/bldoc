package cli

import (
	"strings"
	"testing"
)

func TestRmDep_Valid(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")
	if _, _, err := execute(t, "add-dep", "README", "pyproject.toml"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "rm-dep", "README", "pyproject.toml"); err != nil {
		t.Fatalf("rm-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected no dependencies remaining, got stdout=%q", stdout)
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

func TestRmDep_UnknownTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "rm-dep", "MISSING", "pyproject.toml")
	if err == nil {
		t.Fatal("expected an error for an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestRmDep_UnrecordedDependencyRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")

	_, stderr, err := execute(t, "rm-dep", "README", "other.toml")
	if err == nil {
		t.Fatal("expected an error removing an unrecorded dependency")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}
