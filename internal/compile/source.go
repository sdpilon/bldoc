// Package compile turns a manifest target's recorded dependencies into
// its compiled intermediate: raw-mode concatenation or field-mode JSON.
package compile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// parseSource reads path and parses it as TOML, JSON, or YAML based on
// its file extension, into a generic key/value tree suitable for
// resolveFieldPath.
func parseSource(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var out map[string]interface{}
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".toml":
		if err := toml.Unmarshal(data, &out); err != nil {
			return nil, fmt.Errorf("parsing %s as TOML: %w", path, err)
		}
	case ".json":
		if err := json.Unmarshal(data, &out); err != nil {
			return nil, fmt.Errorf("parsing %s as JSON: %w", path, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &out); err != nil {
			return nil, fmt.Errorf("parsing %s as YAML: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported source format %q for %s (expected .toml, .json, .yaml, or .yml)", ext, path)
	}
	return out, nil
}

// resolveFieldPath walks tree following the dot-separated keys in path,
// requiring the final value to be a scalar (string, number, or bool).
func resolveFieldPath(tree map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	var cur interface{} = tree

	for i, key := range keys {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field-path %q: %q is not a table", path, strings.Join(keys[:i], "."))
		}
		v, ok := m[key]
		if !ok {
			return nil, fmt.Errorf("field-path %q: key %q not found", path, strings.Join(keys[:i+1], "."))
		}
		cur = v
	}

	if !isScalar(cur) {
		return nil, fmt.Errorf("field-path %q resolves to a table or array, not a scalar value", path)
	}
	return cur, nil
}

// isScalar reports whether v is a string, number, or boolean — i.e. not
// a nested table/object or array/list.
func isScalar(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return false
	default:
		return true
	}
}

// requireFlatScalarDoc rejects a parsed document with any top-level
// value that is not a scalar (a nested table/object or array/list),
// returning the first offending key found.
func requireFlatScalarDoc(doc map[string]interface{}) error {
	for key, v := range doc {
		if !isScalar(v) {
			return fmt.Errorf("key %q resolves to a table or array, not a scalar value", key)
		}
	}
	return nil
}
