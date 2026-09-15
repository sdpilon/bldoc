package cli

import (
	"runtime/debug"
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

func TestBuildVersion_NoBuildInfo(t *testing.T) {
	got := buildVersion("0.0.0-dev", nil)
	if got != "0.0.0-dev" {
		t.Fatalf("expected plain version unchanged, got %q", got)
	}
}

func TestBuildVersion_NoVCSSettings(t *testing.T) {
	info := &debug.BuildInfo{}
	got := buildVersion("0.0.0-dev", info)
	if got != "0.0.0-dev" {
		t.Fatalf("expected plain version unchanged, got %q", got)
	}
}

func TestBuildVersion_WithCommitAndTime(t *testing.T) {
	info := &debug.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "vcs.revision", Value: "29e081c95c2287acb4d0be465b20bf334100e7d6"},
		{Key: "vcs.time", Value: "2026-09-15T02:10:19Z"},
		{Key: "vcs.modified", Value: "false"},
	}}
	got := buildVersion("0.0.0-dev", info)
	want := "0.0.0-dev (commit 29e081c95c22, 2026-09-15T02:10:19Z)"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildVersion_ModifiedFlag(t *testing.T) {
	info := &debug.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "vcs.revision", Value: "29e081c95c2287acb4d0be465b20bf334100e7d6"},
		{Key: "vcs.time", Value: "2026-09-15T02:10:19Z"},
		{Key: "vcs.modified", Value: "true"},
	}}
	got := buildVersion("0.0.0-dev", info)
	if !strings.HasSuffix(got, ", modified)") {
		t.Fatalf("expected a trailing modified flag, got %q", got)
	}
}
