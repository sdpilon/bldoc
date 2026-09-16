package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// heading is one parsed ATX heading (level 1-6, from its leading '#'
// count) in a Markdown document: its GitHub-style slug, its ancestor
// chain's slugs (root-first, not including itself — its breadcrumb),
// and the byte offsets bounding its own line in the source (used to
// find a matched heading's content boundaries).
type heading struct {
	level      int
	slug       string
	breadcrumb []string
	lineStart  int
	lineEnd    int // offset of the first byte after the heading line (and its newline)
}

// parseHeadings scans content for ATX headings (# through ######),
// skipping any inside fenced code blocks (``` or ~~~ delimited), and
// returns them in document order with each one's slug and ancestor
// breadcrumb computed.
func parseHeadings(content string) []heading {
	var headings []heading
	var stack []heading // currently open ancestor headings, outermost first

	inFence := false
	offset := 0
	for _, line := range strings.SplitAfter(content, "\n") {
		if line == "" {
			continue
		}
		lineStart := offset
		offset += len(line)
		trimmed := strings.TrimSpace(line)

		if isFenceDelim(trimmed) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}

		level, title, ok := parseATXHeading(strings.TrimRight(line, "\n"))
		if !ok {
			continue
		}

		for len(stack) > 0 && stack[len(stack)-1].level >= level {
			stack = stack[:len(stack)-1]
		}
		breadcrumb := make([]string, len(stack))
		for i, a := range stack {
			breadcrumb[i] = a.slug
		}

		h := heading{
			level:      level,
			slug:       slugify(title),
			breadcrumb: breadcrumb,
			lineStart:  lineStart,
			lineEnd:    offset,
		}
		headings = append(headings, h)
		stack = append(stack, h)
	}
	return headings
}

// isFenceDelim reports whether trimmed (a line with surrounding
// whitespace already stripped) opens or closes a fenced code block.
func isFenceDelim(trimmed string) bool {
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}

// parseATXHeading reports whether line (with its trailing newline, if
// any, already stripped) is an ATX heading, returning its level (1-6)
// and title text (leading/trailing whitespace and any closing '#'
// sequence trimmed) if so.
func parseATXHeading(line string) (level int, title string, ok bool) {
	if !strings.HasPrefix(line, "#") {
		return 0, "", false
	}
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level > 6 {
		return 0, "", false
	}
	rest := line[level:]
	if rest != "" && rest[0] != ' ' && rest[0] != '\t' {
		return 0, "", false
	}
	title = strings.TrimSpace(rest)
	title = strings.TrimRight(title, "#")
	title = strings.TrimSpace(title)
	return level, title, true
}

// slugify computes title's GitHub-style anchor slug: lowercased, with
// letters, digits, spaces, hyphens, and underscores kept (spaces
// becoming hyphens) and every other character dropped.
func slugify(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case r == ' ':
			b.WriteRune('-')
		case r == '-' || r == '_' || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
		default:
			// dropped: punctuation, backticks, etc.
		}
	}
	return b.String()
}

// matchesAnchor reports whether h is addressed by segments — a plain
// anchor's single segment must equal h's own slug; a breadcrumb
// anchor's segments (root-first) must equal h's own slug as the last
// segment, and the nearest N ancestors' slugs (a suffix of h's full
// root-first breadcrumb) as the preceding segments, for as many
// ancestor segments as were given.
func matchesAnchor(h heading, segments []string) bool {
	own := segments[len(segments)-1]
	if h.slug != own {
		return false
	}
	ancestorsWanted := segments[:len(segments)-1]
	if len(ancestorsWanted) > len(h.breadcrumb) {
		return false
	}
	hSuffix := h.breadcrumb[len(h.breadcrumb)-len(ancestorsWanted):]
	for i, want := range ancestorsWanted {
		if hSuffix[i] != want {
			return false
		}
	}
	return true
}

// findHeadingIndices returns the indices into headings of every heading
// matched by anchor (split on '/' into breadcrumb segments; a
// single-segment anchor is a plain slug lookup).
func findHeadingIndices(headings []heading, anchor string) []int {
	segments := strings.Split(anchor, "/")
	var matches []int
	for i, h := range headings {
		if matchesAnchor(h, segments) {
			matches = append(matches, i)
		}
	}
	return matches
}

// resolveSection returns content's raw text belonging to the heading at
// headings[matchIdx] — from just after its own line to the start of the
// next heading at a level <= its own (nested) or the very next heading
// regardless of level (default), or the end of content if there is no
// such next heading.
func resolveSection(content string, headings []heading, matchIdx int, nested bool) string {
	match := headings[matchIdx]
	end := len(content)
	for i := matchIdx + 1; i < len(headings); i++ {
		if !nested || headings[i].level <= match.level {
			end = headings[i].lineStart
			break
		}
	}
	return content[match.lineEnd:end]
}

// isMarkdownSource reports whether path's extension is .md or
// .markdown, case-insensitive.
func isMarkdownSource(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return true
	default:
		return false
	}
}

// ResolveMarkdownAnchor reads path (which must have a .md or .markdown
// extension) and resolves anchor's Markdown section content — the
// matched heading's own text (nested false) or its text plus all nested
// subsections (nested true). It rejects an unsupported source
// extension, an anchor matching no heading, and an anchor matching more
// than one heading (regardless of whether it was a plain or breadcrumb
// anchor).
func ResolveMarkdownAnchor(path, anchor string, nested bool) (string, error) {
	if !isMarkdownSource(path) {
		return "", fmt.Errorf("anchor %q: %s is not a Markdown source (expected .md or .markdown)", anchor, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", path, err)
	}
	content := string(data)
	headings := parseHeadings(content)
	matches := findHeadingIndices(headings, anchor)
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("anchor %q: no heading matches in %s", anchor, path)
	case 1:
		return resolveSection(content, headings, matches[0], nested), nil
	default:
		return "", fmt.Errorf("anchor %q: matches %d headings in %s; disambiguate with a breadcrumb path", anchor, len(matches), path)
	}
}

// AmbiguousAnchorSuggestions reads path and, if the given plain
// (single-segment) anchor matches more than one heading, returns one
// suggested breadcrumb per match (that match's immediate parent slug,
// if it has one, joined with its own slug). It returns a nil slice with
// no error when the anchor is unambiguous (zero or one match) — this is
// a best-effort check for add-dep's ambiguity warning, so callers
// should treat a read error as "skip the check", not a hard failure;
// add-dep never requires a source file to exist.
func AmbiguousAnchorSuggestions(path, anchor string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	headings := parseHeadings(string(data))
	matches := findHeadingIndices(headings, anchor)
	if len(matches) < 2 {
		return nil, nil
	}
	suggestions := make([]string, len(matches))
	for i, idx := range matches {
		h := headings[idx]
		if len(h.breadcrumb) == 0 {
			suggestions[i] = h.slug
			continue
		}
		suggestions[i] = h.breadcrumb[len(h.breadcrumb)-1] + "/" + h.slug
	}
	return suggestions, nil
}
