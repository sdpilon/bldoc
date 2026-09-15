package compile

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEncode_RawPassesThrough(t *testing.T) {
	result := Result{Raw: []byte("hello world")}

	data, err := Encode(result, "yaml")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("expected raw bytes unchanged, got %q", data)
	}
}

func TestEncode_FieldModeJSON(t *testing.T) {
	result := Result{IsField: true, Fields: map[string]Field{"summary": {Value: "notes", Raw: "notes"}}}

	data, err := Encode(result, "")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if !strings.Contains(string(data), `"summary"`) || !strings.Contains(string(data), `"value": "notes"`) {
		t.Fatalf("expected JSON with summary/value, got %q", data)
	}
}

func TestEncode_FieldModeYAML(t *testing.T) {
	result := Result{IsField: true, Fields: map[string]Field{"summary": {Value: "notes", Raw: "notes"}}}

	data, err := Encode(result, "yaml")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	var out map[string]map[string]string
	if err := yaml.Unmarshal(data, &out); err != nil {
		t.Fatalf("compiled output is not valid YAML: %v\noutput:\n%s", err, data)
	}
	if out["summary"]["value"] != "notes" {
		t.Fatalf("expected summary.value == notes, got %+v", out)
	}
}

func TestEncode_ListModeJSON(t *testing.T) {
	result := Result{
		IsList:      true,
		RecordOrder: []string{"bldoc"},
		Records:     map[string]map[string]interface{}{"bldoc": {"path": "~/bldoc"}},
	}

	data, err := Encode(result, "")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(string(data)), "[") {
		t.Fatalf("expected a JSON array, got %q", data)
	}
}

func TestEncode_ListModeYAML(t *testing.T) {
	result := Result{
		IsList:      true,
		RecordOrder: []string{"bldoc", "projx"},
		Records: map[string]map[string]interface{}{
			"bldoc": {"path": "~/bldoc"},
			"projx": {"path": "~/projx"},
		},
	}

	data, err := Encode(result, "yaml")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	var out []map[string]interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		t.Fatalf("compiled output is not valid YAML: %v\noutput:\n%s", err, data)
	}
	if len(out) != 2 || out[0]["path"] != "~/bldoc" || out[1]["path"] != "~/projx" {
		t.Fatalf("expected ordered records, got %+v", out)
	}
}

func TestEncode_ListModeYAMLExtCaseInsensitive(t *testing.T) {
	result := Result{IsList: true, RecordOrder: []string{}, Records: map[string]map[string]interface{}{}}

	data, err := Encode(result, "YML")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	var out []map[string]interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		t.Fatalf("compiled output is not valid YAML: %v\noutput:\n%s", err, data)
	}
}

func TestEncode_EmptyListModeIsEmptyArray(t *testing.T) {
	result := Result{IsList: true, RecordOrder: []string{}, Records: map[string]map[string]interface{}{}}

	jsonData, err := Encode(result, "")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if strings.TrimSpace(string(jsonData)) != "[]" {
		t.Fatalf("expected [], got %q", jsonData)
	}
}
