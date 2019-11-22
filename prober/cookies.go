package prober

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

// RestrictedCookies is a http.CookieJar implementation that only
// stores cookies if the name, domain, and path matches predefined
// expectations.
type RestrictedCookies struct {
	cookies map[string]*http.Cookie // map of names to accept to stored cookies
	domain  string                  // domain to match for
	path    string                  // expected path on domain
}

// SetCookies stores the specified cookies in the jar.
//
// Cookies will only be stored if the domain, path and cookie names
// are what we wanted.
func (cj *RestrictedCookies) SetCookies(u *url.URL, cookies []*http.Cookie) {
	log.Printf("SetCookies(%v, %v)\n", u, cookies)
	if !strings.HasPrefix(u.String(), cj.domain) {
		return
	}
	for _, c := range cookies {
		if c.Path != cj.path {
			continue
		}
		_, ok := cj.cookies[c.Name]
		if !ok {
			continue
		}
		log.Printf("cookie match, storing %v\n", c)
		cj.cookies[c.Name] = c
	}
}

// Cookies returns the cookies for the given domain.
func (cj *RestrictedCookies) Cookies(u *url.URL) []*http.Cookie {
	log.Printf("Cookies(%v) (has %v)\n", u, cj.cookies)
	if !strings.HasPrefix(u.String(), cj.domain) {
		return []*http.Cookie{}
	}
	r := []*http.Cookie{}
	for _, c := range cj.cookies {
		if c != nil {
			r = append(r, c)
		}
	}
	return r
}
