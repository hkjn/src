// scrape_torrent fetches HTML from a site and looks for magnet: links.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"code.google.com/p/go.net/html"
)

// The site to connect to.
const site = "http://thepiratebay.se/search/%s/0/7/0"

var dry_run = flag.Bool("dry_run", false, "if specified, don't actually request URLs")

type magnet struct {
	url string
}

// getBody fetches url and returns a io.ReadCloser for the content.
func getBody(url string) (io.ReadCloser, error) {
	log.Printf("Connecting to %s..\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return &DummyReadCloser{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &DummyReadCloser{}, errors.New(
			fmt.Sprintf("Failed to retrieve %s: %d", url, resp.Status))
	}
	return resp.Body, nil
}

// DummyReadCloser implements io.ReadCloser by returning static data.
type DummyReadCloser struct {
	consumed bool
}

// Read returns a static []byte the first time, then io.EOF.
func (d *DummyReadCloser) Read(b []byte) (n int, err error) {
	if d.consumed {
		return 0, io.EOF
	}
	data := []byte(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Fake</title>
<body>
<div>
	<a href="magnet:?xt=urn:btih:db4c7f5d&dn=something%5D&tr=udp%3A%2F%2Ftracker.foo.com%3A80&tr=udp%3A%2F%2Ftracker.bar.com%3A80&tr=udp%3A%2F%2Ftracker.qux.it%3A6969&tr=udp%3A%2F%2Ftfoobar.de%3A80&tr=udp%3A%2F%2Fopen.foobar.com%3A1337" title="Download this torrent using magnet"><img src="/static/img/icon-magnet.gif" alt="Magnet link" /></a>
</div>
</body>`)
	copy(b, data)
	// No more data from us.
	d.consumed = true
	return len(data), nil
}

// Close doesn't do anything.
func (DummyReadCloser) Close() error { return nil }

// getMagnet returns 'magnet:' link within a html.Token, or empty string.
func getMagnet(tok html.Token) magnet {
	if tok.Type == html.StartTagToken {
		for _, a := range tok.Attr {
			if a.Key == "href" {
				if len(a.Val) > 8 && a.Val[:8] == "magnet:?" {
					return magnet{a.Val}
				}
			}
		}
	}
	return magnet{}
}

// getMagnetLinks retrieves up to n "magnet:" links from url.
func getMagnetLinks(url string, n int, dry_run bool) []magnet {
	var rcloser io.ReadCloser
	var err error
	if dry_run {
		log.Printf("--dry_run\n")
		return getMagnets(&DummyReadCloser{}, n)
	} else {
		log.Printf("--nodry_run\n")
		rcloser, err = getBody(url)
		if err != nil {
			log.Fatal(err)
		}
		return getMagnets(rcloser, n)
	}
}

// getMagnets returns up to n "magnet:" links to be found in the ReadCloser.
func getMagnets(rcloser io.ReadCloser, n int) (result []magnet) {
	z := html.NewTokenizer(rcloser)
	result = make([]magnet, n)
	i := 0
	for {
		_, tok := z.Next(), z.Token()
		if tok.Type == html.ErrorToken {
			if z.Err().Error() == "EOF" {
				break
			}
			log.Fatalf("Parse error: %v", z.Err())
		}
		m := getMagnet(tok)
		if (m != magnet{}) {
			result[i] = m
			i += 1
			if i >= n {
				break // Found all the magnets we wanted.
			}
		}
	}
	rcloser.Close()
	return
}

func main() {
	flag.Parse()
	search_term := flag.Arg(0)
	log.Printf("Search term: %s\n", search_term)
	magnets := getMagnetLinks(fmt.Sprintf(site, search_term), 10, *dry_run)
	for _, m := range magnets {
		if (m == magnet{}) {
			break
		}
		fmt.Printf("%s\n", m.url)
		// TODO: Also print URL-decoded 'dn' param to show title?
		// TODO: Also print # seeds / leaches?
	}
}
