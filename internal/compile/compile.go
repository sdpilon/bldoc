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

// Result is a target's compiled intermediate: RawBytes (raw mode,
// including the empty-target case), Fields (field mode), or an ordered
// set of Records (list mode).
type Result struct {
	IsField     bool
	IsList      bool
	Raw         []byte
	Fields      map[string]Field
	RecordOrder []string
	Records     map[string]map[string]interface{}
}

// Target compiles target into its intermediate. A target with no
// dependencies compiles to an empty result matching its declared Mode
// (field-mode targets get an empty JSON object, list-mode targets get an
// empty array; raw or undeclared get empty raw bytes). A declared
// list-mode target always compiles via compileList. A target with
// dependencies and no declared list mode has its mode determined by
// whether its first dependency names a field, matching the raw/field
// exclusivity the manifest already enforces.
func Target(target manifest.Target) (Result, error) {
	if target.Mode == "list" {
		return compileList(target)
	}
	if len(target.Deps) == 0 {
		if target.Mode == "field" {
			return Result{IsField: true, Fields: map[string]Field{}}, nil
		}
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

// compileList resolves target's list-mode dependencies into an ordered
// set of records: a record-only dependency (bare record, no field) merges
// every top-level key of its parsed source document into that record; a
// record.field dependency resolves its value exactly as a field-mode
// dependency would (whole-file content, or a field-path with an optional
// format template) into that one field, unwrapped (no {value, raw}).
// Records appear in RecordOrder in the order each record's first
// dependency was encountered. A field name duplicated within the same
// record — from any combination of merged and explicit fields — is
// rejected.
func compileList(target manifest.Target) (Result, error) {
	order := []string{}
	records := make(map[string]map[string]interface{})
	seen := make(map[string]map[string]bool)

	ensureRecord := func(record string) {
		if _, ok := records[record]; ok {
			return
		}
		records[record] = make(map[string]interface{})
		seen[record] = make(map[string]bool)
		order = append(order, record)
	}

	setField := func(record, field string, value interface{}) error {
		if seen[record][field] {
			return fmt.Errorf("record %q has more than one dependency naming field %q", record, field)
		}
		seen[record][field] = true
		records[record][field] = value
		return nil
	}

	for _, dep := range target.Deps {
		record, field, ok := manifest.SplitRecordField(dep.Field)
		if !ok {
			return Result{}, fmt.Errorf("target %q: invalid list-mode field %q", target.Name, dep.Field)
		}
		ensureRecord(record)

		if field == "" {
			doc, err := parseSource(dep.Source)
			if err != nil {
				return Result{}, err
			}
			if err := requireFlatScalarDoc(doc); err != nil {
				return Result{}, fmt.Errorf("record %q: %w", record, err)
			}
			for key, value := range doc {
				if err := setField(record, key, value); err != nil {
					return Result{}, err
				}
			}
			continue
		}

		raw, err := resolveDepValue(dep)
		if err != nil {
			return Result{}, err
		}
		var value interface{} = raw
		if dep.Format != "" {
			value = fmt.Sprintf(dep.Format, raw)
		}
		if err := setField(record, field, value); err != nil {
			return Result{}, err
		}
	}

	return Result{IsList: true, RecordOrder: order, Records: records}, nil
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
