// scrape_torrent fetches HTML from a site and looks for magnet: links.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// The site to connect to.
//
// Alternatives:
// domain = "https://unblockpirate.uk/search/%s/0/7/0"
const domain = "https://proxtpb.art"
const searchPattern = "%s/search/%s/0/7/0"

var dry_run = flag.Bool("dry_run", false, "if specified, don't actually request URLs")

type (
	magnet struct {
		url string
	}
	torrentLink struct {
		url string
	}
)

// getBody fetches url and returns a io.ReadCloser for the content.
func getBody(url string) (io.ReadCloser, error) {
	log.Printf("Connecting to %s..\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return &DummyReadCloser{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &DummyReadCloser{}, errors.New(
			fmt.Sprintf("Failed to retrieve %q: %v", url, resp.Status))
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

// getTorrentLink returns '<a href="/torrent/">' element within a html.Token, or empty string if there's none.
func getTorrentLink(tok html.Token) *torrentLink {
	if tok.Type == html.StartTagToken {
		for _, a := range tok.Attr {
			if a.Key == "href" {
				if len(a.Val) > 8 && a.Val[:8] == "/torrent" {
					log.Printf("found href to /torrent: %v\n", a.Val)
					return &torrentLink{a.Val}
				}
			}
		}
	}
	return nil
}

// getMagnet returns 'magnet:' link within a html.Token.
func getMagnet(tok html.Token) *magnet {
	if tok.Type == html.StartTagToken {
		for _, a := range tok.Attr {
			if a.Key == "href" {
				if len(a.Val) > 8 && a.Val[:8] == "magnet:?" {
					return &magnet{a.Val}
				}
			}
		}
	}
	return nil
}

// getMagnetLinks retrieves up to n "magnet:" links from url.
func getMagnetLinks(url string, n int, dry_run bool) []*magnet {
	var body io.ReadCloser
	var err error
	if dry_run {
		log.Printf("--dry_run\n")
		return getMagnets(&DummyReadCloser{}, url, n)
	} else {
		log.Printf("--nodry_run\n")
		body, err = getBody(url)
		if err != nil {
			log.Fatal(err)
		}
		return getMagnets(body, url, n)
	}
}

// scrapeTorrent returns the magnet link found by following the torrentLink.
func scrapeTorrent(link torrentLink, baseUrl string) (*magnet, error) {
	url := fmt.Sprintf("%s%s", domain, link.url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve %q: %v", link, resp.Status)
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	result := &magnet{}
	i := 0
	for {
		_, tok := tokenizer.Next(), tokenizer.Token()
		if tok.Type == html.ErrorToken {
			if tokenizer.Err().Error() == "EOF" {
				break
			}
			log.Fatalf("Parse error: %v", tokenizer.Err())
		}
		m := getMagnet(tok)
		if m != nil {
			result = m
			i += 1
			break
		}
	}
	return result, nil
}

// getMagnets returns up to n "magnet:" links to be found in the body.
func getMagnets(body io.ReadCloser, url string, n int) []*magnet {
	defer body.Close()
	tokenizer := html.NewTokenizer(body)
	result := make([]*magnet, n)
	i := 0
	for {
		_, tok := tokenizer.Next(), tokenizer.Token()
		if tok.Type == html.ErrorToken {
			if tokenizer.Err().Error() == "EOF" {
				break
			}
			log.Fatalf("Parse error: %v", tokenizer.Err())
		}
		link := getTorrentLink(tok)
		if link != nil {
			m, err := scrapeTorrent(*link, url)
			if err != nil {
				log.Printf("failed to fetch magnet: %v\n", err)
				m = nil
			}
			result[i] = m
			i += 1
			if i >= n {
				break // Found all the magnets we wanted.
			}
		}
	}
	return result
}

func main() {
	flag.Parse()
	searchTerm := flag.Arg(0)
	log.Printf("Search term: %s\n", searchTerm)
	magnets := getMagnetLinks(fmt.Sprintf(searchPattern, domain, searchTerm), 10, *dry_run)
	for _, m := range magnets {
		if m == nil {
			break
		}
		fmt.Printf("%s\n", m.url)
	}
}
