package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/terminal"
	"github.com/ongaeshi/gomilk/search"
	"github.com/ongaeshi/gomilk/search/option"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const version = "0.1.0"

var opts option.Option

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gomilk"
	parser.Usage = "[OPTIONS] PATTERN1 [PATTERN2 ..]"

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

	if len(args) >= 2 {
		fmt.Println("AND search does not currently support")
		os.Exit(1)
	}

	var root = "."

	if opts.Directory != "" {
		root = strings.TrimRight(opts.Directory, "\"")
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

	if opts.Update {
		prevDir, _ := filepath.Abs(".")
		os.Chdir(root)

		command := "milk"

		if runtime.GOOS == "windows" {
			command = "milk.bat"
		}

		cmd := exec.Command(command, "update")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}

		os.Chdir(prevDir)
	}

	searcher := search.Searcher{root, pattern, &opts}
	err = searcher.Search()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
