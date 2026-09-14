package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bldoc/internal/manifest"
)

func TestList_ExtraArgumentRejected(t *testing.T) {
	_, stderr, err := execute(t, "list", "README")
	if err == nil {
		t.Fatal("expected a usage error for an extra positional argument")
	}
	if strings.Contains(stderr, "not yet implemented") {
		t.Fatalf("expected a usage error, not the not-yet-implemented stub, got stderr=%q", stderr)
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestList_NoTargets(t *testing.T) {
	t.Chdir(t.TempDir())

	stdout, _, err := execute(t, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected no target names, got stdout=%q", stdout)
	}
}

func TestList_TargetsListed(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")
	newTarget(t, "CHANGELOG", "raw")

	stdout, _, err := execute(t, "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	lines := strings.Fields(stdout)
	if len(lines) != 2 || lines[0] != "README" || lines[1] != "CHANGELOG" {
		t.Fatalf("expected [README CHANGELOG] in declaration order, got %v", lines)
	}
}

func TestList_MalformedManifestRejected(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	path := filepath.Join(dir, manifest.FileName)
	if err := os.WriteFile(path, []byte("not [ valid toml"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, stderr, err := execute(t, "list")
	if err == nil {
		t.Fatal("expected an error reading a malformed manifest")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}
