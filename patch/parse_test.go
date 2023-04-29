package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocument(t *testing.T) {

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
