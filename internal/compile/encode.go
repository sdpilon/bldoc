package compile

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Encode serializes result's structured data (field- or list-mode) into
// its final output bytes, choosing the encoder from ext: YAML when ext
// is "yaml" or "yml" (case-insensitive), JSON otherwise. A raw-mode
// result's bytes are returned unchanged — ext never affects raw mode's
// encoding, only its filename (handled by the caller).
func Encode(result Result, ext string) ([]byte, error) {
	if !result.IsField && !result.IsList {
		return result.Raw, nil
	}

	useYAML := strings.EqualFold(ext, "yaml") || strings.EqualFold(ext, "yml")

	var value interface{}
	if result.IsList {
		list := make([]map[string]interface{}, len(result.RecordOrder))
		for i, record := range result.RecordOrder {
			list[i] = result.Records[record]
		}
		value = list
	} else {
		value = result.Fields
	}

	if useYAML {
		data, err := yaml.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("encoding intermediate as YAML: %w", err)
		}
		return data, nil
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encoding intermediate as JSON: %w", err)
	}
	return data, nil
}
