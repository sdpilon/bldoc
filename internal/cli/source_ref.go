package cli

import (
	"fmt"
	"strings"
)

// SourceRef is a parsed source-ref: `source`, `source:path`, or
// `source#anchor` (anchor either a single GitHub-style slug or a
// `parent/child` breadcrumb path of slugs). Path and Anchor are never
// both set.
type SourceRef struct {
	Source string
	Path   string
	Anchor string
}

func parseSourceRef(s string) (SourceRef, error) {
	hashParts := strings.Split(s, "#")
	switch len(hashParts) {
	case 1:
		parts := strings.Split(s, ":")
		switch len(parts) {
		case 1:
			return SourceRef{Source: parts[0]}, nil
		case 2:
			return SourceRef{Source: parts[0], Path: parts[1]}, nil
		default:
			return SourceRef{}, fmt.Errorf("invalid source-ref %q: at most one ':' is allowed", s)
		}
	case 2:
		source, anchor := hashParts[0], hashParts[1]
		if source == "" {
			return SourceRef{}, fmt.Errorf("invalid source-ref %q: empty source before '#'", s)
		}
		if anchor == "" {
			return SourceRef{}, fmt.Errorf("invalid source-ref %q: empty anchor after '#'", s)
		}
		if strings.Contains(source, ":") {
			return SourceRef{}, fmt.Errorf("invalid source-ref %q: cannot combine ':field-path' and '#anchor'", s)
		}
		return SourceRef{Source: source, Anchor: anchor}, nil
	default:
		return SourceRef{}, fmt.Errorf("invalid source-ref %q: at most one '#' is allowed", s)
	}
}
