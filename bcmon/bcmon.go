// bcmon.go is a tool for pulling out and serving up data from bitcoin-cli for monitoring.
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
)

type (
	cli struct {
		callCounters map[string]prometheus.Counter
	}
	// peerInfo describes one peer from the bitcoin-cli getpeerinfo response.
	peerInfo struct {
		PeerId          string           `json:"id"`
		Addr            string           `json:"addr"`
		AddrLocal       string           `json:"addrlocal"`
		AddrBind        string           `json:"addrbind"`
		Services        string           `json:"services"`
		RelayTxes       bool             `json:"relaytxes"`
		LastSend        int64            `json:"lastsend"`
		LastRecv        int64            `json:"lastrecv"`
		BytesSent       int64            `json:"bytessent"`
		BytesRecv       int64            `json:"bytesrecv"`
		ConnTime        int64            `json:"conntime"`
		TimeOffset      int64            `json:"timeoffset"`
		PingTimeMs      int64            `json:"pingtime"`
		MinPingMs       int64            `json:"minping"`
		Version         int64            `json:"version"`
		Subver          string           `json:"subver"`
		Inbound         bool             `json:"inbound"`
		AddNode         bool             `json:"addnode"`
		StartingHeight  int64            `json:"startingheight"`
		BanScore        int64            `json:"banscore"`
		SyncedHeaders   int64            `json:"synced_headers"`
		SyncedBlocks    int64            `json:"synced_blocks"`
		Inflight        []string         `json:"inflight"`
		Whitelisted     bool             `json:"whitelisted"`
		BytesSentPerMsg map[string]int64 `json:"bytessent_per_msg"`
		BytesRecvPerMsg map[string]int64 `json:"bytesrecv_per_msg"`
	}
	// state describes the last known state.
	state struct {
		// MonVersion is the version of bcmon.
		MonVersion string
		pid        int
		args       []string
		peerInfo   []peerInfo
		gauges     map[string]prometheus.Gauge
	}
	indexHandler struct {
		tmpl template.Template
		s    *state
	}
)

const counterPrefix = "bitcoind"

var (
	// TODO: eliminate global variable
	allState   state
	addr       = os.Getenv("BCMON_ADDR")
	hostname   = os.Getenv("BCMON_HOSTNAME")
	httpPrefix = os.Getenv("BCMON_HTTP_PREFIX")
)

// getFile returns the contents of the specified file.
func getFile(f string) ([]byte, error) {
	// Asset is defined in bindata.go.
	return Asset(f)
}

func (s state) isRunning() bool {
	return s.pid != 0
}

// execCmd executes specified command with arguments and returns the output.
func execCmd(cmd string, arg ...string) (string, error) {
	c := exec.Command(cmd, arg...)
	out := bytes.Buffer{}
	stderr := bytes.Buffer{}
	c.Stdout = &out
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Command exited with non-zero status.
			errstring := stderr.String()
			errmsg := fmt.Sprintf("Command %q exited with non-zero status: %v", fmt.Sprintf("%s %s", cmd, strings.Join(arg, " ")), err)
			if errstring != "" {
				errmsg += fmt.Sprintf(", stderr=%q", stderr.String())
			}
			return "", fmt.Errorf(errmsg)
		}
		return "", err
	}
	return out.String(), nil
}

// newCli returns a new cli.
func newCli() *cli {
	cliCalls := []string{
		"getpeerinfo",
	}
	counters := map[string]prometheus.Counter{}
	for _, call := range cliCalls {
		c := prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: counterPrefix,
				Name:      call + "calls_total",
				Help:      fmt.Sprintf("Number of calls to %q CLI.", call),
			},
		)
		counters[call] = c
		// prometheus.MustRegister(c)
	}
	return &cli{callCounters: counters}
}

func (c cli) exec(cmd string) (string, error) {
	return execCmd("lightning-cli", cmd)
}

func (c cli) incCounter(call string) {
}

