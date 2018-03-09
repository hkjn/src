// lnmon.go is a tool for pulling out and serving up data from lightning-cli for monitoring.
//
// TODO: If we were to track state between CLI polls, we could detect e.g. channel state transitions,
// new channels, etc., to create an event stream.
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/crypto/acme/autocert"
)

type (
	cli         struct{}
	addressInfo struct {
		AddressType string `json:"type"`
		Address     string `json:"address"`
		Port        int    `json:"port"`
	}
	address []addressInfo
	// channeFundlListing describes one channel from listfunds output.
	channelFundListing struct {
		PeerId          nodeId  `json:"peer_id"`
		ShortChannelId  string  `json:"short_channel_id"`
		ChannelSat      satoshi `json:"channel_sat"`
		ChannelTotalSat satoshi `json:"channel_total_sat"`
		FundingTxId     string  `json:"funding_tx_id"`
	}
	// channelListing describes one channel from listchannels output.
	channelListing struct {
		Source          nodeId `json:"source"`
		Destination     nodeId `json:"destination"`
		ShortChannelId  string `json:"short_channel_id"`
		Flags           int64  `json:"flags"`
		Active          bool   `json:"active"`
		Public          bool   `json:"public"`
		LastUpdate      unixTs `json:"last_update"`
		BaseFeeMsats    int64  `json:"base_fee_millisatoshi"`
		FeePerMillionth int64  `json:"fee_per_millionth"`
		Delay           int64  `json:"delay"`
	}
	// channelListings describes several channels from listchannels output.
	channelListings []channelListing

	// msatoshi is number of millisatoshis, the main unit used in the LN protocol.
	msatoshi int64
	// satoshi is number of satoshis, also sometimes used in the LN protocol.
	satoshi int64
	// channel describes an individual channel.
	//
	// For channels to our own node, we know the state and msatoshi details, but
	// for other nodes we only know if they are active and/or public and their short_channel_id.
	channel struct {
		State                     channelState `json:"state"`
		Owner                     nodeId       `json:"owner"`
		ShortChannelId            string       `json:"short_channel_id"`
		FundingTxId               string       `json:"funding_txid"`
		MsatsToUs                 msatoshi     `json:"msatoshi_to_us"`
		MsatsTotal                msatoshi     `json:"msatoshi_total"`
		DustLimitSatoshis         int64        `json:"dust_limit_satoshis"`
		MaxHtlcValueInFlightMsats int64        `json:"max_htlc_value_in_flight_msats"`
		ChannelReserveSatoshis    int64        `json:"channel_reserve_satoshis"`
		HtlcMinimumMsat           msatoshi     `json:"htlc_minimum_msat"`
		ToSelfDelay               int64        `json:"to_self_delay"`
		MaxAcceptedHtlcs          int64        `json:"max_accepted_htlcs"`
		// Active is true if the channel's last known state is active.
		Active bool `json:"active"`
		// Public is true if the channel is public.
		Public bool `json:"public"`
	}
	// alias is an optional human-readable alias for a node.
	alias string
	// channels describes several channel structures.
	channels []channel
	// netaddr describes the network addresses for a peer.
	netaddr []string
	// peer describes a single peer.
	peer struct {
		PeerId    nodeId   `json:"id"`
		Connected bool     `json:"connected"`
		Netaddr   netaddr  `json:"netaddr"`
		Channels  channels `json:"channels"`
	}
	// int64 represents a UNIX timestamp as int.
	unixTs int64
	// nodeId represents a unique id for a LN node.
	nodeId string
	// node describes a single node.
	node struct {
		// isPeer is true if this node is one of our peers.
		isPeer bool
		// Connected is set to true for nodes that are our peers and also are currently connected.
		Connected bool
		// Channels holds any channels that the node has.
		Channels channels `json:"channels"`

		NodeId        nodeId  `json:"nodeid"`
		Alias         alias   `json:"alias"`
		Color         string  `json:"color"`
		LastTimestamp unixTs  `json:"last_timestamp"`
		Addresses     address `json:"addresses"`
	}
	// invoiceRequest is a request to create a new invoice.
	invoiceRequest struct {
		Msatoshi    msatoshi `json:"msatoshi"`
		Label       string   `json:"label"`
		Description string   `json:"description"`
	}
	// invoiceResponse is the response from lightning-cli invoice.
	invoiceResponse struct {
		PaymentHash string `json:"payment_hash"`
		ExpiryTime  unixTs `json:"expiry_time"`
		ExpiresAt   unixTs `json:"expires_at"`
		Bolt11      string `json:"bolt11"`
	}
	// invoiceListing is a single result from listinvoices.
	invoiceListing struct {
		Label            string   `json:"label"`
		PaymentHash      string   `json:"payment_hash"`
		Msatoshi         msatoshi `json:"msatoshi"`
		Status           string   `json:"status"`
		ExpiryTime       unixTs   `json:"expiry_time"`
		ExpiresAt        unixTs   `json:"expires_at"`
		PaidTimestamp    unixTs   `json:"paid_timestamp"`
		PaidAt           unixTs   `json:"paid_at"`
		PayIndex         int      `json:"pay_index"`
		MsatoshiReceived msatoshi `json:"msatoshi_received"`
	}
	// invoiceListings describes several results from listinvoices.
	invoiceListings []invoiceListing
	// nodes describes several nodes.
	nodes []node
	// candidates describes candidate nodes to open channels with.
	candidates nodes
	// output describes an individual output.
	output struct {
		TxId string `json:"txid"`
		// Output is the index of the txo.
		Output int64 `json:"output"`
		Value  int64 `json:"value"`
	}
	// outputs describes several outputs.
	outputs             []output
	channelFundListings []channelFundListing
	fundListing         struct {
		Outputs  outputs             `json:"outputs"`
		Channels channelFundListings `json:"channels"`
	}

	payment struct {
		PaymentId       int64    `json:"id"`
		PaymentHash     string   `json:"payment_hash"`
		Destination     nodeId   `json:"destination"`
		Msatoshi        msatoshi `json:"msatoshi"`
		Timestamp       unixTs   `json:"timestamp"`
		CreatedAt       unixTs   `json:"created_at"`
		Status          string   `json:"status"`
		PaymentPreimage string   `json:"payment_preimage"`
	}
	payments []payment

	// getInfoResponse is the format of the getinfo response from lightning-cli.
	getInfoResponse struct {
		NodeId      nodeId  `json:"id"`
		Port        int     `json:"port"`
		Address     address `json:"address"`
		Version     string  `json:"version"`
		Blockheight int     `json:"blockheight"`
		Network     string  `json:"network"`
	}
	// allNodes is a map from node id to that node.
	allNodes map[nodeId]node
	// state describes the last known state.
	state struct {
		// MonVersion is the version of lnmon.
		MonVersion   string
		pid          int
		args         []string
		Alias        alias
		Info         getInfoResponse
		Nodes        allNodes
		Channels     channelListings
		Payments     payments
		Outputs      fundListing
		Invoices     invoiceListings
		gauges       map[string]prometheus.Gauge
		counterVecs  map[string]*prometheus.CounterVec
		gaugeVecs    map[string]*prometheus.GaugeVec
		callCounters map[string]*prometheus.CounterVec
	}
	channelStateNum int
	channelState    string
	// indexHandler serves the index requests.
	indexHandler struct {
		tmpl template.Template
		s    *state
	}
	// nodeHandler serves /node requests.
	nodeHandler struct {
		tmpl template.Template
		s    *state
	}
	// cmdHandler serves /cmd requests.
	cmdHandler struct {
		s *state
	}
)

