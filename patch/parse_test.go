package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
