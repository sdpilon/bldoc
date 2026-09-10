package cli

import (
	"strings"
	"testing"
)

func newTarget(t *testing.T, name string) {
	t.Helper()
	if _, stderr, err := execute(t, "new", name); err != nil {
		t.Fatalf("new %s: err=%v stderr=%q", name, err, stderr)
	}
}

func TestAddDep_WholeFile(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")

	if _, stderr, err := execute(t, "add-dep", "README", "pyproject.toml"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.TrimSpace(stdout) != "pyproject.toml" {
		t.Fatalf("expected pyproject.toml recorded, got stdout=%q", stdout)
	}
}

func TestAddDep_FieldWithFormat(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")

	_, stderr, err := execute(t, "add-dep", "README:version", "--format", "Python version must be %s.", "pyproject.toml:project.requires-python")
	if err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "version") || !strings.Contains(stdout, "pyproject.toml:project.requires-python") || !strings.Contains(stdout, "Python version must be %s.") {
		t.Fatalf("expected field, source, path, and format recorded, got stdout=%q", stdout)
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

func TestAddDep_UnknownTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "add-dep", "MISSING", "pyproject.toml")
	if err == nil {
		t.Fatal("expected an error for an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestAddDep_FieldModeRejectedOnRawTarget(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")
	if _, _, err := execute(t, "add-dep", "README", "a.toml"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README:version", "--format", "x", "b.toml:x")
	if err == nil {
		t.Fatal("expected an error mixing a field-addressed dep into a raw-mode target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestAddDep_RawModeRejectedOnFieldTarget(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")
	if _, _, err := execute(t, "add-dep", "README:version", "--format", "x", "a.toml:x"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README", "b.toml")
	if err == nil {
		t.Fatal("expected an error mixing a raw dep into a field-mode target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestAddDep_DuplicateRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README")
	if _, _, err := execute(t, "add-dep", "README", "pyproject.toml"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README", "pyproject.toml")
	if err == nil {
		t.Fatal("expected an error adding a duplicate dependency")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}
