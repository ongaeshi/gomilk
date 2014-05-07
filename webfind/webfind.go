package webfind

import (
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root string, pattern *pattern.Pattern) {
	results, err := search([]string{pattern.Pattern})

	if err != nil {
		close(self.Out)
		return
	}

	for _, path := range results {
		fileType := ""
		if self.Option.FilesWithRegexp == "" {
			fileType = file.IdentifyType(path)
			if fileType == file.ERROR || fileType == file.BINARY {
				continue
			}
		}
		self.Out <- &grep.Params{path, fileType, pattern}
	}

	close(self.Out)
}

func search(args []string) ([]string, error) {
	query := strings.Join(args, " ")

	contents, err := readURL(fmt.Sprintf("http://127.0.0.1:9292/gmilk?package=milkode&query=%s", url.QueryEscape(query))) // @todo package, port, address

	if err != nil {
		return []string{}, err
	}
	
	return strings.Fields(contents), nil
}

func readURL(url string) (string, error) {
	response, err := http.Get(url)
	
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	
	if err != nil {
		return "", err
	}

	return string(contents), nil
}


