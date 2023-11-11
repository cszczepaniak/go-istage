package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/cszczepaniak/go-istage/git"
	"github.com/cszczepaniak/go-istage/logging"
	"github.com/cszczepaniak/go-istage/nolibgit"
	"github.com/cszczepaniak/go-istage/services"
	"github.com/cszczepaniak/go-istage/ui"
)

type lines []int

func (l *lines) UnmarshalText(bs []byte) error {
	numStrs := strings.Split(string(bs), `,`)

	ls := make(lines, 0, len(numStrs))
	for _, s := range numStrs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ls = append(ls, n)
	}

	*l = ls
	return nil
}

func (l *lines) MarshalText() ([]byte, error) {
	if l == nil {
		return nil, nil
	}
	res := &strings.Builder{}
	for i, n := range *l {
		res.WriteString(strconv.Itoa(n))
		if i < len(*l)-1 {
			res.WriteString(`,`)
		}
	}
	return []byte(res.String()), nil
}

func main() {
	linesArg := make(lines, 0)

	flag.TextVar(&linesArg, `lines`, &lines{}, ``)

	flag.Parse()

	err := logging.Init(logging.Config{
		OutputPath: `log/debug.log`,
	})
	if err != nil {
		panic(`failed to initialize logging: ` + err.Error())
	}

	gitEnv, err := nolibgit.LoadEnvironment()
	if err != nil {
		logging.Error(`failed to initialize git env`, `err`, err)
	}

	gs, err := git.NewClient(gitEnv)
	if err != nil {
		logging.Error(`failed to initialize git service`, `err`, err)
	}

	ds, err := services.NewDocumentService(gs)
	if err != nil {
		logging.Error(`failed to initialize document service`, `err`, err)
	}

	ps := services.NewPatchingService(gs)

	err = ui.RunUI(ps, ds, gs, gs)
	if err != nil {
		logging.Error(`error during UI runtime`, `err`, err)
	}
}
