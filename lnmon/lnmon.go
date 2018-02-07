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
	cli struct{}
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
	bitcoindState struct {
		pid      int
		args     []string
		peerInfo []peerInfo
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

	// msatoshi is number of millisatoshis, used in the LN protocol.
	msatoshi int64
	// channel describes an individual channel.
	channel struct {
		State                     channelState `json:"state"`
		Owner                     string       `json:"owner"`
		ShortChannelId            string       `json:"short_channel_id"`
		FundingTxId               string       `json:"funding_txid"`
		MsatoshiToUs              msatoshi     `json:"msatoshi_to_us"`
		MsatoshiTotal             msatoshi     `json:"msatoshi_total"`
		DustLimitSatoshis         int64        `json:"dust_limit_satoshis"`
		MaxHtlcValueInFlightMsats int64        `json:"max_htlc_value_in_flight_msats"`
		ChannelReserveSatoshis    int64        `json:"channel_reserve_satoshis"`
		HtlcMinimumMsat           msatoshi     `json:"htlc_minimum_msat"`
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
	peer struct {
		PeerId    string   `json:"id"`
		Connected bool     `json:"connected"`
		Netaddr   netaddr  `json:"netaddr"`
		Channels  channels `json:"channels"`
	}
	// node describes a single node.
	node struct {
		// isPeer is true if this node is one of our peers.
		isPeer bool
		// Connected is set to true for peers that are currently connected.
		Connected bool
		// Channels holds any channels that nodes that are our peers has.
		Channels channels `json:"channels"`

		NodeId        string  `json:"nodeid"`
		Alias         alias   `json:"alias"`
		Color         string  `json:"color"`
		LastTimestamp int64   `json:"last_timestamp"`
		Addresses     address `json:"addresses"`
	}
	// nodes describes several nodes.
	nodes []node
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
	// allNodes is a map from node id to that node.
	allNodes map[string]node
	// lightningState describes the last known state of the lightningd daemon.
	lightningdState struct {
		pid      int
		args     []string
		Alias    alias
		Info     getInfoResponse
		Nodes    allNodes
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

var (
	states = map[channelStateNum]channelState{
		ChanneldNormalState:          "CHANNELD_NORMAL",
		ChanneldAwaitingLockinState:  "CHANNELD_AWAITING_LOCKIN",
		OpeningdState:                "OPENINGD",
		OnchaindTheirUnilateralState: "ONCHAIND_THEIR_UNILATERAL",
		OnchaindOurUnilateralState:   "ONCHAIND_OUR_UNILATERAL",
		ClosingdSigexchangeState:     "CLOSINGD_SIGEXCHANGE",
	}
	// TODO: eliminate global variable
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
	channelCapacities = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "channel_capacities_msatoshi",
			Help:      "Capacity of channels in millisatoshi by name and channel state.",
		},
		[]string{"node_id", "state"},
	)
	channelBalances = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "lightningd",
			Name:      "channel_balances_msatoshi",
			Help:      "Balance to us of channels in millisatoshi by name and channel state.",
		},
		[]string{"node_id", "state", "direction"},
	)
	// TODO: Add metrics showing last seen timestamp from listnodes instead of binary connected/unconnected status, which
	// doesn't form timeseries easily.
	debugging = os.Getenv("LNMON_DEBUGGING") == "1"
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(bitcoindRunning)
	prometheus.MustRegister(lightningdRunning)
	prometheus.MustRegister(availableFunds)
	prometheus.MustRegister(channelCapacities)
	prometheus.MustRegister(numChannels)
	prometheus.MustRegister(numNodes)
	prometheus.MustRegister(numPeers)
	prometheus.MustRegister(channelBalances)
	prometheus.MustRegister(ourChannels)
}

// getFile returns the contents of the specified file.
func getFile(f string) ([]byte, error) {
	// Asset is defined in bindata.go.
	return Asset(f)
}

// String returns a human-readable description of the netaddr.
func (n netaddr) String() string {
	if len(n) < 1 {
		return "netaddr{}"
	}
	return fmt.Sprintf("%s", n[0])
}

