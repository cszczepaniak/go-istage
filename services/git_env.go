package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	git "github.com/libgit2/git2go/v34"
)

type GitEnvironment struct {
	repoPath  string
	pathToGit string
}

func NewGitEnvironment(repoPath string, pathToGit string) (GitEnvironment, error) {
	if pathToGit == `` {
		var err error
		pathToGit, err = resolveGitPath()
		if err != nil {
			return GitEnvironment{}, err
		}
	}
	if repoPath == `` {
		var err error
		repoPath, err = resolveRepoPath()
		if err != nil {
			return GitEnvironment{}, err
		}
	}
	return GitEnvironment{
		repoPath:  repoPath,
		pathToGit: pathToGit,
	}, nil
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
	return git.Discover(d, false, nil)
}
