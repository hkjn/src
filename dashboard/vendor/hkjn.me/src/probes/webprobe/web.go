// Package webprobe implements a HTTP probe.
package webprobe // import "hkjn.me/src/probes/webprobe"

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"hkjn.me/src/prober"
	"hkjn.me/src/probes"
)

const (
	MaxResponseBytes int64 = 10e5 // largest response size accepted
	defaultName            = "WebProber"
)

// WebProber probes a target's HTTP response.
type WebProber struct {
	Target         string // URL to probe
	Method         string // GET, POST, PUT, etc.
	Name           string // name of the prober
	Body           io.Reader
	wantCode       int
	wantInResponse string
	wantHeaders    map[string]string
	alertFn        prober.AlertFn
}

// Name sets the name for the prober.
func Name(name string) func(*WebProber) {
	return func(p *WebProber) {
		p.Name = fmt.Sprintf("%s_%s", defaultName, name)
	}
}

// Body sets the HTTP request body for the prober.
func Body(body io.Reader) func(*WebProber) {
	return func(p *WebProber) {
		p.Body = body
	}
}

// InResponse applies the option that the prober wants given string in the HTTP response.
func InResponse(str string) func(*WebProber) {
	return func(p *WebProber) {
		p.wantInResponse = str
	}
}

// ResponseHeaders adds expected header names and values for the response.
func ResponseHeaders(headers map[string]string) func(*WebProber) {
	return func(p *WebProber) {
		p.wantHeaders = headers
	}
}

// Alert sets a custom alerting function.
//
// If Alert is not called, the probes.SendAlertEmail function is called.
func Alert(fn prober.AlertFn) func(*WebProber) {
	return func(p *WebProber) {
		p.alertFn = fn
	}
}

// New returns a new instance of the web probe with specified options.
func New(target, method string, code int, options ...func(*WebProber)) *prober.Probe {
	return NewWithGeneric(target, method, code, []prober.Option{}, options...)
}

// NewWithGeneric returns a new instance of the web probe with specified options.
//
// NewWithGeneric passes through specified prober.Options, after
// applying the webprobe-specific options.
func NewWithGeneric(target, method string, code int, genericOpts []prober.Option, options ...func(*WebProber)) *prober.Probe {
	name := defaultName
	p := &WebProber{
		Target:      target,
		Name:        name,
		Method:      method,
		wantCode:    code,
		wantHeaders: make(map[string]string),
		alertFn:     probes.SendAlertEmail,
	}
	for _, opt := range options {
		opt(p)
	}
	return prober.NewProbe(p, p.Name, fmt.Sprintf("Probes HTTP response of %s", target), genericOpts...)
}

// Probe verifies that the target's HTTP response is as expected.
func (p WebProber) Probe() prober.Result {
	req, err := http.NewRequest(p.Method, p.Target, p.Body)
	if err != nil {
		return prober.FailedWith(fmt.Errorf("failed to create HTTP request: %v", err))
	}
	// Inform the server that we'd like the connection to be closed once
	// we're done:
	// http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
	req.Header.Set("Connection", "close")
	t := http.Transport{}
	resp, err := t.RoundTrip(req)
	if err != nil {
		return prober.FailedWith(fmt.Errorf("failed to send HTTP request: %v", err))
	}
	defer resp.Body.Close()

	for h, want := range p.wantHeaders {
		got := resp.Header.Get(h)
		if want != got {
			return prober.FailedWith(fmt.Errorf("want header %q=%q in response, got %q", h, want, got))
		}
	}

	if resp.StatusCode != p.wantCode {
		return prober.FailedWith(fmt.Errorf("bad HTTP response status; want %d, got %d", p.wantCode, resp.StatusCode))
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, MaxResponseBytes))
	if err != nil {
		return prober.FailedWith(fmt.Errorf("failed to read HTTP response: %v", err))
	}
	sb := string(body)
	if !strings.Contains(sb, p.wantInResponse) {
		return prober.FailedWith(fmt.Errorf("response doesn't contain %q: \n%v\n", p.wantInResponse, sb))
	}
	return prober.Passed()
}

// Alert calls the prober.AlertFn for the prober.
//
// If no prober.AlertFn was set with the Alert() option,
// probes.SendAlertEmail is used by default.
func (p *WebProber) Alert(name, desc string, badness int, records prober.Records) error {
	return p.alertFn(name, desc, badness, records)
}
