// +build !appengine

// main.go holds bits needed when not serving hkjnweb on appengine.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
	"hkjn.me/hkjnweb"
)

var hosts = []string{
	"hkjn.me",
	"www.hkjn.me",
}

// registerStatic registers handler for static files.
func registerStatic(dir string) {
	if dir == "" {
		dir = "static"
	}
	path := fmt.Sprintf("/%s/", dir)
	h := http.StripPrefix(
		path,
		http.FileServer(http.Dir(dir)))
	http.Handle(path, h)
}

func main() {
	flag.Parse()
	hkjnweb.Register(os.Getenv("PROD") != "")
	registerStatic(os.Getenv("STATIC_DIR"))

	addr := os.Getenv("BIND_ADDR")
	if addr == "" {
		addr = ":12345"
	}
	if os.Getenv("SERVE_HTTP") != "" {
		log.Printf("webserver serving HTTP on %s..\n", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatalf("FATAL: http server exited: %v\n", err)
		}
	} else {
		log.Println("Since SERVE_HTTP isn't set, we should serve https")
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/etc/ssl/hkjn.me"),
			HostPolicy: autocert.HostWhitelist(hosts...),
		}
		s := &http.Server{
			Addr:      addr,
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}
		log.Printf("webserver serving HTTPS on %s..\n", addr)
		err := s.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("FATAL: https server exited: %v\n", err)
		}
	}
}
