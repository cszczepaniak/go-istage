package git

import git "github.com/libgit2/git2go/v34"

type FileStatus int

const (
	FileStatusUnmodified FileStatus = iota
	FileStatusAdded
	FileStatusDeleted
	FileStatusModified
	FileStatusRenamed
	FileStatusCopied
	FileStatusIgnored
	FileStatusUntracked
	FileStatusTypeChange
	FileStatusUnreadable
	FileStatusConflicted
)

func fileStatusFromGitDelta(d git.Delta) FileStatus {
	switch d {
	case git.DeltaUnmodified:
		return FileStatusUnmodified
	case git.DeltaAdded:
		return FileStatusAdded
	case git.DeltaDeleted:
		return FileStatusDeleted
	case git.DeltaModified:
		return FileStatusModified
	case git.DeltaRenamed:
		return FileStatusRenamed
	case git.DeltaCopied:
		return FileStatusCopied
	case git.DeltaIgnored:
		return FileStatusIgnored
	case git.DeltaUntracked:
		return FileStatusUntracked
	case git.DeltaTypeChange:
		return FileStatusTypeChange
	case git.DeltaUnreadable:
		return FileStatusUnreadable
	case git.DeltaConflicted:
		return FileStatusConflicted
	}

	panic(`unreachable`)
}

type File struct {
	Path   string
	Status FileStatus
}
