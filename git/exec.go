package git

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/cszczepaniak/go-istage/nolibgit"
	git "github.com/libgit2/git2go/v34"
)

type GitExecBuilder struct {
	env  nolibgit.Environment
	repo *git.Repository

	stdin      io.Reader
	updateRepo bool
	capture    bool
	args       []string
}

func (gs *Client) Exec(name string) *GitExecBuilder {
	return &GitExecBuilder{
		env:  gs.env,
		repo: gs.repo,

		updateRepo: true,
		capture:    true,
		args:       []string{name},
	}
}

func Exec(env nolibgit.Environment, repo *git.Repository, name string) *GitExecBuilder {
	return &GitExecBuilder{
		env:  env,
		repo: repo,

		updateRepo: true,
		capture:    true,
		args:       []string{name},
	}
}

func (eb *GitExecBuilder) WithStdin(r io.Reader) *GitExecBuilder {
	eb.stdin = r
	return eb
}

func (eb *GitExecBuilder) SkipUpdate() *GitExecBuilder {
	eb.updateRepo = false
	return eb
}

func (eb *GitExecBuilder) SkipCapture() *GitExecBuilder {
	eb.capture = false
	return eb
}

func (eb *GitExecBuilder) WithArgs(a ...string) *GitExecBuilder {
	eb.args = append(eb.args, a...)
	return eb
}

func (eb *GitExecBuilder) Run() error {
	cmd := exec.Command(eb.env.GitExecutable, eb.args...)
	cmd.Dir = eb.env.WorkingDir

	var out strings.Builder
	if eb.capture {
		cmd.Stdout = &out
		cmd.Stderr = &out
	}
	if eb.stdin != nil {
		cmd.Stdin = eb.stdin
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing %+v:\n%s", eb.args, out.String())
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
		repo, err := git.OpenRepository(eb.env.RepoDir)
		if err != nil {
			return err
		}

		if repo == nil {
			return errors.New(`nil repository returned`)
		}
		*eb.repo = *repo
	}
	return nil
}
