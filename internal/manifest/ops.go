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

// AddTarget records a new target named name. It rejects a name that
// already exists.
func AddTarget(m *Manifest, name string) error {
	if _, err := findTargetIndex(m, name); err == nil {
		return fmt.Errorf("target %q already exists", name)
	}
	m.Targets = append(m.Targets, Target{Name: name})
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
// unknown target, a dep whose field-addressing (dep.Field set or not)
// doesn't match every dependency the target already has, and a dep
// whose (Source, Path) pair is already recorded on the target.
func AddDep(m *Manifest, target string, dep Dep) error {
	idx, err := findTargetIndex(m, target)
	if err != nil {
		return err
	}
	t := &m.Targets[idx]

	if len(t.Deps) > 0 {
		existingFieldMode := t.Deps[0].Field != ""
		newFieldMode := dep.Field != ""
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
