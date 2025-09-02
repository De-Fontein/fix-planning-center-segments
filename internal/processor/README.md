This package is used to clean up ChordPro files which have been parsed into segments.
It has a couple of rules:
1. Remove all `(x2)`, `(3x)`, `(×4)`, `x3`, `2x`, etc. references (2 being variable).
2. Remove references like `(To Chorus)`, `(To Outro)`, `(To End)`, `(Naar Refrein)`, `(Naar Slot)`, etc. which may appear at the end of a line.
3. Surround "naked chords" (chords without brackets) with square brackets `[]` for consistency on chord-only lines.
4. Remove consecutive empty lines, singular empty lines are fine.
