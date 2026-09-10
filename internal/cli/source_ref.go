package cli

import (
	"fmt"
	"strings"
)

// SourceRef is a parsed source-ref: `source` or `source:path`.
type SourceRef struct {
	Source string
	Path   string
}

func parseSourceRef(s string) (SourceRef, error) {
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		return SourceRef{Source: parts[0]}, nil
	case 2:
		return SourceRef{Source: parts[0], Path: parts[1]}, nil
	default:
		return SourceRef{}, fmt.Errorf("invalid source-ref %q: at most one ':' is allowed", s)
	}
}
