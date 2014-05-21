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
	"os"
	"path/filepath"
	"regexp"
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
	query := strings.Join(args, " ")
	path, _ := filepath.Abs(root)
	url := fmt.Sprintf("http://127.0.0.1:%d/gomilk?dir=%s&query=%s", self.Option.Port, url.QueryEscape(path), url.QueryEscape(query)) // @todo address

	if self.Option.All {
		url += "&all=1"
	}

	contents, err := readURL(url)

	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Need \"milk web --gomilk\"\n")
		os.Exit(1)
	}

	// Get absolute path array from 'milk web -g'
	apaths := strings.Split(contents, "\n")

	// Error?
	if len(apaths) > 0 {
		submatch := regexp.MustCompile("Error: (.*)").FindSubmatch([]byte(apaths[0]))
		if len(submatch) > 0 && submatch[1] != nil {
			fmt.Printf("Get %s: response is error\n", url)
			fmt.Printf("%s\n", submatch[1])
			os.Exit(1)
		}
	}

	// Return match paths
	if self.Option.ExpandPath {
		// abs -> abs
		return apaths, nil

	} else {
		// abs -> relative
		currentDir, _ := filepath.Abs(".")
		rpaths := make([]string, len(apaths))

		for i, apath := range apaths {
			rpath, err := filepath.Rel(currentDir, apath)

			if err == nil {
				rpaths[i] = rpath
			} else {
				rpaths[i] = apaths[i]
			}
		}

		return rpaths, nil
	}
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
