package manifest

import "strings"

// SplitRecordField splits a list-mode dependency's target-side field
// string into its record and field components. Zero dots names a
// record-only dependency (field is ""); exactly one dot names a
// record.field dependency. An empty field, or a field with more than one
// dot, is invalid.
func SplitRecordField(field string) (record, subfield string, ok bool) {
	if field == "" {
		return "", "", false
	}
	parts := strings.Split(field, ".")
	switch len(parts) {
	case 1:
		return parts[0], "", true
	case 2:
		return parts[0], parts[1], true
	default:
		return "", "", false
	}
}
