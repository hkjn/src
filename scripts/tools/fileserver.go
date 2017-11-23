// Fileserver is a simple HTTP fileserver.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
)

var fileName = flag.String("name", "", "File to serve")

func init() {
	flag.Parse()
}

func main() {
	http.HandleFunc("/file", serveFile)
	log.Fatalln(http.ListenAndServe(":1234", nil))
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	glog.Infof("serving /file with %q..\n", *fileName)
	http.ServeFile(w, r, *fileName)
}
