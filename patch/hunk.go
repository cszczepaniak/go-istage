package patch

type Hunk struct {
	Offset    int
	Length    int
	OldStart  int
	OldLength int
	NewStart  int
	NewLength int
}
