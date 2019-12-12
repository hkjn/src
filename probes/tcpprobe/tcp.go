// Package tcpprobe implements a TCP probe.
package tcpprobe // import "hkjn.me/src/probes/tcpprobe"

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
	defaultName            = "TcpProber"
	defaultTimeout         = time.Second * 5
)

// TcpProber probes a target's HTTP response.
type WebProber struct {
	Target  string        // host to probe
	Name    string        // name of the prober
	Timeout time.Duration // maximum duration to be allowed
	alertFn prober.AlertFn
}

// Name sets the name for the prober.
func Name(name string) func(*TcpProber) {
	return func(p *TcpProber) {
		p.Name = fmt.Sprintf("%s_%s", defaultName, name)
	}
}

// Timeout sets the timeout for the prober.
func Timeout(timeout time.Duration) func(*TcpProber) {
	return func(p *TcpProber) {
		p.Timeout = timeout
	}
}

// Alert sets a custom alerting function.
//
// If Alert is not called, the probes.SendAlertEmail function is called.
func Alert(fn prober.AlertFn) func(*TcpProber) {
	return func(p *TcpProber) {
		p.alertFn = fn
	}
}

// New returns a new tcp probe with specified options.
func New(target, options ...func(*WebProber)) *prober.Probe {
	return NewWithGeneric(target, method, code, []prober.Option{}, options...)
}

// NewWithGeneric returns a new instance of the tcp probe with specified options.
//
// NewWithGeneric passes through specified prober.Options, after
// applying the tcpprobe-specific options.
func NewWithGeneric(target, genericOpts []prober.Option, options ...func(*TcpProber)) *prober.Probe {
	p := &TcpProber{
		Target:  target,
		Name:    defaultName,
		Timeout: defaultTimeout,
		alertFn: probes.SendAlertEmail,
	}
	for _, opt := range options {
		opt(p)
	}
	return prober.NewProbe(
		p,
		p.Name,
		fmt.Sprintf("Probes HTTP response of %s", target),
		genericOpts...,
	)
}

// Probe verifies that the target can be reached within the timeout.
func (p TcpProber) Probe() prober.Result {
	conn, err := net.DialTimeout("tcp", p.Target, p.Timeout)
	if err != nil {
		return prober.FailedWith(fmt.Errorf("failed to dial tcp: %v", err))
	}
	defer conn.Close()
	return prober.Passed()
}

// Alert calls the prober.AlertFn for the prober.
//
// If no prober.AlertFn was set with the Alert() option,
// probes.SendAlertEmail is used by default.
func (p *WebProber) Alert(name, desc string, badness int, records prober.Records) error {
	return p.alertFn(name, desc, badness, records)
}
