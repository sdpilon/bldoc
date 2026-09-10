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
	if err := AddTarget(m, "README"); err != nil {
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

	if err := AddTarget(m, "README"); err == nil {
		t.Fatal("expected an error adding a duplicate target, got nil")
	}
	if len(m.Targets) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets)
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

func TestAddDepRejectsDuplicate(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README", Deps: []Dep{{Source: "pyproject.toml"}}}}}

	if err := AddDep(m, "README", Dep{Source: "pyproject.toml"}); err == nil {
		t.Fatal("expected an error adding a duplicate dependency, got nil")
	}
	if len(m.Targets[0].Deps) != 1 {
		t.Fatalf("expected manifest unchanged, got %+v", m.Targets[0].Deps)
	}
}

func TestRemoveDepUnknownTarget(t *testing.T) {
	m := &Manifest{}

	if err := RemoveDep(m, "MISSING", "pyproject.toml", ""); err == nil {
		t.Fatal("expected an error removing a dep from an unknown target, got nil")
	}
}

func TestRemoveDepUnrecorded(t *testing.T) {
	m := &Manifest{Targets: []Target{{Name: "README"}}}

	if err := RemoveDep(m, "README", "other.toml", ""); err == nil {
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

	if err := RemoveDep(m, "README", "b.toml", ""); err != nil {
		t.Fatalf("RemoveDep: %v", err)
	}
	got := m.Targets[0].Deps
	if len(got) != 2 || got[0].Source != "a.toml" || got[1].Source != "c.toml" {
		t.Fatalf("expected [a.toml, c.toml] preserving order, got %+v", got)
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
