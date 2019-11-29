// Package dashboard implements a web dashboard for monitoring.
package dashboard // import "hkjn.me/src/dashboard"

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"hkjn.me/src/config"
	"hkjn.me/src/dashboard/gen"
	"hkjn.me/src/prober"
	"hkjn.me/src/probes"
)

var (
	emailTemplate = `{{define "email"}}
The probe <a href="http://j.mp/hkjndash#{{.Name}}">{{.Name}}</a> failed enough that this alert fired, as the arbitrary metric of 'badness' is {{.Badness}}, which we can all agree is a big number.<br/>
The description of the probe is: &ldquo;{{.Desc}}&rdquo;<br/>
Failure details follow:<br/>
{{range $r := .Records.RecentFailures}}
  <h2>{{$r.Timestamp}} ({{$r.Ago}})</h2>
  <p>{{$r.Result.Info}}</p>
{{end}}
{{end}}`
	probecfg = struct {
		WebProbes []struct {
			Target, Want, Name string
			WantStatus         int
		}
		VarsProbes []struct {
			Target, Name, Key, WantValue string
		}
		DnsProbes []struct {
			Target  string
			Records struct {
				Cname string
				A     []string
				Mx    []struct {
					Host string
					Pref uint16
				}
				Ns  []string
				Txt []string
			}
		}
	}{}
	loadConfigOnce = sync.Once{}
)

type Config struct {
	Debug            bool `default:"true"`
	BindAddr         string
	AllowedGoogleIds []string
	SendgridToken    string
	EmailSender      string
	EmailRecipient   string
}

// setProbeCfg sets the config values.
func setProbesCfg(conf Config, emailTemplate string) error {
	if conf.Debug {
		log.Printf("Starting in debug mode..")
		return nil
	}
	// TODO(hkjn): Unify probes.Config vs dashboard.Config.
	probes.Config.SendgridToken = conf.SendgridToken
	if conf.SendgridToken == "" {
		return errors.New("no DASHBOARD_SENDGRIDTOKEN specified")
	}

	if conf.EmailSender == "" {
		return errors.New("no DASHBOARD_EMAILSENDER specified")
	}
	if conf.Debug {
		log.Printf(
			"Sending any alert emails from %q to %q\n",
			conf.EmailSender,
			conf.EmailRecipient,
		)
	}
	probes.Config.Alert.Sender = conf.EmailSender
	if conf.EmailRecipient == "" {
		return errors.New("no DASHBOARD_EMAILRECIPIENT specified")
	}
	probes.Config.Alert.Recipient = conf.EmailRecipient
	if emailTemplate == "" {
		return errors.New("no email template")
	}
	probes.Config.Template = emailTemplate

	if conf.Debug {
		log.Printf("These Google+ IDs are allowed access: %q\n", conf.AllowedGoogleIds)
	}
	return nil
}

// getIndexData returns the data for the index page.
func getIndexData(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data := struct {
		Version string
		Links   []struct {
			Name, URL string
		}
		Probes         []*prober.Probe
		ProberDisabled bool
	}{}
	data.Version = gen.Version
	data.Probes = getProbes()
	data.ProberDisabled = *proberDisabled
	return data, nil
}

// Start returns the HTTP routes for the dashboard.
func Start(conf Config) *mux.Router {
	r := func(filename string) ([]byte, error) {
		return ioutil.ReadFile(filename)
	}
	if !conf.Debug {
		r = func(filename string) ([]byte, error) {
			return gen.Asset(filename)
		}
	}
	config.MustLoadNameFrom("probes.yaml", &probecfg, r)

	ps := getProbes()
	log.Printf("Starting %d probes..\n", len(ps))
	for _, p := range ps {
		go p.Run()
	}

	if err := setProbesCfg(conf, emailTemplate); err != nil {
		log.Fatalf("FATAL: Couldn't set probes config: %v\n", err)
	}
	return newRouter(conf.Debug)
}
