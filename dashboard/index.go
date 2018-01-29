package dashboard

import (
	"net/http"
	"os"

	"hkjn.me/src/prober"
)

func getVersion() string {
	v := os.Getenv("DASHBOARD_VERSION")
	if v == "" {
		v = "<unknown DASHBOARD_VERSION>"
	}
	return v
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
	data.Version = getVersion()
	data.Probes = getProbes()
	data.ProberDisabled = *proberDisabled
	return data, nil
}
