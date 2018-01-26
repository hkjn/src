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
	cli          struct{}
	bitcoindInfo struct {
		pid  int
		args []string
	}
	addressInfo struct {
		AddressType string `json:"type"`
		Address     string `json:"address"`
		Port        int    `json:"port"`
	}
	// getinfo is the format of the getinfo response from lightning-cli.
	getinfo struct {
		NodeId      string        `json:"id"`
		Port        int           `json:"port"`
		Address     []addressInfo `json:"address"`
		Version     string        `json:"version"`
		Blockheight int           `json:"blockheight"`
	}

	peer struct {
		State     string   `json:"state"`
		Netaddr   []string `json:"netaddr"`
		PeerId    string   `json:"peerid"`
		Connected bool     `json:"connected"`
		Owner     string   `json:"owner"`
	}
	// listpeers is the format of the listpeers response from lightning-cli.
	listpeers struct {
		Peers []peer `json:"peers"`
	}
	lightningdInfo struct {
		pid   int
		args  []string
		info  getinfo
		peers listpeers
	}
	info struct {
		bitcoind   bitcoindInfo
		lightningd lightningdInfo
	}
)

var (
	state info
	// btcInfo bitcoindInfo
	// lnInfo  lightningdInfo
)

func (info bitcoindInfo) String() string {
	if info.pid == 0 {
		return "bitcoindInfo{not running}"
	} else {
		return fmt.Sprintf("bitcoindInfo{pid: %d, args: %q}", info.pid, strings.Join(info.args, " "))
	}
}

func (info lightningdInfo) String() string {
	if info.pid == 0 {
		return "lightningdInfo{not running}"
	} else {
		return fmt.Sprintf("lightningdInfo{pid: %d, args: %q}", info.pid, strings.Join(info.args, " "))
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

func (c cli) GetInfo() (*getinfo, error) {
	infostring, err := c.exec("getinfo")
	if err != nil {
		return nil, err
	}
	var info getinfo
	if err := json.Unmarshal([]byte(infostring), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ListPeers returns the lightning-cli response to listpeers.
func (c cli) ListPeers() (*listpeers, error) {
	respstring, err := c.exec("listpeers")
	if err != nil {
		return nil, err
	}
	var peers listpeers
	if err := json.Unmarshal([]byte(respstring), &peers); err != nil {
		return nil, err
	}
	return &peers, nil
}

func (c cli) GetNodes() (string, error) {
	return "", nil
}

func getBtcInfo() (*bitcoindInfo, error) {
	s, err := execCmd("pgrep", "-a", "bitcoind")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(s, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success, seems like there's still >= 1 parts..
	if len(parts) < 1 || len(parts[0]) == 0 {
		return &bitcoindInfo{}, nil
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	info := bitcoindInfo{
		pid:  pid,
		args: []string{},
	}
	for _, arg := range parts[1:] {
		info.args = append(info.args, arg)
	}
	// log.Println(btcInfo)
	return &info, nil
}

func getLnInfo() (*lightningdInfo, error) {
	s, err := execCmd("pgrep", "-a", "lightningd")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(s, " ")
	// Note: seems to get >= 1 parts even if pgrep returns non-success.
	if len(parts) < 1 || len(parts[0]) == 0 {
		return &lightningdInfo{}, nil
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	info := lightningdInfo{
		pid:  pid,
		args: []string{},
	}
	for _, arg := range parts[1:] {
		info.args = append(info.args, arg)
	}
	c := cli{}
	getinfo, err := c.GetInfo()
	if err != nil {
		return nil, err
	}
	info.info = *getinfo
	log.Printf("lightningd getinfo response: %+v\n", getinfo)

	peers, err := c.ListPeers()
	if err != nil {
		return nil, err
	}
	info.peers = *peers
	log.Printf("lightningd listpeers response: %+v\n", peers)
	return &info, nil
}

func refresh() {
	for {
		btcInfo, err := getBtcInfo()
		if err != nil {
			log.Fatalf("Failed to get bitcoind info: %v\n", err)
		}
		state.bitcoind = *btcInfo

		lnInfo, err := getLnInfo()
		if err != nil {
			log.Fatalf("Failed to get lightningd info: %v\n", err)
		}
		state.lightningd = *lnInfo

		// lightning-cli getinfo
		// lightning-cli getpeers
		// lightning-cli getnodes
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
	if state.bitcoind.pid == 0 {
		s += `<p><code>bitcoind</code> is <strong>not running</strong>.</p>`
	} else {
		s += `<p><code>bitcoind</code> is <strong>running</strong>.</p>`
	}
	s += `<h2><code>lightningd</code> info</h2>`
	if state.lightningd.pid == 0 {
		s += `<p><code>lightningd</code> is <strong>not running</strong>.</p>`
	} else {
		s += `<p><code>lightningd</code> is <strong>running</strong>.</p>`
		s += fmt.Sprintf(`<p>Our id is <code>%s</code>.</p>`, state.lightningd.info.NodeId)
		s += fmt.Sprintf(`<p>Our address is is <code>%s:%d</code>.</p>`, state.lightningd.info.Address[0].Address, state.lightningd.info.Address[0].Port)
		s += fmt.Sprintf(`<p>Our version is is <code>%s</code>.</p>`, state.lightningd.info.Version)
		s += fmt.Sprintf(`<p>Our blockheight is is <code>%d</code>.</p>`, state.lightningd.info.Blockheight)

		s += fmt.Sprintf(`<h3>We have %d lightning peers</h3>`, len(state.lightningd.peers.Peers))
		for _, peer := range state.lightningd.peers.Peers {
			if len(peer.Netaddr) > 0 {
				s += fmt.Sprintf(`<p><code>%s</code> running at <code>%s</code> is in state <code>%s</code>.</p>`, peer.PeerId, peer.Netaddr[0], peer.State)
			} else {
				s += fmt.Sprintf(`<p><code>%s</code> is in state <code>%s</code>.</p>`, peer.PeerId, peer.State)
			}
		}
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
