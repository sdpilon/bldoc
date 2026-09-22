package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddDepBatch_TwoRawPairs(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")

	if _, stderr, err := execute(t, "add-dep", "README", "pyproject.toml", "go.mod"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "pyproject.toml") || !strings.Contains(stdout, "go.mod") {
		t.Fatalf("expected both dependencies recorded, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_TwoFieldPairs(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README", "version:pyproject.toml@project.version", "summary:spec.md#purpose"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "version <- pyproject.toml@project.version") {
		t.Fatalf("expected version field recorded, got stdout=%q", stdout)
	}
	if !strings.Contains(stdout, "summary <- spec.md#purpose") {
		t.Fatalf("expected summary field recorded, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_ListModeRecordWithTwoFieldPairs(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "PROJECTS", "list")

	if _, stderr, err := execute(t, "add-dep", "PROJECTS", "bldoc", "description:entry-description.md", "path:project-headers/bldoc.yaml"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "PROJECTS")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "bldoc.description <- entry-description.md") {
		t.Fatalf("expected record field 'description' recorded, got stdout=%q", stdout)
	}
	if !strings.Contains(stdout, "bldoc.path <- project-headers/bldoc.yaml") {
		t.Fatalf("expected record field 'path' recorded, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_ListModeRecordOnlyPairAlongsideFieldPair(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "PROJECTS", "list")

	if _, stderr, err := execute(t, "add-dep", "PROJECTS", "bldoc", "defaults.yaml", "description:entry-description.md"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "PROJECTS")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "bldoc <- defaults.yaml") {
		t.Fatalf("expected record-only dependency recorded, got stdout=%q", stdout)
	}
	if !strings.Contains(stdout, "bldoc.description <- entry-description.md") {
		t.Fatalf("expected record field dependency recorded, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_AtInTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")

	_, stderr, err := execute(t, "add-dep", "README@version", "pyproject.toml", "go.mod")
	if err == nil {
		t.Fatal("expected an error for a target argument containing '@'")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestAddDepBatch_FewerThanTwoPairsRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "PROJECTS", "list")

	_, stderr, err := execute(t, "add-dep", "PROJECTS", "bldoc", "description:entry-description.md")
	if err == nil {
		t.Fatal("expected an error for fewer than two pairs after the record name")
	}
	if stderr == "" {
		t.Fatal("expected a usage error message on stderr")
	}
}

func TestAddDepBatch_UnknownTargetRejected(t *testing.T) {
	t.Chdir(t.TempDir())

	_, stderr, err := execute(t, "add-dep", "MISSING", "a.toml", "b.toml")
	if err == nil {
		t.Fatal("expected an error for an unknown target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}
}

func TestAddDepBatch_ModeExclusivityRejectsWholeBatch(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	stdoutBefore, _, _ := execute(t, "show", "README")

	_, stderr, err := execute(t, "add-dep", "README", "version:pyproject.toml@project.version", "go.mod")
	if err == nil {
		t.Fatal("expected a mode-exclusivity error mixing a field pair and a bare pair on a field-mode target")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}

	stdoutAfter, _, _ := execute(t, "show", "README")
	if stdoutBefore != stdoutAfter {
		t.Fatalf("expected manifest unchanged after a failed batch call; before=%q after=%q", stdoutBefore, stdoutAfter)
	}
}

func TestAddDepBatch_DuplicateLaterPairAbortsEntireBatch(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	newTarget(t, "README", "raw")
	if _, _, err := execute(t, "add-dep", "README", "go.mod"); err != nil {
		t.Fatalf("seed add-dep: %v", err)
	}

	manifestBefore, err := os.ReadFile(filepath.Join(dir, "bldoc.toml"))
	if err != nil {
		t.Fatalf("reading manifest before: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README", "pyproject.toml", "go.mod")
	if err == nil {
		t.Fatal("expected a duplicate-dependency error")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}

	manifestAfter, err := os.ReadFile(filepath.Join(dir, "bldoc.toml"))
	if err != nil {
		t.Fatalf("reading manifest after: %v", err)
	}
	if string(manifestBefore) != string(manifestAfter) {
		t.Fatalf("expected manifest byte-for-byte unchanged after a failed batch call\nbefore=%q\nafter=%q", manifestBefore, manifestAfter)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Contains(stdout, "pyproject.toml") {
		t.Fatalf("expected pyproject.toml not recorded despite being individually valid, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_DuplicatePairWithinSameCallRejected(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "raw")

	_, stderr, err := execute(t, "add-dep", "README", "pyproject.toml", "pyproject.toml")
	if err == nil {
		t.Fatal("expected a duplicate-dependency error for the same source-ref listed twice")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Contains(stdout, "pyproject.toml") {
		t.Fatalf("expected neither occurrence recorded, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_AmbiguousAnchorWarnsPerPair(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	newTarget(t, "README", "field")
	specContent := "# Title\n\n### Requirement: A\n\n#### Scenario: Dup\nbody\n\n### Requirement: B\n\n#### Scenario: Dup\nbody\n\n### Requirement: C\n\n#### Scenario: Other\nbody\n"
	if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, stderr, err := execute(t, "add-dep", "README", "first:spec.md#scenario-dup", "second:spec.md#scenario-other")
	if err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}
	if !strings.Contains(stderr, "scenario-dup") {
		t.Fatalf("expected an ambiguity warning for the ambiguous pair, got stderr=%q", stderr)
	}
	if strings.Contains(stderr, "scenario-other") {
		t.Fatalf("expected no warning for the unambiguous pair, got stderr=%q", stderr)
	}
}

func TestAddDepBatch_NestedAppliesToEveryPair(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README", "section:spec.md#requirements", "sub:spec.md#configuration", "--nested"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Count(stdout, "(nested)") != 2 {
		t.Fatalf("expected both dependencies recorded with --nested, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_FormatAppliesToEveryPair(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README", "--format", "%s (pinned)", "version:pyproject.toml@project.version", "tool:pyproject.toml@project.name"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Count(stdout, `(format: "%s (pinned)")`) != 2 {
		t.Fatalf("expected both dependencies recorded with the same format template, got stdout=%q", stdout)
	}
}

func TestAddDepBatch_FormatRejectedOnRecordOnlyPairInBatch(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "PROJECTS", "list")

	_, stderr, err := execute(t, "add-dep", "PROJECTS", "bldoc", "--format", "%s", "defaults.yaml", "description:entry-description.md")
	if err == nil {
		t.Fatal("expected an error: --format is rejected on a record-only pair within a batch call")
	}
	if stderr == "" {
		t.Fatal("expected an error message on stderr")
	}

	stdout, _, err := execute(t, "show", "PROJECTS")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if strings.Contains(stdout, "defaults.yaml") || strings.Contains(stdout, "entry-description.md") {
		t.Fatalf("expected manifest unchanged, got stdout=%q", stdout)
	}
}

func TestAddDep_TwoArgumentFormStillUsesSingleDependencyGrammar(t *testing.T) {
	t.Chdir(t.TempDir())
	newTarget(t, "README", "field")

	if _, stderr, err := execute(t, "add-dep", "README@version", "pyproject.toml@project.version"); err != nil {
		t.Fatalf("add-dep: err=%v stderr=%q", err, stderr)
	}

	stdout, _, err := execute(t, "show", "README")
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(stdout, "version <- pyproject.toml@project.version") {
		t.Fatalf("expected the two-argument grammar to still apply, got stdout=%q", stdout)
	}
}