const (
	// Channel states are enumerated here. Note that states with lower numbers sort before higher.
	ChannelUnknownState channelStateNum = iota
	ChanneldNormalState
	ChanneldAwaitingLockinState
	OpeningdState
	OnchaindTheirUnilateralState
	OnchaindOurUnilateralState
	ClosingdSigexchangeState

	// Prometheus monitoring prefix.
	counterPrefix = "lightningd"
)

var (
	// allState state

	states = map[channelStateNum]channelState{
		ChannelUnknownState:          "<unknown channel state>",
		ChanneldNormalState:          "CHANNELD_NORMAL",
		ChanneldAwaitingLockinState:  "CHANNELD_AWAITING_LOCKIN",
		OpeningdState:                "OPENINGD",
		OnchaindTheirUnilateralState: "ONCHAIND_THEIR_UNILATERAL",
		OnchaindOurUnilateralState:   "ONCHAIND_OUR_UNILATERAL",
		ClosingdSigexchangeState:     "CLOSINGD_SIGEXCHANGE",
	}
	lnmonVersion string
	debugging    = os.Getenv("LNMON_DEBUGGING") == "1"
	addr         = os.Getenv("LNMON_ADDR")
	hostname     = os.Getenv("LNMON_HOSTNAME")
	httpPrefix   = os.Getenv("LNMON_HTTP_PREFIX")
)

// getFile returns the contents of the specified file.
func getFile(f string) ([]byte, error) {
	// Asset is defined in bindata.go.
	return Asset(f)
}

// WithEscapedEntities returns the escaped HTML entities for the alias.
//
// TODO: We added this to attempt to correctly render utf8 characters instead of '?'
// in HTML, but using alias field of node id
// 03939ff69d65a13c4bb2585042e7eb7e75a7c77289ab5794d1b973721d86c6839c
// as an example, it seems that either lightning-cli or us is mangling the bytes:
// Via sites rendering the node's alias correctly,, byte stream should be 43 6F 63 6F 50 69 E2 9A A1
// We are for some reason getting actual '?' / 3F characters:             43 6F 63 6F 50 69 3F 3F 3F
// Also, we probably shouldn't (need to) render actual HTML entities, since that could lead to XSS vulns.
func (a alias) WithEscapedEntities() string {
	result := []string{}
	s := string(a)
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		result = append(result, fmt.Sprintf(`&#x%X;`, r))
	}
	return strings.Join(result, " ")
}

// Short returns the first few characters of the node id.
func (nid nodeId) Short() string {
	return string(nid[:6]) + "[..]"
}

// Time returns the unixTs converted to regular time.Time format.
func (ts unixTs) Time() time.Time {
	return time.Unix(int64(ts), 0)
}

// Since returns a description of the duration passed since the timestamp.
func (ts unixTs) Since() string {
	d := time.Since(ts.Time())
	hrs := d.Hours()
	if hrs > 24*365 {
		// More than one year? Probably start of epoch.
		return "never seen"
	}
	if hrs > 24.0 {
		return fmt.Sprintf("%.2f days", hrs/24)
	}
	if d > time.Hour {
		return fmt.Sprintf("%.2f hrs", hrs)
	}
	if d > time.Minute {
		return fmt.Sprintf("%.2f min", d.Minutes())
	}
	if d > time.Second {
		return fmt.Sprintf("%.2f sec", d.Seconds())
	}
	return fmt.Sprintf("%.2f ms", float64(d.Nanoseconds()/1e6))
}

