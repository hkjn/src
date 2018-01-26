// lnmon.go is a simple wrapper around c-lightning's lightning-cli for monitoring state.
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	// getInfoResponse is the format of the getinfo response from lightning-cli.
	getInfoResponse struct {
		NodeId      string        `json:"id"`
		Port        int           `json:"port"`
		Address     []addressInfo `json:"address"`
		Version     string        `json:"version"`
		Blockheight int           `json:"blockheight"`
	}

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
	peer struct {
		PeerId    string    `json:"id"`
		Connected bool      `json:"connected"`
		Netaddr   []string  `json:"netaddr"`
		Channels  []channel `json:"channels"`
	}
	// listPeersResponse is the format of the listpeers response from lightning-cli.
	listPeersResponse struct {
		Peers []peer `json:"peers"`
	}
	node struct {
		NodeId        string        `json:"nodeid"`
		Alias         string        `json:"alias"`
		Color         string        `json:"color"`
		LastTimestamp int64         `json:"last_timestamp"`
		Addresses     []addressInfo `json:"addresses"`
	}
	// listnNodesResponse si the format of the listnodes response from lightning-cli.
	listNodesResponse struct {
		Nodes []node `json:"nodes"`
	}
	lightningdState struct {
		pid   int
		args  []string
		info  getInfoResponse
		peers listPeersResponse
		nodes listNodesResponse
	}
	state struct {
		// aliases maps LN node ids to their human-readable aliases
		aliases    map[string]string
		bitcoind   bitcoindState
		lightningd lightningdState
	}
)

// TODO: eliminate global variable
var allState state

func (s bitcoindState) String() string {
	if s.pid == 0 {
		return "bitcoindState{not running}"
	} else {
		return fmt.Sprintf("bitcoindState{pid: %d, args: %q}", s.pid, strings.Join(s.args, " "))
	}
}

func (s lightningdState) String() string {
	if s.pid == 0 {
		return "lightningdState{not running}"
	} else {
		return fmt.Sprintf("lightningdState{pid: %d, args: %q}", s.pid, strings.Join(s.args, " "))
	}
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
			log.Printf("Command %q exited with non-zero status: %v, stderr=%s\n", fmt.Sprintf(cmd, strings.Join(arg, " ")), err, stderr.String())
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

func (c cli) GetNodes() (string, error) {
	return "", nil
}

// getBtcState returns the current bitcoind state.
func getBtcState() (*bitcoindState, error) {
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

func getLnState() (*lightningdState, error) {
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
	s.info = *info
	log.Printf("lightningd getinfo response: %+v\n", info)

	peers, err := c.ListPeers()
	if err != nil {
		return nil, err
	}
	s.peers = *peers
	// log.Printf("lightningd listpeers response: %+v\n", peers)

	nodes, err := c.ListNodes()
	if err != nil {
		return nil, err
	}
	s.nodes = *nodes
	// log.Printf("lightningd listnodes response: %+v\n", nodes)
	return &s, nil
}

func refresh() {
	allState.aliases = map[string]string{}
	for {
		btcState, err := getBtcState()
		if err != nil {
			log.Fatalf("Failed to get bitcoind state: %v\n", err)
		}
		allState.bitcoind = *btcState

		lnState, err := getLnState()
		if err != nil {
			log.Fatalf("Failed to get lightningd state: %v\n", err)
		}
		allState.lightningd = *lnState

		// lightning-cli getinfo
		// lightning-cli getchannels

		// lightning-cli listfunds
		// lightning-cli listinvoice
		// lightning-cli listpayments
		time.Sleep(time.Minute)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	s := "<html>"
	s += "<h1>Hello!</h1>"
	s += `<p>This is an experiment running at ln.hkjn.me with Lightning Network and Bitcoin by <a href="https://hkjn.me">hkjn</a>.</p>`
	s += `<h2><code>bitcoind</code> info</h2>`
	if allState.bitcoind.pid == 0 {
		s += `<p><code>bitcoind</code> is <strong>not running</strong>.</p>`
	} else {
		s += `<p><code>bitcoind</code> is <strong>running</strong>.</p>`
	}
	s += `<h2><code>lightningd</code> info</h2>`
	if allState.lightningd.pid == 0 {
		s += `<p><code>lightningd</code> is <strong>not running</strong>.</p>`
	} else {
		s += `<p><code>lightningd</code> is <strong>running</strong>.</p>`
		s += fmt.Sprintf(`<p>Our id is <code>%s</code>.</p>`, allState.lightningd.info.NodeId)
		alias, exists := allState.aliases[allState.lightningd.info.NodeId]
		if exists {
			s += fmt.Sprintf(`<p>Our alias is <code>%s</code>.</p>`, alias)
		}
		s += fmt.Sprintf(`<p>Our address is is <code>%s:%d</code>.</p>`, allState.lightningd.info.Address[0].Address, allState.lightningd.info.Address[0].Port)
		s += fmt.Sprintf(`<p>Our version is is <code>%s</code>.</p>`, allState.lightningd.info.Version)
		s += fmt.Sprintf(`<p>Our blockheight is is <code>%d</code>.</p>`, allState.lightningd.info.Blockheight)

		s += fmt.Sprintf(`<h3>We have %d lightning peers:</h3>`, len(allState.lightningd.peers.Peers))
		s += fmt.Sprintf(`<ul>`)
		for _, peer := range allState.lightningd.peers.Peers {
			alias, exists := allState.aliases[peer.PeerId]
			if exists {
				s += fmt.Sprintf("<li><code>%s</code> (<code>%s</code>) is ", alias, peer.PeerId)
			} else {
				s += fmt.Sprintf(`<li><code>%s</code> is `, peer.PeerId)
			}
			if peer.Connected {
				s += fmt.Sprintf(`<strong>connected</strong> at <code>%s</code>.`, peer.Netaddr[0])
				s += fmt.Sprintf(`<ul>`)
				if len(peer.Channels) > 0 {
					for _, channel := range peer.Channels {
						s += fmt.Sprintf(`<li><code>%s</code>: channel id <code>%s</code>, funding tx id <code>%s</code></li>`, channel.State, channel.ShortChannelId, channel.FundingTxId)
					}
				} else {
					s += fmt.Sprintf(`<li>No channels.</li>`)
				}
				s += fmt.Sprintf(`</ul>`)
			} else {
				s += "not connected."
			}
			s += fmt.Sprintf(`</li>`)
		}
		s += fmt.Sprintf(`</ul>`)

		s += fmt.Sprintf(`<h3>We know of %d lightning nodes:</h3>`, len(allState.lightningd.nodes.Nodes))
		s += fmt.Sprintf(`<ul>`)
		for _, node := range allState.lightningd.nodes.Nodes {
			s += fmt.Sprintf(`<li><code>%s</code>: <code>%s</code></li>`, node.Alias, node.NodeId)
			_, exists := allState.aliases[node.NodeId]
			if !exists {
				log.Printf("Learned alias %q for node %q\n", node.Alias, node.NodeId)
				allState.aliases[node.NodeId] = node.Alias
			}
		}
		s += fmt.Sprintf(`</ul>`)
	}

	s += "</html>"
	fmt.Fprintf(w, s)
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
