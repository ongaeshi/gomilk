package search

import (
	"github.com/ongaeshi/gomilk/search/grep"
	"github.com/ongaeshi/gomilk/search/option"
	"github.com/ongaeshi/gomilk/search/pattern"
	"github.com/ongaeshi/gomilk/search/print"
	"github.com/ongaeshi/gomilk/find"
)

type Searcher struct {
	Root, Pattern string
	Option        *option.Option
}

func (self *Searcher) Search() error {
	pattern, err := self.pattern()
	if err != nil {
		return err
	}
	grep := make(chan *grep.Params, self.Option.Proc)
	match := make(chan *print.Params, self.Option.Proc)
	done := make(chan bool)
	go self.find(grep, pattern)
	go self.grep(grep, match)
	go self.print(match, done)
	<-done
	return nil
}

func (self *Searcher) pattern() (*pattern.Pattern, error) {
	fileRegexp := self.Option.FileSearchRegexp
	if self.Option.FilesWithRegexp != "" {
		fileRegexp = self.Option.FilesWithRegexp
	}
	return pattern.NewPattern(
		self.Pattern,
		fileRegexp,
		self.Option.SmartCase,
		self.Option.IgnoreCase,
	)
}

func (self *Searcher) find(out chan *grep.Params, pattern *pattern.Pattern) {
	finder := find.Finder{out, self.Option}
	finder.Find(self.Root, pattern)
}

func (self *Searcher) grep(in chan *grep.Params, out chan *print.Params) {
	grepper := grep.Grepper{in, out, self.Option}
	grepper.ConcurrentGrep()
}

func (self *Searcher) print(in chan *print.Params, done chan bool) {
	printer := print.NewPrinter(in, done, self.Option)
	printer.Print()
}

