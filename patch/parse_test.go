package patch

import (
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDocument(t *testing.T) {
	const patchDivider = `/* DIVIDER */`
	bs, err := os.ReadFile(`parse_test_data.txt`)
	require.NoError(t, err)

	patches := strings.Split(string(bs), patchDivider)
	for i, p := range patches {
		patches[i] = strings.TrimLeftFunc(p, unicode.IsSpace)
	}
	require.Len(t, patches, 4)

	entry0 := Entry{
		Offset: 0,
		Length: 48,
		Hunks: []Hunk{{
			Offset:    4,
			Length:    13,
			OldStart:  2,
			OldLength: 6,
			NewStart:  2,
			NewLength: 12,
		}, {
			Offset:    17,
			Length:    9,
			OldStart:  11,
			OldLength: 7,
			NewStart:  17,
			NewLength: 7,
		}, {
			Offset:    26,
			Length:    15,
			OldStart:  23,
			OldLength: 11,
			NewStart:  29,
			NewLength: 9,
		}, {
			Offset:    41,
			Length:    7,
			OldStart:  35,
			OldLength: 4,
			NewStart:  39,
			NewLength: 4,
		}},
		Changes: Changes{
			Path:    `change_file.txt`,
			OldPath: `change_file.txt`,
			Mode:    ``,
			OldMode: ``,
		},
	}
	entry1 := Entry{
		Offset: 48,
		Length: 3,
		Hunks:  []Hunk{},
		Changes: Changes{
			Mode:    `100755`,
			OldMode: `100644`,
		},
	}
	entry2 := Entry{
		Offset: 51,
		Length: 10,
		Hunks: []Hunk{{
			Offset:    56,
			Length:    5,
			OldStart:  0,
			OldLength: 0,
			NewStart:  1,
			NewLength: 4,
		}},
		Changes: Changes{
			Path:    `new_file.txt`,
			OldPath: ``,
			Mode:    `100644`,
			OldMode: ``,
		},
	}
	entry3 := Entry{
		Offset: 61,
		Length: 11,
		Hunks: []Hunk{{
			Offset:    66,
			Length:    6,
			OldStart:  1,
			OldLength: 5,
		}},
		Changes: Changes{
			Path:    ``,
			OldPath: `deleted_file.txt`,
			Mode:    ``,
			OldMode: `100755`,
		},
	}

	expLines := []Line{{
		Kind: DiffLine,
		Text: `diff --git a/change_file.txt b/change_file.txt`,
	}, {
		Kind: HeaderLine,
		Text: `index 2674c46..6381df1 100644`,
	}, {
		Kind: HeaderLine,
		Text: `--- a/change_file.txt`,
	}, {
		Kind: HeaderLine,
		Text: `+++ b/change_file.txt`,
	}, {
		Kind: HunkLine,
		Text: `@@ -2,6 +2,12 @@`,
	}, {
		Kind: ContextLine,
		Text: ` 1`,
	}, {
		Kind: ContextLine,
		Text: ` 2`,
	}, {
		Kind: ContextLine,
		Text: ` 3`,
	}, {
		Kind: AdditionLine,
		Text: `+a`,
	}, {
		Kind: AdditionLine,
		Text: `+b`,
	}, {
		Kind: AdditionLine,
		Text: `+c`,
	}, {
		Kind: AdditionLine,
		Text: `+d`,
	}, {
		Kind: AdditionLine,
		Text: `+e`,
	}, {
		Kind: AdditionLine,
		Text: `+f`,
	}, {
		Kind: ContextLine,
		Text: ` 4`,
	}, {
		Kind: ContextLine,
		Text: ` 5`,
	}, {
		Kind: ContextLine,
		Text: ` 6`,
	}, {
		Kind: HunkLine,
		Text: `@@ -11,7 +17,7 @@`,
	}, {
		Kind: ContextLine,
		Text: ` 10`,
	}, {
		Kind: ContextLine,
		Text: ` 11`,
	}, {
		Kind: ContextLine,
		Text: ` 12`,
	}, {
		Kind: RemovalLine,
		Text: `-13`,
	}, {
		Kind: AdditionLine,
		Text: `+thirteen`,
	}, {
		Kind: ContextLine,
		Text: ` 14`,
	}, {
		Kind: ContextLine,
		Text: ` 15`,
	}, {
		Kind: ContextLine,
		Text: ` 16`,
	}, {
		Kind: HunkLine,
		Text: `@@ -23,11 +29,9 @@`,
	}, {
		Kind: ContextLine,
		Text: ` 22`,
	}, {
		Kind: ContextLine,
		Text: ` 23`,
	}, {
		Kind: ContextLine,
		Text: ` 24`,
	}, {
		Kind: RemovalLine,
		Text: `-25`,
	}, {
		Kind: RemovalLine,
		Text: `-26`,
	}, {
		Kind: RemovalLine,
		Text: `-27`,
	}, {
		Kind: RemovalLine,
		Text: `-28`,
	}, {
		Kind: RemovalLine,
		Text: `-29`,
	}, {
		Kind: AdditionLine,
		Text: `+x`,
	}, {
		Kind: AdditionLine,
		Text: `+y`,
	}, {
		Kind: AdditionLine,
		Text: `+z`,
	}, {
		Kind: ContextLine,
		Text: ` 30`,
	}, {
		Kind: ContextLine,
		Text: ` 31`,
	}, {
		Kind: ContextLine,
		Text: ` 32`,
	}, {
		Kind: HunkLine,
		Text: `@@ -35,4 +39,4 @@`,
	}, {
		Kind: ContextLine,
		Text: ` 34`,
	}, {
		Kind: ContextLine,
		Text: ` 35`,
	}, {
		Kind: ContextLine,
		Text: ` 36`,
	}, {
		Kind: RemovalLine,
		Text: `-37`,
	}, {
		Kind: AdditionLine,
		Text: `+37`,
	}, {
		Kind: NoEndOfLineLine,
		Text: `\ No newline at end of file`,
	}, {
		Kind: DiffLine,
		Text: `diff --git a/change_mode.txt b/change_mode.txt`,
	}, {
		Kind: HeaderLine,
		Text: `old mode 100644`,
	}, {
		Kind: HeaderLine,
		Text: `new mode 100755`,
	}, {
		Kind: DiffLine,
		Text: `diff --git a/new_file.txt b/new_file.txt`,
	}, {
		Kind: HeaderLine,
		Text: `new file mode 100644`,
	}, {
		Kind: HeaderLine,
		Text: `index 0000000..1fa3451`,
	}, {
		Kind: HeaderLine,
		Text: `--- /dev/null`,
	}, {
		Kind: HeaderLine,
		Text: `+++ b/new_file.txt`,
	}, {
		Kind: HunkLine,
		Text: `@@ -0,0 +1,4 @@`,
	}, {
		Kind: AdditionLine,
		Text: `+abc`,
	}, {
		Kind: AdditionLine,
		Text: `+abc`,
	}, {
		Kind: AdditionLine,
		Text: `+abc`,
	}, {
		Kind: AdditionLine,
		Text: `+abc`,
	}, {
		Kind: DiffLine,
		Text: `diff --git a/deleted_file.txt b/deleted_file.txt`,
	}, {
		Kind: HeaderLine,
		Text: `deleted file mode 100755`,
	}, {
		Kind: HeaderLine,
		Text: `index e51e7d2..0000000`,
	}, {
		Kind: HeaderLine,
		Text: `--- a/deleted_file.txt`,
	}, {
		Kind: HeaderLine,
		Text: `+++ /dev/null`,
	}, {
		Kind: HunkLine,
		Text: `@@ -1,5 +0,0 @@`,
	}, {
		Kind: RemovalLine,
		Text: `-asdasdasd`,
	}, {
		Kind: RemovalLine,
		Text: `-`,
	}, {
		Kind: RemovalLine,
		Text: `-def`,
	}, {
		Kind: RemovalLine,
		Text: `-`,
	}, {
		Kind: RemovalLine,
		Text: `-ghi`,
	}}

	doc := ParseDocument(patches)

	require.Len(t, doc.Lines, len(expLines))

	for i, exp := range expLines {
		l := doc.Lines[i]
		l.LineBreak = ``
		assert.Equal(t, exp, l, `line %d was wrong`, i)
	}

	require.Len(t, doc.Entries, 4)
	for i, exp := range []Entry{entry0, entry1, entry2, entry3} {
		act := doc.Entries[i]
		assert.Equal(t, exp.Changes, act.Changes, `changes for entry %d`, i)
		assert.Equal(t, exp.Length, act.Length, `length for entry %d`, i)
		assert.Equal(t, exp.Offset, act.Offset, `offset for entry %d`, i)

		require.Len(t, act.Hunks, len(exp.Hunks))
		for j, hunk := range exp.Hunks {
			assert.Equal(t, hunk, act.Hunks[j], `entry %d hunk %d`, i, j)
		}
	}
}

func TestParseLines(t *testing.T) {
	tests := []struct {
		desc       string
		inputLines []line
		expLines   []Line
	}{{
		desc: `header with no hunk`,
		inputLines: []line{{
			text:      `literally does not matter`,
			lineBreak: "\n",
		}, {
			text:      `like at all`,
			lineBreak: "\n",
		}, {
			text:      `diff --git foobar`,
			lineBreak: "\n",
		}},
		expLines: []Line{{
			Kind:      HeaderLine,
			Text:      `literally does not matter`,
			LineBreak: "\n",
		}, {
			Kind:      HeaderLine,
			Text:      `like at all`,
			LineBreak: "\n",
		}, {
			Kind:      DiffLine,
			Text:      `diff --git foobar`,
			LineBreak: "\n",
		}},
	}, {
		desc: `hunk`,
		inputLines: []line{{
			text: `@@ does not have to be a well-formed hunk`,
		}, {
			text: `+ addition`,
		}, {
			text: `- removal`,
		}, {
			text: `context`,
		}, {
			text: `context`,
		}, {
			text: `\ no end of line`,
		}},
		expLines: []Line{{
			Kind: HunkLine,
			Text: `@@ does not have to be a well-formed hunk`,
		}, {
			Kind: AdditionLine,
			Text: `+ addition`,
		}, {
			Kind: RemovalLine,
			Text: `- removal`,
		}, {
			Kind: ContextLine,
			Text: `context`,
		}, {
			Kind: ContextLine,
			Text: `context`,
		}, {
			Kind: NoEndOfLineLine,
			Text: `\ no end of line`,
		}},
	}}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			ls := ParseLines(tc.inputLines)
			assert.Equal(t, tc.expLines, ls)
		})
	}
}

