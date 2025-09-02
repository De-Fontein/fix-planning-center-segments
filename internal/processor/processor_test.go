package processor

import "testing"

func TestRemoveRepeats(t *testing.T) {
	t.Parallel()
	in := "Play this (x2)\nPlay this ( 3x )\nPlay this ×2\nPlay this x3 then stop\nPlay this 2x\n"
	want := "Play this\nPlay this\nPlay this\nPlay this then stop\nPlay this\n"

	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestRemoveToRefEnd(t *testing.T) {
	t.Parallel()
	in := "Line end (To Chorus)\nAndere regel (naar refrein)\nKeep (parenthetical info) in middle\nFinal (To Outro)\n"
	want := "Line end\nAndere regel\nKeep (parenthetical info) in middle\nFinal\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestWrapNakedChords_ChordOnlyLine(t *testing.T) {
	t.Parallel()
	in := "C G Am F\nC | G Am | F\n[C] G Am\n"
	want := "[C] [G] [Am] [F]\n[C] | [G] [Am] | [F]\n[C] [G] [Am]\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestNoWrapOnLyricLine(t *testing.T) {
	t.Parallel()
	in := "C G Amazing grace\nI am free\nA song of joy\n"
	// Do not wrap because lines include lyrics/words
	want := "C G Amazing grace\nI am free\nA song of joy\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestCollapseMultipleBlankLines(t *testing.T) {
	t.Parallel()
	in := "Line 1\n\n\n\nLine 2\n\n\nLine 3\n"
	want := "Line 1\n\nLine 2\n\nLine 3\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestTrailingNewlineNormalization(t *testing.T) {
	t.Parallel()
	in := "Line 1\nLine 2"
	want := "Line 1\nLine 2\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected trailing newline:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestTrimTrailingBlankLines(t *testing.T) {
	t.Parallel()
	in := "Line 1\n\n\n"
	want := "Line 1\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected trimming of trailing blanks:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestDoNotStripToRefInMiddle(t *testing.T) {
	t.Parallel()
	in := "Go (To Chorus) now please\n"
	want := "Go (To Chorus) now please\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected removal of middle parenthetical:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestCaseInsensitiveToRefRemoval(t *testing.T) {
	t.Parallel()
	in := "Lead in (TO BRIDGE)\nOutro (NAAR SLOT)\n"
	want := "Lead in\nOutro\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected case-insensitive To/Naar removal:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestRepeatTokensAtStartAndUnicode(t *testing.T) {
	t.Parallel()
	in := "x3 Play this\n×2 Start now\n2x begin\n"
	want := "Play this\nStart now\nbegin\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected repeat-token removal at start/unicode:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestSpacesCollapseAfterRepeatRemoval(t *testing.T) {
	t.Parallel()
	in := "Play   this   x3  now\n"
	want := "Play this now\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected space collapse:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestAdvancedChordWrapping(t *testing.T) {
	t.Parallel()
	in := "F#m D/E G#maj7 Cb\n"
	want := "[F#m] [D/E] [G#maj7] [Cb]\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected advanced chord wrapping:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestNoWrapOnNonChordLetters(t *testing.T) {
	t.Parallel()
	in := "H I J\n"
	want := "H I J\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected wrap on non-chord letters:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func TestExistingBracketsPreserved(t *testing.T) {
	t.Parallel()
	in := "[C] [G]\n"
	want := "[C] [G]\n"
	got := CleanText(in)
	if got != want {
		t.Fatalf("unexpected change to already-bracketed chords:\n--- got ---\n%q\n--- want ---\n%q", got, want)
	}
}

func Test_wrapChordsIfChordLine_SimpleChords(t *testing.T) {
	t.Parallel()
	in := "C G Am F"
	want := "[C] [G] [Am] [F]"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_WithBars(t *testing.T) {
	t.Parallel()
	in := "C | G Am | F"
	want := "[C] | [G] [Am] | [F]"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_AdvancedChords(t *testing.T) {
	t.Parallel()
	in := "F#m D/E G#maj7 Cb"
	want := "[F#m] [D/E] [G#maj7] [Cb]"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_AlreadyBracketed(t *testing.T) {
	t.Parallel()
	in := "[C] [G] | [Am] [F]"
	want := "[C] [G] | [Am] [F]"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_MixedBracketedAndNaked(t *testing.T) {
	t.Parallel()
	in := "[C] G"
	want := "[C] [G]"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_LyricsPresentShouldNotWrap(t *testing.T) {
	t.Parallel()
	in := "C G Amazing grace"
	want := "C G Amazing grace"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_NonChordLettersShouldNotWrap(t *testing.T) {
	t.Parallel()
	in := "H I J"
	want := "H I J"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_EmptyString(t *testing.T) {
	t.Parallel()
	in := ""
	want := ""
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_SpacesOnlyRemainUnchanged(t *testing.T) {
	t.Parallel()
	in := "   "
	want := "   "
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_LeadingTrailingSpacesNormalizedWhenOnlyChords(t *testing.T) {
	t.Parallel()
	in := "   C   G   "
	want := "[C] [G]"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}

func Test_wrapChordsIfChordLine_OnlyBarToken(t *testing.T) {
	t.Parallel()
	in := "|"
	want := "|"
	got := wrapChordsIfChordLine(in)
	if got != want {
		t.Fatalf("unexpected:\n--- in ---\n%q\n--- got ---\n%q\n--- want ---\n%q", in, got, want)
	}
}
