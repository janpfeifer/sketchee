package main

import (
	"flag"
	"fmt"
	"net/http"

	// Populate registers itself as Game.Populator, which starts a game.
	"github.com/golang/glog"
)

var flagPort = flag.Int("port", 9200, "Port number where to start the server.")
var flagStaticFilesPath = flag.String("static",
	"",
	"Path to static files, like javascripts libraries. Must be set, for instance: --static=`pwd`")

func checkFlags() bool {
	if *flagStaticFilesPath == "" {
		return false
	}
	return true
}

func main() {
	//profiler.SetFlags()
	flag.Parse()
	if !checkFlags() {
		glog.Fatal("Bad flags")
	}

	// Registers static content (files): Use the dotFileHidingFileSystem only because the `http.Dir()` implementation
	// somehow doesn't seem to work when running in Windows (in WSL).
	//http.Handle("/", http.FileServer(http.Dir(*flagStaticFilesPath)))
	fs := dotFileHidingFileSystem{http.Dir(*flagStaticFilesPath)}
	http.Handle("/", &glogHandler{http.FileServer(fs)})
	glog.V(1).Infof("glogHandler set up.")

	// Loops over serving requests.
	fmt.Printf("start listening at port %d\n", *flagPort)
	fmt.Printf("Static path: %s\n", *flagStaticFilesPath)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *flagPort), nil)
	if err != nil {
		fmt.Printf("ListenAndServe error :" + err.Error())
	}
}

type glogHandler struct {
	handler http.Handler
}

func (gh *glogHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.V(1).Infof("ServeHTTP(request=%q)", req.RequestURI)
	gh.handler.ServeHTTP(w, req)
}