// String returns the name of the state, e.g. "OPENINGD".
func (s channelStateNum) String() string {
	if ChanneldNormalState <= s && s <= ClosingdSigexchangeState {
		return string(states[s])

	}
	return "Invalid channelstate " + string(s)
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

// updateChannel updates the node's channels to take into account the new channelListing.
func (n *node) updateChannel(cl channelListing) {
	// TODO: store channels for one node as map?
	found := false
	for _, c := range n.Channels {
		if c.ShortChannelId == cl.ShortChannelId {
			// TODO: might need to update state on known channel?
			found = true
			break
		}
	}
	if !found {
		n.Channels = append(n.Channels, channel{
			// State: ""
			ShortChannelId: cl.ShortChannelId,
			Active:         cl.Active,
			Public:         cl.Public,
		})
	}
}

// String returns a human-readable description of the address.
func (addr address) String() string {
	if len(addr) == 0 {
		return ""
	}
	if len(addr) != 1 {
		return "<unsupported address>"
	}
	return fmt.Sprintf("%s:%d", addr[0].Address, addr[0].Port)
}

func (n node) KnownAddress() bool {
	if len(n.Addresses) != 1 {
		// TODO: support > 1 addresses, if this can occur.
		return false
	}
	return n.Addresses[0].Address != "" && n.Addresses[0].Port != 0
}

// isChannelCandidate returns true if we could open a channel to this node.
func (n node) isChannelCandidate() bool {
	if !n.KnownAddress() {
		return false
	}
	return true
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
	// TODO: could have more smartness here.
	return cs[i].State < cs[j].State
}

func (ns candidates) Len() int      { return len(ns) }
func (ns candidates) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns candidates) Less(i, j int) bool {
	if ns[i].KnownAddress() && !ns[j].KnownAddress() {
		// Nodes with known addresses are "more" than ones we don't know the address for.
		return false
	}
	if !ns[i].KnownAddress() && ns[j].KnownAddress() {
		// Nodes without known addresses are "less" than ones than ones we know the address for.
		return true
	}
	if len(ns[i].Channels) < len(ns[j].Channels) {
		// Peers with fewer channels are "less" than ones with more of them.
		return true
	}
	if len(ns[i].Channels) > len(ns[j].Channels) {
		// Peers with more channels are "more" than ones with fewer of them.
		return false
	}
	if !ns[i].Channels.AnyWithUs() && ns[j].Channels.AnyWithUs() {
		// Nodes without channels to us are "more" than nodes with channels to us.
		return false
	}
	if ns[i].Channels.AnyWithUs() && !ns[j].Channels.AnyWithUs() {
		// Nodes with channels to us are "less" than nodes with no channels to us.
		return true
	}
	if ns[i].isPeer && !ns[j].isPeer {
		// Peers that are our peer are "more" than nodes that are not currently our peer.
		return false
	}
	if !ns[i].isPeer && ns[j].isPeer {
		// Peers that are not our peer are "less" than nodes that are currently our peer.
		return true
	}
	// TODO: should compare ns[i].Channels and ns[j].Channels.

	// Tie-breaker: alphabetic ordering of peer id.
	return ns[i].NodeId < ns[j].NodeId
}

// ChannelCandidates returns nodes that are good candidates for new channels.
//
// The number of nodes returned may be less than the full set of nodes.
func (ns nodes) ChannelCandidates() candidates {
	// TODO: We could use much better metrics to decide candidate nodes to open channels with, like:
	// 1. High availability should rank higher
	// 2. High connectivity (many other nodes having channels to it, with large balances and high volume of HTLCs going between them)
	// 3. Improving network graph structure
	// 4. Fast response time
	// 5. Long lifetime
	result := candidates{}
	for _, n := range ns {
		if !n.isChannelCandidate() {
			continue
		}
		result = append(result, n)
	}
	sort.Sort(sort.Reverse(result))
	max := 50
	if len(result) < 50 {
		max = len(result)
	}
	return result[:max]
}

// AsSat returns a description
func (msat msatoshi) AsSat() string {
	return fmt.Sprintf("%d sat", int64(msat)/1e3)
}

// String returns a string formatting the msatoshi amount as BTC, mBTC or sat as appropriate.
func (msat msatoshi) String() string {
	if msat == 0 {
		return ""
	}
	sat := float64(int64(msat)) / 1.0e3
	if sat < 1.0 {
		return fmt.Sprintf("%d msat", msat)
	}
	if sat < 1.0e4 {
		return fmt.Sprintf("%.3f sat", sat)
	}
	btc := sat / 1.0e8
	if btc < 0.005 {
		return fmt.Sprintf("%.4f mBTC", btc*1000.0)
	}
	return fmt.Sprintf("%.5f BTC", btc)
}

// updateNodes updates the nodes with new node information.
//
// TODO: If we get an error from lightningd we do s.reset(), but we might want to
// forget nodes that disappear are from listnodes output too.
func (s state) updateNodes(newNodes nodes) {
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
				if len(nn.Channels) == 0 && len(on.Channels) > 0 {
					nn.Channels = on.Channels
				}
			}
			// TODO: may or may not want to update on.Addresses.
			if on.LastTimestamp.Time().After(nn.LastTimestamp.Time()) {
				nn.LastTimestamp = on.LastTimestamp
			}
			s.Nodes[nn.NodeId] = nn
		}
	}
}

// AnyWithUs returns true if any of the channels are with us.
func (cs channels) AnyWithUs() bool {
	for _, c := range cs {
		if c.WithUs() {
			return true
		}
	}
	return false
}

