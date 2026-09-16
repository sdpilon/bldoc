package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTarget(t *testing.T, name, mode string) {
	t.Helper()
	if _, stderr, err := execute(t, "new", name, "--mode", mode); err != nil {
		t.Fatalf("new %s --mode %s: err=%v stderr=%q", name, mode, err, stderr)
	}
}

func TestAddDep_WholeFile(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")

	if _, stderr, err := execute(t, "add-dep", "README", "pyproject.toml"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "pyproject.toml") {
		t.Fatalf("expected pyproject.toml recorded, got stdout=%q", stdout)
	}
}

func TestAddDep_FieldWithFormat(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

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
	newTarget(t, "README", "raw")
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
	newTarget(t, "README", "field")
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

func TestAddDep_AnchorDependency(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README:summary", "spec.md#purpose"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "spec.md#purpose") {
		t.Fatalf("expected anchor dependency recorded, got stdout=%q", stdout)
	}
}

func TestAddDep_AnchorBreadcrumbDependency(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README:x", "spec.md#requirement-a/scenario-b"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "spec.md#requirement-a/scenario-b") {
		t.Fatalf("expected breadcrumb anchor dependency recorded, got stdout=%q", stdout)
	}
}

func TestAddDep_AnchorAndPathCombinedRejected(t *testing.T) {
	_, stderr, err := execute(t, "add-dep", "README:x", "spec.md:some.path#purpose")
	if err == nil {
		t.Fatal("expected a usage error combining ':field-path' and '#anchor'")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestAddDep_NestedAcceptedWithAnchorAndField(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README:section", "spec.md#requirements", "--nested"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "(nested)") {
		t.Fatalf("expected nested flag recorded, got stdout=%q", stdout)
	}
}

func TestAddDep_NestedRejectedWithoutAnchor(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	_, stderr, err := execute(t, "add-dep", "README:section", "spec.md", "--nested")
	if err == nil {
		t.Fatal("expected an error using --nested without an anchor")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestAddDep_NestedRejectedOnRawModeDependency(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")

	_, stderr, err := execute(t, "add-dep", "README", "spec.md#requirements", "--nested")
	if err == nil {
		t.Fatal("expected an error using --nested on a whole-file (raw-mode) dependency")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestAddDep_AmbiguousAnchorWarnsButRecords(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	newTarget(t, "README", "field")
	specContent := "# Title\n\n### Requirement: A\n\n#### Scenario: Dup\nbody\n\n### Requirement: B\n\n#### Scenario: Dup\nbody\n"
	if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	stdout, stderr, err := execute(t, "add-dep", "README:x", "spec.md#scenario-dup")
	if err != nil {
		t.Fatalf("add-dep: err=%v stdout=%q stderr=%q", err, stdout, stderr)
	}
	if !strings.Contains(stderr, "warning") || !strings.Contains(stderr, "scenario-dup") {
		t.Fatalf("expected an ambiguity warning on stderr, got stderr=%q", stderr)
	}

	showStdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(showStdout, "spec.md#scenario-dup") {
		t.Fatalf("expected the dependency still recorded despite the warning, got stdout=%q", showStdout)
	}
}

func TestAddDep_UnambiguousAnchorNoWarning(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	newTarget(t, "README", "field")
	if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Title\n\n## Purpose\nbody\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README:x", "spec.md#purpose")
	if err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected no warning for an unambiguous anchor, got stderr=%q", stderr)
	}
}

func TestAddDep_BreadcrumbAnchorNeverWarns(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	newTarget(t, "README", "field")
	specContent := "# Title\n\n### Requirement: A\n\n#### Scenario: Dup\nbody\n\n### Requirement: B\n\n#### Scenario: Dup\nbody\n"
	if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README:x", "spec.md#requirement-a/scenario-dup")
	if err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected no ambiguity check for a breadcrumb anchor, got stderr=%q", stderr)
	}
}

func TestAddDep_AnchorUnreadableSourceSkipsCheckSilently(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	_, stderr, err := execute(t, "add-dep", "README:x", "missing.md#purpose")
	if err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected no warning or error for an unreadable source, got stderr=%q", stderr)
	}
}

func TestAddDep_DuplicateRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
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
