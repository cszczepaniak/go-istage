func ChangesFromHeader(header []Line) Changes {
		if l.Kind != HeaderLine && l.Kind != DiffLine {
			return ch
			ch.OldPath = strings.TrimPrefix(l.Text, `--- a/`)
		if strings.HasPrefix(l.Text, `+++ b/`) {
			ch.Path = strings.TrimPrefix(l.Text, `+++ b/`)
	return ch