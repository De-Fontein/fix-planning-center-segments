package parser

import (
	"reflect"
	"testing"
)

func TestParse_HeaderInBrackets(t *testing.T) {
	txt := "[Chorus]\nLa la"
	got := Parse(txt)
	if len(got) != 1 {
		t.Fatalf("expected 1 section, got %d", len(got))
	}
	if got[0].Header != "CHORUS" {
		t.Fatalf("expected header 'CHORUS', got %q", got[0].Header)
	}
	want := []string{"La la"}
	if !reflect.DeepEqual(got[0].Content, want) {
		t.Fatalf("content mismatch: want %#v, got %#v", want, got[0].Content)
	}
}

func TestParse_HeaderInBoldTags_WithNumber(t *testing.T) {
	txt := "<b>Verse 2</b>:\nLine"
	got := Parse(txt)
	if len(got) != 1 {
		t.Fatalf("expected 1 section, got %d", len(got))
	}
	if got[0].Header != "VERSE 2" {
		t.Fatalf("expected header 'VERSE 2', got %q", got[0].Header)
	}
}

func TestNormalizeNewlines_LFOnly(t *testing.T) {
	t.Parallel()

	in := "Line 1\nLine 2\nLine 3"
	want := "Line 1\nLine 2\nLine 3"
	if got := normalizeNewlines(in); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeNewlines_CRLF(t *testing.T) {
	t.Parallel()

	in := "Line 1\r\nLine 2\r\nLine 3"
	want := "Line 1\nLine 2\nLine 3"
	if got := normalizeNewlines(in); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeNewlines_CROnly(t *testing.T) {
	t.Parallel()

	in := "Line 1\rLine 2\rLine 3"
	want := "Line 1\nLine 2\nLine 3"
	if got := normalizeNewlines(in); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeNewlines_Mixed(t *testing.T) {
	t.Parallel()

	in := "Line 1\r\nLine 2\rLine 3\nLine 4\r\n"
	want := "Line 1\nLine 2\nLine 3\nLine 4\n"
	if got := normalizeNewlines(in); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParse_NoHeaders_AllGeneral(t *testing.T) {
	t.Parallel()

	txt := "Line 1\nLine 2\n\nLine 3"
	got := Parse(txt)
	if len(got) != 1 {
		t.Fatalf("expected 1 section, got %d", len(got))
	}
	if got[0].Header != "GENERAL" {
		t.Fatalf("expected header 'GENERAL', got %q", got[0].Header)
	}
	wantContent := []string{"Line 1", "Line 2", "", "Line 3"}
	if !reflect.DeepEqual(got[0].Content, wantContent) {
		t.Fatalf("content mismatch:\nwant: %#v\n got: %#v", wantContent, got[0].Content)
	}
}

func TestParse_SimpleHeaders_EN(t *testing.T) {
	t.Parallel()

	txt := "Verse\nA\nB\nChorus\nC\nD"
	got := Parse(txt)
	if len(got) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(got))
	}
	if got[0].Header != "VERSE" {
		t.Fatalf("expected first header 'VERSE', got %q", got[0].Header)
	}
	if got[1].Header != "CHORUS" {
		t.Fatalf("expected second header 'CHORUS', got %q", got[1].Header)
	}
}

func TestParse_HeaderVariants(t *testing.T) {
	t.Parallel()

	txt := "pre chorus\nA\nPreChorus\nB\nPRE-CHORUS\nC"
	got := Parse(txt)
	if len(got) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(got))
	}
	// Duplicates should be numbered 1..3
	if got[0].Header != "PRE-CHORUS 1" || got[1].Header != "PRE-CHORUS 2" || got[2].Header != "PRE-CHORUS 3" {
		t.Fatalf("unexpected headers: %#v", []string{got[0].Header, got[1].Header, got[2].Header})
	}
}

func TestParse_NumberedHeaders_And_AutoNumbering(t *testing.T) {
	t.Parallel()

	txt := "Verse 1:\nA\nVerse\nB\nVerse 1\nC\nVerse 3\nD"
	got := Parse(txt)
	if len(got) != 4 {
		t.Fatalf("expected 4 sections, got %d", len(got))
	}
	// Expect unique sequential-ish numbering without collisions.
	// First: VERSE 1 (explicit)
	// Second: VERSE 2 (auto, since another VERSE exists)
	// Third: VERSE 3 (explicit 1 collides, bump to next available)
	// Fourth: VERSE 4 (explicit 3 accepted -> next becomes 3, then next available for following would be 4)
	want := []string{"VERSE 1", "VERSE 2", "VERSE 3", "VERSE 4"}
	var gotHeaders []string
	for _, s := range got {
		gotHeaders = append(gotHeaders, s.Header)
	}
	if !reflect.DeepEqual(gotHeaders, want) {
		t.Fatalf("headers mismatch:\nwant: %#v\n got: %#v", want, gotHeaders)
	}
}

func TestParse_DutchKeywords(t *testing.T) {
	t.Parallel()

	txt := "Refrein 1:\nLine\nCouplet 1\nLine\nBrug 1\nLine\nUitro 1\nLine"
	got := Parse(txt)
	want := []string{"CHORUS 1", "VERSE 1", "BRIDGE 1", "OUTRO 1"}
	var gotHeaders []string
	for _, s := range got {
		gotHeaders = append(gotHeaders, s.Header)
	}
	if !reflect.DeepEqual(gotHeaders, want) {
		t.Fatalf("headers mismatch:\nwant: %#v\n got: %#v", want, gotHeaders)
	}
}

func TestParse_SingletonDoesNotForceNumber(t *testing.T) {
	t.Parallel()

	txt := "Bridge\nA\nB"
	got := Parse(txt)
	if len(got) != 1 {
		t.Fatalf("expected 1 section, got %d", len(got))
	}
	// Only one BRIDGE -> keep as "BRIDGE" (no numbering)
	if got[0].Header != "BRIDGE" {
		t.Fatalf("expected header 'BRIDGE', got %q", got[0].Header)
	}
}

func TestParse_DuplicateExplicitNumbersAvoidCollision(t *testing.T) {
	t.Parallel()

	txt := "Chorus 1\nA\nChorus 1\nB\nChorus\nC"
	got := Parse(txt)
	want := []string{"CHORUS 1", "CHORUS 2", "CHORUS 3"}
	var gotHeaders []string
	for _, s := range got {
		gotHeaders = append(gotHeaders, s.Header)
	}
	if !reflect.DeepEqual(gotHeaders, want) {
		t.Fatalf("headers mismatch:\nwant: %#v\n got: %#v", want, gotHeaders)
	}
}
