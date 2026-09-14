package manifest

import "fmt"

func findTargetIndex(m *Manifest, name string) (int, error) {
	for i, t := range m.Targets {
		if t.Name == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("target %q not found", name)
}

// AddTarget records a new target named name, with the given mode ("raw",
// "field", or "" for no declared mode) and ext (meaningful for "raw"
// only; "" for none). It rejects a name that already exists. Callers are
// responsible for validating mode and the mode/ext combination before
// calling this — AddTarget stores whatever it's given.
func AddTarget(m *Manifest, name, mode, ext string) error {
	if _, err := findTargetIndex(m, name); err == nil {
		return fmt.Errorf("target %q already exists", name)
	}
	m.Targets = append(m.Targets, Target{Name: name, Mode: mode, Ext: ext})
	return nil
}

// RemoveTarget deletes target and all of its recorded dependencies. It
// rejects a name that does not exist.
func RemoveTarget(m *Manifest, name string) error {
	idx, err := findTargetIndex(m, name)
	if err != nil {
		return err
	}
	m.Targets = append(m.Targets[:idx], m.Targets[idx+1:]...)
	return nil
}

// ListTargets returns every target's name, in declaration order.
func ListTargets(m *Manifest) []string {
	names := make([]string, len(m.Targets))
	for i, t := range m.Targets {
		names[i] = t.Name
	}
	return names
}

// ShowTarget returns the named target, including its recorded
// dependencies in the order they were added. It rejects a name that
// does not exist.
func ShowTarget(m *Manifest, name string) (Target, error) {
	idx, err := findTargetIndex(m, name)
	if err != nil {
		return Target{}, err
	}
	return m.Targets[idx], nil
}

// AddDep appends dep to target's recorded dependencies. It rejects an
// unknown target and a dep whose (Source, Path) pair is already recorded
// on the target. For mode-exclusivity: a target with a declared Mode
// rejects any dependency (including the first) whose field-addressing
// doesn't match that mode; a target with no declared Mode rejects a dep
// whose field-addressing doesn't match every dependency it already has,
// once it has at least one.
func AddDep(m *Manifest, target string, dep Dep) error {
	idx, err := findTargetIndex(m, target)
	if err != nil {
		return err
	}
	t := &m.Targets[idx]

	newFieldMode := dep.Field != ""
	if t.Mode != "" {
		if wantField := t.Mode == "field"; wantField != newFieldMode {
			return fmt.Errorf("target %q is declared mode %q, which doesn't accept this dependency's field-addressing", target, t.Mode)
		}
	} else if len(t.Deps) > 0 {
		existingFieldMode := t.Deps[0].Field != ""
		if existingFieldMode != newFieldMode {
			return fmt.Errorf("target %q mixes raw and field-mode dependencies", target)
		}
	}

	for _, d := range t.Deps {
		if d.Source == dep.Source && d.Path == dep.Path {
			return fmt.Errorf("dependency %q already recorded on target %q", dep.Source, target)
		}
	}

	t.Deps = append(t.Deps, dep)
	return nil
}

// RenameTarget renames the target named oldName to newName, preserving
// its recorded dependencies, their order, its declared Mode, and its
// Ext unchanged. It rejects an oldName that does not exist and a newName
// that already names a target.
func RenameTarget(m *Manifest, oldName, newName string) error {
	idx, err := findTargetIndex(m, oldName)
	if err != nil {
		return err
	}
	if _, err := findTargetIndex(m, newName); err == nil {
		return fmt.Errorf("target %q already exists", newName)
	}
	m.Targets[idx].Name = newName
	return nil
}

// RemoveDep removes the dependency matching (source, path) from
// target's recorded dependencies, preserving the order of the rest. It
// rejects an unknown target and a (source, path) not recorded on it.
func RemoveDep(m *Manifest, target, source, path string) error {
	idx, err := findTargetIndex(m, target)
	if err != nil {
		return err
	}
	t := &m.Targets[idx]

	for i, d := range t.Deps {
		if d.Source == source && d.Path == path {
			t.Deps = append(t.Deps[:i], t.Deps[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("dependency %q not recorded on target %q", source, target)
}
