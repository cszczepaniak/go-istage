package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/cszczepaniak/go-istage/patch"
	git "github.com/libgit2/git2go/v34"
)

type GitService struct {
	repo *git.Repository
	env  GitEnvironment

	repoChanged func()
}

func NewGitService(env GitEnvironment) (*GitService, error) {
	gs := &GitService{
		env: env,
	}
	err := gs.UpdateRepository()
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func (gs *GitService) UpdateRepository() error {
	repo, err := git.OpenRepository(gs.env.repoPath)
	if err != nil {
		return err
	}

	gs.repo = repo

	if gs.repoChanged != nil {
		gs.repoChanged()
	}

	return nil
}

func (gs *GitService) ApplyPatch(patchContents string, dir patch.Direction) error {
	path, err := writePatchToFile(patchContents)
	if err != nil {
		return err
	}
	defer os.Remove(path)

	isUndo := dir == patch.Reset || dir == patch.Unstage

	b := gs.Exec(`apply`)
	b.WithArgs(`-v`, `--whitespace=nowarn`)

	if isUndo {
		b.WithArgs(`--reverse`)
	}
	if dir != patch.Reset {
		b.WithArgs(`--cached`)
	}
	b.WithArgs(fmt.Sprintf(`%q`, path))

	return b.Run()
}

func writePatchToFile(contents string) (string, error) {
	patchFile, err := os.CreateTemp(``, ``)
	if err != nil {
		return ``, err
	}
	defer patchFile.Close()

	_, err = io.WriteString(patchFile, contents)
	if err != nil {
		return ``, err
	}

	return patchFile.Name(), nil
}

type gitExecBuilder struct {
	gs *GitService

	updateRepo bool
	capture    bool
	args       []string
}

func (gs *GitService) Exec(name string) *gitExecBuilder {
	return &gitExecBuilder{
		gs: gs,

		updateRepo: true,
		capture:    true,
		args:       []string{name},
	}
}

func (eb *gitExecBuilder) SkipUpdate() *gitExecBuilder {
	eb.updateRepo = false
	return eb
}

func (eb *gitExecBuilder) SkipCapture() *gitExecBuilder {
	eb.capture = false
	return eb
}

func (eb *gitExecBuilder) WithArgs(a ...string) *gitExecBuilder {
	eb.args = append(eb.args, a...)
	return eb
}

func (eb *gitExecBuilder) Run() error {
	cmd := exec.Command(eb.gs.env.pathToGit, eb.args...)

	cmd.Dir = eb.gs.repo.Workdir()

	var out strings.Builder
	if eb.capture {
		cmd.Stdout = &out
		cmd.Stderr = &out
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	if eb.capture {
		output := out.String()
		sc := bufio.NewScanner(strings.NewReader(output))
		for sc.Scan() {
			txt := sc.Text()
			if strings.HasPrefix(txt, `fatal:`) || strings.HasPrefix(txt, `error:`) {
				return fmt.Errorf("error executing %+v:\n%s", eb.args, output)
			}
		}
	}

	if eb.updateRepo {
		return eb.gs.UpdateRepository()
	}
	return nil
}

func (gs *GitService) OnRepoChanged(fn func()) {
	if fn == nil {
		return
	}
	if gs.repoChanged == nil {
		gs.repoChanged = fn
		return
	}
	old := gs.repoChanged
	gs.repoChanged = func() {
		fn()
		old()
	}
}