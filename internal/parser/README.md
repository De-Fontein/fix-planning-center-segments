This package is used to parse raw ChordPro files into raw segments.

1. Look for section headers using English or Dutch case-insensitive keywords (Verse, Chorus, Refrain, Pre-Chorus, Bridge, Intro, Outro, Ending, Instrumental, Interlude, Tag, Turnaround, Vamp, Refrain, PreChorus, PostChorus, Post-Chorus, Breakdown, Verse 1, Verse 2, Chorus 1, Chorus 2, Intro, Uitro, Refrein, Couplet, Brug, etc.). They should be the first word at the start of a line, but the numbers are important too.
2. If no keywords (section headers) are found, treat all lyrics as one big section and name this section "General".
3. Make sure to use an array while collecting, since duplicate keywords (section headers) do exist and the order is important.
4. Make sure each keyword (section header) is unique at the end by adding/incrementing a number to make them unique
