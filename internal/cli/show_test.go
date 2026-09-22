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

func TestShow_UnknownTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "show", "MISSING")
	if err == nil {
		t.Fatal("expected an error for an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestShow_DependenciesShownInAddedOrder(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	if _, _, err := execute(t, "add-dep", "README", "a.toml"); err != nil {
		t.Fatalf("add-dep a.toml: %v", err)
	}
	if _, _, err := execute(t, "add-dep", "README", "b.toml"); err != nil {
		t.Fatalf("add-dep b.toml: %v", err)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	aIdx := strings.Index(stdout, "a.toml")
	bIdx := strings.Index(stdout, "b.toml")
	if aIdx == -1 || bIdx == -1 || aIdx > bIdx {
		t.Fatalf("expected a.toml before b.toml in added order, got stdout=%q", stdout)
	}
}

func TestShow_DeclaredModeShown(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "PROJECTS", "field")

	stdout, _, err := execute(t, "show", "PROJECTS")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "mode: field") {
		t.Fatalf("expected declared mode field to be shown, got stdout=%q", stdout)
	}
}

func TestShow_ComputedModeShownForLegacyTarget(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	if _, _, err := execute(t, "add-dep", "README", "a.toml"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "mode: raw") {
		t.Fatalf("expected computed mode raw to be shown, got stdout=%q", stdout)
	}
}

func TestShow_ExtensionShown(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, _, err := execute(t, "new", "PROJECTS", "--mode", "raw", "--ext", "yaml"); err != nil {
		t.Fatalf("new: %v", err)
	}

	stdout, _, err := execute(t, "show", "PROJECTS")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "ext: yaml") {
		t.Fatalf("expected extension yaml to be shown, got stdout=%q", stdout)
	}
}

func TestShow_AnchorDependencyShown(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")
	if _, _, err := execute(t, "add-dep", "README@summary", "spec.md#purpose"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "summary <- spec.md#purpose") {
		t.Fatalf("expected anchor dependency shown, got stdout=%q", stdout)
	}
}

func TestShow_NestedAnchorDependencyShown(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")
	if _, _, err := execute(t, "add-dep", "README@section", "spec.md#requirements", "--nested"); err != nil {
		t.Fatalf("add-dep: %v", err)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "spec.md#requirements") || !strings.Contains(stdout, "(nested)") {
		t.Fatalf("expected anchor and nested flag shown, got stdout=%q", stdout)
	}
}

func TestShow_NoExtensionNotShown(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Contains(stdout, "ext:") {
		t.Fatalf("expected no extension line, got stdout=%q", stdout)
	}
}
