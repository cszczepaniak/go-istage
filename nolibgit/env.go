package nolibgit

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Environment struct {
	RepoDir       string
	WorkingDir    string
	GitExecutable string
}

func LoadEnvironment() (Environment, error) {
	path, err := resolveRepoPath()
	if err != nil {
		return Environment{}, err
	}

	git, err := resolveGitPath()
	if err != nil {
		return Environment{}, err
	}

	return Environment{
		RepoDir:       path.dir,
		WorkingDir:    path.workDir,
		GitExecutable: git,
	}, nil
}

type repoPath struct {
	dir     string
	workDir string
}

func resolveRepoPath() (repoPath, error) {
	d, err := os.Getwd()
	if err != nil {
		return repoPath{}, err
	}

	return resolveRepoPathFrom(d)
}

func resolveRepoPathFrom(start string) (repoPath, error) {
	entries, err := os.ReadDir(start)
	if err != nil {
		return repoPath{}, err
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasSuffix(e.Name(), `.git`) {
			return repoPath{
				dir:     start,
				workDir: e.Name(),
			}, nil
		}
	}
	return resolveRepoPathFrom(path.Join(`..`, start))
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
