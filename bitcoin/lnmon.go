// lnmon.go is a simple wrapper around c-lightning's lightning-cli for monitoring state.
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"expvar"
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
	// getInfoResponse is the format of the getinfo response from lightning-cli.
	getInfoResponse struct {
		NodeId      string  `json:"id"`
		Port        int     `json:"port"`
		Address     address `json:"address"`
		Version     string  `json:"version"`
		Blockheight int     `json:"blockheight"`
	}

	// channel describes an individual channel.
	channel struct {
		State                     string `json:"state"`
		Owner                     string `json:"owner"`
		ShortChannelId            string `json:"short_channel_id"`
		FundingTxId               string `json:"funding_txid"`
		MsatoshiToUs              int64  `json:"msatoshi_to_us"`
		MsatoshiTotal             int64  `json:"msatoshi_total"`
		DustLimitSatoshis         int64  `json:"dust_limit_satoshis"`
		MaxHtlcValueInFlightMsats int64  `json:"max_htlc_value_in_flight_msats"`
		ChannelReserveSatoshis    int64  `json:"channel_reserve_satoshis"`
		HtlcMinimumMsat           int64  `json:"htlc_minimum_msat"`
		ToSelfDelay               int64  `json:"to_self_delay"`
		MaxAcceptedHtlcs          int64  `json:"max_accepted_htlcs"`
	}
	// channels describes several channel structures.
	channels []channel
	// peers describes several peers.
	peers []peer
	// peer describes a single peer.
	peer struct {
		PeerId    string   `json:"id"`
		Connected bool     `json:"connected"`
		Netaddr   []string `json:"netaddr"`
		Channels  channels `json:"channels"`
	}
	// listPeersResponse is the format of the listpeers response from lightning-cli.
	listPeersResponse struct {
		Peers peers `json:"peers"`
	}
	// node describes a single node.
	node struct {
		NodeId        string        `json:"nodeid"`
		Alias         string        `json:"alias"`
		Color         string        `json:"color"`
		LastTimestamp int64         `json:"last_timestamp"`
		Addresses     []addressInfo `json:"addresses"`
	}
	// nodes describes several nodes.
	nodes []node
	// listNodesResponse si the format of the listnodes response from lightning-cli.
	listNodesResponse struct {
		Nodes nodes `json:"nodes"`
	}
	lightningdState struct {
		pid   int
		args  []string
		Alias string
		Info  getInfoResponse
		Peers listPeersResponse
		Nodes listNodesResponse
	}
	state struct {
		// aliases maps LN node ids to their human-readable aliases
		aliases    map[string]string
		Bitcoind   bitcoindState
		Lightningd lightningdState
	}
)

// TODO: eliminate global variable
var (
	allState state
	counts   = expvar.NewMap("counters")
)

// getFile returns the contents of the specified file.
func getFile(f string) ([]byte, error) {
	// Asset is defined in bindata.go.
	return Asset(f)
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
	// Note: there's several more states, we just order the ones  we care about here.
	statePrio := map[string]int{
		"CHANNELD_NORMAL":          0,
		"CHANNELD_AWAITING_LOCKIN": 1,
	}
	getPrio := func(s string) int {
		prio, exists := statePrio[s]
		if !exists {
			// Unknown state sorts after known ones.
			return 10
		}
		return prio
	}
	return getPrio(cs[i].State) < getPrio(cs[j].State)
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
			log.Printf("Command %q exited with non-zero status: %v, stderr=%s\n", fmt.Sprintf("%s %s", cmd, strings.Join(arg, " ")), err, stderr.String())
			return stderr.String(), nil
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
	btcState, err := execCmd("pgrep", "-a", "bitcoind")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(btcState, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success, seems like there's still >= 1 parts..
	if len(parts) < 1 || len(parts[0]) == 0 {
		return &bitcoindState{}, nil
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
	lightningState, err := execCmd("pgrep", "-a", "lightningd")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(lightningState, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success.
	if len(parts) < 1 || len(parts[0]) == 0 {
		return &lightningdState{}, nil
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	s := lightningdState{
		pid:  pid,
		args: []string{},
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

	peers, err := c.ListPeers()
	if err != nil {
		return nil, err
	}
	s.Peers = *peers
	sort.Sort(sort.Reverse(s.Peers.Peers))
	// log.Printf("lightningd listpeers response: %+v\n", peers)

	nodes, err := c.ListNodes()
	if err != nil {
		return nil, err
	}
	s.Nodes = *nodes
	// log.Printf("lightningd listnodes response: %+v\n", nodes)
	return &s, nil
}

func refresh() {
	allState.aliases = map[string]string{}
	bitcoindRunning := expvar.NewInt("bitcoind.running")
	lightningdRunning := expvar.NewInt("lightningd.running")
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
			for _, node := range allState.Lightningd.Nodes.Nodes {
				_, exists := allState.aliases[node.NodeId]
				if !exists {
					log.Printf("Learned alias %q for node %q\n", node.Alias, node.NodeId)
					allState.aliases[node.NodeId] = node.Alias
				}
				if node.NodeId == allState.Lightningd.Info.NodeId {
					allState.Lightningd.Alias = node.Alias
				}
			}
		}
		if allState.Lightningd.IsRunning() {
			lightningdRunning.Set(1)
		} else {
			lightningdRunning.Set(0)
		}

		// lightning-cli getchannels
		// lightning-cli listfunds
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
	}
	tmpl, err := template.New("index").Parse(string(s))
	if err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Fatalf("Failed to parse .tmpl: %v\n", err)
	}

	if err := tmpl.Execute(w, allState); err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Fatalf("Failed to execute template: %v\n", err)
	}
}

func main() {
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
		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(s.ListenAndServeTLS("", ""))
	} else {
		fmt.Printf("Serving plaintext HTTP on %s..\n", addr)
		log.Fatal(s.ListenAndServe())
	}
}
