// lnmon.go is a simple wrapper around c-lightning's lightning-cli for monitoring state.
package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type (
	cli          struct{}
	bitcoindInfo struct {
		pid  int
		args []string
	}
)

var btcInfo bitcoindInfo

func execCmd(cmd string, arg ...string) (string, error) {
	c := exec.Command(cmd, arg...)
	out := bytes.Buffer{}
	c.Stdout = &out
	if err := c.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (c cli) exec(cmd string) (string, error) {
	return execCmd("lightning-cli", cmd)
}

func (c cli) GetInfo() (string, error) {
	return c.exec("getinfo")
}

func (c cli) GetPeers() (string, error) {
	return "", nil
}

func (c cli) GetNodes() (string, error) {
	return "", nil
}

func updateBtcInfo() error {
	s, err := execCmd("pgrep", "-a", "bitcoind")
	if err != nil {
		return err
	}
	log.Fatalf("FIXMEH: bitcoind info: %q\n", s)
	parts := strings.Split(s, " ")

	btcInfo.pid = -1
	btcInfo.args = parts[1:]
	return nil
}

func main() {
	if err := updateBtcInfo(); err != nil {
		log.Fatalf("Failed to get bitcoind info: %v\n", err)
	}

	log.Printf("Running lightning-cli getinfo..\n")
	c := cli{}
	info, err := c.GetInfo()
	if err != nil {
		log.Fatalf("Failed to get info: %v\n", err)
	}
	log.Printf("lightningd info: %s\n", info)
	// lightning-cli getinfo
	// lightning-cli getpeers
	// lightning-cli getnodes
	// lightning-cli getchannels

	// lightning-cli listfunds
	// lightning-cli listinvoice
	// lightning-cli listpayments
}
