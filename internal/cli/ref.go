package cli

import (
	"fmt"
	"strings"
)

// TargetRef is a parsed target-ref: `target` or `target:field`.
type TargetRef struct {
	Target string
	Field  string
}

func parseTargetRef(s string) (TargetRef, error) {
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		return TargetRef{Target: parts[0]}, nil
	case 2:
		return TargetRef{Target: parts[0], Field: parts[1]}, nil
	default:
		return TargetRef{}, fmt.Errorf("invalid target-ref %q: at most one ':' is allowed", s)
	}
}
