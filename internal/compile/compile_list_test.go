package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bldoc/internal/manifest"
)

func TestTarget_EmptyDeclaredListModeCompilesToEmptyArray(t *testing.T) {
	result, err := Target(manifest.Target{Name: "PROJECTS", Mode: "list"})
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if !result.IsList {
		t.Fatal("expected list-mode result for an empty declared-list target")
	}
	if len(result.RecordOrder) != 0 || len(result.Records) != 0 {
		t.Fatalf("expected no records, got order=%v records=%+v", result.RecordOrder, result.Records)
	}
}

func TestTarget_ListModeMergesWholeDocument(t *testing.T) {
	dir := t.TempDir()
	headerPath := filepath.Join(dir, "bldoc.yaml")
	if err := os.WriteFile(headerPath, []byte("path: ~/Projects/_Claude/bldoc\nremote: git@github.com:sdpilon/bldoc.git\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: headerPath, Field: "bldoc"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	record, ok := result.Records["bldoc"]
	if !ok {
		t.Fatalf("expected record %q, got %+v", "bldoc", result.Records)
	}
	if record["path"] != "~/Projects/_Claude/bldoc" || record["remote"] != "git@github.com:sdpilon/bldoc.git" {
		t.Fatalf("expected merged keys, got %+v", record)
	}
}

func TestTarget_ListModeResolvesRecordField(t *testing.T) {
	dir := t.TempDir()
	descPath := filepath.Join(dir, "entry-description.md")
	if err := os.WriteFile(descPath, []byte("a Go CLI tool"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: descPath, Field: "bldoc.description"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	record := result.Records["bldoc"]
	if record["description"] != "a Go CLI tool" {
		t.Fatalf("expected description field, got %+v", record)
	}
}

func TestTarget_ListModeResolvesRecordFieldAnchor(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Title\n\n## Purpose\n\na Go CLI tool\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: specPath, Anchor: "purpose", Field: "bldoc.description"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	record := result.Records["bldoc"]
	got, _ := record["description"].(string)
	if !strings.Contains(got, "a Go CLI tool") {
		t.Fatalf("expected description field from the anchor's section, got %+v", record)
	}
}

func TestTarget_ListModeAppliesFormatTemplate(t *testing.T) {
	dir := t.TempDir()
	versionPath := filepath.Join(dir, "version.toml")
	if err := os.WriteFile(versionPath, []byte("version = \"1.2\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: versionPath, Path: "version", Field: "bldoc.version", Format: "v%s"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if got := result.Records["bldoc"]["version"]; got != "v1.2" {
		t.Fatalf("expected v1.2, got %v", got)
	}
}

func TestTarget_ListModeFieldsAreUnwrapped(t *testing.T) {
	dir := t.TempDir()
	descPath := filepath.Join(dir, "entry-description.md")
	if err := os.WriteFile(descPath, []byte("a Go CLI tool"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: descPath, Field: "bldoc.description"}},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if _, ok := result.Records["bldoc"]["description"].(Field); ok {
		t.Fatal("expected a plain value, not a wrapped Field{value, raw}")
	}
}

func TestTarget_ListModeOrdersRecordsByFirstAppearance(t *testing.T) {
	dir := t.TempDir()
	descPath := filepath.Join(dir, "entry-description.md")
	if err := os.WriteFile(descPath, []byte("desc"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{
			{Source: descPath, Field: "bldoc.description"},
			{Source: descPath, Field: "projx.description"},
		},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if len(result.RecordOrder) != 2 || result.RecordOrder[0] != "bldoc" || result.RecordOrder[1] != "projx" {
		t.Fatalf("expected order [bldoc projx], got %v", result.RecordOrder)
	}
}

func TestTarget_ListModeRejectsNonScalarMerge(t *testing.T) {
	dir := t.TempDir()
	headerPath := filepath.Join(dir, "bldoc.yaml")
	if err := os.WriteFile(headerPath, []byte("nested:\n  x: y\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: headerPath, Field: "bldoc"}},
	}
	if _, err := Target(target); err == nil {
		t.Fatal("expected an error merging a document with a non-scalar top-level value")
	}
}

func TestTarget_ListModeRejectsMalformedRecordOnlySource(t *testing.T) {
	dir := t.TempDir()
	headerPath := filepath.Join(dir, "bldoc.yaml")
	if err := os.WriteFile(headerPath, []byte(": not valid yaml :::"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{{Source: headerPath, Field: "bldoc"}},
	}
	if _, err := Target(target); err == nil {
		t.Fatal("expected an error for a malformed record-only source")
	}
}

func TestTarget_ListModeRejectsDuplicateFieldWithinRecord(t *testing.T) {
	dir := t.TempDir()
	descPath := filepath.Join(dir, "entry-description.md")
	if err := os.WriteFile(descPath, []byte("desc"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{
			{Source: descPath, Field: "bldoc.description"},
			{Source: descPath, Field: "bldoc.description"},
		},
	}
	if _, err := Target(target); err == nil {
		t.Fatal("expected an error for a duplicate field name within the same record")
	}
}

func TestTarget_ListModeRejectsMergeCollisionWithExplicitField(t *testing.T) {
	dir := t.TempDir()
	headerPath := filepath.Join(dir, "bldoc.yaml")
	if err := os.WriteFile(headerPath, []byte("description: from header\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	descPath := filepath.Join(dir, "entry-description.md")
	if err := os.WriteFile(descPath, []byte("from file"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{
			{Source: headerPath, Field: "bldoc"},
			{Source: descPath, Field: "bldoc.description"},
		},
	}
	if _, err := Target(target); err == nil {
		t.Fatal("expected an error when a merged key collides with an explicit field")
	}
}

func TestTarget_ListModeAllowsSameFieldNameInDifferentRecords(t *testing.T) {
	dir := t.TempDir()
	descPath := filepath.Join(dir, "entry-description.md")
	if err := os.WriteFile(descPath, []byte("desc"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	target := manifest.Target{
		Name: "PROJECTS",
		Mode: "list",
		Deps: []manifest.Dep{
			{Source: descPath, Field: "bldoc.description"},
			{Source: descPath, Field: "projx.description"},
		},
	}
	result, err := Target(target)
	if err != nil {
		t.Fatalf("Target: %v", err)
	}
	if result.Records["bldoc"]["description"] != "desc" || result.Records["projx"]["description"] != "desc" {
		t.Fatalf("expected both records to carry description, got %+v", result.Records)
	}
}
