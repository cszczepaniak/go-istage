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

	doc := ParseDocument(patches)

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
