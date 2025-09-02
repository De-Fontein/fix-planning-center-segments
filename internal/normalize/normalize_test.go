package normalize

import "testing"

func TestNewlines_CRLFToLF(t *testing.T) {
	t.Parallel()
	in, want := "a\r\nb\r\nc", "a\nb\nc"
	if got := Newlines(in); got != want {
		t.Fatalf("Newlines(%q) = %q, want %q", in, got, want)
	}
}

func TestNewlines_CRToLF(t *testing.T) {
	t.Parallel()
	in, want := "a\rb\rc", "a\nb\nc"
	if got := Newlines(in); got != want {
		t.Fatalf("Newlines(%q) = %q, want %q", in, got, want)
	}
}

func TestNewlines_MixedCRAndCRLF(t *testing.T) {
	t.Parallel()
	in, want := "a\rb\r\nc\n", "a\nb\nc\n"
	if got := Newlines(in); got != want {
		t.Fatalf("Newlines(%q) = %q, want %q", in, got, want)
	}
}

func TestNewlines_AlreadyLF(t *testing.T) {
	t.Parallel()
	in, want := "a\nb\nc\n", "a\nb\nc\n"
	if got := Newlines(in); got != want {
		t.Fatalf("Newlines(%q) = %q, want %q", in, got, want)
	}
}

func TestNewlines_Empty(t *testing.T) {
	t.Parallel()
	in, want := "", ""
	if got := Newlines(in); got != want {
		t.Fatalf("Newlines(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_BracketsOnly(t *testing.T) {
	t.Parallel()
	in, want := "[Verse]", "Verse"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_BoldOnly(t *testing.T) {
	t.Parallel()
	in, want := "<b>Chorus</b>", "Chorus"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_BoldOnlyWithSpaces(t *testing.T) {
	t.Parallel()
	in, want := "  <b> Chorus </b>  ", "Chorus"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_UppercaseBoldTags(t *testing.T) {
	t.Parallel()
	in, want := "<B>Tag</B>", "Tag"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_BracketsOutsideBold(t *testing.T) {
	t.Parallel()
	in, want := "[<b>Bridge</b>]", "Bridge"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

// Note: Bold outside brackets does NOT strip inner brackets with current implementation.
func TestStripDecorations_BoldOutsideBrackets(t *testing.T) {
	t.Parallel()
	in, want := "<b>[Chorus]</b>", "[Chorus]"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_MismatchedBoldNoClose(t *testing.T) {
	t.Parallel()
	in, want := "<b>Tag", "<b>Tag"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_NonWrappingBracket(t *testing.T) {
	t.Parallel()
	in, want := "Intro]", "Intro]"
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_Empty(t *testing.T) {
	t.Parallel()
	in, want := "", ""
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestStripDecorations_SpacesOnly(t *testing.T) {
	t.Parallel()
	in, want := "   ", ""
	if got := StripDecorations(in); got != want {
		t.Fatalf("StripDecorations(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_PreChorus(t *testing.T) {
	t.Parallel()
	in, want := "Pre-Chorus", "prechorus"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_Spaces(t *testing.T) {
	t.Parallel()
	in, want := " Pre   Chorus ", "prechorus"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_Underscore(t *testing.T) {
	t.Parallel()
	in, want := "pre_chorus", "prechorus"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_DigitsPreserved(t *testing.T) {
	t.Parallel()
	in, want := "Chorus2", "chorus2"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_Mixed(t *testing.T) {
	t.Parallel()
	in, want := "  Pre-_ Chorus_2  ", "prechorus2"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_Empty(t *testing.T) {
	t.Parallel()
	in, want := "", ""
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_PunctuationPreserved(t *testing.T) {
	t.Parallel()
	in, want := "Verse (Alt) - 2", "verse(alt)2"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_OnlySeparators(t *testing.T) {
	t.Parallel()
	in, want := "---___   ", ""
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_MultipleMixedSeparators(t *testing.T) {
	t.Parallel()
	in, want := " -_A-_-B_-- ", "ab"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestKey_UnicodeLettersAndDiacritics(t *testing.T) {
	t.Parallel()
	in, want := "Pré-Chörus", "préchörus"
	if got := Key(in); got != want {
		t.Fatalf("Key(%q) = %q, want %q", in, got, want)
	}
}

func TestIsDigits_Positive123(t *testing.T) {
	t.Parallel()
	in, want := "123", true
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_Zero(t *testing.T) {
	t.Parallel()
	in, want := "0", false
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_LeadingZerosPositive(t *testing.T) {
	t.Parallel()
	in, want := "01", true
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_AllZeros(t *testing.T) {
	t.Parallel()
	in, want := "000", false
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_Empty(t *testing.T) {
	t.Parallel()
	in, want := "", false
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_Spaces(t *testing.T) {
	t.Parallel()
	in, want := " 123 ", false
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_Negative(t *testing.T) {
	t.Parallel()
	in, want := "-1", false
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}

func TestIsDigits_NonDigits(t *testing.T) {
	t.Parallel()
	in, want := "12a3", false
	if got := IsDigits(in); got != want {
		t.Fatalf("IsDigits(%q) = %v, want %v", in, got, want)
	}
}
