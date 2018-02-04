// lnmon.go is a simple wrapper around bitcoin-cli and lightning-cli for monitoring their state.
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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
)

type (
	cli           struct{}
	bitcoindState struct {
		pid  int
		args []string
	}
	addressInfo struct {
		AddressType string `json:"type"`
		Address     string `json:"address"`
		Port        int    `json:"port"`
	}
	address []addressInfo
	// channelListing describes one channel from listchannels output.
	channelListing struct {
		Source               string `json:"source"`
		Destination          string `json:"destination"`
		ShortChannelId       string `json:"short_channel_id"`
		Flags                int64  `json:"flags"`
		Active               bool   `json:"active"`
		Public               bool   `json:"public"`
		LastUpdate           int64  `json:"last_update"`
		BaseFeeMillisatoshis int64  `json:"base_fee_millisatoshi"`
		FeePerMillionth      int64  `json:"fee_per_millionth"`
		Delay                int64  `json:"delay"`
	}
	// channelListings describes several channels from listchannels output.
	channelListings []channelListing

	// channel describes an individual channel.
	channel struct {
		State                     channelState `json:"state"`
		Owner                     string       `json:"owner"`
		ShortChannelId            string       `json:"short_channel_id"`
		FundingTxId               string       `json:"funding_txid"`
		MsatoshiToUs              int64        `json:"msatoshi_to_us"`
		MsatoshiTotal             int64        `json:"msatoshi_total"`
		DustLimitSatoshis         int64        `json:"dust_limit_satoshis"`
		MaxHtlcValueInFlightMsats int64        `json:"max_htlc_value_in_flight_msats"`
		ChannelReserveSatoshis    int64        `json:"channel_reserve_satoshis"`
		HtlcMinimumMsat           int64        `json:"htlc_minimum_msat"`
		ToSelfDelay               int64        `json:"to_self_delay"`
		MaxAcceptedHtlcs          int64        `json:"max_accepted_htlcs"`
	}
	// alias is an optional human-readable alias for a node.
	alias string
	// channels describes several channel structures.
	channels []channel
	// netaddr describes the network addresses for a peer.
	netaddr []string
	// peer describes a single peer.
	//
	// TODO: finish unifying peer -> node types.
	peer struct {
		PeerId    string   `json:"id"`
		Connected bool     `json:"connected"`
		Netaddr   netaddr  `json:"netaddr"`
		Channels  channels `json:"channels"`
	}
	// peers describes several peers.
	peers []peer
	// node describes a single node.
	node struct {
		// isPeer is true if this node is one of our peers.
		isPeer    bool
		NodeId    string `json:"nodeid"`
		Connected bool   `json:"connected"
		// Netaddr is present for listpeer responses. Should unify with addresses from listnodes..`
		Netaddr netaddr `json:"netaddr"`
		// Channels holds the channels for peers.
		Channels      channels      `json:"channels"`
		Alias         alias         `json:"alias"`
		Color         string        `json:"color"`
		LastTimestamp int64         `json:"last_timestamp"`
		Addresses     []addressInfo `json:"addresses"`
	}
	// nodes describes several nodes.
	nodes map[alias]node
	// output describes an individual output.
	output struct {
		TxId string `json:"txid"`
		// Output is the index of the txo.
		Output int64 `json:"output"`
		Value  int64 `json:"value"`
	}
	// outputs describes several outputs.
	outputs []output

	// getInfoResponse is the format of the getinfo response from lightning-cli.
	getInfoResponse struct {
		NodeId      string  `json:"id"`
		Port        int     `json:"port"`
		Address     address `json:"address"`
		Version     string  `json:"version"`
		Blockheight int     `json:"blockheight"`
	}
	// listChannelsResponse is the format of the listchannels response from lightning-cli.
	listChannelsResponse struct {
		Channels []channelListing `json:"channels"`
	}
	// listFundsResponse is the format of the listfunds response from lightning-cli.
	listFundsResponse struct {
		Outputs outputs `json:"outputs"`
	}
	// listPeersResponse is the format of the listpeers response from lightning-cli.
	listPeersResponse struct {
		Peers peers `json:"peers"`
	}
	// listNodesResponse is the format of the listnodes response from lightning-cli.
	listNodesResponse struct {
		Nodes []node `json:"nodes"`
	}
	// lightningState describes the last known state of the lightningd daemon.
	lightningdState struct {
		pid      int
		args     []string
		Alias    alias
		Info     getInfoResponse
		Peers    peers
		Nodes    nodes
		Channels channelListings
		Outputs  outputs
	}
	channelStateNum int
	channelState    string
	state           struct {
		Bitcoind   bitcoindState
		Lightningd lightningdState
	}
)

