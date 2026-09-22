package cli

import (
	"fmt"
	"strings"
)

// pair is one entry in the variadic `add-dep` form: a bare source-ref
// (Field empty — a whole-file, anchor-addressed, or list-mode
// record-only dependency), or `field:source-ref` (a field-mode
// dependency, or a single field within a list-mode record). Because a
// source-ref never itself contains ':', splitting on the first ':' is
// unambiguous regardless of the target's mode; whether a given pair's
// shape is valid for that mode is left to manifest.AddDep, exactly as it
// already is for a single add-dep call.
type pair struct {
	Field  string
	Source SourceRef
}

// validateBareTarget rejects a target argument carrying a '@field'
// suffix, which has no meaning in the variadic add-dep form — every
// pair carries its own field instead.
func validateBareTarget(s string) error {
	if strings.Contains(s, "@") {
		return fmt.Errorf("invalid target %q: the target argument must be bare (no '@field') in this form", s)
	}
	return nil
}

// validateRecordName rejects a record-name argument containing any of
// the symbols with meaning elsewhere in the grammar ('@' field
// addressing, ':' pair mapping, '.' record.field nesting).
func validateRecordName(s string) error {
	if strings.ContainsAny(s, "@:.") {
		return fmt.Errorf("invalid record name %q: must be bare (no '@', ':', or '.')", s)
	}
	return nil
}

func parsePair(s string) (pair, error) {
	idx := strings.Index(s, ":")
	if idx == -1 {
		sourceRef, err := parseSourceRef(s)
		if err != nil {
			return pair{}, err
		}
		return pair{Source: sourceRef}, nil
	}
	field, rest := s[:idx], s[idx+1:]
	if field == "" {
		return pair{}, fmt.Errorf("invalid pair %q: empty field before ':'", s)
	}
	sourceRef, err := parseSourceRef(rest)
	if err != nil {
		return pair{}, err
	}
	return pair{Field: field, Source: sourceRef}, nil
}
