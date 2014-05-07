package webfind

import (
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root string, pattern *pattern.Pattern) {
	close(self.Out)
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

func Search(args []string) {
	query := strings.Join(args, " ")

	contents, err := readURL(fmt.Sprintf("http://127.0.0.1:9292/gmilk?package=milkode&query=%s", url.QueryEscape(query))) // @todo package, port, address

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Println(contents)
}