// updateChannels updates the state based on the channelListings from listchannels.
func (s state) updateChannels(cls channelListings) {
	for _, cl := range cls {
		sn, exists := s.Nodes[cl.Source]
		if !exists {
			fmt.Printf("no such source node %s\n", cl.Source)
			continue
		}
		dn, exists := s.Nodes[cl.Destination]
		if !exists {
			fmt.Printf("no such dest node %s\n", cl.Destination)
			continue
		}

		sn.updateChannel(cl)
		dn.updateChannel(cl)

		s.Nodes[sn.NodeId] = sn
		s.Nodes[dn.NodeId] = dn
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

// ChannelCandidates returns nodes that are good candidates for new channels.
//
// The number of nodes returned may be less than the full set of nodes.
func (ns allNodes) ChannelCandidates() candidates {
	result := make(nodes, len(ns), len(ns))
	i := 0
	for _, n := range ns {
		result[i] = n
		i += 1
	}
	return result.ChannelCandidates()
}

// NumNormalChannels returns the number of CHANNELD_NORMAL channels with our node.
func (ns nodes) NumNormalChannels() int {
	sum := 0
	for _, n := range ns {
		if !n.isPeer {
			continue
		}
		for _, c := range n.Channels {
			if !c.WithUs() {
				continue
			}
			if c.State != states[ChanneldNormalState] {
				continue
			}
			sum += 1
		}
	}
	return sum
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
	if !ns[i].Channels.AnyWithUs() && ns[j].Channels.AnyWithUs() {
		// Nodes without channels to us are "less" than nodes with channels to us.
		return true
	}
	if ns[i].Channels.AnyWithUs() && !ns[j].Channels.AnyWithUs() {
		// Nodes with channels to us are "more" than nodes with no channels to us.
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
func (fs fundListing) String() string {
	return fmt.Sprintf("outputs totalling %v sat", fs.Sum())
}

// Sum returns the total value of all the outputs in the fund listing.
func (fs fundListing) Sum() int64 {
	sum := int64(0)
	for _, o := range fs.Outputs {
		sum += o.Value
	}
	return sum
}

// NumChannelsByState returns a map from channel state to number of channels in that state.
func (ns nodes) NumChannelsByState() map[channelState]int64 {
	byState := map[channelState]int64{}
	for _, n := range ns {
		if !n.isPeer {
			// Nodes that are not our peers can't have channels with us.
			continue
		}
		for _, c := range n.Channels {
			byState[c.State] += 1
		}
	}
	return byState
}

// String returns a human-readable description of the channel listing.
func (cl channelListing) String() string {
	parts := []string{
		fmt.Sprintf("source: %s", cl.Source),
		fmt.Sprintf("destination: %s", cl.Destination),
		fmt.Sprintf("short_channel_id: %s", cl.ShortChannelId),
	}
	return fmt.Sprintf("channelListing{%s}", strings.Join(parts, ", "))
}

// String returns a human-readable description of the channel listings.
func (cls channelListings) String() string {
	return fmt.Sprintf("%d channels", len(cls))
}

// WithUs returns true if the channel is to our node.
func (c channel) WithUs() bool {
	// TODO: could use a cleaner check than deducing that if we know the capacity or balance it must be our channel..
	return c.MsatsTotal > msatoshi(0) || c.MsatsToUs > msatoshi(0)
}

// BalanceToUs returns the total msatoshi to us in the channels.
func (cs channels) BalanceToUs() msatoshi {
	total := msatoshi(0)
	for _, c := range cs {
		total += c.MsatsToUs
	}
	return total
}

// BalanceToThem returns the total msatoshi to them in the channels.
func (cs channels) BalanceToThem() msatoshi {
	total := msatoshi(0)
	for _, c := range cs {
		total += (c.MsatsTotal - c.MsatsToUs)
	}
	return total
}

func (cs channels) DescBalance() string {
	us := cs.BalanceToUs()
	them := cs.BalanceToThem()
	if us == 0 && them == 0 {
		return "no balance in channel!?"
	}
	usDesc := us.String()
	themDesc := them.String()
	if us == 0 {
		usDesc = "0"
	}
	if them == 0 {
		themDesc = "0"
	}
	return fmt.Sprintf("%s / %s", usDesc, themDesc)
}

// OurChannels returns the channels that are with our node.
func (cs channels) OurChannels() channels {
	result := channels{}
	for _, c := range cs {
		if !c.WithUs() {
			continue
		}
		result = append(result, c)
	}
	return result
}

// DescState describes the state of the channels.
func (cs channels) DescState() string {
	if len(cs) < 1 {
		return ""
	}
	desc := []string{}
	withUs := 0
	for _, c := range cs {
		if c.WithUs() {
			desc = append(desc, string(c.State))
			withUs += 1
		}
	}
	if withUs > 0 {
		if withUs == 1 {
			return fmt.Sprintf(
				"Channel with us (%s): %s.",
				strings.Join(desc, ", "),
				cs.DescBalance(),
			)
		} else {
			return fmt.Sprintf(
				"%d channels with us, in state %s.",
				withUs,
				strings.Join(desc, ", "),
			)
		}
	} else {
		return "No channels with us."
	}
}

// String returns a human-readable description of the channel.
func (c channel) String() string {
	parts := []string{
		fmt.Sprintf("state: %s", c.State),
	}
	if c.FundingTxId != "" {
		parts = append(parts, fmt.Sprintf("funding_txid: %s", c.FundingTxId))
	}
	if c.MsatsToUs != 0 {
		parts = append(parts, fmt.Sprintf("msatoshi_to_us: %d", c.MsatsToUs))
	}
	if c.MsatsTotal != 0 {
		parts = append(parts, fmt.Sprintf("msatoshi_total: %d", c.MsatsTotal))
	}
	return fmt.Sprintf("channel{%s}", strings.Join(parts, ", "))
}

// IsRunning returns true if lightningd is running.
func (s state) IsRunning() bool {
	return s.pid != 0
}

// String returns a human-readable description of the invoice request.
func (ir invoiceRequest) String() string {
	return fmt.Sprintf("invoiceRequest{msatoshi: %v, label: %q, description: %q}", ir.Msatoshi, ir.Label, ir.Description)
}

// String returns a human-readable description of the invoice listing.
func (il invoiceListings) String() string {
	return fmt.Sprintf("%d invoices", len(il))
}

// String returns a human-readable description of the payments.
func (ps payments) String() string {
	return fmt.Sprintf("%d payments", len(ps))
}

// execCmd executes specified command with arguments and returns the output.
func execCmd(cmd string, arg ...string) (string, error) {
	c := exec.Command(cmd, arg...)
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Command exited with non-zero status.
			errmsg := fmt.Sprintf("Command %q exited with non-zero status: %v", fmt.Sprintf("%s %s", cmd, strings.Join(arg, " ")), err)
			errstring := stderr.String()
			if errstring != "" {
				errmsg += fmt.Sprintf(", stderr=%s", stderr.String())
			}
			outstring := stdout.String()
			if outstring != "" {
				errmsg += fmt.Sprintf(", stdout=%s", stdout.String())
			}
			return "", fmt.Errorf(errmsg)
		}
		return "", err
	}
	return stdout.String(), nil
}

// exec returns output when executing specified lightning-cli command.
func (c cli) exec(cmd string, args ...string) (string, error) {
	allArgs := append([]string{cmd}, args...)
	return execCmd("lightning-cli", allArgs...)
}

// GetInfo returns the getinfo response.
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

// Invoice returns the lightning-cli response to invoice.
func (c cli) Invoice(r invoiceRequest) (*invoiceResponse, error) {
	respstring, err := c.exec("invoice", fmt.Sprintf("%d", r.Msatoshi), r.Label, r.Description)
	if err != nil {
		return nil, err
	}
	resp := invoiceResponse{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
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
func (c cli) ListFunds() (*fundListing, error) {
	respstring, err := c.exec("listfunds")
	if err != nil {
		return nil, err
	}
	resp := fundListing{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
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

// ListInvoice returns the listinvoices response.
func (c cli) ListInvoices() (*invoiceListings, error) {
	respstring, err := c.exec("listinvoices")
	if err != nil {
		return nil, err
	}
	resp := struct {
		Invoices invoiceListings `json:"invoices"`
	}{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp.Invoices, nil
}

// ListPayments returns the listpayments response.
func (c cli) ListPayments() (*payments, error) {
	respstring, err := c.exec("listpayments")
	if err != nil {
		return nil, err
	}
	resp := struct {
		Payments payments `json:"payments"`
	}{}
	if err := json.Unmarshal([]byte(respstring), &resp); err != nil {
		return nil, err
	}
	return &resp.Payments, nil
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

// Equal returns true if the addresses are the same.
func (addr address) Equal(other address) bool {
	if len(addr) != len(other) {
		return false
	}
	for i, a := range addr {
		if a.AddressType != other[i].AddressType {
			return false
		}
		if a.Address != other[i].Address {
			return false
		}
		if a.Port != other[i].Port {
			return false
		}
	}
	return true
}

// Equal returns true if the info responses are the same.
func (info getInfoResponse) Equal(other getInfoResponse) bool {
	if info.NodeId != other.NodeId {
		return false
	}
	if info.Port != other.Port {
		return false
	}
	if !info.Address.Equal(other.Address) {
		return false
	}
	if info.Version != other.Version {
		return false
	}
	if info.Blockheight != other.Blockheight {
		return false
	}
	return true
}

// Equal returns true if the fundlistings are the same.
func (fl fundListing) Equal(other fundListing) bool {
	if !fl.Outputs.Equal(other.Outputs) {
		return false
	}
	if !fl.Channels.Equal(other.Channels) {
		return false
	}
	return true
}

// Equal returns true if the invoice listings are the same.
func (ils invoiceListings) Equal(other invoiceListings) bool {
	if len(ils) != len(other) {
		return false
	}
	for i, il := range ils {
		if !il.Equal(other[i]) {
			return false
		}
	}
	return true
}

// Equal returns true if the invoice listings are the same.
func (il invoiceListing) Equal(other invoiceListing) bool {
	if il.Label != other.Label {
		return false
	}
	if il.PaymentHash != other.PaymentHash {
		return false
	}
	if il.Msatoshi != other.Msatoshi {
		return false
	}
	if il.Status != other.Status {
		return false
	}
	if il.ExpiryTime != other.ExpiryTime {
		return false
	}
	if il.ExpiresAt != other.ExpiresAt {
		return false
	}
	if il.PaidTimestamp != other.PaidTimestamp {
		return false
	}
	if il.PaidAt != other.PaidAt {
		return false
	}
	if il.PayIndex != other.PayIndex {
		return false
	}
	if il.MsatoshiReceived != other.MsatoshiReceived {
		return false
	}
	return true
}

// Equal returns true if the payments are the same.
func (ps payments) Equal(other payments) bool {
	if len(ps) != len(other) {
		return false
	}
	for i, p := range ps {
		if !p.Equal(other[i]) {
			return false
		}
	}
	return true
}

// Equal returns true if the payments are the same.
func (p payment) Equal(other payment) bool {
	if p.PaymentId != other.PaymentId {
		return false
	}
	if p.PaymentHash != other.PaymentHash {
		return false
	}
	if p.Destination != other.Destination {
		return false
	}
	if p.Msatoshi != other.Msatoshi {
		return false
	}
	if p.Timestamp != other.Timestamp {
		return false
	}
	if p.CreatedAt != other.CreatedAt {
		return false
	}
	if p.Status != other.Status {
		return false
	}
	if p.PaymentPreimage != other.PaymentPreimage {
		return false
	}
	return true
}

// Equal returns true if the outputs are the same.
func (outs outputs) Equal(other outputs) bool {
	if len(outs) != len(other) {
		return false
	}
	for i, o := range outs {
		if o.TxId != other[i].TxId {
			return false
		}
		if o.Output != other[i].Output {
			return false
		}
		if o.Value != other[i].Value {
			return false
		}
	}
	return true
}

// Equal returns true if the channel fund listings are the same.
func (cs channelFundListings) Equal(other channelFundListings) bool {
	if len(cs) != len(other) {
		return false
	}
	for i, l := range cs {
		if !l.Equal(other[i]) {
			return false
		}
	}
	return true
}

// Equal returns true if the channel fund listings are the same.
func (cfl channelFundListing) Equal(other channelFundListing) bool {
	if cfl.PeerId != other.PeerId {
		return false
	}
	if cfl.ShortChannelId != other.ShortChannelId {
		return false
	}
	if cfl.ChannelSat != other.ChannelSat {
		return false
	}
	if cfl.ChannelTotalSat != other.ChannelTotalSat {
		return false
	}
	if cfl.FundingTxId != other.FundingTxId {
		return false
	}
	return true
}

// update brings the state in sync with the lightning-cli responses.
//
// Note that we reset all state between lightning-cli calls, to make sure we're not presenting stale data from earlier.
// This means that any failure to fetch new state from cli will result in empty state.
func (s *state) update() error {
	s.MonVersion = lnmonVersion
	s.Nodes = allNodes{}
	// TODO: grab mutex here to avoid data race when we write and ServeHTTP may read.
	ps, err := execCmd("pgrep", "-a", "lightningd")
	if err != nil {
		return err
	}
	parts := strings.Split(ps, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success.
	if len(parts) < 1 || len(parts[0]) == 0 {
		return fmt.Errorf("failed to parse lightningd status: %v", ps)
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	s.pid = pid
	for _, arg := range parts[1:] {
		s.args = append(s.args, arg)
	}
	s.gauges["running"].Set(1)

	c := &cli{}
	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "getinfo"}).Inc()
	info, err := c.GetInfo()
	if err != nil {
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "getinfo"}).Inc()
		return err
	}
	if !s.Info.Equal(*info) {
		log.Printf("lightningd getinfo response changed: %+v\n", info)
		s.Info = *info
		s.counterVecs["info"].Reset()
		s.counterVecs["info"].With(
			prometheus.Labels{
				"lnmon_version":      s.MonVersion,
				"lightningd_version": s.Info.Version,
				"network":            s.Info.Network,
			},
		).Set(1.0)
		s.gauges["blockheight"].Set(float64(s.Info.Blockheight))
	}

	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "listchannels"}).Inc()
	channels, err := c.ListChannels()
	if err != nil {
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "listchannels"}).Inc()
		return err
	}
	if len(s.Channels) != len(*channels) {
		s.Channels = *channels
		log.Printf("We now know of %d channels.\n", len(s.Channels))
	}
	s.gauges["num_channels"].Set(float64(len(s.Channels)))

	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "listfunds"}).Inc()
	outputs, err := c.ListFunds()
	if err != nil {
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "listfunds"}).Inc()
		return err
	}
	if !s.Outputs.Equal(*outputs) {
		s.Outputs = *outputs
		log.Printf("We now know of %d %v.\n", len(s.Outputs.Outputs), s.Outputs)
	}
	s.gauges["total_funds"].Set(float64(s.Outputs.Sum()))

	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "listpeers"}).Inc()
	peerNodes, err := c.ListPeers()
	if err != nil {
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "listpeers"}).Inc()
		return err
	}
	s.updateNodes(*peerNodes)

	peers := s.Nodes.Peers()
	s.gaugeVecs["num_peers"].With(prometheus.Labels{"connected": "connected"}).Set(float64(peers.NumConnected()))
	s.gaugeVecs["num_peers"].With(prometheus.Labels{"connected": "unconnected"}).Set(float64(len(peers) - peers.NumConnected()))

	// Delete all old metrics in vectors, so we don't accidentally persist e.g gauges measuring earlier
	// seen channel states. This is unfortunate, but seems necessary unless we want to track and
	// mutate earlier states.
	s.gaugeVecs["our_channels"].Reset()
	for state, n := range s.Nodes.Peers().NumChannelsByState() {
		// log.Printf("We have %d channels in state %q\n", n, state)
		s.gaugeVecs["our_channels"].With(prometheus.Labels{"state": string(state)}).Set(float64(n))
	}

	// We merge in the output from listpeers with the one from listpeers above.
	// log.Printf("lightningd listpeers response: %+v\n", peers)
	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "listnodes"}).Inc()
	nodes, err := c.ListNodes()
	if err != nil {
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "listnodes"}).Inc()
		return err
	}
	s.updateNodes(*nodes)
	s.gauges["num_nodes"].Set(float64(len(s.Nodes)))
	// log.Printf("lightningd listnodes response: %+v\n", nodes)

	// TODO: remove s.Channels entirely?
	s.updateChannels(s.Channels)

	s.counterVecs["aliases"].Reset()
	s.gaugeVecs["channel_capacities_msatoshi"].Reset()
	s.gaugeVecs["channel_balances_msatoshi"].Reset()
	for _, n := range s.Nodes {
		s.counterVecs["aliases"].With(
			prometheus.Labels{
				"node_id": string(n.NodeId),
				"alias":   string(n.Alias),
			}).Set(float64(n.LastTimestamp))

		for _, c := range n.Channels.OurChannels() {
			if c.MsatsTotal > msatoshi(0) {
				s.gaugeVecs["channel_capacities_msatoshi"].With(
					prometheus.Labels{
						"node_id": string(n.NodeId),
						"state":   string(c.State),
					}).Set(float64(c.MsatsTotal))
			}
			if c.MsatsToUs > msatoshi(0) {
				s.gaugeVecs["channel_balances_msatoshi"].With(
					prometheus.Labels{
						"node_id":   string(n.NodeId),
						"state":     string(c.State),
						"direction": "to_us",
					}).Set(float64(c.MsatsToUs))
			}
			if c.MsatsTotal-c.MsatsToUs > msatoshi(0) {
				s.gaugeVecs["channel_balances_msatoshi"].With(
					prometheus.Labels{
						"node_id":   string(n.NodeId),
						"state":     string(c.State),
						"direction": "to_them",
					}).Set(float64(c.MsatsTotal - c.MsatsToUs))
			}
		}
	}

	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "listinvoices"}).Inc()
	invoices, err := c.ListInvoices()
	if err != nil {
		// TODO: Set gauge with # of invoices (and payments), as well as total value in msats.
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "listinvoices"}).Inc()
		return err
	}
	if !s.Invoices.Equal(*invoices) {
		s.Invoices = *invoices
		log.Printf("We learned of new invoices: %v\n", s.Invoices)
	}

	s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "listpayments"}).Inc()
	payments, err := c.ListPayments()
	if err != nil {
		s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "listpayments"}).Inc()
		return err
	}
	if !s.Payments.Equal(*payments) {
		s.Payments = *payments
		log.Printf("We learned of new payments: %v\n", s.Payments)
	}

	n, exists := s.Nodes[s.Info.NodeId]
	if exists && s.Alias != n.Alias {
		s.Alias = n.Alias
		log.Printf("Found that our own alias is %q.\n", s.Alias)
	}
	return nil
}

