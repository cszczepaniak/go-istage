package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/cszczepaniak/go-istage/patch"
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

var (
	applyPatch = flag.Bool(`apply`, false, `apply`)
)

func main() {
	linesArg := make(lines, 0)

	flag.TextVar(&linesArg, `lines`, &lines{}, ``)

	flag.Parse()

	gitEnv, err := services.NewGitEnvironment(`/home/connor/src/go-istage`, ``)
	if err != nil {
		log.Fatalln(err)
	}

	gs, err := services.NewGitService(gitEnv)
	if err != nil {
		log.Fatalln(err)
	}

	ds, err := services.NewDocumentService(gs)
	if err != nil {
		log.Fatalln(err)
	}

	ps := services.NewPatchingService(gs, ds)

	if *applyPatch {
		err = ps.ApplyPatch(patch.Stage, false, linesArg)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err = ui.RunUI(ds.Document, ps, &docUpdater{ds: ds})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

type docUpdater struct {
	ds *services.DocumentService
}

func (du *docUpdater) UpdateDocument() (patch.Document, error) {
	err := du.ds.UpdateDocument()
	if err != nil {
		return patch.Document{}, err
	}
	return du.ds.Document, nil
}
