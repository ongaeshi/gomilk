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
	results, err := self.search(root, []string{pattern.Pattern})

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

func (self *Finder) search(root string, args []string) ([]string, error) {
	query   := strings.Join(args, " ")
	path, _ := filepath.Abs(root)
	url     := fmt.Sprintf("http://127.0.0.1:9292/gomilk?dir=%s&query=%s", url.QueryEscape(path), url.QueryEscape(query)) // @todo port, address

	if (self.Option.All) {
		url += "&all=1"
	}

	contents, err := readURL(url)

	if err != nil {
		return []string{}, err
	}
	
	// Get absolute path array from 'milk web -F'
	apaths := strings.Fields(contents)

	// abs -> abs
	// return apaths, nil

	// abs -> relative
	currentDir, _ := filepath.Abs(".")
	rpaths := make([]string, len(apaths))

	for i, apath := range apaths {
		rpath, err := filepath.Rel(currentDir, apath)

		if (err == nil) {
			rpaths[i] = rpath
		} else {
			rpaths[i] = apaths[i]
		}
	}

	return rpaths, nil
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


