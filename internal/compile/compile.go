package compile

import (
	"bytes"
	"fmt"
	"os"

	"bldoc/internal/manifest"
)

// Field is one resolved field-mode entry in a target's compiled
// intermediate.
type Field struct {
	Value string `json:"value"`
	Raw   string `json:"raw"`
}

// Result is a target's compiled intermediate: either RawBytes (raw
// mode, including the empty-target case) or Fields (field mode).
type Result struct {
	IsField bool
	Raw     []byte
	Fields  map[string]Field
}

// Target compiles target into its intermediate. A target with no
// dependencies compiles to an empty raw-mode result; otherwise its mode
// is determined by whether its first dependency names a field, matching
// the raw/field exclusivity the manifest already enforces.
func Target(target manifest.Target) (Result, error) {
	if len(target.Deps) == 0 {
		return Result{Raw: []byte{}}, nil
	}
	if target.Deps[0].Field == "" {
		raw, err := compileRaw(target)
		if err != nil {
			return Result{}, err
		}
		return Result{Raw: raw}, nil
	}
	fields, err := compileFields(target)
	if err != nil {
		return Result{}, err
	}
	return Result{IsField: true, Fields: fields}, nil
}

// compileRaw concatenates target's whole-file dependencies' bytes, in
// the order they were added.
func compileRaw(target manifest.Target) ([]byte, error) {
	var buf bytes.Buffer
	for _, dep := range target.Deps {
		data, err := os.ReadFile(dep.Source)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", dep.Source, err)
		}
		buf.Write(data)
	}
	return buf.Bytes(), nil
}

// compileFields resolves target's field-mode dependencies into a map of
// field name to Field, after rejecting duplicate field names.
func compileFields(target manifest.Target) (map[string]Field, error) {
	if err := checkDuplicateFields(target); err != nil {
		return nil, err
	}

	fields := make(map[string]Field, len(target.Deps))
	for _, dep := range target.Deps {
		raw, err := resolveDepValue(dep)
		if err != nil {
			return nil, err
		}
		value := raw
		if dep.Format != "" {
			value = fmt.Sprintf(dep.Format, raw)
		}
		fields[dep.Field] = Field{Value: value, Raw: raw}
	}
	return fields, nil
}

func checkDuplicateFields(target manifest.Target) error {
	seen := make(map[string]bool, len(target.Deps))
	for _, dep := range target.Deps {
		if seen[dep.Field] {
			return fmt.Errorf("target %q has more than one dependency naming field %q", target.Name, dep.Field)
		}
		seen[dep.Field] = true
	}
	return nil
}

// resolveDepValue returns dep's raw string value: the whole source
// file's content when dep has no field-path, or the field-path's
// resolved scalar value when it does.
func resolveDepValue(dep manifest.Dep) (string, error) {
	if dep.Path == "" {
		data, err := os.ReadFile(dep.Source)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", dep.Source, err)
		}
		return string(data), nil
	}

	tree, err := parseSource(dep.Source)
	if err != nil {
		return "", err
	}
	v, err := resolveFieldPath(tree, dep.Path)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(v), nil
}
