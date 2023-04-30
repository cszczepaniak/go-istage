//go:build !prod

package git

import (
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"testing"

	git "github.com/libgit2/git2go/v34"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRepo struct {
	t    testing.TB
	env  Environment
	repo *git.Repository
	path string
}

func NewTestRepo(t testing.TB) testRepo {
	t.Helper()

	tempDir, err := os.MkdirTemp(``, ``)
	require.NoError(t, err, `failed to make temporary directory`)

	t.Cleanup(func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	})

	require.NoError(t, os.Chdir(tempDir), `failed to change to temporary directory`)

	out, err := exec.Command(`git`, `init`, `-b`, `master`).CombinedOutput()
	require.NoError(t, err, "failed to init git repo:\n %s\n", out)

	out, err = exec.Command(`git`, `commit`, `--allow-empty`, `-m`, `initial commit`).CombinedOutput()
	require.NoError(t, err, "failed to make initial git commit:\n %s\n", out)

	env, err := NewEnvironment(``, ``)
	require.NoError(t, err, `failed to init git environment`)

	repo, err := git.OpenRepository(env.repoPath)
	require.NoError(t, err, `failed to open git repository`)

	return testRepo{
		t:    t,
		env:  env,
		repo: repo,
		path: tempDir,
	}
}

func (tr testRepo) Add(path string) {
	err := tr.env.Exec(tr.repo, `add`).WithArgs(path).Run()
	require.NoError(tr.t, err)
}

func (tr testRepo) AddAll() {
	tr.Add(`.`)
}

func (tr testRepo) ClearUnstagedChanges() {
	err := tr.env.Exec(tr.repo, `checkout`).WithArgs(`.`).Run()
	require.NoError(tr.t, err)
}

func (tr testRepo) Commit(msg string) {
	err := tr.env.Exec(tr.repo, `commit`).WithArgs(`-m`, msg).Run()
	require.NoError(tr.t, err)
}

func (tr testRepo) MakeFile(t testing.TB, path string) *fileBuilder {
	return defaultFileBuilder(t, tr, path)
}

type fileAdditionMode int

const (
	unstaged fileAdditionMode = iota
	staged
	committed
)

type fileBuilder struct {
	sb *strings.Builder

	t  testing.TB
	tr testRepo

	path         string
	additionMode fileAdditionMode
	commitMsg    string
	newline      string
	filemode     fs.FileMode
}

func defaultFileBuilder(t testing.TB, tr testRepo, path string) *fileBuilder {
	return &fileBuilder{
		sb: &strings.Builder{},

		tr:           tr,
		path:         path,
		additionMode: unstaged,
		newline:      "\n",
		filemode:     0o644,
	}
}

func (fb *fileBuilder) ShouldStage() *fileBuilder {
	fb.additionMode = staged
	return fb
}

func (fb *fileBuilder) ShouldCommit(msg string) *fileBuilder {
	fb.additionMode = committed
	fb.commitMsg = msg
	return fb
}

func (fb *fileBuilder) WithNewline(newline string) *fileBuilder {
	fb.newline = newline
	return fb
}

func (fb *fileBuilder) WithFileMode(mode fs.FileMode) *fileBuilder {
	fb.filemode = mode
	return fb
}

func (fb *fileBuilder) AddLine(content string) *fileBuilder {
	fb.sb.WriteString(content + fb.newline)
	return fb
}

func (fb *fileBuilder) Add(content string) *fileBuilder {
	fb.sb.WriteString(content)
	return fb
}

func (fb *fileBuilder) Write(b []byte) (int, error) {
	return fb.sb.Write(b)
}

type fileBuildResult struct {
	t    testing.TB
	path string
}

func (fb *fileBuilder) Build() fileBuildResult {
	err := os.WriteFile(fb.path, []byte(fb.sb.String()), fb.filemode)
	require.NoError(fb.t, err)

	if fb.additionMode > unstaged {
		fb.tr.Add(fb.path)
	}
	if fb.additionMode > staged {
		fb.tr.Commit(fb.commitMsg)
	}

	return fileBuildResult{
		t:    fb.t,
		path: fb.path,
	}
}

func (fbr fileBuildResult) Append(s string) {
	f, err := os.OpenFile(fbr.path, os.O_APPEND|os.O_WRONLY, 0o644)
	require.NoError(fbr.t, err)
	defer f.Close()

	_, err = f.Write([]byte(s))
	require.NoError(fbr.t, err)
}

func (fbr fileBuildResult) Replace(s string) {
	f, err := os.OpenFile(fbr.path, os.O_TRUNC|os.O_WRONLY, 0o644)
	require.NoError(fbr.t, err)
	defer f.Close()

	_, err = f.Write([]byte(s))
	require.NoError(fbr.t, err)
}

func (fbr fileBuildResult) Remove() {
	err := os.Remove(fbr.path)
	require.NoError(fbr.t, err)
}

func (fbr fileBuildResult) Rename(name string) {
	require.NotEqual(fbr.t, fbr.path, name, `cannot rename file to itself`)

	bs, err := os.ReadFile(fbr.path)
	require.NoError(fbr.t, err)

	err = os.WriteFile(name, bs, 0o644)
	require.NoError(fbr.t, err)

	err = os.Remove(fbr.path)
	require.NoError(fbr.t, err)
}
