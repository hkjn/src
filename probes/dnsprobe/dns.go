// Package dnsprobe implements a DNS probe.
package dnsprobe // import "hkjn.me/src/probes/dnsprobe"

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"time"

	"hkjn.me/src/prober"
	"hkjn.me/src/probes"
)

// DnsProber probes a target host's DNS records.
type DnsProber struct {
	Target     string // host to probe
	name, desc string // name and description of the prober
	wantMX     mxRecords
	wantA      []string
	wantNS     nsRecords
	wantCNAME  string
	wantTXT    []string
	alertFn    prober.AlertFn
}

// Alert sets a custom alerting function.
//
// If Alert is not called, the probes.SendAlertEmail function is called.
func Alert(fn prober.AlertFn) func(*DnsProber) {
	return func(p *DnsProber) {
		p.alertFn = fn
	}
}

// New returns a new instance of the DNS probe with specified options.
func New(target string, options ...func(*DnsProber)) *prober.Probe {
	return NewWithGeneric(target, []prober.Option{}, options...)
}

// NewWithGeneric returns a new instance of the DNS probe with specified options.
//
// NewWithGeneric passes through specified prober.Options, after
// applying the dnsprobe-specific options.
func NewWithGeneric(target string, genericOpts []prober.Option, options ...func(*DnsProber)) *prober.Probe {
	p := &DnsProber{
		Target:  target,
		alertFn: probes.SendAlertEmail,
	}
	for _, opt := range options {
		opt(p)
	}

	// Set a default name if none was specified.
	if p.name == "" {
		p.name = fmt.Sprintf("DnsProber_%s", target)
	}
	// Set a default desc if none was specified.
	if p.desc == "" {
		p.desc = fmt.Sprintf("Probes DNS records of %s", target)
	}
	// TODO(hkjn): This currently doesn't let prober.Interval be
	// overridden, because even if it's supplied in genericOpts we
	// override it ourselves..
	return prober.NewProbe(p, p.name, p.desc,
		append(genericOpts, prober.Interval(time.Minute*5), prober.FailurePenalty(5))...)
}

// Name sets specified name.
func Name(name string) func(*DnsProber) {
	return func(p *DnsProber) {
		p.name = name
	}
}

// Desc sets specified description.
func Desc(desc string) func(*DnsProber) {
	return func(p *DnsProber) {
		p.desc = desc
	}
}

// MX sets expected MX records.
func MX(mx []*net.MX) func(*DnsProber) {
	return func(p *DnsProber) {
		wantMX := mxRecords(mx)
		sort.Sort(wantMX)
		p.wantMX = wantMX
	}
}

// A sets expected A records.
func A(a []string) func(*DnsProber) {
	return func(p *DnsProber) {
		sort.Strings(a)
		p.wantA = a
	}
}

// NS sets expected NS records.
func NS(ns []*net.NS) func(*DnsProber) {
	return func(p *DnsProber) {
		nsRec := nsRecords(ns)
		sort.Sort(nsRec)
		p.wantNS = nsRec
	}
}

// CNAME sets expected 1 CNAME record.
func CNAME(cname string) func(*DnsProber) {
	return func(p *DnsProber) {
		p.wantCNAME = cname
	}
}

// TXT applies the option that the prober wants specific TXT records.
func TXT(txt []string) func(*DnsProber) {
	return func(p *DnsProber) {
		sort.Strings(txt)
		p.wantTXT = txt
	}
}

// Probe verifies that the target's DNS records are as expected.
func (p *DnsProber) Probe() prober.Result {
	if len(p.wantMX) > 0 {
		log.Printf("Checking %d MX records..\n", len(p.wantMX))
		if err := p.checkMX(); err != nil {
			return prober.FailedWith(err)
		}
	}
	if len(p.wantA) > 0 {
		log.Printf("Checking %d A records..\n", len(p.wantA))
		if err := p.checkA(); err != nil {
			return prober.FailedWith(err)
		}
	}
	if len(p.wantNS) > 0 {
		log.Printf("Checking %d NS records..\n", len(p.wantNS))
		if err := p.checkNS(); err != nil {
			return prober.FailedWith(err)
		}
	}
	if p.wantCNAME != "" {
		log.Printf("Checking CNAME record..\n")
		if err := p.checkCNAME(); err != nil {
			return prober.FailedWith(err)
		}
	}
	if len(p.wantTXT) > 0 {
		log.Printf("Checking %d TXT records..\n", len(p.wantTXT))
		if err := p.checkTXT(); err != nil {
			return prober.FailedWith(err)
		}
	}
	return prober.PassedWith(p.String(), p.Target)
}

// String returns the human-readable description of the prober.
func (p DnsProber) String() string {
	var parts []string
	if len(p.wantMX) > 0 {
		parts = append(parts, fmt.Sprintf("MX=%s", p.wantMX))
	}
	if len(p.wantA) > 0 {
		parts = append(parts, fmt.Sprintf("A=%+v", p.wantA))
	}
	if len(p.wantNS) > 0 {
		parts = append(parts, fmt.Sprintf("NS=%+v", p.wantNS))
	}
	if p.wantCNAME != "" {
		parts = append(parts, fmt.Sprintf("CNAME=%s", p.wantCNAME))
	}
	if len(p.wantTXT) > 0 {
		parts = append(parts, fmt.Sprintf("TXT=%+v", p.wantTXT))
	}
	s := fmt.Sprintf("%s: %s", p.name, p.desc)
	if len(parts) > 0 {
		s += fmt.Sprintf(" (checks %s)", strings.Join(parts, ", "))
	}
	return s
}

