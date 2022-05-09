// wasm main is started as soon as the page is loaded.
//
// There are several parameters that can be passed by the URL,
// they are described in warmap/wasm/j/params.go.
package main

import (
	"flag"
	"github.com/golang/glog"
	coreJS "github.com/gowebapi/webapi/core/js" // proxy for syscall/js ?
	"math/rand"
	skJS "sketchee/wasm/js"
	"time"
)

var (
	_ = coreJS.Func{}
)

func main() {
	// Set and parse flags needed by glog package.
	_ = flag.Set("logtostderr", "true")

	// Set -vmodule flag for glog.
	vmodule := skJS.URI.Query().Get(skJS.Q_VMODULE)
	_ = flag.Set("vmodule", vmodule)

	flag.Parse()
	glog.Infof("Sketchee's WebAssembly client started.")
	if glog.V(1) || vmodule != "" {
		glog.Infof("vmodule=%s", vmodule)
	}

	// Trivial randomness source.
	rand.Seed(time.Now().UTC().UnixNano())

	// Connect to server with websocket (TODO).
	//address := connection.DefaultWebServerLocation()
	//glog.V(1).Infof("  Server address: %s\n", address)
	//conn, err := connection.New(context.Background(), address,
	//	func() {
	//		// On connection closing.
	//		glog.Fatalf("Connection (websocket) to server closed, stopping client.")
	//	})
	//if err != nil {
	//	glog.Fatalf("Failed to connect to server, you will have to reload page: %v", err)
	//}

	// Loop indefinitely, leave Go subsystem running.
	select {}
}
