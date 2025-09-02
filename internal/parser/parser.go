package parser

import (
	"chordparser/internal/normalize"
	"strconv"
	"strings"
)

// Canonical, case-insensitive keywords (EN + NL) with variant normalization.
// Variants normalize by removing spaces/hyphens/underscores and lowercasing.
var keywordCanonical = map[string]string{
	// English
	"verse":        "VERSE",
	"chorus":       "CHORUS",
	"refrain":      "REFRAIN",
	"prechorus":    "PRE-CHORUS",
	"pre-chorus":   "PRE-CHORUS",
	"postchorus":   "POST-CHORUS",
	"post-chorus":  "POST-CHORUS",
	"bridge":       "BRIDGE",
	"intro":        "INTRO",
	"outro":        "OUTRO",
	"ending":       "ENDING",
	"instrumental": "INSTRUMENTAL",
	"interlude":    "INTERLUDE",
	"tag":          "TAG",
	"turnaround":   "TURNAROUND",
	"vamp":         "VAMP",
	"breakdown":    "BREAKDOWN",

	// Dutch (mapped to English)
	"refrein":       "CHORUS",
	"couplet":       "VERSE",
	"brug":          "BRIDGE",
	"uitro":         "OUTRO",
	"intermezzo":    "INTERLUDE",
	"instrumentaal": "INSTRUMENTAL",
	"slot":          "ENDING",
}

// Section represents a parsed section.
type Section struct {
	Header  string   `json:"header"`
	Content []string `json:"content"`
}

func Parse(text string) []Section {
	text = normalize.Newlines(text)
	lines := strings.Split(text, "\n")

	var sections []Section
	header := "GENERAL"
	content := []string{}
	foundAnyHeader := false

	for _, line := range lines {
		if base, num, ok := detectHeader(line); ok {
			// flush previous content if any
			if len(content) > 0 {
				sections = append(sections, Section{Header: header, Content: content})
			}
			foundAnyHeader = true
			if num > 0 {
				header = base + " " + strconv.Itoa(num)
			} else {
				header = base
			}
			content = nil
			continue
		}
		content = append(content, line)
	}

	// flush last accumulated content
	if len(content) > 0 {
		sections = append(sections, Section{Header: header, Content: content})
	}

	// If at least one header was found, ensure uniqueness across duplicates.
	// If none were found, keep the single "General" section unnumbered.
	if foundAnyHeader {
		sections = makeUniqueHeaders(sections)
	}

	return sections
}

// Parse splits a raw chord/lyrics text into ordered sections.
// - Detects headers using EN/NL keywords (case-insensitive, supports variants).
// - Understands numbered headers (e.g., "Verse 1").
// - If no headers exist, returns one "General" section.
// - Ensures duplicate headers are made unique by numbering at the end.
func normalizeNewlines(s string) string {
	// Normalize Windows CRLF and old Mac CR to LF
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// detectHeader attempts to parse the given line as a section header.
// It returns the canonical base header, an explicit number if present (0 if not),
// and whether the line is a recognized header.
func detectHeader(line string) (string, int, bool) {
	// Strip simple bold tags to support headers like <b>Verse 2</b>:
	line = strings.ReplaceAll(line, "<b>", "")
	line = strings.ReplaceAll(line, "</b>", "")
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "", 0, false
	}

	// Drop trailing colons (common in headers like "Verse 1:")
	for strings.HasSuffix(trimmed, ":") {
		trimmed = strings.TrimSuffix(trimmed, ":")
		trimmed = strings.TrimSpace(trimmed)
	}

	if trimmed == "" {
		return "", 0, false
	}

	fields := strings.Fields(trimmed)

	// Detect trailing number (e.g., "Chorus 2")
	var num int
	var err error
	last := fields[len(fields)-1]

	hasNumber := false
	if normalize.IsDigits(last) {
		hasNumber = true
		num, err = strconv.Atoi(last)
		if err != nil {
			// Shouldn't happen due to isDigits, but guard anyway.
			hasNumber = false
			num = 0
		}
	}

	// Base name to match against keywords (may contain decorations)
	name := trimmed
	if hasNumber {
		name = strings.TrimSpace(strings.TrimSuffix(trimmed, last))
	}

	// Try direct match on full name (original and decoration-stripped)
	if canon, ok := keywordCanonical[normalize.Key(name)]; ok {
		return canon, num, true
	}

	// Finally, inspect the first "word" only, allowing decorations and "includes" matching.
	first := fields[0]
	first = normalize.StripDecorations(first)
	firstNorm := normalize.Key(first)

	if canon, ok := keywordCanonical[firstNorm]; ok {
		return canon, num, true
	}
	// If the first token contains a keyword (e.g., "[Pre-Chorus]"), accept it.
	for k, canon := range keywordCanonical {
		if strings.Contains(firstNorm, k) {
			return canon, num, true
		}
	}

	return "", 0, false
}

// makeUniqueHeaders ensures headers are unique by adding/incrementing numbers
// only for bases that appear multiple times. Single occurrences are left as-is.
func makeUniqueHeaders(sections []Section) []Section {
	// Count occurrences per base (ignore any existing numbers).
	baseCounts := make(map[string]int)
	for _, s := range sections {
		base, _ := splitBaseAndNumber(s.Header)
		baseCounts[base]++
	}

	// For bases with duplicates, assign sequence numbers in order.
	next := make(map[string]int)
	for i := range sections {
		base, n := splitBaseAndNumber(sections[i].Header)
		count := baseCounts[base]

		if count <= 1 {
			// Leave singletons as-is (preserve explicit numbering if the author provided it).
			continue
		}

		if n == 0 {
			// No explicit number: assign next
			next[base]++
			n = next[base]
		} else if n <= next[base] {
			// Explicit number collides with or precedes the sequence; bump to next
			next[base]++
			n = next[base]
		} else {
			// Explicit number is ahead; accept and move the cursor
			next[base] = n
		}

		sections[i].Header = base + " " + strconv.Itoa(n)
	}

	return sections
}

// splitBaseAndNumber splits "Header 2" -> ("Header", 2). If no trailing number, returns (Header, 0).
func splitBaseAndNumber(header string) (string, int) {
	h := strings.TrimSpace(header)
	if h == "" {
		return "", 0
	}
	parts := strings.Fields(h)
	if len(parts) < 2 {
		return h, 0
	}
	last := parts[len(parts)-1]
	if normalize.IsDigits(last) {
		n, err := strconv.Atoi(last)
		if err == nil {
			return strings.TrimSpace(strings.TrimSuffix(h, last)), n
		}
	}
	return h, 0
}