// GetPeerInfo returns the getinfo response.
func (c cli) GetPeerInfo() (*[]peerInfo, error) {
	c.incCounter("getpeerinfo")

	respstring, err := c.exec("getpeerinfo")
	if err != nil {
		return nil, err
	}
	resp := []peerInfo{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// update refreshes the state.
func (s *state) update() error {
	ps, err := execCmd("pgrep", "-a", "bitcoind")
	if err != nil {
		return err
	}
	parts := strings.Split(ps, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success, seems like there's still >= 1 parts..
	if len(parts) < 1 || len(parts[0]) == 0 {
		return fmt.Errorf("failed to parse bitcoind status: %v", ps)
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	s.pid = pid
	for _, arg := range parts[1:] {
		s.args = append(s.args, arg)
	}
	return nil
}

// reset forgets all lightningd state.
func (s *state) reset() {
	s.pid = 0
	s.args = []string{}
	s.peerInfo = []peerInfo{}
}

// newIndexHandler returns a new http handler for the index page.
func newIndexHandler(s *state) (*indexHandler, error) {
	b, err := getFile("index.tmpl")
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("index").Parse(string(b))
	if err != nil {
		return nil, err
	}
	return &indexHandler{tmpl: *tmpl, s: s}, nil
}

// ServeHTTP serves the index page.
func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v] HTTP %s %s\n", r.RemoteAddr, r.Method, r.URL)

	data := struct {
		IsRunning     bool
		DashboardLink string
	}{
		IsRunning: allState.isRunning(),
	}

	if os.Getenv("BCMON_LINK_DASHBOARD") != "" {
		data.DashboardLink = os.Getenv("BCMON_LINK_DASHBOARD")
	}
	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Printf("Failed to execute template: %v\n", err)
		return
	}
}

// registerMetrics registers the Prometheus monitoring metrics.
func (h indexHandler) registerMetrics() {
	for _, m := range h.s.gauges {
		prometheus.MustRegister(m)
	}
}

func refresh(s *state) {
	// TODO: Maybe don't assume that bitcoind always is in "pid" namespace..
	namespace := "pid"
	registeredBitcoin := false
	for {
		if err := s.update(); err != nil {
			// TODO: increment counter here, so we can alert on possible lightningd crashes.
			log.Printf("Failed to get state: %v\n", err)
			s.reset()
		}
		if s.isRunning() {
			if !registeredBitcoin {
				lc := prometheus.NewProcessCollector(s.pid, namespace)
				prometheus.MustRegister(lc)
				registeredBitcoin = true
				log.Printf("Registered ProcessCollector for bitcoind pid %d in namespace %q\n", s.pid, namespace)
			}
			s.gauges["running"].Set(1)
		} else {
			s.gauges["running"].Set(0)
		}

		time.Sleep(time.Minute)
	}
}

// newRouter returns a new http router.
func newRouter(ih http.Handler, prefix string) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", ih).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	if prefix != "" {
		log.Printf("Serving resources with prefix %q..\n", prefix)
		sr := r.PathPrefix(prefix).Subrouter()
		sr.Handle("/", ih).Methods("GET")
		sr.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}
	return r
}

// newState returns a new state.
func newState() *state {
	return &state{
		gauges: map[string]prometheus.Gauge{
			"running": prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: counterPrefix,
				Name:      "running",
				Help:      "Whether bitcoind process is running (1) or not (0).",
			}),
		},
	}
}

func main() {
	s := newState()
	ih, err := newIndexHandler(s)
	if err != nil {
		log.Fatalf("Failed to create http handlers: %v\n", err)
	}
	ih.registerMetrics()
	router := newRouter(ih, httpPrefix)
	http.Handle("/", router)
	go refresh(s)

	if addr == "" {
		addr = ":9740"
	}
	server := &http.Server{Addr: addr}
	if addr == ":443" {
		fmt.Printf("Serving TLS at %q as %q..\n", addr, hostname)
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/etc/secrets/acme/"),
			HostPolicy: autocert.HostWhitelist(hostname),
		}
		// Configure extra tcp/80 server for http-01 challenge:
		// https://godoc.org/golang.org/x/crypto/acme/autocert#Manager.HTTPHandler
		httpServer := &http.Server{
			Handler: m.HTTPHandler(nil),
			Addr:    ":80",
		}
		go httpServer.ListenAndServe()

		server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		fmt.Printf("Serving plaintext HTTP on %s..\n", addr)
		log.Fatal(server.ListenAndServe())
	}
}