// String returns the name of the state, e.g. "OPENINGD".
func (s channelStateNum) String() string {
	if ChanneldNormalState <= s && s <= ClosingdSigexchangeState {
		return string(states[s])

	}
	return fmt.Sprintf("Invalid channelstate %v", s)
}

const (
	// Channel states are enumerated here. Note that states with lower numbers sort before higher.
	ChanneldNormalState channelStateNum = 1 + iota
	ChanneldAwaitingLockinState
	OpeningdState
	OnchaindTheirUnilateralState
	OnchaindOurUnilateralState
	ClosingdSigexchangeState
)

// TODO: eliminate global variable
var (
	states = map[channelStateNum]channelState{
		ChanneldNormalState:          "CHANNELD_NORMAL",
		ChanneldAwaitingLockinState:  "CHANNELD_AWAITING_LOCKIN",
		OpeningdState:                "OPENINGD",
		OnchaindTheirUnilateralState: "ONCHAIND_THEIR_UNILATERAL",
		OnchaindOurUnilateralState:   "ONCHAIND_OUR_UNILATERAL",
		ClosingdSigexchangeState:     "CLOSINGD_SIGEXCHANGE",
	}
	allState        state
	bitcoindRunning = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "bitcoind",
		Name:      "running",
		Help:      "Whether bitcoind process is running (1) or not (0).",
	})
	lightningdRunning = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "lightningd",
		Name:      "running",
		Help:      "Whether lightningd process is running (1) or not (0).",
	})
	availableFunds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "total_funds",
			Help:      "Sum of all funds available for opening channels.",
		},
	)
	numChannels = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "num_channels",
			Help:      "Number of Lightning channels this node knows about.",
		},
	)
	numPeers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "num_peers",
			Help:      "Number of Lightning peers of this node.",
		},
		[]string{"connected"},
	)
	numNodes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "num_nodes",
			Help:      "Number of Lightning nodes known by this node.",
		},
	)
	ourChannels = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "our_channels",
			Help:      "Number of channels per state to and from our node.",
		},
		[]string{"state"},
	)
	channelCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			// TODO: update console to no longer expect old name lightningd_total_channel_capacity_msatoshi
			Name: "channel_capacity_msatoshi",
			Help: "Capacity of channels in millisatoshi by direction and state.",
		},
		[]string{"state"},
	)
	channelToUsBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "channel_to_us_balance_msatoshi",
			Help:      "Balance to us of channels in millisatoshi by state.",
		},
		[]string{"state"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(bitcoindRunning)
	prometheus.MustRegister(lightningdRunning)
	prometheus.MustRegister(availableFunds)
	prometheus.MustRegister(channelCapacity)
	prometheus.MustRegister(channelToUsBalance)
	prometheus.MustRegister(numChannels)
	prometheus.MustRegister(numNodes)
	prometheus.MustRegister(numPeers)
	prometheus.MustRegister(ourChannels)
}

// getFile returns the contents of the specified file.
func getFile(f string) ([]byte, error) {
	// Asset is defined in bindata.go.
	return Asset(f)
}

// String returns a human-readable description of the peer.
func (n netaddr) String() string {
	if len(n) < 1 {
		return "netaddr{}"
	}
	return fmt.Sprintf("%s", n[0])
}

// String returns a human-readable description of the peer.
func (p peer) String() string {
	parts := []string{
		fmt.Sprintf("id: %s", p.PeerId),
		fmt.Sprintf("connected: %v", p.Connected),
	}
	if p.Connected {
		parts = append(parts, fmt.Sprintf("netaddr: %s", p.Netaddr))
	}
	if len(p.Channels) > 0 {
		parts = append(parts, fmt.Sprintf("channels: %s", p.Channels))
	}
	return fmt.Sprintf(
		"peer{%s}",
		strings.Join(parts, ", "),
	)
}

