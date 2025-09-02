package processor

import (
	"regexp"
	"strings"
)

var (
	// Remove trailing parenthetical repeats like "(x2)" or "( 3x )" at end of line
	reParenRepeatEnd = regexp.MustCompile(`\s*\(\s*(?:\d+\s*[x×]|[x×]\s*\d+)\s*\)\s*$`)
	// Remove trailing parenthetical directives to sections, e.g., "(To Chorus)", "(naar refrein)" at end of line
	reToRefEnd = regexp.MustCompile(`(?i)\s*\(\s*(?:to|naar)\b[^)]*\)\s*$`)
	// Remove standalone repeat tokens like "x3", "3x", "×2" anywhere
	reRepeatToken = regexp.MustCompile(`(?:^|\s)(?:\d+\s*[x×]|[x×]\s*\d+)(?:\s|$)`)
	// Collapse multiple spaces
	reMultiSpaces = regexp.MustCompile(` {2,}`)
	// Chord token matcher (simple but covers common chord formats)
	reChord = regexp.MustCompile(`^(?:[A-G](?:#|b)?(?:(?:maj|min|m|dim|aug|sus)\d*|\d*)?(?:/[A-G](?:#|b)?)?)$`)
)

// CleanText normalizes a song text by removing repeat notations, trailing section directives,
// wrapping naked chord-only lines in brackets, and collapsing multiple blank lines.
func CleanText(in string) string {
	lines := strings.Split(in, "\n")

	out := make([]string, 0, len(lines))
	prevBlank := false

	for _, line := range lines {
		s := line

		// Remove trailing "(To ...)" or "(naar ...)" parentheticals and repeat markers
		s = reToRefEnd.ReplaceAllString(s, "")
		s = reParenRepeatEnd.ReplaceAllString(s, "")

		// Remove standalone repeat tokens like "x3", "3x", "×2"
		s = reRepeatToken.ReplaceAllStringFunc(s, func(m string) string {
			// keep a single leading space if there was one
			if strings.HasPrefix(m, " ") || strings.HasPrefix(m, "\t") {
				return " "
			}
			return ""
		})

		// Tidy spaces
		s = strings.TrimSpace(reMultiSpaces.ReplaceAllString(s, " "))

		// Wrap chords on chord-only lines
		s = wrapChordsIfChordLine(s)

		// Collapse multiple blank lines
		if s == "" {
			if !prevBlank && (len(out) == 0 || out[len(out)-1] != "") {
				out = append(out, "")
			}
			prevBlank = true
		} else {
			out = append(out, s)
			prevBlank = false
		}
	}

	// Remove trailing blank lines
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}

	// Ensure exactly one trailing newline
	return strings.Join(out, "\n") + "\n"
}

func wrapChordsIfChordLine(s string) string {
	if s == "" {
		return s
	}

	tokens := strings.Fields(s)
	if len(tokens) == 0 {
		return s
	}

	onlyChords := true
	for _, t := range tokens {
		if t == "|" {
			continue
		}
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			continue
		}
		if reChord.MatchString(t) {
			continue
		}
		onlyChords = false
		break
	}

	if !onlyChords {
		return s
	}

	for i, t := range tokens {
		if t == "|" {
			continue
		}
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			continue
		}
		if reChord.MatchString(t) {
			tokens[i] = "[" + t + "]"
		}
	}

	return strings.Join(tokens, " ")
}
