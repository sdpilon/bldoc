package cli

import (
	"strings"
	"testing"
)

func TestRename_MissingArguments(t *testing.T) {
	for _, args := range [][]string{{"rename"}, {"rename", "README"}} {
		_, stderr, err := execute(t, args...)
		if err == nil {
			t.Fatalf("expected a usage error for %v", args)
		}
		if stderr == "" {
			t.Fatalf("expected a usage error message on stderr for %v", args)
		}
	}
}

func TestRename_ExtraArgumentRejected(t *testing.T) {
	_, stderr, err := execute(t, "rename", "README", "NOTES", "EXTRA")
	if err == nil {
		t.Fatal("expected a usage error for an extra positional argument")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestRename_UnknownOldTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "rename", "MISSING", "NEW")
	if err == nil {
		t.Fatal("expected an error renaming an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestRename_DuplicateNewNameRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	newTarget(t, "NOTES", "raw")

	_, stderr, err := execute(t, "rename", "README", "NOTES")
	if err == nil {
		t.Fatal("expected an error renaming onto an existing target name")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestRename_Success(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	if _, _, err := execute(t, "add-dep", "README", "a.toml"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	if _, stderr, err := execute(t, "rename", "README", "NOTES"); err != nil {
		t.Fatalf("rename: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if strings.TrimSpace(stdout) != "NOTES" {
		t.Fatalf("expected NOTES to be recorded in place of README, got stdout=%q", stdout)
	}

	stdout, _, err = execute(t, "show", "NOTES")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "a.toml") {
		t.Fatalf("expected renamed target to keep its dependencies, got stdout=%q", stdout)
	}
}