// ToNode returns the node representation of this peer. All peers are nodes, but not all nodes are peers.
func (p peer) ToNode() node {
	addrinfo := []addressInfo{}
	if p.Connected {
		// TODO: Find out difference between netaddr (values like [::ffff:138.68.252.183]:58638, presumably source addr of tcp socket for our peer,
		// and addresses field from listnodes, presumably announced service addr for accepting p2p traffic.
		// addr := p.Netaddr.String()
		// parts := strings.Split(addr, ":")
		// addrinfo = append(addrinfo, addressInfo{AddressType: "notimplemented", Address: p.Netaddr.String(), Port: 9735})
	}
	return node{
		isPeer:    true,
		Connected: p.Connected,
		Addresses: addrinfo,
		Channels:  p.Channels,
		NodeId:    p.PeerId,
	}
}

// String returns a human-readable description of the node.
func (n node) String() string {
	parts := []string{
		fmt.Sprintf("id: %s", n.NodeId),
		fmt.Sprintf("isPeer: %v", n.isPeer),
		fmt.Sprintf("Connected: %v", n.Connected),
	}
	if n.Alias != "" {
		parts = append(parts, fmt.Sprintf("Alias: %s", n.Alias))
	}
	if n.Color != "" {
		parts = append(parts, fmt.Sprintf("Color: %s", n.Color))
	}
	if n.Connected {
		parts = append(parts, fmt.Sprintf("Addresses: %s", n.Addresses))
	}
	if len(n.Channels) > 0 {
		parts = append(parts, fmt.Sprintf("Channels: %s", n.Channels))
	}
	return fmt.Sprintf(
		"node{%s}",
		strings.Join(parts, ", "),
	)
}

// String returns a human-readable description of the address.
func (addr address) String() string {
	if len(addr) < 1 {
		return fmt.Sprintf("address{}")
	}
	parts := make([]string, len(addr), len(addr))
	for i, a := range addr {
		parts[i] = fmt.Sprintf("%s:%d", a.Address, a.Port)
	}
	return strings.Join(parts, ", ")
}

// String returns a human-readable description of the nodes.
func (ns nodes) String() string {
	desc := make([]string, len(ns), len(ns))
	for i, n := range ns {
		desc[i] = fmt.Sprintf("%d: %s", i, n)
	}
	return fmt.Sprintf("%d nodes: %s", len(ns), strings.Join(desc, ", "))
}

// Desc returns a description of the nodes.
func (ns nodes) Desc() string {
	return fmt.Sprintf("%d nodes", len(ns))
}

// Implement sort.Interface for channels to sort them in reasonable order.
func (cs channels) Len() int      { return len(cs) }
func (cs channels) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs channels) Less(i, j int) bool {
	return cs[i].State < cs[j].State
}

// AsBTC returns a string formatting the msatoshi amount as BTC.
func (msat msatoshi) AsBTC() string {
	return fmt.Sprintf("%.5f BTC", float64(int64(msat))/1.0e11)
}

func (ns allNodes) getAlias(nodeId string) (alias, error) {
	n, exists := ns[nodeId]
	if !exists {
		return "", fmt.Errorf("no such node %q", nodeId)
	}
	return n.Alias, nil
}

// updateNodes updates the nodes with new node information.
func (s lightningdState) updateNodes(newNodes nodes) {
	for _, nn := range newNodes {
		on, exists := s.Nodes[nn.NodeId]
		if !exists {
			if debugging {
				fmt.Printf("Learned about new node %v\n", nn.NodeId)
			}
			s.Nodes[nn.NodeId] = nn
		} else {
			if debugging {
				fmt.Printf("Updating any stale info we had on node %v\n", nn.NodeId)
				fmt.Printf("Updating any stale info we had on old node %v to new %v\n", on, nn)
			}
			if on.Alias != "" && nn.Alias == "" {
				// Preserve alias if we knew it.
				nn.Alias = on.Alias
			}
			if on.Color != "" && nn.Color == "" {
				// Preserve color if we knew it.
				nn.Color = on.Color
			}
			if on.isPeer {
				// Note: We avoid marking peers as non-peers when listnodes returns with isPeer=false. This could be less hacky.
				nn.isPeer = true
				nn.Connected = on.Connected
				nn.Channels = on.Channels
			}
			// TODO: may or may not want to update on.Addresses and on.last_timestamp
			s.Nodes[nn.NodeId] = nn
		}
	}
}

