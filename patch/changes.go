package patch

import (
	"errors"
	"strings"
)

type Changes struct {
	Path    string
	OldPath string

	Mode    string
	OldMode string
}

func ChangesFromHeader(header []Line) (Changes, error) {
	ch := Changes{}

	for _, l := range header {
		if l.Kind != HeaderLine {
			return Changes{}, errors.New(`dev error: non-header line found`)
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

		if strings.HasPrefix(l.Text, `--- a/`) && !strings.Contains(l.Text, `/dev/null`) {
			ch.OldPath = strings.TrimPrefix(`--- a/`, l.Text)
			continue
		}

		if strings.HasPrefix(l.Text, `--- b/`) {
			ch.Path = strings.TrimPrefix(`--- b/`, l.Text)
			continue
		}
	}

	return ch, nil
}