func TestGetLines(t *testing.T) {
	tests := []struct {
		desc     string
		content  string
		expLines []line
	}{{
		desc:    `LF no trailing`,
		content: "abc\ndef",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\n",
		}, {
			text:      `def`,
			lineBreak: ``,
		}},
	}, {
		desc:    `LF with trailing`,
		content: "abc\ndef\n",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\n",
		}, {
			text:      `def`,
			lineBreak: "\n",
		}},
	}, {
		desc:    `CR no trailing`,
		content: "abc\rdef",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\r",
		}, {
			text:      `def`,
			lineBreak: "",
		}},
	}, {
		desc:    `CR with trailing`,
		content: "abc\rdef\r",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\r",
		}, {
			text:      `def`,
			lineBreak: "\r",
		}},
	}, {
		desc:    `CRLF no trailing`,
		content: "abc\r\ndef",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\r\n",
		}, {
			text:      `def`,
			lineBreak: "",
		}},
	}, {
		desc:    `CRLF with trailing`,
		content: "abc\r\ndef\r\n",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\r\n",
		}, {
			text:      `def`,
			lineBreak: "\r\n",
		}},
	}, {
		desc:    `kitchen sink`,
		content: "abc\r\ndef\nghi\rḦeḼḼo unicode",
		expLines: []line{{
			text:      `abc`,
			lineBreak: "\r\n",
		}, {
			text:      `def`,
			lineBreak: "\n",
		}, {
			text:      `ghi`,
			lineBreak: "\r",
		}, {
			text:      `ḦeḼḼo unicode`,
			lineBreak: "",
		}},
	}}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expLines, GetLines(tc.content))
		})
	}
}
