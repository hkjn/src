package dashboard

import (
	"flag"
	"net"
	"sort"
	"sync"

	"github.com/golang/glog"

	"hkjn.me/src/prober"
	"hkjn.me/src/probes/dnsprobe"
	"hkjn.me/src/probes/webprobe"
)

// TODO(hkjn): Add support for sending POST requests in webprobe.
var (
	proberDisabled = flag.Bool("no_probes", false, "disables probes")
	allProbes      = prober.Probes{}
	createOnce     = sync.Once{}
)

func getWebProbes() prober.Probes {
	probes := prober.Probes{}
	for _, p := range probecfg.WebProbes {
		probes = append(probes,
			webprobe.New(
				p.Target,
				"GET",
				p.WantStatus,
				webprobe.Name(p.Name),
				webprobe.InResponse(p.Want)))
	}
	return probes
}

func getDnsProbes() prober.Probes {
	probes := prober.Probes{}
	for _, p := range probecfg.DnsProbes {
		mxRecords := []*net.MX{}
		for _, mx := range p.Records.Mx {
			mxRecords = append(mxRecords, &net.MX{
				Host: mx.Host,
				Pref: mx.Pref,
			})
		}
		nsRecords := []*net.NS{}
		for _, ns := range p.Records.Ns {
			nsRecords = append(nsRecords, &net.NS{Host: ns})
		}
		probes = append(
			probes,
			dnsprobe.New(
				p.Target, dnsprobe.MX(mxRecords), dnsprobe.A(p.Records.A),
				dnsprobe.NS(nsRecords), dnsprobe.CNAME(p.Records.Cname), dnsprobe.TXT(p.Records.Txt)))
	}
	return probes
}

// getProbes returns all probes in the dashboard.
func getProbes() prober.Probes {
	createOnce.Do(func() {
		if !flag.Parsed() {
			flag.Parse()
		}
		if *proberDisabled {
			glog.Infof("Probes are disabled with -no_probes\n")
		} else {
			allProbes = append(getDnsProbes(), getWebProbes()...)
		}
	})
	sort.Sort(allProbes)
	return allProbes
}