// reset forgets all lightningd state.
func (s *state) reset() {
	s.pid = 0
	s.args = []string{}
	s.Alias = alias("")
	s.Info = getInfoResponse{}
	s.Nodes = allNodes{}
	s.Channels = channelListings{}
	s.Payments = payments{}
	s.Outputs = fundListing{}
	s.counterVecs["aliases"].Reset()
	s.counterVecs["info"].Reset()
	s.gauges["blockheight"].Set(0.0)
	s.gauges["num_channels"].Set(0.0)
	s.gauges["total_funds"].Set(0.0)
	s.gauges["running"].Set(0)
	s.gaugeVecs["num_peers"].Reset()
	s.gaugeVecs["our_channels"].Reset()
	s.gaugeVecs["channel_capacities_msatoshi"].Reset()
	s.gaugeVecs["channel_balances_msatoshi"].Reset()
}

// refresh updates the state.
func refresh(s *state) {
	// TODO: Maybe don't assume that lightningd always is in "pid" namespace..
	namespace := "pid"
	registeredLn := false
	for {
		if err := s.update(); err != nil {
			// TODO: increment counter here, so we can alert on possible lightningd crashes.
			log.Printf("Failed to update state: %v\n", err)
			s.reset()
		}
		if s.IsRunning() {
			if !registeredLn {
				// TODO: Need to handle case where we registered collector to pid #1, then
				// lightningd crashed and restarted with pid #2. Since we use --pid=ln
				// when operating inside container, our process will die if lightningd dies anyway
				// currently, so this is less of an issue in practice.
				lc := prometheus.NewProcessCollector(s.pid, namespace)
				prometheus.MustRegister(lc)
				registeredLn = true
				log.Printf("Registered ProcessCollector for lightningd pid %d in namespace %s\n", s.pid, namespace)
			}
		}
		time.Sleep(time.Minute)
	}
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

// newNodeHandler returns a new http handler for the /node page.
func newNodeHandler(s *state) (*nodeHandler, error) {
	b, err := getFile("node.tmpl")
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("node").Parse(string(b))
	if err != nil {
		return nil, err
	}
	return &nodeHandler{tmpl: *tmpl, s: s}, nil
}

// registerMetrics registers the Prometheus monitoring metrics.
func (h indexHandler) registerMetrics() {
	for _, m := range h.s.gauges {
		prometheus.MustRegister(m)
	}
	for _, c := range h.s.callCounters {
		prometheus.MustRegister(c)
	}
	for _, m := range h.s.counterVecs {
		prometheus.MustRegister(m)
	}
	for _, m := range h.s.gaugeVecs {
		prometheus.MustRegister(m)
	}
}

// ServeHTTP serves the index page.
func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v] HTTP %s %s\n", r.RemoteAddr, r.Method, r.URL)
	h.s.counterVecs["http_calls"].With(prometheus.Labels{"call": "index"}).Inc()
	data := struct {
		IsRunning         bool
		MonVersion        string
		Alias             alias
		Info              getInfoResponse
		NumNodes          int
		NumChannels       int
		ChannelCandidates candidates
		Peers             nodes
		Invoices          []invoiceListing
		Payments          payments
	}{
		IsRunning:         h.s.IsRunning(),
		MonVersion:        h.s.MonVersion,
		Alias:             h.s.Alias,
		Info:              h.s.Info,
		NumNodes:          len(h.s.Nodes),
		NumChannels:       len(h.s.Channels),
		ChannelCandidates: h.s.Nodes.ChannelCandidates(),
		Peers:             h.s.Nodes.Peers(),
		Invoices:          h.s.Invoices,
		Payments:          h.s.Payments,
	}
	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Printf("Failed to execute template: %v\n", err)
		return
	}
}

