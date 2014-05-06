package main
 
import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/terminal"
	"github.com/monochromegane/the_platinum_searcher/search"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	// "github.com/ongaeshi/gomilk/search"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
)

const version = "0.1.0"

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
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
 
func main() {
	query := strings.Join(os.Args[1:], " ")
	contents, err := readURL(fmt.Sprintf("http://127.0.0.1:9292/gmilk?package=milkode&query=%s", url.QueryEscape(query)))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Println(contents)

	// test
	var opts option.Option

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gomilk"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"

	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if len(args) == 0 && opts.FilesWithRegexp == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	var root = "."
	if len(args) == 2 {
		root = strings.TrimRight(args[1], "\"")
		_, err := os.Lstat(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	opts.Proc = runtime.NumCPU()

	if !terminal.IsTerminal(os.Stdout) {
		opts.NoColor = true
		opts.NoGroup = true
	}

	if opts.Context > 0 {
		opts.Before = opts.Context
		opts.After = opts.Context
	}

	if opts.Context > 0 {
		opts.Before = opts.Context
		opts.After = opts.Context
	}

	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	searcher := search.Searcher{root, pattern, &opts}
	err = searcher.Search()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
