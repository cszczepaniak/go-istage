package git

import (
	"strings"

	"github.com/cszczepaniak/go-istage/patch"
	git "github.com/libgit2/git2go/v34"
)

type Client struct {
	repo *git.Repository
	env  Environment
}

func NewClient(env Environment) (*Client, error) {
	c := &Client{
		env: env,
	}
	err := c.UpdateRepository()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) UpdateRepository() error {
	repo, err := git.OpenRepository(c.env.repoPath)
	if err != nil {
		return err
	}

	c.repo = repo

	return nil
}

func (c *Client) ApplyPatch(patchContents string, dir patch.Direction) error {
	isUndo := dir == patch.Reset || dir == patch.Unstage

	b := c.Exec(`apply`).WithStdin(strings.NewReader(patchContents))

	b.WithArgs(`-v`)
	if dir != patch.Reset {
		b.WithArgs(`--cached`)
	}
	if isUndo {
		b.WithArgs(`--reverse`)
	}
	b.WithArgs(`--whitespace=nowarn`)

	return b.Run()
}

func (c *Client) UnstagedFiles() ([]File, error) {
	opts := &git.StatusOptions{
		Show:  git.StatusShowWorkdirOnly,
		Flags: git.StatusOptIncludeUntracked | git.StatusOptRecurseUntrackedDirs | git.StatusOptRenamesIndexToWorkdir,
	}
	sl, err := c.repo.StatusList(opts)
	if err != nil {
		return nil, err
	}

	n, err := sl.EntryCount()
	if err != nil {
		return nil, err
	}

	res := make([]File, 0, n)
	for i := 0; i < n; i++ {
		e, err := sl.ByIndex(i)
		if err != nil {
			return nil, err
		}
		res = append(res, File{
			Path:   e.IndexToWorkdir.NewFile.Path,
			Status: fileStatusFromGitDelta(e.IndexToWorkdir.Status),
		})
	}

	return res, nil
}

func (c *Client) UnstagedChanges() ([]string, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	opts.Flags |= git.DiffShowUntrackedContent
	opts.Flags |= git.DiffRecurseUntracked

	diff, err := c.repo.DiffIndexToWorkdir(nil, &opts)
	if err != nil {
		return nil, err
	}

	err = handleRenames(diff)
	if err != nil {
		return nil, err
	}

	return patchesFromDiff(diff)
}

func (c *Client) StagedChanges() ([]string, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}

	headRef, err := c.repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := c.repo.LookupCommit(headRef.Target())
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	diff, err := c.repo.DiffTreeToIndex(tree, nil, &opts)
	if err != nil {
		return nil, err
	}

	err = handleRenames(diff)
	if err != nil {
		return nil, err
	}

	return patchesFromDiff(diff)
}

func patchesFromDiff(diff *git.Diff) ([]string, error) {
	ndl, err := diff.NumDeltas()
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, ndl)
	for i := 0; i < ndl; i++ {
		p, err := diff.Patch(i)
		if err != nil {
			return nil, err
		}

		txt, err := p.String()
		if err != nil {
			return nil, err
		}

		if txt != `` {
			res = append(res, txt)
		}
	}

	return res, nil
}

func handleRenames(diff *git.Diff) error {
	findOpts, err := git.DefaultDiffFindOptions()
	if err != nil {
		return err
	}
	findOpts.Flags |= git.DiffFindRenames
	findOpts.Flags |= git.DiffFindForUntracked

	return diff.FindSimilar(&findOpts)
}
