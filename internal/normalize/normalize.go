package normalize

import (
	"strconv"
	"strings"
)

// Newlines normalizes CRLF/CR to LF.
func Newlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// StripDecorations removes simple wrappers like [ ... ] and <b>...</b>.
func StripDecorations(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if s[0] == '[' && s[len(s)-1] == ']' {
		s = strings.TrimSpace(s[1 : len(s)-1])
	}
	l := strings.ToLower(s)
	if strings.HasPrefix(l, "<b>") && strings.HasSuffix(l, "</b>") {
		s = strings.TrimSpace(s[3 : len(s)-4])
	}
	return s
}

// Key normalizes a string for matching: remove spaces/hyphens/underscores and lowercase.
func Key(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, " ", "")
	return strings.ToLower(s)
}

// IsDigits reports whether s consists entirely of digits and represents a positive integer.
func IsDigits(s string) bool {
	if s == "" {
		return false
	}
	n, err := strconv.Atoi(s)

	return err == nil && n > 0
}