// Peers returns all nodes that are our peers.
func (ns allNodes) Peers() nodes {
	result := nodes{}
	for _, n := range ns {
		if n.isPeer {
			result = append(result, n)
		}
	}
	sort.Sort(sort.Reverse(result))
	return result
}

// Implement sort.Interface for nodes to sort them in reasonable order.
func (ns nodes) Len() int      { return len(ns) }
func (ns nodes) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodes) Less(i, j int) bool {
	if !ns[i].isPeer && ns[j].isPeer {
		// Non-peer nodes are "less" than peer nodes.
		return true
	}
	if ns[i].isPeer && !ns[j].isPeer {
		// Peer nodes are "more" than non-peer nodes.
		return false
	}
	if len(ns[i].Channels) < len(ns[j].Channels) {
		// Peers with fewer channels are "less" than ones with more of them.
		return true
	}
	if len(ns[i].Channels) > len(ns[j].Channels) {
		// Peers with more channels are never "less" than ones with fewer of them.
		return false
	}
	if !ns[i].Connected && ns[j].Connected {
		// Unconnected peers are "less" than connected ones.
		return true
	}
	if ns[i].Connected && !ns[j].Connected {
		// Connected peers are never "less" than unconnected ones.
		return false
	}
	if len(ns[i].Channels) > 1 && len(ns[j].Channels) > 1 {
		// If we and the other peer has at least one channel, let us be "less" than
		// our peer if our first channel is "less" than theirs.
		cs := channels{ns[i].Channels[0], ns[j].Channels[0]}
		sort.Sort(cs)
		return cs.Less(0, 1)
	}
	// Tie-breaker: alphabetic ordering of peer id.
	return ns[i].NodeId < ns[j].NodeId
}

