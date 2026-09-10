package cli

import (
	"strings"
	"testing"
)

func TestRoot_VersionFlag(t *testing.T) {
	stdout, _, err := execute(t, "--version")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(stdout, "bldoc") {
		t.Fatalf("expected version output to mention bldoc, got %q", stdout)
	}
}
