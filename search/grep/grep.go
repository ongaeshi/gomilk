package grep

import (
	"bufio"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"github.com/ongaeshi/gomilk/search/file"
	"github.com/ongaeshi/gomilk/search/match"
	"github.com/ongaeshi/gomilk/search/option"
	"github.com/ongaeshi/gomilk/search/pattern"
	"github.com/ongaeshi/gomilk/search/print"
	"os"
	"sync"
)

type Params struct {
	Path, Encode string
	Pattern      *pattern.Pattern
}

type Grepper struct {
	In     chan *Params
	Out    chan *print.Params
	Option *option.Option
}

func (self *Grepper) ConcurrentGrep() {
	var wg sync.WaitGroup
	sem := make(chan bool, self.Option.Proc)
	for arg := range self.In {
		sem <- true
		wg.Add(1)
		go func(self *Grepper, arg *Params, sem chan bool) {
			defer wg.Done()
			self.Grep(arg.Path, arg.Encode, arg.Pattern, sem)
		}(self, arg, sem)
	}
	wg.Wait()
	close(self.Out)
}

func getDecoder(encode string) transform.Transformer {
	switch encode {
	case file.EUCJP:
		return japanese.EUCJP.NewDecoder()
	case file.SHIFTJIS:
		return japanese.ShiftJIS.NewDecoder()
	}
	return nil
}

func (self *Grepper) Grep(path, encode string, pattern *pattern.Pattern, sem chan bool) {
	if self.Option.FilesWithRegexp != "" {
		self.Out <- &print.Params{pattern, path, nil}
		<-sem
		return
	}

	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var f *bufio.Reader
	if dec := getDecoder(encode); dec != nil {
		f = bufio.NewReader(transform.NewReader(fh, dec))
	} else {
		f = bufio.NewReader(fh)
	}

	var buf []byte
	matches := make([]*match.Match, 0)
	m := match.NewMatch(self.Option.Before, self.Option.After)
	var lineNum = 1
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			break
		}
		if newMatch, ok := m.IsMatch(pattern, lineNum, string(buf)); ok {
			matches = append(matches, m)
			m = newMatch
		}
		lineNum++
	}
	if m.Matched {
		matches = append(matches, m)
	}
	self.Out <- &print.Params{pattern, path, matches}
	fh.Close()
	<-sem
}
