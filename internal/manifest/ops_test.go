package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddTargetCreatesManifestFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)

	m, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := AddTarget(m, "README", "", ""); err != nil {
		t.Fatalf("AddTarget: %v", err)
	}
	if err := Save(path, m); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected manifest file to be created: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load after save: %v", err)
	}
	if len(got.Targets) != 1 || got.Targets[0].Name != "README" {
		t.Fatalf("expected one target %q, got %+v", "README", got.Targets)
	}
}

func TestAddTargetRejectsDuplicate(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}}}

	if err := AddTarget(m, "README", "", ""); err == nil {
		t.Fatal("expected an error adding a duplicate target, got nil")
	}
	if len(m.Targets) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets)
	}
}

func TestAddTargetRecordsModeAndExt(t *testing.T) {
	m := &Manifest{}

	if err := AddTarget(m, "PROJECTS", "raw", "yaml"); err != nil {
		t.Fatalf("AddTarget: %v", err)
	}
	if got := m.Targets[0]; got.Mode != "raw" || got.Ext != "yaml" {
		t.Fatalf("expected mode %q and ext %q, got %+v", "raw", "yaml", got)
	}
}

func TestAddTargetNoModeOrExtByDefault(t *testing.T) {
	m := &Manifest{}

	if err := AddTarget(m, "README", "", ""); err != nil {
		t.Fatalf("AddTarget: %v", err)
	}
	if got := m.Targets[0]; got.Mode != "" || got.Ext != "" {
		t.Fatalf("expected no mode or ext, got %+v", got)
	}
}

func TestAddDepUnknownTarget(t *testing.T) {
	m := &Manifest{}

	if err := AddDep(m, "MISSING", Dep{Source: "pyproject.toml"}); err == nil {
		t.Fatal("expected an error adding a dep to an unknown target, got nil")
	}
}