// ServeHTTP serves /node requests.
func (h nodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nid := vars["id"]

	log.Printf("[%v] nodeHandler serving HTTP for %s %s\n", r.RemoteAddr, r.Method, r.URL)
	h.s.counterVecs["http_calls"].With(prometheus.Labels{"call": "node"}).Inc()
	if len(nid) != 66 {
		log.Printf("Serving 400 for HTTP %s %q\n", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	n, exists := h.s.Nodes[nodeId(nid)]
	if !exists {
		log.Printf("No such node %v, serving 404 for HTTP %s %q\n", nid, r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Bad Request")
		return
	}
	if err := h.tmpl.Execute(w, n); err != nil {
		http.Error(w, "Well, that's embarrassing. Please try again later.", http.StatusInternalServerError)
		log.Printf("Failed to execute template: %v\n", err)
		return
	}
}

// ServeHTTP serves /cmd requests.
//
// TODO: if more DoS resistance is needed, could cap per-RemoteAddr requests per second/minute.
// TODO: allow connecting to new peers via lightning-cli connect?
func (h cmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	log.Printf("[%s] cmdHandler serving HTTP for %s %s\n", r.RemoteAddr, r.Method, r.URL)
	h.s.counterVecs["http_calls"].With(prometheus.Labels{"call": "cmd"}).Inc()
	if action != "invoice" {
		log.Printf("Bad action=%q, serving 400 for HTTP %s %q\n", action, r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	h.s.counterVecs["http_calls"].With(prometheus.Labels{"call": "cmd_invoice"}).Inc()

	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Serving 400 for HTTP %s %q: %q\n", r.Method, r.URL.Path, string(b))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	var ir invoiceRequest
	if err := json.Unmarshal(b, &ir); err != nil {
		log.Printf("Failed to unmarshal as cmdRequest, serving 400 for HTTP %s %q: %q\n", r.Method, r.URL.Path, string(b))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	if ir.Msatoshi == msatoshi(0) {
		log.Printf("Missing %q value in request, serving 400 for HTTP %s %q: %q\n", "msatoshi", r.Method, r.URL.Path, string(b))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	if ir.Label == "" {
		log.Printf("Missing %q value in request, serving 400 for HTTP %s %q: %q\n", "label", r.Method, r.URL.Path, string(b))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	if ir.Description == "" {
		log.Printf("Missing %q value in request, serving 400 for HTTP %s %q: %q\n", "description", r.Method, r.URL.Path, string(b))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}

	c := &cli{}
	h.s.counterVecs["cli_calls"].With(prometheus.Labels{"call": "invoice"}).Inc()
	log.Printf("Creating invoice from %v..\n", ir)
	iv, err := c.Invoice(ir)
	if err != nil {
		h.s.counterVecs["cli_failures"].With(prometheus.Labels{"call": "invoice"}).Inc()
		// TODO: maybe return actual error to user?
		log.Printf("Failed to generate invoice: %v\n", err)
		log.Printf("Serving 400 for HTTP %s %q: %q\n", r.Method, r.URL.Path, string(b))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	log.Printf("lightning-cli invoice response: %+v\n", iv)

	fmt.Fprintf(w, iv.Bolt11)
	log.Printf("Wrote bolt11 as response: %s\n", iv.Bolt11)
}

// newRouter returns a new http router.
func newRouter(s *state, prefix string) (*mux.Router, error) {
	ih, err := newIndexHandler(s)
	if err != nil {
		return nil, err
	}
	ih.registerMetrics()

	nh, err := newNodeHandler(s)
	if err != nil {
		return nil, err
	}

	ch := cmdHandler{s}

	r := mux.NewRouter()
	r.Handle("/", ih).Methods("GET")
	r.Handle("/node/{id:[a-f0-9]+}", nh).Methods("GET")
	r.Handle("/cmd/{action:[a-z]+}", ch).Methods("POST")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	if prefix != "" {
		log.Printf("Serving resources with prefix %q..\n", prefix)
		sr := r.PathPrefix(prefix).Subrouter()
		sr.Handle("/", ih).Methods("GET")
		sr.Handle("/node/{id:[a-f0-9]+}", nh).Methods("GET")
		sr.Handle("/cmd/{action:[a-z]+}", ch).Methods("POST")
		sr.Handle("/metrics", promhttp.Handler()).Methods("GET")
	}

	return r, nil
}

// newState returns a new state.
func newState() *state {
	s := state{
		Nodes: allNodes{},
		gauges: map[string]prometheus.Gauge{
			"blockheight": prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: counterPrefix,
				Name:      "blockheight",
				Help:      "Reported blockheight from lightningd.",
			}),
			"running": prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: counterPrefix,
				Name:      "running",
				Help:      "Whether lightningd process is running (1) or not (0).",
			}),
			"num_channels": prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "num_channels",
					Help:      "Number of Lightning channels this node knows about.",
				},
			),
			"num_nodes": prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "num_nodes",
					Help:      "Number of Lightning nodes known by this node.",
				},
			),
			"total_funds": prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "total_funds",
					Help:      "Sum of all funds available for opening channels.",
				},
			),
		},
		counterVecs: map[string]*prometheus.CounterVec{
			"aliases": prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: counterPrefix,
					Name:      "aliases",
					Help:      "Alias for each node id.",
				},
				[]string{"node_id", "alias"},
			),
			"info": prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: counterPrefix,
					Name:      "info",
					Help:      "Info of lightningd and lnmon version.",
				},
				[]string{"lnmon_version", "lightningd_version", "network"},
			),
			"http_calls": prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: counterPrefix,
					Name:      "http_calls_total",
					Help:      "Number of calls to http handlers.",
				},
				[]string{"call"},
			),
			"cli_calls": prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: counterPrefix,
					Name:      "cli_calls_total",
					Help:      "Number of calls to cli command via CLI.",
				},
				[]string{"call"},
			),
			"cli_failures": prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: counterPrefix,
					Name:      "cli_failures_total",
					Help:      "Number of failed calls to cli command via CLI.",
				},
				[]string{"call"},
			),
		},
		gaugeVecs: map[string]*prometheus.GaugeVec{
			"num_peers": prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "num_peers",
					Help:      "Number of Lightning peers of this node.",
				},
				[]string{"connected"},
			),
			"our_channels": prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "our_channels",
					Help:      "Number of channels per state to and from our node.",
				},
				[]string{"state"},
			),
			"channel_capacities_msatoshi": prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "channel_capacities_msatoshi",
					Help:      "Capacity of channels in millisatoshi by name and channel state.",
				},
				[]string{"node_id", "state"},
			),
			"channel_balances_msatoshi": prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: counterPrefix,
					Name:      "channel_balances_msatoshi",
					Help:      "Balance to us of channels in millisatoshi by name and channel state.",
				},
				[]string{"node_id", "state", "direction"},
			),
		},
	}
	return &s
}

func main() {
	log.Printf("lnmon version %q starting..\n", lnmonVersion)
	s := newState()
	router, err := newRouter(s, httpPrefix)
	if err != nil {
		log.Fatalf("Failed to create router: %v\n", err)
	}
	http.Handle("/", router)
	go refresh(s)

	if addr == "" {
		addr = ":8380"
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
