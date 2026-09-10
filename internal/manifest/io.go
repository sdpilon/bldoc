package manifest

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// FileName is the manifest's fixed file name, resolved relative to the
// current working directory.
const FileName = "bldoc.toml"

// Load reads and parses the manifest at path. A missing file is not an
// error: it yields an empty Manifest, so callers that only read (list,
// show) behave correctly against a project with no manifest yet, and
// callers that mutate (new, add-dep) build on top of it via Save.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Manifest{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading manifest %s: %w", path, err)
	}
	var m Manifest
	if _, err := toml.Decode(string(data), &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %s: %w", path, err)
	}
	return &m, nil
}

// Save writes m to path, creating or truncating it as needed.
func Save(path string, m *Manifest) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("writing manifest %s: %w", path, err)
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(m); err != nil {
		return fmt.Errorf("encoding manifest %s: %w", path, err)
	}
	return nil
}