// mxRecords is a collection of MX records, implementing sort.Interface.
//
// We need this custom order since the sort order in net.LookupMX
// randomizes records with the same preference value, but we'd like a stable
// ordering so we can compared expected vs actual records.
type mxRecords []*net.MX

func (r mxRecords) Len() int { return len(r) }

func (r mxRecords) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r mxRecords) Less(i, j int) bool {
	if r[i].Pref == r[j].Pref {
		return r[i].Host < r[j].Host
	}
	return r[i].Pref < r[j].Pref
}

// String returns a readable description of the MX records.
func (r mxRecords) String() string {
	s := ""
	for i, r := range r {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%s (%d)", r.Host, r.Pref)
	}
	return s
}

// checkMX verifies that the target has expected MX records.
func (p *DnsProber) checkMX() error {
	mx, err := net.LookupMX(p.Target)
	if err != nil {
		return fmt.Errorf("failed to look up MX records for %s: %v", p.Target, err)
	}
	mxRec := mxRecords(mx)
	if len(mxRec) != len(p.wantMX) {
		return fmt.Errorf("want %d MX records, got %d: %s", len(p.wantMX), len(mxRec), mxRec)
	}
	sort.Sort(mxRec)
	for i, r := range mxRec {
		if !strings.EqualFold(p.wantMX[i].Host, r.Host) {
			return fmt.Errorf("bad host %q for MX record #%d; want %q", r.Host, i, p.wantMX[i].Host)
		}
		if p.wantMX[i].Pref != r.Pref {
			return fmt.Errorf("bad prio %d for MX record #%d; want %d", i, r.Pref, p.wantMX[i].Pref)
		}
	}
	return nil
}

// checkA verifies that the target has expected A records.
func (p *DnsProber) checkA() error {
	addr, err := net.LookupHost(p.Target)
	if err != nil {
		return fmt.Errorf("failed to look up A records for %s: %v", p.Target, err)
	}
	if len(addr) != len(p.wantA) {
		return fmt.Errorf("got %d A records, want %d: %v", len(addr), len(p.wantA), addr)
	}
	sort.Strings(addr)
	for i, a := range addr {
		if p.wantA[i] != a {
			return fmt.Errorf("bad A record %q at #%d; want %q", a, i, p.wantA[i])
		}
	}
	return nil
}

// nsRecords is a collection of NS records, implementing sort.Interface.
type nsRecords []*net.NS

func (r nsRecords) Len() int { return len(r) }

func (r nsRecords) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r nsRecords) Less(i, j int) bool { return r[i].Host < r[j].Host }

// String returns a readable description of the NS records.
func (ns nsRecords) String() string {
	s := ""
	for i, r := range ns {
		if i > 0 {
			s += ", "
		}
		s += r.Host
	}
	return s
}

// checkNS verifies that the target has expected NS records.
func (p *DnsProber) checkNS() error {
	ns, err := net.LookupNS(p.Target)
	if err != nil {
		return fmt.Errorf("failed to look up NS records for %s: %v", p.Target, err)
	}
	nsRec := nsRecords(ns)
	if len(nsRec) != len(p.wantNS) {
		return fmt.Errorf("want %d NS records, got %d: %s", len(p.wantNS), len(nsRec), nsRec)
	}
	sort.Sort(nsRec)
	for i, n := range ns {
		if p.wantNS[i].Host != n.Host {
			return fmt.Errorf("bad NS record %q at #%d; want %q", n, i, p.wantNS[i])
		}
	}
	return nil
}

// checkCNAME verifies that the target has expected CNAME record.
func (p *DnsProber) checkCNAME() error {
	cname, err := net.LookupCNAME(p.Target)
	if err != nil {
		return fmt.Errorf("failed to look up CNAME record for %s: %v", p.Target, err)
	}
	if cname != p.wantCNAME {
		return fmt.Errorf("bad CNAME record %q; want %q", cname, p.wantCNAME)
	}
	return nil
}

// checkTXT verifies that the target has expected TXT records.
func (p *DnsProber) checkTXT() error {
	txt, err := net.LookupTXT(p.Target)
	if err != nil {
		return err
	}
	if len(txt) != len(p.wantTXT) {
		return fmt.Errorf("want %d TXT records, got %d: %s", len(p.wantTXT), len(txt), txt)
	}
	sort.Strings(txt)
	for i, t := range txt {
		if p.wantTXT[i] != t {
			return fmt.Errorf("bad TXT record %q at #%d; want %q", t, i, p.wantTXT[i])
		}
	}
	return nil
}

// Alert calls the prober.AlertFn for the prober.
//
// If no prober.AlertFn was set with the Alert() option,
// probes.SendAlertEmail is used by default.
func (p *DnsProber) Alert(name, desc string, badness int, records prober.Records) error {
	return p.alertFn(name, desc, badness, records)
}
