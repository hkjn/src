package secretservice

import (
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/crypto/acme/autocert"
)

const (
	// saltFile is the path to the secretservice salt file.
	saltFile = "/etc/secrets/secretservice/salt"
	// seedFile is the path to the secretservice seed file.
	seedFile = "/etc/secrets/secretservice/seed"
	tpl      = `
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

type Config struct {
	FilesDir string
	Addr     string
	Domain   string
}

// lookup returns the unique prefix to use for given key.
func lookup(key, salt string) string {
	val := fmt.Sprintf("%s|%s\n", key, salt)
	digest := sha512.Sum512([]byte(val))
	return fmt.Sprintf("%x", digest)
}

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

// Serve serves secrets.
func Serve() error {
	var conf Config
	if err := envconfig.Process("SECRETSERVICE", &conf); err != nil {
		return fmt.Errorf("envconfig: %v\n", err)
	}
	if conf.Domain == "" {
		return fmt.Errorf("no SECRETSERVICE_DOMAIN\n")
	}
	if conf.FilesDir == "" {
		conf.FilesDir = "/var/www/secretservice"
	}
	if conf.Addr == "" {
		conf.Addr = ":443"
	}
	hash, err := GetHash()
	if err != nil {
		return err
	}
	baseUri := fmt.Sprintf("/%s/", hash)
	filesUri := fmt.Sprintf("%s%s/", baseUri, "files")
	log.Printf("Secretservice serving from base %q, bound to %q\n", conf.FilesDir, filesUri)
	http.HandleFunc(baseUri, indexHandler)

	fs := http.FileServer(http.Dir(conf.FilesDir))
	fs = http.StripPrefix(filesUri, fs)
	http.Handle(filesUri, fs)
	s := &http.Server{Addr: conf.Addr}
	if conf.Addr == ":443" {
		log.Println("Serving TLS..")
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/etc/secrets/acme/"),
			HostPolicy: autocert.HostWhitelist(conf.Domain),
		}
		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		return s.ListenAndServeTLS("", "")
	} else {
		log.Printf("Serving plaintext HTTP on %s..\n", conf.Addr)
		return s.ListenAndServe()
	}
}

// getSecretServiceHash returns the secret service hash read from files.
func GetHash() (string, error) {
	salt, err := ioutil.ReadFile(saltFile)
	if err != nil {
		return "", err
	}
	seed, err := ioutil.ReadFile(seedFile)
	if err != nil {
		return "", err
	}
	seed = []byte(strings.TrimSpace(string(seed)))
	salt = []byte(strings.TrimSpace(string(salt)))
	val := fmt.Sprintf("%s|%s\n", seed, salt)
	digest := sha512.Sum512([]byte(val))
	return fmt.Sprintf("%x", digest), nil
}