func TestAddDepWholeFile(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}}}

	if err := AddDep(m, "README", Dep{Source: "pyproject.toml"}); err != nil {
		t.Fatalf("AddDep: %v", err)
	}
	if len(m.Targets[0].Deps) != 1 || m.Targets[0].Deps[0].Source != "pyproject.toml" {
		t.Fatalf("expected pyproject.toml recorded, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepFieldAddressedWithFormat(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}}}
	dep := Dep{
		Source: "pyproject.toml",
		Path:   "project.requires-python",
		Field:  "version",
		Format: "Python version must be %s to run this project.",
	}

	if err := AddDep(m, "README", dep); err != nil {
		t.Fatalf("AddDep: %v", err)
	}
	got := m.Targets[0].Deps[0]
	if got != dep {
		t.Fatalf("recorded dep mismatch: got %+v, want %+v", got, dep)
	}
}

func TestAddDepRejectsFieldModeOnRawTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README", Deps: []Dep{{Source: "a.toml"}}}}}

	err := AddDep(m, "README", Dep{Source: "b.toml", Path: "x", Field: "version"})
	if err == nil {
		t.Fatal("expected an error mixing field-mode dep into a raw-mode target, got nil")
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsRawModeOnFieldTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README", Deps: []Dep{{Source: "a.toml", Path: "x", Field: "version"}}}}}

	err := AddDep(m, "README", Dep{Source: "b.toml"})
	if err == nil {
		t.Fatal("expected an error mixing a raw dep into a field-mode target, got nil")
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsFieldModeAsFirstDepOnDeclaredRawTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "raw"}}}

	err := AddDep(m, "PROJECTS", Dep{Source: "entry.yaml", Path: "name", Field: "name"})
	if err == nil {
		t.Fatal("expected an error adding a field-addressed dep as the first dep of a declared raw-mode target, got nil")
	}
	if len(m.Targets[0].Deps) != 0 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsRawModeAsFirstDepOnDeclaredFieldTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "field"}}}

	err := AddDep(m, "PROJECTS", Dep{Source: "entry.yaml"})
	if err == nil {
		t.Fatal("expected an error adding a whole-file dep as the first dep of a declared field-mode target, got nil")
	}
	if len(m.Targets[0].Deps) != 0 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepMatchesDeclaredModeAccepted(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "field"}}}

	if err := AddDep(m, "PROJECTS", Dep{Source: "entry.yaml", Path: "name", Field: "name"}); err != nil {
		t.Fatalf("AddDep: %v", err)
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected one dep recorded, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsBareRefOnListModeTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "list"}}}

	err := AddDep(m, "PROJECTS", Dep{Source: "entry.yaml"})
	if err == nil {
		t.Fatal("expected an error adding a bare (non-record-scoped) dep to a list-mode target, got nil")
	}
	if len(m.Targets[0].Deps) != 0 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepAcceptsRecordOnlyOnListModeTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "list"}}}

	if err := AddDep(m, "PROJECTS", Dep{Source: "project-headers/bldoc.yaml", Field: "bldoc"}); err != nil {
		t.Fatalf("AddDep: %v", err)
	}
	if len(m.Targets[0].Deps) != 1 || m.Targets[0].Deps[0].Field != "bldoc" {
		t.Fatalf("expected record-only dep recorded, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepAcceptsRecordFieldOnListModeTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "list"}}}

	if err := AddDep(m, "PROJECTS", Dep{Source: "entry-description.md", Field: "bldoc.description"}); err != nil {
		t.Fatalf("AddDep: %v", err)
	}
	if len(m.Targets[0].Deps) != 1 || m.Targets[0].Deps[0].Field != "bldoc.description" {
		t.Fatalf("expected record.field dep recorded, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsOverQualifiedFieldOnListModeTarget(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "list"}}}

	err := AddDep(m, "PROJECTS", Dep{Source: "entry-description.md", Field: "bldoc.description.extra"})
	if err == nil {
		t.Fatal("expected an error adding an over-qualified field to a list-mode target, got nil")
	}
	if len(m.Targets[0].Deps) != 0 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsFormatOnRecordOnlyDep(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "list"}}}

	err := AddDep(m, "PROJECTS", Dep{Source: "project-headers/bldoc.yaml", Field: "bldoc", Format: "%s"})
	if err == nil {
		t.Fatal("expected an error adding --format to a record-only dep, got nil")
	}
	if len(m.Targets[0].Deps) != 0 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepAcceptsFormatOnRecordFieldDep(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "PROJECTS", Mode: "list"}}}

	if err := AddDep(m, "PROJECTS", Dep{Source: "version.txt", Field: "bldoc.version", Format: "v%s"}); err != nil {
		t.Fatalf("AddDep: %v", err)
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected one dep recorded, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsDuplicate(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README", Deps: []Dep{{Source: "pyproject.toml"}}}}}

	if err := AddDep(m, "README", Dep{Source: "pyproject.toml"}); err == nil {
		t.Fatal("expected an error adding a duplicate dependency, got nil")
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepRejectsDuplicateAnchorRegardlessOfNested(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README", Deps: []Dep{
		{Source: "spec.md", Anchor: "requirements", Field: "section"},
	}}}}

	if err := AddDep(m, "README", Dep{Source: "spec.md", Anchor: "requirements", Field: "section", Nested: true}); err == nil {
		t.Fatal("expected an error adding a duplicate anchor dependency that only differs by --nested, got nil")
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestAddDepAcceptsDifferentAnchorsOnSameSource(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README", Deps: []Dep{
		{Source: "spec.md", Anchor: "purpose", Field: "a"},
	}}}}

	if err := AddDep(m, "README", Dep{Source: "spec.md", Anchor: "requirements", Field: "b"}); err != nil {
		t.Fatalf("expected a different anchor on the same source to be accepted, got %v", err)
	}
	if len(m.Targets[0].Deps) != 2 {
		t.Fatalf("expected both dependencies recorded, got %+v", m.Targets[0].Deps)
	}
}

func TestRemoveDepUnknownTarget(t *testing.T) {
	m := &Manifest{}

	if err := RemoveDep(m, "MISSING", "pyproject.toml", "", ""); err == nil {
		t.Fatal("expected an error removing a dep from an unknown target, got nil")
	}
}

func TestRemoveDepUnrecorded(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}}}

	if err := RemoveDep(m, "README", "other.toml", "", ""); err == nil {
		t.Fatal("expected an error removing an unrecorded dependency, got nil")
	}
}

func TestRemoveDepPreservesOrder(t *testing.T) {
	m := &Manifest{Targets: []Target{{
		Name: "README",
		Deps: []Dep{
			{Source: "a.toml"},
			{Source: "b.toml"},
			{Source: "c.toml"},
		},
	}}}

	if err := RemoveDep(m, "README", "b.toml", "", ""); err != nil {
		t.Fatalf("RemoveDep: %v", err)
	}
	got := m.Targets[0].Deps
	if len(got) != 2 || got[0].Source != "a.toml" || got[1].Source != "c.toml" {
		t.Fatalf("expected [a.toml, c.toml] preserving order, got %+v", got)
	}
}

