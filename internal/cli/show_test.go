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
	newTarget(t, "README")
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
	lines := strings.Fields(stdout)
	if len(lines) != 2 || lines[0] != "a.toml" || lines[1] != "b.toml" {
		t.Fatalf("expected [a.toml b.toml] in added order, got %v", lines)
	}
}
