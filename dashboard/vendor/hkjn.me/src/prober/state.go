package prober

import (
	"net/http"
	"net/url"
)

// PageState describes values that a probe may use internally to pass
// state around within a single Probe() execution, e.g. keep log-in
// cookies.
type PageState struct {
	Vals      url.Values
	CookieJar http.CookieJar
}

// NewPageState returns a new page state.
//
// Cookies will only be stored if the domain, path and cookie names
// match.
func NewPageState(domain, path string, cookieNames ...string) *PageState {
	ps := &PageState{}
	cookies := map[string]*http.Cookie{}
	for _, n := range cookieNames {
		cookies[n] = nil
	}
	ps.CookieJar = &RestrictedCookies{
		domain:  domain,
		path:    path,
		cookies: cookies,
	}
	return ps
}
