package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type options struct {
	Threads      int
	Wordlist     string
	Timeout      int
	FilterStatus []int
	FilterSize   []int
	FilterHeader []string
	FilterString []string
	IgnoreDir    []string
	IgnoreExt    []string
	IgnoreFile   []string
	NoRecursion  bool
	UseHTTPS     bool
	CustomHeader []string
	CustomData   string
}

func (o *options) setDefaults() {
	o.Threads = 20
	o.Wordlist = "wordlist.txt"
	o.Timeout = 10
	o.FilterStatus = []int{200, 204, 301, 302, 307, 401, 403}
	o.FilterSize = []int{500, 1000, 1500}
	o.FilterHeader = []string{"Content-Type: text/html"}
	o.FilterString = []string{}
	o.IgnoreDir = []string{"/admin/", "/backup/", "/cgi-bin/"}
	o.IgnoreExt = []string{".png", ".jpg", ".gif", ".js", ".css"}
	o.IgnoreFile = []string{}
	o.NoRecursion = false
	o.UseHTTPS = false
	o.CustomHeader = []string{}
	o.CustomData = ""
}

func parseOptions() *options {
	var o options

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] URL\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&o.Threads, "t", 20, "Number of concurrent threads")
	flag.StringVar(&o.Wordlist, "w", "wordlist.txt", "Path to the wordlist file")
	flag.IntVar(&o.Timeout, "timeout", 10, "Timeout in seconds")
	filterStatus := flag.String("c", "200,204,301,302,307,401,403", "Filter by status code")
	filterSize := flag.String("s", "500,1000,1500", "Filter by content length")
	filterHeader := flag.String("H", "Content-Type: text/html", "Filter by header")
	filterString := flag.String("w", "", "Filter by string")
	ignoreDir := flag.String("i", "/admin/,/backup/,/cgi-bin/", "Ignore directories")
	ignoreExt := flag.String("e", ".png,.jpg,.gif,.js,.css", "Ignore file extensions")
	ignoreFile := flag.String("ignore", "", "Ignore specific files")
	noRecursion := flag.Bool("no-recursion", false, "Disable recursion")
	useHTTPS := flag.Bool("https", false, "Use HTTPS instead of HTTP")
	customHeader := flag.String("header", "", "Custom header")
	customData := flag.String("data", "", "Custom data")
	flag.Parse()

	o.FilterStatus = splitInt(*filterStatus)
	o.FilterSize = splitInt(*filterSize)
	o.FilterHeader = splitString(*filterHeader)
	o.FilterString = splitString(*filterString)
	o.IgnoreDir = splitString(*ignoreDir)
	o.IgnoreExt = splitString(*ignoreExt)
	o.IgnoreFile = splitString(*ignoreFile)
	o.NoRecursion = *noRecursion
	o.UseHTTPS = *useHTTPS
	o.CustomHeader = splitString(*customHeader)
	o.CustomData = *customData

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	return &o
}