// String returns a human-readable description of the address.
func (addr address) String() string {
	if len(addr) != 1 {
		return fmt.Sprintf("<unsupported address of len %d: %v>", len(addr), addr)
	}
	return fmt.Sprintf("%s:%d", addr[0].Address, addr[0].Port)
}

// String returns a human-readable description of the nodes.
func (ns nodes) String() string {
	return fmt.Sprintf("%d nodes", len(ns))
}

// Implement sort.Interface for channels to sort them in reasonable order.
func (cs channels) Len() int      { return len(cs) }
func (cs channels) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs channels) Less(i, j int) bool {
	return cs[i].State < cs[j].State
}

// Implement sort.Interface for peers to sort them in reasonable order.
func (ps peers) Len() int      { return len(ps) }
func (ps peers) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }
func (ps peers) Less(i, j int) bool {
	if !ps[i].Connected && ps[j].Connected {
		// Unconnected peers are "less" than connected ones.
		return true
	}
	if ps[i].Connected && !ps[j].Connected {
		// Connected peers are never "less" than unconnected ones.
		return false
	}
	if len(ps[i].Channels) < len(ps[j].Channels) {
		// Peers with fewer channels are "less" than ones with more of them.
		return true
	}
	if len(ps[i].Channels) > len(ps[j].Channels) {
		// Peers with more channels are never "less" than ones with fewer of them.
		return false
	}
	if len(ps[i].Channels) > 1 && len(ps[j].Channels) > 1 {
		// If we and the other peer has at least one channel, let us be "less" than
		// our peer if our first channel is "less" than theirs.
		cs := channels{ps[i].Channels[0], ps[j].Channels[0]}
		sort.Sort(cs)
		return cs.Less(0, 1)
	}
	// Tie-breaker: alphabetic ordering of peer id.
	return ps[i].PeerId < ps[j].PeerId
}

// NumConnected returns the number of connected peers.
func (ps peers) NumConnected() int {
	n := 0
	for _, p := range ps {
		if p.Connected {
			n += 1
		}
	}
	return n
}

// String returns a human-readable description of the outputs.
func (outs outputs) String() string {
	return fmt.Sprintf("outputs totalling %v sat", outs.Sum())
}

// Sum returns the total value of all the outputs.
func (outs outputs) Sum() int64 {
	sum := int64(0)
	for _, o := range outs {
		sum += o.Value
	}
	return sum
}

// NumChannelsByState returns a map from channel state to number of channels in that state.
func (ps peers) NumChannelsByState() map[channelState]int64 {
	byState := map[channelState]int64{}
	for _, p := range ps {
		for _, c := range p.Channels {
			byState[c.State] += 1
		}
	}
	return byState
}

// TotalChannelCapacity returns the total capacity of channels by state (e.g CHANNELD_NORMAL).
func (ps peers) TotalChannelCapacity() map[channelState]int64 {
	result := map[channelState]int64{}
	for _, p := range ps {
		for _, c := range p.Channels {
			result[c.State] += c.MsatoshiTotal
		}
	}
	return result
}

// ToUsChannelBalance returns the balance of channels by state (e.g CHANNELD_NORMAL).
func (ps peers) ToUsChannelBalance() map[channelState]int64 {
	result := map[channelState]int64{}
	for _, p := range ps {
		for _, c := range p.Channels {
			result[c.State] += c.MsatoshiToUs
		}
	}
	return result
}

// String returns a human-readable description of the channel listings.
func (cls channelListings) String() string {
	return fmt.Sprintf("%d channels", len(cls))
}

// String returns a human-readable description of the channel listing.
func (cl channelListing) String() string {
	return "not implemented"
}

// String returns a human-readable description of the channels.
func (cs channels) String() string {
	if len(cs) == 0 {
		return "channels{}"
	}
	if len(cs) > 1 {
		// TODO: Find how this is supported by protocol.
		return "<unsupported multiple channels>"
	}
	return cs[0].String()
}

