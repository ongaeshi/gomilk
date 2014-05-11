package find

import (
	"fmt"
	"github.com/ongaeshi/gomilk/search/file"
	"github.com/ongaeshi/gomilk/search/grep"
	"github.com/ongaeshi/gomilk/search/option"
	"github.com/ongaeshi/gomilk/search/pattern"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root string, pattern *pattern.Pattern) {
	results, err := search(root, []string{pattern.Pattern})

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

func search(root string, args []string) ([]string, error) {
	query   := strings.Join(args, " ")
	path, _ := filepath.Abs(root)
	url     := fmt.Sprintf("http://127.0.0.1:9292/gmilk?dir=%s&query=%s", url.QueryEscape(path), url.QueryEscape(query))  // @todo package, port, address

	contents, err := readURL(url)

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


