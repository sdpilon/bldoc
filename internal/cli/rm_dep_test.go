package cli

import (
	"strings"
	"testing"
)

func TestRmDep_Valid(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
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
	if strings.Contains(stdout, "pyproject.toml") {
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

func TestRmDep_AnchorDependency(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")
	if _, _, err := execute(t, "add-dep", "README:a", "spec.md#purpose"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README:b", "spec.md#requirements"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "rm-dep", "README:a", "spec.md#purpose"); err != nil {
		t.Fatalf("rm-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Contains(stdout, "purpose") {
		t.Fatalf("expected the purpose anchor dependency removed, got stdout=%q", stdout)
	}
	if !strings.Contains(stdout, "requirements") {
		t.Fatalf("expected the requirements anchor dependency to remain, got stdout=%q", stdout)
	}
}

func TestRmDep_NestedFlagRejected(t *testing.T) {
	_, stderr, err := execute(t, "rm-dep", "README", "--nested", "spec.md#purpose")
	if err == nil {
		t.Fatal("expected an error: rm-dep does not accept --nested")
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
	newTarget(t, "README", "raw")

	_, stderr, err := execute(t, "rm-dep", "README", "other.toml")
	if err == nil {
		t.Fatal("expected an error removing an unrecorded dependency")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}
