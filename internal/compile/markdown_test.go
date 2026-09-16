package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleDoc = `# Doc Title

## Purpose

Purpose text.

## Requirements

### Requirement: Manifest file location
Requirement body.

#### Scenario: Manifest path
- **WHEN** a thing happens
- **THEN** another thing happens

### Requirement: Second one
More body.

#### Scenario: Manifest path
- **WHEN** a different thing happens
- **THEN** yet another thing happens

## Impact

Impact text.
`

func TestParseHeadings_LevelsAndSlugs(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	want := []struct {
		level int
		slug  string
	}{
		{1, "doc-title"},
		{2, "purpose"},
		{2, "requirements"},
		{3, "requirement-manifest-file-location"},
		{4, "scenario-manifest-path"},
		{3, "requirement-second-one"},
		{4, "scenario-manifest-path"},
		{2, "impact"},
	}
	if len(headings) != len(want) {
		t.Fatalf("got %d headings, want %d: %+v", len(headings), len(want), headings)
	}
	for i, w := range want {
		if headings[i].level != w.level || headings[i].slug != w.slug {
			t.Fatalf("heading %d: got {level:%d slug:%q}, want {level:%d slug:%q}", i, headings[i].level, headings[i].slug, w.level, w.slug)
		}
	}
}

func TestParseHeadings_Breadcrumb(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	var scenario heading
	found := 0
	for _, h := range headings {
		if h.slug == "scenario-manifest-path" && h.breadcrumb[len(h.breadcrumb)-1] == "requirement-manifest-file-location" {
			scenario = h
			found++
		}
	}
	if found != 1 {
		t.Fatalf("expected exactly one scenario nested under requirement-manifest-file-location, found %d", found)
	}
	got := scenario.breadcrumb
	if len(got) != 3 || got[0] != "doc-title" || got[1] != "requirements" || got[2] != "requirement-manifest-file-location" {
		t.Fatalf("got breadcrumb %+v", got)
	}
}

