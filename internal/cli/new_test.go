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

	if _, stderr, err := execute(t, "new", "README", "--mode", "raw"); err != nil {
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
	if _, _, err := execute(t, "new", "README", "--mode", "raw"); err != nil {
		t.Fatalf("first new: %v", err)
	}

	_, stderr, err := execute(t, "new", "README", "--mode", "raw")
	if err == nil {
		t.Fatal("expected an error creating a duplicate target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestNew_MissingModeRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "new", "README")
	if err == nil {
		t.Fatal("expected an error when --mode is omitted")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}

	stdout, _, _ := execute(t, "list")
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected no target recorded, got stdout=%q", stdout)
	}
}

func TestNew_InvalidModeRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "new", "README", "--mode", "nonsense")
	if err == nil {
		t.Fatal("expected an error for an invalid --mode value")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestNew_ExtRequiresRawMode(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "new", "PROJECTS", "--mode", "field", "--ext", "yaml")
	if err == nil {
		t.Fatal("expected an error combining --ext with --mode field")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}

	stdout, _, _ := execute(t, "list")
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected no target recorded, got stdout=%q", stdout)
	}
}

func TestNew_ExtRecordedWithRawMode(t *testing.T) {
	t.Chdir(t.TempDir())

	if _, stderr, err := execute(t, "new", "PROJECTS", "--mode", "raw", "--ext", "yaml"); err != nil {
		t.Fatalf("expected new to succeed, got err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "PROJECTS")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "yaml") {
		t.Fatalf("expected extension yaml to be shown, got stdout=%q", stdout)
	}
}
