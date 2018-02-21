// Fileserver is a minimal service for to serving contents from the file system over HTTP.

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	filesDir := os.Getenv("FILESERVER_DIR")
	if filesDir == "" {
		filesDir = "/var/www"
	}
	fs := http.FileServer(http.Dir(filesDir))
	http.Handle("/", fs)

	addr := os.Getenv("FILESERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	s := &http.Server{
		Addr: addr,
	}
	if addr == ":443" {
		host := os.Getenv("FILESERVER_HOST")
		if host == "" {
			log.Fatalf("FILESERVER_HOST must be set to serve TLS.")
		}
		fmt.Printf("Serving TLS as %s..\n", host)
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/etc/secrets/acme/"),
			HostPolicy: autocert.HostWhitelist(host),
		}
		// Configure extra tcp/80 server for http-01 challenge:
		// https://godoc.org/golang.org/x/crypto/acme/autocert#Manager.HTTPHandler
		httpServer := &http.Server{
			Handler: m.HTTPHandler(nil),
			Addr:    ":80",
		}
		go httpServer.ListenAndServe()

		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(s.ListenAndServeTLS("", ""))
	} else {
		fmt.Printf("Serving plaintext HTTP on %s..\n", addr)
		log.Fatal(s.ListenAndServe())
	}
}