func TestParseHeadings_SkipsFencedCode(t *testing.T) {
	doc := "# Title\n\n```\n# not a heading\n```\n\n## Real Heading\n"
	headings := parseHeadings(doc)
	if len(headings) != 2 {
		t.Fatalf("expected 2 headings (fenced '#' skipped), got %d: %+v", len(headings), headings)
	}
	if headings[1].slug != "real-heading" {
		t.Fatalf("got %+v", headings[1])
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Purpose":                             "purpose",
		"Requirement: Manifest file location": "requirement-manifest-file-location",
		"`--format` requires a field-addressed target": "--format-requires-a-field-addressed-target",
		"Scenario: Unknown target rejected":            "scenario-unknown-target-rejected",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Fatalf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFindHeadingIndices_PlainAnchorAmbiguous(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	matches := findHeadingIndices(headings, "scenario-manifest-path")
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches for the duplicated slug, got %d", len(matches))
	}
}

func TestFindHeadingIndices_BreadcrumbDisambiguates(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	matches := findHeadingIndices(headings, "requirement-second-one/scenario-manifest-path")
	if len(matches) != 1 {
		t.Fatalf("expected exactly 1 match via breadcrumb, got %d", len(matches))
	}
	got := headings[matches[0]]
	if got.breadcrumb[len(got.breadcrumb)-1] != "requirement-second-one" {
		t.Fatalf("breadcrumb match resolved to the wrong heading: %+v", got)
	}
}

func TestFindHeadingIndices_NotFound(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	matches := findHeadingIndices(headings, "does-not-exist")
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func TestResolveSection_DefaultStopsAtNextHeadingAnyLevel(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	idx := -1
	for i, h := range headings {
		if h.slug == "requirements" {
			idx = i
		}
	}
	if idx == -1 {
		t.Fatal("requirements heading not found")
	}
	got := resolveSection(sampleDoc, headings, idx, false)
	if strings.Contains(got, "Requirement: Manifest file location") {
		t.Fatalf("default capture should stop before the first nested heading, got %q", got)
	}
}

func TestResolveSection_NestedIncludesSubsectionsStopsAtSibling(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	idx := -1
	for i, h := range headings {
		if h.slug == "requirements" {
			idx = i
		}
	}
	got := resolveSection(sampleDoc, headings, idx, true)
	if !strings.Contains(got, "Requirement: Manifest file location") || !strings.Contains(got, "Requirement: Second one") {
		t.Fatalf("nested capture should include all subsections, got %q", got)
	}
	if strings.Contains(got, "Impact text") {
		t.Fatalf("nested capture should stop at the next ##-or-shallower heading, got %q", got)
	}
}

func TestResolveSection_ReachesEndOfFile(t *testing.T) {
	headings := parseHeadings(sampleDoc)
	idx := -1
	for i, h := range headings {
		if h.slug == "impact" {
			idx = i
		}
	}
	got := resolveSection(sampleDoc, headings, idx, false)
	if !strings.Contains(got, "Impact text.") {
		t.Fatalf("expected the last heading's content to reach end of file, got %q", got)
	}
}

func TestResolveMarkdownAnchor_UnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(path, []byte("# Purpose\nbody\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := ResolveMarkdownAnchor(path, "purpose", false); err == nil {
		t.Fatal("expected an error for a non-Markdown extension")
	}
}

func TestResolveMarkdownAnchor_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(path, []byte("# Title\n\n## Purpose\nbody\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := ResolveMarkdownAnchor(path, "missing", false); err == nil {
		t.Fatal("expected an error for an anchor with no matching heading")
	}
}

func TestResolveMarkdownAnchor_Ambiguous(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(path, []byte(sampleDoc), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := ResolveMarkdownAnchor(path, "scenario-manifest-path", false); err == nil {
		t.Fatal("expected an error for an ambiguous plain anchor")
	}
}

func TestResolveMarkdownAnchor_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(path, []byte(sampleDoc), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := ResolveMarkdownAnchor(path, "purpose", false)
	if err != nil {
		t.Fatalf("ResolveMarkdownAnchor: %v", err)
	}
	if !strings.Contains(got, "Purpose text.") {
		t.Fatalf("got %q", got)
	}
}

func TestAmbiguousAnchorSuggestions_UnreadableFileErrors(t *testing.T) {
	dir := t.TempDir()
	_, err := AmbiguousAnchorSuggestions(filepath.Join(dir, "missing.md"), "purpose")
	if err == nil {
		t.Fatal("expected an error for a missing file (callers treat this as skip-the-check)")
	}
}

func TestAmbiguousAnchorSuggestions_Unambiguous(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(path, []byte(sampleDoc), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := AmbiguousAnchorSuggestions(path, "purpose")
	if err != nil {
		t.Fatalf("AmbiguousAnchorSuggestions: %v", err)
	}
	if got != nil {
		t.Fatalf("expected no suggestions for an unambiguous anchor, got %+v", got)
	}
}

func TestAmbiguousAnchorSuggestions_SuggestsBreadcrumbs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "spec.md")
	if err := os.WriteFile(path, []byte(sampleDoc), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	got, err := AmbiguousAnchorSuggestions(path, "scenario-manifest-path")
	if err != nil {
		t.Fatalf("AmbiguousAnchorSuggestions: %v", err)
	}
	want := []string{
		"requirement-manifest-file-location/scenario-manifest-path",
		"requirement-second-one/scenario-manifest-path",
	}
	if len(got) != len(want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	}
}

// TestParseHeadings_RepoFixtureHasKnownDuplicates grounds the ambiguity
// mechanism in this repo's own real specs, which already contain
// duplicate headings (see the markdown-section-deps proposal).
func TestParseHeadings_RepoFixtureHasKnownDuplicates(t *testing.T) {
	cases := []struct {
		path string
		slug string
	}{
		{"../../openspec/specs/cli-shell/spec.md", "scenario-extra-argument-rejected"},
		{"../../openspec/specs/manifest-store/spec.md", "scenario-unknown-target-rejected"},
	}
	for _, c := range cases {
		data, err := os.ReadFile(c.path)
		if err != nil {
			t.Fatalf("ReadFile %s: %v", c.path, err)
		}
		headings := parseHeadings(string(data))
		matches := findHeadingIndices(headings, c.slug)
		if len(matches) < 2 {
			t.Fatalf("%s: expected at least 2 headings matching %q, got %d", c.path, c.slug, len(matches))
		}
	}
}
