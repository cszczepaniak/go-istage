package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/cszczepaniak/go-istage/logging"
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

	err := logging.Init()
	if err != nil {
		panic(`failed to initialize logging: ` + err.Error())
	}

	gitEnv, err := services.NewGitEnvironment(``, ``)
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

	ps := services.NewPatchingService(gs)

	doc, err := ds.UnstagedChanges()
	if err != nil {
		log.Fatalln(err)
	}

	err = ui.RunUI(doc, ps, ds)
	if err != nil {
		log.Fatalln(err)
	}
}
