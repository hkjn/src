// secretservice is a fileserver serving up secrets over HTTPS.
//
package main

import (
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/crypto/acme/autocert"

	"hkjn.me/hkjninfra/secretservice"
)

type Config struct {
	Seed     string
	FilesDir string
	Addr     string
	Domain   string
}

const (
	salt = "planetary location or range, biology, and behaviors"
	tpl  = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>{{.Title}}</title>
</head>
<body>
<h1>{{.Title}}</h1>

<ul>
  <li><a href="files">View</a></li>
</ul>
</body>
</html>`
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.RequestURI, "/")
	if len(parts) != 3 || parts[0] != "" || parts[2] != "" {
		log.Printf("[%q] Requests %q, returning 404\n", r.RemoteAddr, r.RequestURI)
		http.NotFound(w, r)
		return
	}
	log.Printf("[%q] Serving page template for %q\n", r.RemoteAddr, r.RequestURI)
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		http.Error(w, "Oops.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "File server",
	}

	if err = t.Execute(w, data); err != nil {
		http.Error(w, "Oops.", http.StatusInternalServerError)
		return
	}
}

// lookup returns the unique prefix to use for given key.
func lookup(key string) string {
	val := fmt.Sprintf("%s|%s\n", key, salt)
	digest := sha512.Sum512([]byte(val))
	return fmt.Sprintf("%x", digest)
}

func main() {
	var c Config
	if err := envconfig.Process("SECRETSERVICE", &c); err != nil {
		log.Fatalf("envconfig: %v\n", err)
	}
	if c.Domain == "" {
		log.Fatalf("no SECRETSERVICE_DOMAIN\n")
	}
	log.Printf("Using %d character seed\n", len(c.Seed))
	if c.Seed == "" {
		log.Fatalf("no SECRETSERVICE_SEED\n")
	}
	hash := lookup(c.Seed)
	if c.FilesDir == "" {
		c.FilesDir = "/var/www/secretservice"
	}
	if c.Addr == "" {
		c.Addr = ":443"
	}
	baseUri := fmt.Sprintf("/%s/", hash)
	filesUri := fmt.Sprintf("%s%s/", baseUri, "files")
	log.Printf("Serving from base %q, bound to %q\n", c.FilesDir, filesUri)
	http.HandleFunc(baseUri, indexHandler)

	fs := http.StripPrefix(filesUri, http.FileServer(http.Dir(c.FilesDir)))
	http.Handle(filesUri, fs)
	s := &http.Server{Addr: c.Addr}
	if c.Addr == ":443" {
		log.Println("Serving TLS..")
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/etc/secrets/acme/"),
			HostPolicy: autocert.HostWhitelist(c.Domain),
		}
		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(s.ListenAndServeTLS("", ""))
	} else {
		log.Printf("Serving plaintext HTTP on %s..\n", c.Addr)
		log.Fatal(s.ListenAndServe())
	}
}
