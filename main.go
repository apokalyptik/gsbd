package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"code.google.com/p/log4go"

	"github.com/rjohnsondev/go-safe-browsing-api"
)

var apiKey string
var dataDir string
var pidFile string
var verbosity int
var httpHost string
var httpPort int
var gsb *safebrowsing.SafeBrowsing
var logger log4go.Logger

func init() {
	flag.StringVar(&apiKey, "key", "", "Google Safe Browsing api key")
	flag.StringVar(&dataDir, "dir", "/var/lib/safebrowsing/", "path to store GSB data")
	flag.StringVar(&pidFile, "pid", "/var/run/gsbd.pid", "path to PID file")
	flag.StringVar(&httpHost, "listen", "0.0.0.0", "address on which to listen for HTTP requests")
	flag.IntVar(&httpPort, "port", 8888, "port on which to listen for HTTP requests")
	flag.IntVar(&verbosity, "verbose", 0, "0: Critical+, 1:Warning+, 2: Info+, 3: Debug")
}

func main() {
	var err error
	flag.Parse()

	if apiKey == "" {
		log.Fatal("Please profide an API key")
	}

	if pidFile != "" {
		handlePidFile()
	}

	if err = os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error initializing %s: %s", dataDir, err.Error())
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	switch verbosity {
	case 0:
		logger = log4go.NewDefaultLogger(log4go.CRITICAL)
	case 1:
		logger = log4go.NewDefaultLogger(log4go.WARNING)
	case 2:
		logger = log4go.NewDefaultLogger(log4go.INFO)
	case 3:
		logger = log4go.NewDefaultLogger(log4go.DEBUG)
	}
	safebrowsing.Logger = logger
	if gsb, err = safebrowsing.NewSafeBrowsing(apiKey, dataDir); err != nil {
		log.Fatalf("Error initializing GSB api client: %s", err.Error())
	}

	if httpHost != "" && httpPort != 0 {
		go handleHTTP()
	}
	// in theory we may want to add other transports here at some point
	// like raw dns, or zeromq, etc

	var wait chan struct{}
	<-wait
}