func TestRemoveDepMatchesByAnchor(t *testing.T) {
	m := &Manifest{Targets: []Target{{
		Name: "README",
		Deps: []Dep{
			{Source: "spec.md", Anchor: "purpose", Field: "a"},
			{Source: "spec.md", Anchor: "requirements", Field: "b"},
		},
	}}}

	if err := RemoveDep(m, "README", "spec.md", "", "purpose"); err != nil {
		t.Fatalf("RemoveDep: %v", err)
	}
	got := m.Targets[0].Deps
	if len(got) != 1 || got[0].Anchor != "requirements" {
		t.Fatalf("expected only the requirements-anchor dependency to remain, got %+v", got)
	}
}

func TestRemoveTargetUnknown(t *testing.T) {
	m := &Manifest{}

	if err := RemoveTarget(m, "MISSING"); err == nil {
		t.Fatal("expected an error removing an unknown target, got nil")
	}
}

func TestRemoveTargetDeletesTargetAndDeps(t *testing.T) {
	m := &Manifest{Targets: []Target{
		{Name: "README", Deps: []Dep{{Source: "a.toml"}}},
		{Name: "CHANGELOG"},
	}}

	if err := RemoveTarget(m, "README"); err != nil {
		t.Fatalf("RemoveTarget: %v", err)
	}
	if len(m.Targets) != 1 || m.Targets[0].Name != "CHANGELOG" {
		t.Fatalf("expected only CHANGELOG to remain, got %+v", m.Targets)
	}
}

func TestListTargetsDeclarationOrder(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}, {Name: "CHANGELOG"}}}

	got := ListTargets(m)
	if len(got) != 2 || got[0] != "README" || got[1] != "CHANGELOG" {
		t.Fatalf("expected [README, CHANGELOG], got %v", got)
	}
}

func TestListTargetsEmpty(t *testing.T) {
	m := &Manifest{}

	got := ListTargets(m)
	if len(got) != 0 {
		t.Fatalf("expected no targets, got %v", got)
	}
}

func TestRenameTargetUnknownOld(t *testing.T) {
	m := &Manifest{}

	if err := RenameTarget(m, "MISSING", "NEW"); err == nil {
		t.Fatal("expected an error renaming an unknown target, got nil")
	}
}

func TestRenameTargetDuplicateNewRejected(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}, {Name: "NOTES"}}}

	if err := RenameTarget(m, "README", "NOTES"); err == nil {
		t.Fatal("expected an error renaming onto an existing target name, got nil")
	}
	if m.Targets[0].Name != "README" || m.Targets[1].Name != "NOTES" {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets)
	}
}

func TestRenameTargetPreservesDepsOrderModeExt(t *testing.T) {
	m := &Manifest{Targets: []Target{{
		Name: "PROJECTS",
		Mode: "raw",
		Ext:  "yaml",
		Deps: []Dep{{Source: "a.yaml"}, {Source: "b.yaml"}},
	}}}

	if err := RenameTarget(m, "PROJECTS", "REGISTRY"); err != nil {
		t.Fatalf("RenameTarget: %v", err)
	}
	if len(m.Targets) != 1 {
		t.Fatalf("expected exactly one target, got %+v", m.Targets)
	}
	got := m.Targets[0]
	if got.Name != "REGISTRY" || got.Mode != "raw" || got.Ext != "yaml" {
		t.Fatalf("expected renamed target with mode/ext preserved, got %+v", got)
	}
	if len(got.Deps) != 2 || got.Deps[0].Source != "a.yaml" || got.Deps[1].Source != "b.yaml" {
		t.Fatalf("expected deps preserved in order, got %+v", got.Deps)
	}
}

func TestShowTargetUnknown(t *testing.T) {
	m := &Manifest{}

	if _, err := ShowTarget(m, "MISSING"); err == nil {
		t.Fatal("expected an error showing an unknown target, got nil")
	}
}

func TestShowTargetReturnsDepsInAddedOrder(t *testing.T) {
	m := &Manifest{Targets: []Target{{
		Name: "README",
		Deps: []Dep{
			{Source: "a.toml"},
			{Source: "b.toml", Path: "x", Field: "version", Format: "v%s"},
		},
	}}}

	got, err := ShowTarget(m, "README")
	if err != nil {
		t.Fatalf("ShowTarget: %v", err)
	}
	if len(got.Deps) != 2 || got.Deps[0].Source != "a.toml" || got.Deps[1].Source != "b.toml" {
		t.Fatalf("expected deps in added order, got %+v", got.Deps)
	}
}
