package git

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cszczepaniak/go-istage/nolibgit"
)

type Environment struct {
	repoPath  string
	pathToGit string
}

func NewEnvironment(repoPath string, pathToGit string) (Environment, error) {
	if pathToGit == `` {
		var err error
		pathToGit, err = resolveGitPath()
		if err != nil {
			return Environment{}, err
		}
	}
	if repoPath == `` {
		var err error
		repoPath, err = resolveRepoPath()
		if err != nil {
			return Environment{}, err
		}
	}
	return Environment{
		repoPath:  repoPath,
		pathToGit: pathToGit,
	}, nil
}

func newEnvironmentFromNoLibGit(env nolibgit.Environment) Environment {
	return Environment{
		repoPath:  env.RepoDir,
		pathToGit: env.GitExecutable,
	}
}

func resolveGitPath() (string, error) {
	path := os.Getenv("PATH")
	if path == `` {
		return ``, nil
	}

	paths := strings.Split(path, string(os.PathListSeparator))
	lookingFor := []string{"git", "git.exe"}

	var res string
	for _, p := range paths {
		for _, exe := range lookingFor {
			res = filepath.Join(p, exe)
			_, err := os.Stat(res)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return ``, err
			}

			return res, nil
		}
	}
	return ``, errors.New(`could not find git`)
}

func resolveRepoPath() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return ``, err
	}

	return resolveRepoPathFrom(d)
}

func resolveRepoPathFrom(start string) (string, error) {
	entries, err := os.ReadDir(start)
	if err != nil {
		return ``, err
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasSuffix(e.Name(), `.git`) {
			path, err := filepath.Abs(filepath.Join(start, e.Name()))
			if err != nil {
				return ``, err
			}
			return path, nil
		}
	}
	return resolveRepoPathFrom(path.Join(`..`, start))
}
