# go-istage

`go-istage` is a TUI utility based on Immo Landwerth's [git-istage](https://github.com/terrajobst/git-istage).

His tool is much more fully featured than this one, so you should definitely check it out.

## libgit2

This project currently uses git2go to perform certain Git operations. git2go depends on libgit2, so in order to build 
this project, you need to have the correct version of libgit2 installed (which is 1.5.0).

To be honest, I don't really know the _right_ wait to do this. Here's how I did it:
- After running `go mod download` here, find out where git2go is installed (should be `$(go env GOPATH)/pkg/mod/github.com/libgit2/git2go/v34@v34.0.0)`)
- Head to that directory
- Setup libgit2; all of these commands might need `sudo`
  - `mkdir vendor`
  - `cd vendor`
  - `git clone -b v1.5.0 https://github.com/libgit2/libgit2.git`
  - `cd -`
  - `chmod +x ./script/build-libgit2-static.sh`
  - `chmod +x ./script/build-libgit2.sh`
  - `make install-static`
- Now you can head back to this repo, and `make build`, `make test`, etc. should all work to produce static executables of this project