// String returns a human-readable description of the channel.
func (c channel) String() string {
	parts := []string{
		fmt.Sprintf("state: %s", c.State),
	}
	if c.FundingTxId != "" {
		parts = append(parts, fmt.Sprintf("funding_txid: %s", c.FundingTxId))
	}
	return fmt.Sprintf("channel{%s}", strings.Join(parts, ", "))
}

// String returns a human-readable description of the peers.
func (ps peers) String() string {
	return fmt.Sprintf("%d peers", len(ps))
}

// String returns a human-readable description of the bitcoind state.
func (s bitcoindState) String() string {
	if s.pid == 0 {
		return "bitcoindState{not running}"
	} else {
		return fmt.Sprintf("bitcoindState{pid: %d, args: %q}", s.pid, strings.Join(s.args, " "))
	}
}

func (s bitcoindState) IsRunning() bool {
	return s.pid != 0
}

func (s lightningdState) String() string {
	if s.pid == 0 {
		return "lightningdState{not running}"
	} else {
		return fmt.Sprintf("lightningdState{pid: %d, args: %q}", s.pid, strings.Join(s.args, " "))
	}
}

func (s lightningdState) IsRunning() bool {
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

func (c cli) exec(cmd string) (string, error) {
	return execCmd("lightning-cli", cmd)
}

func (c cli) GetInfo() (*getInfoResponse, error) {
	infostring, err := c.exec("getinfo")
	if err != nil {
		return nil, err
	}
	resp := getInfoResponse{}
	if err := json.Unmarshal([]byte(infostring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListChannels returns the lightning-cli response to listchannels.
func (c cli) ListChannels() (*listChannelsResponse, error) {
	respstring, err := c.exec("listchannels")
	if err != nil {
		return nil, err
	}
	resp := listChannelsResponse{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListFunds returns the lightning-cli response to listfunds.
func (c cli) ListFunds() (*listFundsResponse, error) {
	respstring, err := c.exec("listfunds")
	if err != nil {
		return nil, err
	}
	resp := listFundsResponse{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListNodes returns the lightning-cli response to listnodes.
func (c cli) ListNodes() (*listNodesResponse, error) {
	respstring, err := c.exec("listnodes")
	if err != nil {
		return nil, err
	}
	resp := listNodesResponse{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListPeers returns the lightning-cli response to listpeers.
func (c cli) ListPeers() (*listPeersResponse, error) {
	respstring, err := c.exec("listpeers")
	if err != nil {
		return nil, err
	}
	resp := listPeersResponse{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// getBitcoindState returns the current bitcoind state.
func getBitcoindState() (*bitcoindState, error) {
	ps, err := execCmd("pgrep", "-a", "bitcoind")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(ps, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success, seems like there's still >= 1 parts..
	if len(parts) < 1 || len(parts[0]) == 0 {
		return nil, fmt.Errorf("failed to parse bitcoind status: %v", ps)
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	s := bitcoindState{
		pid:  pid,
		args: []string{},
	}
	for _, arg := range parts[1:] {
		s.args = append(s.args, arg)
	}
	return &s, nil
}

// getLightningState returns the current lightningd state.
func getLightningdState() (*lightningdState, error) {
	ps, err := execCmd("pgrep", "-a", "lightningd")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(ps, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success.
	if len(parts) < 1 || len(parts[0]) == 0 {
		return nil, fmt.Errorf("failed to parse lightningd status: %v", ps)
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	s := lightningdState{
		pid:   pid,
		args:  []string{},
		Nodes: map[alias]node{},
	}
	for _, arg := range parts[1:] {
		s.args = append(s.args, arg)
	}
	c := cli{}
	info, err := c.GetInfo()
	if err != nil {
		return nil, err
	}
	s.Info = *info
	log.Printf("lightningd getinfo response: %+v\n", info)

	channels, err := c.ListChannels()
	if err != nil {
		return nil, err
	}
	s.Channels = channels.Channels
	numChannels.Set(float64(len(s.Channels)))
	log.Printf("Learned of %d channels.\n", len(s.Channels))

	funds, err := c.ListFunds()
	if err != nil {
		return nil, err
	}
	s.Outputs = funds.Outputs
	availableFunds.Set(float64(s.Outputs.Sum()))
	log.Printf("Learned of %d %v.\n", len(s.Outputs), s.Outputs)

	peers, err := c.ListPeers()
	if err != nil {
		return nil, err
	}
	s.Peers = peers.Peers
	log.Printf("Learned of %d peers.\n", len(s.Peers))

	// TODO: need to either combine node/peer types, or set peer.alias based on node.alias.
	sort.Sort(sort.Reverse(s.Peers))
	numPeers.With(prometheus.Labels{"connected": "connected"}).Set(float64(s.Peers.NumConnected()))
	numPeers.With(prometheus.Labels{"connected": "unconnected"}).Set(float64(len(s.Peers) - s.Peers.NumConnected()))
	for state, n := range s.Peers.NumChannelsByState() {
		// log.Printf("We have %d channels in state %q\n", n, state)
		ourChannels.With(prometheus.Labels{"state": string(state)}).Set(float64(n))
	}
	totalCap := s.Peers.TotalChannelCapacity()
	for state, cap := range totalCap {
		channelCapacity.With(prometheus.Labels{"state": string(state)}).Set(float64(cap))
	}
	toUsBalance := s.Peers.ToUsChannelBalance()
	for state, balance := range toUsBalance {
		channelToUsBalance.With(prometheus.Labels{"state": string(state)}).Set(float64(balance))
	}
	// log.Printf("lightningd listpeers response: %+v\n", peers)

	iodes, err := c.ListNodes()
	if err != nil {
		return nil, err
	}
	for _, n := range nodes.Nodes {
		s.Nodes[n.Alias] = n
	}
	numNodes.Set(float64(len(s.Nodes)))
	log.Printf("Learned of %d nodes.\n", len(s.Nodes))
	// log.Printf("lightningd listnodes response: %+v\n", nodes)

	return &s, nil
}

func refresh() {
	for {
		btcState, err := getBitcoindState()
		if err != nil {
			log.Printf("Failed to get bitcoind state: %v\n", err)
			allState.Bitcoind = bitcoindState{}
		} else {
			allState.Bitcoind = *btcState
		}
		if allState.Bitcoind.IsRunning() {
			bitcoindRunning.Set(1)
		} else {
			bitcoindRunning.Set(0)
		}

		lnState, err := getLightningdState()
		if err != nil {
			log.Printf("Failed to get lightningd state: %v\n", err)
			allState.Lightningd = lightningdState{}
		} else {
			allState.Lightningd = *lnState
		}
		if allState.Lightningd.IsRunning() {
			lightningdRunning.Set(1)
		} else {
			lightningdRunning.Set(0)
		}

		// lightning-cli listinvoice
		// lightning-cli listpayments
		time.Sleep(time.Minute)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v] HTTP %s %s\n", r.RemoteAddr, r.Method, r.URL)
	if r.Method != "GET" {
		log.Printf("Serving 400 for HTTP %s %q\n", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	if r.URL.Path != "/" {
		log.Printf("Serving 404 for GET %q\n", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Page Not Found")
		return
	}
	// TODO: read and parse .tmpl once on startup
	s, err := getFile("lnmon.tmpl")
	if err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Fatalf("Failed to read lnmon.tmpl: %v\n", err)
		return
	}
	tmpl, err := template.New("index").Parse(string(s))
	if err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Fatalf("Failed to parse .tmpl: %v\n", err)
		return
	}

	if err := tmpl.Execute(w, allState); err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Printf("Failed to execute template: %v\n", err)
		return
	}
}

func main() {
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())

	go refresh()

	http.HandleFunc("/", indexHandler)

	addr := ":80"
	if os.Getenv("LNMON_ADDR") != "" {
		addr = os.Getenv("LNMON_ADDR")
	}
	hostname := "ln.hkjn.me"

	fmt.Printf("Serving TLS at %q as %q..\n", addr, hostname)
	s := &http.Server{
		Addr: addr,
	}
	if addr == ":443" {
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

		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(s.ListenAndServeTLS("", ""))
	} else {
		fmt.Printf("Serving plaintext HTTP on %s..\n", addr)
		log.Fatal(s.ListenAndServe())
	}
}
