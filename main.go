package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tpmanc/myhttp/internal"
	"github.com/tpmanc/myhttp/internal/config"
	"github.com/tpmanc/myhttp/internal/hash"
	appHttp "github.com/tpmanc/myhttp/internal/http"
)

func main() {
	var parallelCount int
	var isHelp bool

	flag.Usage = help
	flag.BoolVar(&isHelp, "help", false, "Display this help message")
	flag.IntVar(&parallelCount, "parallel", config.ParallelCountDefault, "Number of parallel requests")
	flag.Parse()

	if isHelp {
		flag.Usage()
		os.Exit(0)
	}

	if parallelCount <= 0 {
		flag.Usage()
		panic("-parallel param must be grater than 0")
	}

	urls := flag.Args()
	if len(urls) == 0 {
		flag.Usage()
		panic("URLs list can't be empty")
	}

	logger := log.Default()

	httpClient := appHttp.New(http.DefaultClient)

	hashService, err := hash.New(hash.NewMD5())
	if err != nil {
		panic(err)
	}

	app, err := internal.NewApp(logger, parallelCount, hashService, httpClient)
	if err != nil {
		panic(err)
	}

	res := app.Process(urls)
	defer app.Stop()

	for _, p := range res {
		fmt.Println(fmt.Sprintf("%s %s", p.URL, p.Hash))
	}
}

// help prints information about usage
func help() {
	fmt.Println(`myhttp:
	Makes HTTP request to the provided URLs and prints result in following format:
		URL_1 hash_string_1
		URL_2 hash_string_2
		...
	Usage:
	    myhttp [option] [url_1 url_2 ...]
	Options:
	    -parallel [int]  Set number of parallel requests, should be grater than 0
	    -help            Display help message
	Examples:
        ./myhttp http://facebook.com http://google.com
        ./myhttp -parallel 5 http://facebook.com http://google.com
		./myhttp -help`)
}
