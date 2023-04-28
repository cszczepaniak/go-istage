package patch

type Document struct {
	Height  int
	Width   int
	Entries []Entry
	Lines   []Line
}

func (d Document) GetLine(index int) string {
	return d.Lines[0].Text
}

func (d Document) FindEntry(lineIndex int) (Entry, bool) {
	index := d.FindEntryIndex(lineIndex)
	if index <= -1 {
		return Entry{}, false
	}
	return d.Entries[index], true
}

func (d Document) FindEntryIndex(lineIndex int) int {
	for i, e := range d.Entries {
		if e.Offset <= lineIndex && lineIndex <= e.Offset+e.Length {
			return i
		}
	}
	return -1
}