// NumConnected returns the number of connected nodes.
func (ns nodes) NumConnected() int {
	num := 0
	for _, n := range ns {
		if n.Connected {
			num += 1
		}
	}
	return num
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
func (ns nodes) NumChannelsByState() map[channelState]int64 {
	byState := map[channelState]int64{}
	for _, n := range ns {
		if !n.isPeer {
			// Nodes that are not our peers can't have cbannels with us.
			continue
		}
		for _, c := range n.Channels {
			byState[c.State] += 1
		}
	}
	return byState
}

// TotalChannelCapacity returns the total capacity of channels by state (e.g CHANNELD_NORMAL).
func (ns nodes) TotalChannelCapacity() map[channelState]msatoshi {
	result := map[channelState]msatoshi{}
	for _, n := range ns {
		for _, c := range n.Channels {
			if !n.isPeer {
				// Nodes that are not our peers can't have cbannels towards us.
				continue
			}
			result[c.State] += c.MsatoshiTotal
		}
	}
	return result
}

// ToUsChannelBalance returns the balance of channels by state (e.g CHANNELD_NORMAL).
func (ns nodes) ToUsChannelBalance() map[channelState]msatoshi {
	result := map[channelState]msatoshi{}
	for _, n := range ns {
		if !n.isPeer {
			// Nodes that are not our peers can't have cbannels towards us.
			continue
		}
		for _, c := range n.Channels {
			result[c.State] += c.MsatoshiToUs
		}
	}
	return result
}

// String returns a human-readable description of the channel listings.
func (cls channelListings) String() string {
	return fmt.Sprintf("%d channels", len(cls))
}

// String returns a human-readable description of the channels.
func (cs channels) String() string {
	if len(cs) == 0 {
		return ""
	}
	if len(cs) > 1 {
		// TODO: Find how this is supported by protocol.
		return "<unsupported multiple channels>"
	}
	return cs[0].String()
}

func (cs channels) State() string {
	if len(cs) != 1 {
		return ""
	}
	return string(cs[0].State)
}

func (cs channels) MilliSatoshiToUs() msatoshi {
	if len(cs) != 1 {
		return msatoshi(-1)
	}
	return cs[0].MsatoshiToUs
}

func (cs channels) MilliSatoshiTotal() msatoshi {
	if len(cs) != 1 {
		return msatoshi(-1)
	}
	return cs[0].MsatoshiTotal
}

func (cs channels) MilliSatoshiToThem() msatoshi {
	return cs.MilliSatoshiTotal() - cs.MilliSatoshiToUs()
}

// String returns a human-readable description of the channel.
func (c channel) String() string {
	parts := []string{
		fmt.Sprintf("state: %s", c.State),
	}
	if c.FundingTxId != "" {
		parts = append(parts, fmt.Sprintf("funding_txid: %s", c.FundingTxId))
	}
	if c.MsatoshiToUs != 0 {
		parts = append(parts, fmt.Sprintf("msatoshi_to_us: %d", c.MsatoshiToUs))
	}
	if c.MsatoshiTotal != 0 {
		parts = append(parts, fmt.Sprintf("msatoshi_total: %d", c.MsatoshiTotal))
	}
	return fmt.Sprintf("channel{%s}", strings.Join(parts, ", "))
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

// TODO: expose counters per cli command.
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
func (c cli) ListChannels() (*channelListings, error) {
	respstring, err := c.exec("listchannels")
	if err != nil {
		return nil, err

	}
	resp := struct {
		Channels channelListings `json:"channels"`
	}{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp.Channels, nil
}

// ListFunds returns the lightning-cli response to listfunds.
func (c cli) ListFunds() (*outputs, error) {
	respstring, err := c.exec("listfunds")
	if err != nil {
		return nil, err
	}
	resp := struct {
		Outputs outputs `json:"outputs"`
	}{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp.Outputs, nil
}

// ListNodes returns the lightning-cli response to listnodes.
func (c cli) ListNodes() (*nodes, error) {
	respstring, err := c.exec("listnodes")
	if err != nil {
		return nil, err
	}
	resp := struct {
		Nodes nodes `json:"nodes"`
	}{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp.Nodes, nil
}

// ListPeers returns the nodes returned by lightning-cli listpeers.
func (c cli) ListPeers() (*nodes, error) {
	respstring, err := c.exec("listpeers")
	if err != nil {
		return nil, err
	}
	resp := struct {
		Peers []peer `json:"peers"`
	}{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	result := make(nodes, len(resp.Peers), len(resp.Peers))
	for i, p := range resp.Peers {
		result[i] = p.ToNode()
	}
	return &result, nil
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
		Nodes: allNodes{},
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
	s.Channels = *channels
	numChannels.Set(float64(len(s.Channels)))
	log.Printf("Learned of %d channels.\n", len(s.Channels))

	outputs, err := c.ListFunds()
	if err != nil {
		return nil, err
	}
	s.Outputs = *outputs
	availableFunds.Set(float64(s.Outputs.Sum()))
	log.Printf("Learned of %d %v.\n", len(s.Outputs), s.Outputs)
	if debugging {
		log.Printf("Before doing listpeers we know about %d nodes.\n", len(s.Nodes))
	}
	peerNodes, err := c.ListPeers()
	if err != nil {
		return nil, err
	}
	if debugging {
		fmt.Printf("Doing listpeers we learned about %d peer nodes:\n", len(*peerNodes))
		for _, p := range *peerNodes {
			fmt.Printf("  %s\n", p.NodeId)
		}
	}
	s.updateNodes(*peerNodes)
	if debugging {
		log.Printf("After updating with listpeers results, we now know of %d nodes:\n", len(s.Nodes))
		for k, _ := range s.Nodes {
			fmt.Printf("  %s\n", k)
		}
	}

	peers := s.Nodes.Peers() // TODO: could do the below with entire s.Nodes too; methods filter out non-peers where applicable.
	numPeers.With(prometheus.Labels{"connected": "connected"}).Set(float64(peers.NumConnected()))
	numPeers.With(prometheus.Labels{"connected": "unconnected"}).Set(float64(len(peers) - peers.NumConnected()))
	for state, n := range s.Nodes.Peers().NumChannelsByState() {
		// log.Printf("We have %d channels in state %q\n", n, state)
		ourChannels.With(prometheus.Labels{"state": string(state)}).Set(float64(n))
	}
	// log.Printf("lightningd listpeers response: %+v\n", peers)
	if debugging {
		log.Printf("Before doing listnodes we know about %d nodes.\n", len(s.Nodes))
	}
	nodes, err := c.ListNodes()
	if err != nil {
		return nil, err
	}
	if debugging {
		log.Printf("Doing listnodes we learned about %d nodes.\n", len(*nodes))
	}
	s.updateNodes(*nodes)
	if debugging {
		log.Printf("After updating with listnodes results, we now know of %d nodes.\n", len(s.Nodes))
	}
	numNodes.Set(float64(len(s.Nodes)))
	// log.Printf("lightningd listnodes response: %+v\n", nodes)

	for _, n := range s.Nodes {
		if n.isPeer && len(n.Channels) >= 1 {
			if n.Channels.MilliSatoshiTotal() > 0 {
				channelCapacities.With(
					prometheus.Labels{
						"node_id": n.NodeId,
						"state":   n.Channels.State(),
					}).Set(float64(n.Channels.MilliSatoshiTotal()))
			}
			if n.Channels.MilliSatoshiToUs() > 0 {
				channelBalances.With(
					prometheus.Labels{
						"node_id":   n.NodeId,
						"state":     n.Channels.State(),
						"direction": "to_us",
					}).Set(float64(n.Channels.MilliSatoshiToUs()))
			}
			if n.Channels.MilliSatoshiToThem() > 0 {
				channelBalances.With(
					prometheus.Labels{
						"node_id":   n.NodeId,
						"state":     n.Channels.State(),
						"direction": "to_them",
					}).Set(float64(n.Channels.MilliSatoshiToThem()))
			}
		}
	}

	s.Alias, err = s.Nodes.getAlias(s.Info.NodeId)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func refresh() {
	// TODO: Maybe don't assume that bitcoind/lightningd always is in "pid" namespace..
	namespace := "pid"
	registeredBitcoin := false
	registeredLn := false
	for {
		btcState, err := getBitcoindState()
		if err != nil {
			log.Printf("Failed to get bitcoind state: %v\n", err)
		} else {
			allState.Bitcoind = *btcState
		}
		// TODO: Find way to avoid below panicking on "duplicate metrics collector registration attemped".. we'd like
		// to just export process metrics for both bitcoind and lightningd, but doesn't seem to be supported in current
		// prometheus client lib.. can either patch + upstream, or run one collector process per bitcoind/lightningd..
		// or put the processes in different namespaces may also make the metrics unique.
		if allState.Bitcoind.IsRunning() {
			if false && !registeredBitcoin {
				lc := prometheus.NewProcessCollector(allState.Bitcoind.pid, namespace)
				prometheus.MustRegister(lc)
				registeredBitcoin = true
				log.Printf("Registered ProcessCollector for bitcoind pid %d in namespace %q\n", allState.Bitcoind.pid, namespace)
			}
			bitcoindRunning.Set(1)
		} else {
			bitcoindRunning.Set(0)
		}

		// TODO: need to persist lnState.Nodes if we want to persist info we find between polls.
		lnState, err := getLightningdState()
		if err != nil {
			log.Printf("Failed to get lightningd state: %v\n", err)
		} else {
			allState.Lightningd = *lnState
		}
		if allState.Lightningd.IsRunning() {
			if !registeredLn {
				lc := prometheus.NewProcessCollector(allState.Lightningd.pid, namespace)
				prometheus.MustRegister(lc)
				registeredLn = true
				log.Printf("Registered ProcessCollector for lightningd pid %d in namespace %s\n", allState.Lightningd.pid, namespace)
			}
			lightningdRunning.Set(1)
		} else {
			lightningdRunning.Set(0)
		}

		// lightning-cli listinvoice
		// lightning-cli listpayments
		time.Sleep(time.Minute)
	}
}

// indexHandler serves the index page.
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
