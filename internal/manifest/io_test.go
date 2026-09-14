package manifest

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadSaveRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	want := &Manifest{
		Targets: []Target{
			{
				Name: "README",
				Deps: []Dep{
					{Source: "pyproject.toml"},
					{Source: "pyproject.toml", Path: "project.version", Field: "version", Format: "v%s"},
				},
			},
		},
	}

	if err := Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip mismatch:\n got: %+v\nwant: %+v", got, want)
	}
}

func TestLoadSaveRoundTrip_ModeAndExt(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	want := &Manifest{
		Targets: []Target{
			{Name: "PROJECTS", Mode: "raw", Ext: "yaml", Deps: []Dep{{Source: "a.yaml"}}},
			{Name: "README", Deps: []Dep{{Source: "pyproject.toml"}}},
		},
	}

	if err := Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip mismatch:\n got: %+v\nwant: %+v", got, want)
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Targets) != 0 {
		t.Fatalf("expected empty manifest for missing file, got %+v", got)
	}
}

func TestLoadMalformedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	if err := os.WriteFile(path, []byte("this is not [ valid toml"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("expected an error parsing a malformed manifest, got nil")
	}
}
