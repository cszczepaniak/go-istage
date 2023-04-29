package patch

import (
	"strings"
)

type Changes struct {
	Path    string
	OldPath string

	Mode    string
	OldMode string
}

func ChangesFromHeader(header []Line) Changes {
	ch := Changes{}

	for _, l := range header {
		if l.Kind != HeaderLine && l.Kind != DiffLine {
			return ch
		}

		if strings.HasPrefix(l.Text, `old mode `) {
			ch.OldMode = strings.TrimPrefix(l.Text, `old mode `)
			continue
		}

		if strings.HasPrefix(l.Text, `new mode `) {
			ch.Mode = strings.TrimPrefix(l.Text, `new mode `)
			continue
		}

		if strings.HasPrefix(l.Text, `new file mode `) {
			ch.Mode = strings.TrimPrefix(l.Text, `new file mode `)
			continue
		}

		if strings.HasPrefix(l.Text, `deleted file mode `) {
			ch.OldMode = strings.TrimPrefix(l.Text, `deleted file mode `)
			continue
		}

		if strings.HasPrefix(l.Text, `--- a/`) && !strings.Contains(l.Text, `/dev/null`) {
			ch.OldPath = strings.TrimPrefix(l.Text, `--- a/`)
			continue
		}

		if strings.HasPrefix(l.Text, `+++ b/`) {
			ch.Path = strings.TrimPrefix(l.Text, `+++ b/`)
			continue
		}
	}

	return ch
}
