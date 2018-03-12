// Package dashboard implements a web dashboard for monitoring.
package dashboard // import "hkjn.me/src/dashboard"

import (
	"errors"
	"sync"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"hkjn.me/src/config"
	"hkjn.me/src/googleauth"
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
	Debug            bool
	BindAddr         string
	AllowedGoogleIds []string
	GoogleServiceId  string
	GoogleSecret     string
	SendgridToken    string
	EmailSender      string
	EmailRecipient   string
}

// setProbeCfg sets the config values.
func setProbesCfg(conf Config, emailTemplate string) error {
	if conf.Debug {
		glog.Infoln("Starting in debug mode (no auth)..")
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
	glog.V(1).Infof(
		"Sending any alert emails from %q to %q\n",
		conf.EmailSender,
		conf.EmailRecipient,
	)
	probes.Config.Alert.Sender = conf.EmailSender
	if conf.EmailRecipient == "" {
		return errors.New("no DASHBOARD_EMAILRECIPIENT specified")
	}
	probes.Config.Alert.Recipient = conf.EmailRecipient
	if emailTemplate == "" {
		return errors.New("no email template")
	}
	probes.Config.Template = emailTemplate

	if conf.GoogleServiceId == "" {
		return errors.New("no service ID")
	}
	glog.V(1).Infof("Our Google service ID is %q\n", conf.GoogleServiceId)
	if conf.GoogleSecret == "" {
		return errors.New("no service secret")
	}
	googleauth.BaseURL = getHttpPrefix()
	googleauth.SetCredentials(conf.GoogleServiceId, conf.GoogleSecret)
	if len(conf.AllowedGoogleIds) == 0 {
		return errors.New("no allowed IDs")
	}
	glog.V(1).Infof("These Google+ IDs are allowed access: %q\n", conf.AllowedGoogleIds)
	googleauth.SetGatingFunc(func(id string) bool {
		for _, aid := range conf.AllowedGoogleIds {
			if id == aid {
				return true
			}
		}
		return false
	})
	return nil
}

// Start returns the HTTP routes for the dashboard.
func Start(conf Config) *mux.Router {
	config.MustLoad(&probecfg, config.Name("probes.yaml"))

	ps := getProbes()
	glog.Infof("Starting %d probes..\n", len(ps))
	for _, p := range ps {
		go p.Run()
	}

	if err := setProbesCfg(conf, emailTemplate); err != nil {
		glog.Fatalf("FATAL: Couldn't set probes config: %v\n", err)
	}
	return newRouter(conf.Debug)
}
